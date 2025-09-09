package generator

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"go-DBmodeler/pkg/logger"
	"strings"
)

// JavaScriptProcessor 表示JavaScript脚本处理器
type JavaScriptProcessor struct {
	log *logger.Logger
	vm  *goja.Runtime
}

// NewJavaScriptProcessor 创建一个新的JavaScript处理器
func NewJavaScriptProcessor(log *logger.Logger) *JavaScriptProcessor {
	return &JavaScriptProcessor{
		log: log,
		vm:  goja.New(),
	}
}

// Process 使用JavaScript脚本处理表结构数据和生成的TypeScript代码
func (p *JavaScriptProcessor) Process(tsCode string, script string, metadata interface{}) (string, error) {
	if strings.TrimSpace(script) == "" {
		return tsCode, nil // 如果没有脚本，直接返回原代码
	}

	// 将Go结构体转换为JSON字符串，然后传递给JavaScript环境
	var inputJSON string
	if metadata != nil {
		jsonBytes, err := json.Marshal(metadata)
		if err != nil {
			return "", fmt.Errorf("JSON序列化错误: %v", err)
		}
		inputJSON = string(jsonBytes)
	} else {
		inputJSON = "{}"
	}

	// 设置全局变量
	p.vm.Set("inputJSON", inputJSON) // 传递JSON字符串
	p.vm.Set("tsCode", tsCode)       // 同时传递生成的TypeScript代码

	// 自动解析JSON并设置input变量
	var parsedInput interface{}
	if err := json.Unmarshal([]byte(inputJSON), &parsedInput); err != nil {
		p.log.Errorf("JSON解析错误: %v", err)
		parsedInput = map[string]interface{}{}
	}
	p.vm.Set("input", parsedInput)
	p.vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			p.log.Infof("JS Console: %v", args...)
		},
		"warn": func(args ...interface{}) {
			p.log.Warnf("JS Console: %v", args...)
		},
		"error": func(args ...interface{}) {
			p.log.Errorf("JS Console: %v", args...)
		},
	})

	// 执行脚本
	_, err := p.vm.RunString(script)
	if err != nil {
		return "", fmt.Errorf("JavaScript执行错误: %v", err)
	}

	// 获取处理后的结果
	outputVal := p.vm.Get("output")
	if outputVal == nil || goja.IsUndefined(outputVal) {
		return tsCode, nil // 如果没有输出，返回原代码
	}

	output := outputVal.String()
	if output == "" {
		return tsCode, nil // 如果输出为空，返回原代码
	}

	return output, nil
}

// ProcessWithDefaultScript 使用默认脚本处理TypeScript代码
func (p *JavaScriptProcessor) ProcessWithDefaultScript(tsCode string, metadata interface{}) (string, error) {
	defaultScript := `
// 默认脚本：添加文件头注释
const header = ` + "`" + `/**
 * 自动生成的TypeScript模型
 * 生成时间: ${new Date().toLocaleString()}
 * 不要手动编辑此文件
 */

` + "`" + `;

// 处理代码
let processedCode = tsCode;

// 如果没有文件头注释，添加一个
if (!processedCode.includes('/**')) {
	processedCode = header + processedCode;
}

// 移除多余的空行
processedCode = processedCode.replace(/\n{3,}/g, '\n\n');

// 设置输出
output = processedCode;
`

	return p.Process(tsCode, defaultScript, metadata)
}
