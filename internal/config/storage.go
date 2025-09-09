package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConnectionConfig represents database connection configuration
type ConnectionConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database"`
}

// AppConfig represents application configuration
type AppConfig struct {
	Connections []ConnectionConfig `json:"connections"`
	Templates   map[string]string  `json:"templates"`
	Scripts     map[string]string  `json:"scripts"`
}

// Storage represents configuration storage
type Storage struct {
	config     AppConfig
	configDir  string
	configFile string
}

// NewStorage creates a new configuration storage
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".godbmodeler")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	configFile := filepath.Join(configDir, "config.json")

	storage := &Storage{
		configDir:  configDir,
		configFile: configFile,
		config: AppConfig{
			Connections: make([]ConnectionConfig, 0),
			Templates:   map[string]string{"default": DefaultTemplate()},
			Scripts:     map[string]string{},
		},
	}

	storage.config.Scripts["camelCase"] = DefaultCamelCaseScript()
	storage.config.Scripts["addHeader"] = DefaultHeaderScript()
	storage.config.Scripts["formatCode"] = DefaultFormatScript()
	storage.config.Scripts["addImports"] = DefaultImportScript()

	if err := storage.Load(); err != nil {
		if os.IsNotExist(err) {
			if err := storage.Save(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return storage, nil
}

// Load loads configuration from file
func (s *Storage) Load() error {
	data, err := os.ReadFile(s.configFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &s.config); err != nil {
		return err
	}

	return nil
}

// Save saves configuration to file
func (s *Storage) Save() error {
	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.configFile, data, 0644); err != nil {
		return err
	}

	return nil
}

// GetConnections gets all connection configurations
func (s *Storage) GetConnections() []ConnectionConfig {
	return s.config.Connections
}

// AddConnection adds a connection configuration
func (s *Storage) AddConnection(config ConnectionConfig) error {
	for _, conn := range s.config.Connections {
		if conn.Name == config.Name {
			return fmt.Errorf("connection name '%s' already exists", config.Name)
		}
	}

	encryptedConfig, err := EncryptConnectionPassword(config)
	if err != nil {
		return err
	}

	s.config.Connections = append(s.config.Connections, encryptedConfig)
	return s.Save()
}

// UpdateConnection updates a connection configuration
func (s *Storage) UpdateConnection(config ConnectionConfig) error {
	for i, conn := range s.config.Connections {
		if conn.Name == config.Name {
			encryptedConfig, err := EncryptConnectionPassword(config)
			if err != nil {
				return err
			}

			s.config.Connections[i] = encryptedConfig
			return s.Save()
		}
	}

	return fmt.Errorf("connection '%s' does not exist", config.Name)
}

// DeleteConnection deletes a connection configuration
func (s *Storage) DeleteConnection(name string) error {
	for i, conn := range s.config.Connections {
		if conn.Name == name {
			s.config.Connections = append(s.config.Connections[:i], s.config.Connections[i+1:]...)
			return s.Save()
		}
	}

	return fmt.Errorf("connection '%s' does not exist", name)
}

// GetTemplates gets all templates
func (s *Storage) GetTemplates() map[string]string {
	return s.config.Templates
}

// GetTemplate gets a template
func (s *Storage) GetTemplate(name string) (string, error) {
	if template, ok := s.config.Templates[name]; ok {
		return template, nil
	}

	return "", fmt.Errorf("template '%s' does not exist", name)
}

// SetTemplate sets a template
func (s *Storage) SetTemplate(name, template string) error {
	s.config.Templates[name] = template
	return s.Save()
}

// DeleteTemplate deletes a template
func (s *Storage) DeleteTemplate(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot delete default template")
	}

	if _, ok := s.config.Templates[name]; !ok {
		return fmt.Errorf("template '%s' does not exist", name)
	}

	delete(s.config.Templates, name)
	return s.Save()
}

// GetScripts gets all scripts
func (s *Storage) GetScripts() map[string]string {
	// 合并存储的脚本和导入的脚本
	allScripts := make(map[string]string)

	// 添加存储的脚本
	for name, script := range s.config.Scripts {
		allScripts[name] = script
	}

	// 添加导入的脚本
	importedScripts := s.GetImportedScripts()
	for name, script := range importedScripts {
		allScripts[name] = script
	}

	return allScripts
}

// GetImportedScripts gets scripts from imported directory
func (s *Storage) GetImportedScripts() map[string]string {
	importedScripts := make(map[string]string)

	// 创建导入目录
	importDir := filepath.Join(s.configDir, "scripts", "imported")
	if err := os.MkdirAll(importDir, 0755); err != nil {
		return importedScripts
	}

	// 读取所有 .js 文件
	files, err := os.ReadDir(importDir)
	if err != nil {
		return importedScripts
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".js" {
			filePath := filepath.Join(importDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err == nil {
				scriptName := file.Name()[:len(file.Name())-3] // 移除 .js 扩展名
				importedScripts[scriptName] = string(content)
			}
		}
	}

	return importedScripts
}

// GetScript gets a script
func (s *Storage) GetScript(name string) (string, error) {
	if script, ok := s.config.Scripts[name]; ok {
		return script, nil
	}

	return "", fmt.Errorf("script '%s' does not exist", name)
}

// SetScript sets a script
func (s *Storage) SetScript(name, script string) error {
	s.config.Scripts[name] = script
	return s.Save()
}

// DeleteScript deletes a script
func (s *Storage) DeleteScript(name string) error {
	defaultScripts := []string{"camelCase", "addHeader", "formatCode", "addImports"}
	for _, defaultName := range defaultScripts {
		if name == defaultName {
			return fmt.Errorf("cannot delete default script")
		}
	}

	if _, ok := s.config.Scripts[name]; !ok {
		return fmt.Errorf("script '%s' does not exist", name)
	}

	delete(s.config.Scripts, name)
	return s.Save()
}

// GetConfigDir returns the configuration directory path
func (s *Storage) GetConfigDir() string {
	return s.configDir
}

// DefaultTemplate returns the default TypeScript template
func DefaultTemplate() string {
	return "export interface {{.TableName}} {\n{{range .Fields}}  /** {{.Comment}} */\n  {{.Name}}: {{.TsType}};\n{{end}}\n}\n"
}

// DefaultCamelCaseScript returns the default camel case conversion script
func DefaultCamelCaseScript() string {
	return `// Camel case conversion script
function toCamelCase(str) {
    str = String(str);
    return str.replace(/_([a-z])/g, (g) => g[1].toUpperCase())
              .replace(/^[A-Z]/, (g) => g.toLowerCase());
}

// input is already parsed JSON object
let result = "export interface " + input.tableName + " {\n";

// Process each field
for (const field of input.fields) {
    const fieldName = String(field.name || '');
    const camelCaseName = toCamelCase(fieldName);

    if (field.comment) {
        result += "  /** " + field.comment + " */\n";
    }
    result += "  " + camelCaseName + ": " + field.tsType + ";\n";
}

result += "}\n";

// Set output result
output = result;`
}

// DefaultHeaderScript returns the default header comment script
func DefaultHeaderScript() string {
	return "// Header comment script\n" +
		"// input is already parsed JSON object\n\n" +
		"// Generate header comment\n" +
		"const header = `/**\n" +
		" * ${input.tableName} model\n" +
		" * Generated at: ${new Date().toLocaleString()}\n" +
		" * Table: ${input.tableName}\n" +
		" * Field count: ${input.fields.length}\n" +
		" * Do not edit this file manually\n" +
		" */\n\n`;\n\n" +
		"// Get current generated code\n" +
		"let result = tsCode;\n\n" +
		"// If no header comment, add one\n" +
		"if (!result.includes('/**')) {\n" +
		"    result = header + result;\n" +
		"}\n\n" +
		"// Set output result\n" +
		"output = result;"
}

// DefaultFormatScript returns the default code formatting script
func DefaultFormatScript() string {
	return `// Code formatting script
// input is already parsed JSON object

// Get current generated code
let result = tsCode;

// Format code - add proper indentation and line breaks
result = result
    .replace(/\{\s*\n/g, '{\n  ')
    .replace(/;\s*\n/g, ';\n  ')
    .replace(/\n\s*\}/g, '\n}')
    .replace(/\n\s*\n/g, '\n')
    .trim();

// Ensure there's a newline at the end
if (!result.endsWith('\n')) {
    result += '\n';
}

// Set output result
output = result;`
}

// DefaultImportScript returns the default type import script
func DefaultImportScript() string {
	return "// Type import script\n" +
		"// input is already parsed JSON object\n\n" +
		"// Get current generated code\n" +
		"let result = tsCode;\n\n" +
		"// Collect types that need to be imported\n" +
		"const importTypes = new Set();\n\n" +
		"// Check field types to determine what needs to be imported\n" +
		"for (const field of input.fields) {\n" +
		"    const type = field.tsType.toLowerCase();\n" +
		"    if (type.includes('date') || type.includes('time')) {\n" +
		"        importTypes.add('Date');\n" +
		"    }\n" +
		"}\n\n" +
		"// If there are types to import, add import statements\n" +
		"if (importTypes.size > 0) {\n" +
		"    let imports = '// Auto-generated import statements\\n';\n" +
		"    for (const type of importTypes) {\n" +
		"        imports += '// import { ' + type + ' } from \\'./types\\'\\n';\n" +
		"    }\n" +
		"    imports += '\\n';\n\n" +
		"    // Add import statements at the beginning of the code\n" +
		"    result = imports + result;\n" +
		"}\n\n" +
		"// Set output result\n" +
		"output = result;"
}
