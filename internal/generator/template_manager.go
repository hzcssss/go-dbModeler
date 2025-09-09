package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go-DBmodeler/internal/config"
	"go-DBmodeler/pkg/logger"
)

// TemplateManager 管理模板文件的存储和加载
type TemplateManager struct {
	log         *logger.Logger
	templateDir string
}

// NewTemplateManager 创建一个新的模板管理器
func NewTemplateManager(log *logger.Logger, templateDir string) *TemplateManager {
	return &TemplateManager{
		log:         log,
		templateDir: templateDir,
	}
}

// InitializeDefaultTemplates 初始化默认模板到文件系统
func (tm *TemplateManager) InitializeDefaultTemplates() error {
	// 确保模板目录存在
	if err := os.MkdirAll(tm.templateDir, 0755); err != nil {
		return fmt.Errorf("创建模板目录失败: %v", err)
	}

	// 创建默认模板文件
	defaultTemplate := config.DefaultTemplate()
	if err := tm.SaveTemplateToFile("default", defaultTemplate); err != nil {
		return fmt.Errorf("保存默认模板失败: %v", err)
	}

	tm.log.Info("默认模板已初始化到文件系统")
	return nil
}

// SaveTemplateToFile 将模板保存到文件
func (tm *TemplateManager) SaveTemplateToFile(filename, content string) error {
	// 确保文件路径在模板目录内
	fullPath := filepath.Join(tm.templateDir, filename+".tpl")

	// 安全检查：防止目录遍历攻击
	if !strings.HasPrefix(fullPath, tm.templateDir) {
		return fmt.Errorf("无效的模板文件路径")
	}

	// 确保模板目录存在
	if err := os.MkdirAll(tm.templateDir, 0755); err != nil {
		return fmt.Errorf("创建模板目录失败: %v", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("保存模板文件失败: %v", err)
	}

	return nil
}

// LoadTemplateFromFile 从文件加载模板
func (tm *TemplateManager) LoadTemplateFromFile(filename string) (string, error) {
	fullPath := filepath.Join(tm.templateDir, filename+".tpl")

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("模板文件不存在: %s", filename)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %v", err)
	}

	return string(content), nil
}

// DeleteTemplateFile 删除模板文件
func (tm *TemplateManager) DeleteTemplateFile(filename string) error {
	fullPath := filepath.Join(tm.templateDir, filename+".tpl")

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 文件不存在，无需删除
	}

	// 删除文件
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("删除模板文件失败: %v", err)
	}

	return nil
}

// GetAllTemplateFiles 获取所有模板文件
func (tm *TemplateManager) GetAllTemplateFiles() ([]string, error) {
	// 确保模板目录存在
	if err := os.MkdirAll(tm.templateDir, 0755); err != nil {
		return nil, fmt.Errorf("创建模板目录失败: %v", err)
	}

	// 读取目录内容
	files, err := ioutil.ReadDir(tm.templateDir)
	if err != nil {
		return nil, fmt.Errorf("读取模板目录失败: %v", err)
	}

	var templateNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tpl") {
			// 去除文件扩展名
			name := strings.TrimSuffix(file.Name(), ".tpl")
			templateNames = append(templateNames, name)
		}
	}

	return templateNames, nil
}

// GetTemplateDir 获取模板目录路径
func (tm *TemplateManager) GetTemplateDir() string {
	return tm.templateDir
}
