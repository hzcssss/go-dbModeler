package generator

import (
	"bytes"
	"fmt"
	"go-DBmodeler/internal/db/connector"
	"go-DBmodeler/pkg/logger"
	"text/template"
)

// TemplateData 表示模板数据
type TemplateData struct {
	TableName string      `json:"tableName"`
	Fields    []FieldData `json:"fields"`
}

// FieldData 表示字段数据
type FieldData struct {
	Name    string `json:"name"`
	TsType  string `json:"tsType"`
	Comment string `json:"comment"`
}

// Generator 表示TypeScript模型生成器
type Generator struct {
	mapper        TypeMapper
	template      *template.Template
	log           *logger.Logger
	script        string         // JavaScript处理脚本
	scriptManager *ScriptManager // 脚本管理器
}

// NewGenerator 创建一个新的生成器
func NewGenerator(dbType string, templateStr string, log *logger.Logger) (*Generator, error) {
	// 创建类型映射器
	mapper := NewTypeMapper(dbType)

	// 解析模板
	tmpl, err := template.New("tsmodel").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	// 创建脚本管理器
	scriptManager := NewScriptManager(log, "scripts/imported")

	return &Generator{
		mapper:        mapper,
		template:      tmpl,
		log:           log,
		scriptManager: scriptManager,
	}, nil
}

// SetScript 设置JavaScript处理脚本
func (g *Generator) SetScript(script string) {
	g.script = script
}

// SetScriptFromFile 从文件设置JavaScript处理脚本
func (g *Generator) SetScriptFromFile(filename string) error {
	if g.scriptManager == nil {
		return fmt.Errorf("脚本管理器未初始化")
	}

	script, err := g.scriptManager.LoadScriptFromFile(filename)
	if err != nil {
		return err
	}

	g.script = script
	return nil
}

// GetScriptContent 获取当前设置的脚本内容
func (g *Generator) GetScriptContent() (string, error) {
	if g.script == "" {
		return "", fmt.Errorf("没有设置脚本")
	}
	return g.script, nil
}

// Generate 生成TypeScript模型
func (g *Generator) Generate(metadata *connector.TableMetadata) (string, error) {
	// 准备模板数据
	data := TemplateData{
		TableName: metadata.Name,
		Fields:    make([]FieldData, 0, len(metadata.Fields)),
	}

	// 转换字段数据
	for _, field := range metadata.Fields {
		data.Fields = append(data.Fields, FieldData{
			Name:    field.Name,
			TsType:  g.mapper.Map(field.Type),
			Comment: field.Comment,
		})
	}

	// 执行模板
	var buf bytes.Buffer
	if err := g.template.Execute(&buf, data); err != nil {
		return "", err
	}

	tsCode := buf.String()

	// 如果有JavaScript脚本，进行处理
	if g.script != "" {
		processor := NewJavaScriptProcessor(g.log)
		// 将表结构数据和生成的TypeScript代码传递给处理器
		processedCode, err := processor.Process(tsCode, g.script, data)
		if err != nil {
			g.log.Warnf("JavaScript处理失败: %v", err)
			return tsCode, err // 返回错误，让调用者处理
		}
		return processedCode, nil
	}

	return tsCode, nil
}

// DefaultTemplate 返回默认的TypeScript模板
func DefaultTemplate() string {
	return `export interface {{.TableName}} {
{{range .Fields}}  /** {{.Comment}} */
  {{.Name}}: {{.TsType}};
{{end}}
}
`
}
