package pages

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"go-DBmodeler/internal/config"
	"go-DBmodeler/pkg/logger"
	"os"
	"path/filepath"
)

// ScriptManagerPage 表示脚本管理页面
type ScriptManagerPage struct {
	container *fyne.Container
	log       *logger.Logger
	storage   *config.Storage

	// UI组件
	scriptList  *widget.List
	addBtn      *widget.Button
	editBtn     *widget.Button
	deleteBtn   *widget.Button
	importBtn   *widget.Button
	previewArea *widget.Entry

	// 数据
	scripts        map[string]string
	scriptNames    []string
	selectedScript string
}

// NewScriptManagerPage 创建一个新的脚本管理页面
func NewScriptManagerPage(log *logger.Logger, storage *config.Storage) *fyne.Container {
	page := &ScriptManagerPage{
		log:     log,
		storage: storage,
		scripts: storage.GetScripts(),
	}

	// 构建UI并返回容器
	container := page.buildUI()

	return container
}

// buildUI 构建脚本管理页面的UI
func (p *ScriptManagerPage) buildUI() *fyne.Container {
	// 更新脚本名称列表
	p.updateScriptNames()

	// 创建脚本列表
	p.scriptList = widget.NewList(
		func() int {
			return len(p.scriptNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("模板脚本名称")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(p.scriptNames[i])
		},
	)

	// 设置列表选择事件
	p.scriptList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(p.scriptNames) {
			p.selectedScript = p.scriptNames[id]
			p.showScriptPreview(p.selectedScript)
			p.editBtn.Enable()
			p.deleteBtn.Enable()
		}
	}

	p.scriptList.OnUnselected = func(id widget.ListItemID) {
		p.selectedScript = ""
		p.previewArea.SetText("")
		p.editBtn.Disable()
		p.deleteBtn.Disable()
	}

	// 创建按钮
	p.addBtn = widget.NewButton("新增脚本", p.onAddClicked)
	p.editBtn = widget.NewButton("编辑脚本", p.onEditClicked)
	p.editBtn.Disable()
	p.deleteBtn = widget.NewButton("删除脚本", p.onDeleteClicked)
	p.deleteBtn.Disable()
	p.importBtn = widget.NewButton("导入脚本", p.onImportClicked)

	// 创建预览区域
	p.previewArea = widget.NewMultiLineEntry()
	p.previewArea.SetPlaceHolder("选择脚本查看预览...")
	p.previewArea.Disable()
	p.previewArea.Wrapping = fyne.TextWrapOff

	// 创建按钮容器
	buttonContainer := container.NewHBox(
		p.addBtn,
		p.editBtn,
		p.deleteBtn,
		p.importBtn,
		layout.NewSpacer(),
	)

	// 创建左侧面板
	leftPanel := container.NewBorder(
		widget.NewLabelWithStyle("📝 脚本列表", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		buttonContainer,
		nil,
		nil,
		p.scriptList,
	)

	// 创建右侧预览面板
	rightPanel := container.NewBorder(
		widget.NewLabelWithStyle("👀 脚本预览", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewVScroll(p.previewArea),
	)

	// 创建分割布局
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.3 // 左侧占30%，右侧占70%

	// 创建主容器
	p.container = container.NewPadded(split)

	return p.container
}

// updateScriptNames 更新脚本名称列表
func (p *ScriptManagerPage) updateScriptNames() {
	p.scriptNames = make([]string, 0, len(p.scripts))
	for name := range p.scripts {
		p.scriptNames = append(p.scriptNames, name)
	}
}

// showScriptPreview 显示脚本预览
func (p *ScriptManagerPage) showScriptPreview(scriptName string) {
	if content, exists := p.scripts[scriptName]; exists {
		p.previewArea.SetText(content)
	}
}

// onAddClicked 处理新增按钮点击事件
func (p *ScriptManagerPage) onAddClicked() {
	p.showScriptDialog("", "", true)
}

// onEditClicked 处理编辑按钮点击事件
func (p *ScriptManagerPage) onEditClicked() {
	if p.selectedScript == "" {
		return
	}

	content := p.scripts[p.selectedScript]
	p.showScriptDialog(p.selectedScript, content, false)
}

// onDeleteClicked 处理删除按钮点击事件
func (p *ScriptManagerPage) onDeleteClicked() {
	if p.selectedScript == "" {
		return
	}

	// 显示确认对话框
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("确认删除",
		fmt.Sprintf("确定要删除脚本 '%s' 吗？此操作不可恢复。", p.selectedScript),
		func(confirmed bool) {
			if confirmed {
				p.storage.DeleteScript(p.selectedScript)
				p.scripts = p.storage.GetScripts()
				p.updateScriptNames()
				p.scriptList.UnselectAll()
				p.scriptList.Refresh()
				p.log.Infof("已删除脚本: %s", p.selectedScript)
			}
		}, w)
}

// onImportClicked 处理导入按钮点击事件
func (p *ScriptManagerPage) onImportClicked() {
	w := fyne.CurrentApp().Driver().AllWindows()[0]

	// 创建文件选择对话框
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return // 用户取消了选择
		}
		defer reader.Close()

		// 读取文件内容
		content, err := os.ReadFile(reader.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("读取文件失败: %v", err), w)
			return
		}

		// 获取文件名（不含扩展名）
		fileName := filepath.Base(reader.URI().Name())
		if ext := filepath.Ext(fileName); ext != "" {
			fileName = fileName[:len(fileName)-len(ext)]
		}

		// 显示导入确认对话框
		p.showImportDialog(fileName, string(content))
	}, w)

	// 设置文件过滤器
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".js"}))
	fileDialog.SetFileName("script.js")
	fileDialog.Show()
}

// showImportDialog 显示导入确认对话框
func (p *ScriptManagerPage) showImportDialog(name, content string) {
	w := fyne.CurrentApp().Driver().AllWindows()[0]

	// 创建名称输入框
	nameEntry := widget.NewEntry()
	nameEntry.SetText(name)
	nameEntry.SetPlaceHolder("输入脚本名称")

	// 创建内容预览
	previewEntry := widget.NewMultiLineEntry()
	previewEntry.SetText(content)
	previewEntry.Disable()
	previewEntry.Wrapping = fyne.TextWrapOff
	previewEntry.SetMinRowsVisible(10)

	// 创建表单
	form := widget.NewForm(
		widget.NewFormItem("脚本名称", nameEntry),
		widget.NewFormItem("脚本内容预览", previewEntry),
	)

	// 创建对话框
	dialog := dialog.NewCustomConfirm("导入脚本", "导入", "取消", form,
		func(confirmed bool) {
			if confirmed {
				scriptName := nameEntry.Text
				if scriptName == "" {
					dialog.ShowError(fmt.Errorf("脚本名称不能为空"), w)
					return
				}

				// 保存脚本到导入目录
				importDir := filepath.Join(p.storage.GetConfigDir(), "scripts", "imported")
				if err := os.MkdirAll(importDir, 0755); err != nil {
					dialog.ShowError(fmt.Errorf("创建导入目录失败: %v", err), w)
					return
				}

				filePath := filepath.Join(importDir, scriptName+".js")
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					dialog.ShowError(fmt.Errorf("保存脚本文件失败: %v", err), w)
					return
				}

				// 刷新脚本列表
				p.scripts = p.storage.GetScripts()
				p.updateScriptNames()
				p.scriptList.Refresh()

				p.log.Infof("已导入脚本: %s", scriptName)
				dialog.ShowInformation("成功", fmt.Sprintf("脚本 '%s' 导入成功！\n文件已保存到: %s", scriptName, filePath), w)
			}
		}, w)

	dialog.Resize(fyne.NewSize(600, 500))
	dialog.Show()
}

// showScriptDialog 显示脚本编辑对话框
func (p *ScriptManagerPage) showScriptDialog(name, content string, isNew bool) {
	w := fyne.CurrentApp().Driver().AllWindows()[0]

	// 创建名称输入框
	nameEntry := widget.NewEntry()
	nameEntry.SetText(name)
	nameEntry.SetPlaceHolder("输入脚本名称")

	// 创建内容编辑器
	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetText(content)
	contentEntry.SetPlaceHolder("输入JavaScript脚本内容...\n\n" + p.getScriptExamples())
	contentEntry.Wrapping = fyne.TextWrapOff
	contentEntry.SetMinRowsVisible(15)

	// 创建复制按钮
	copyExampleBtn := widget.NewButton("复制示例", func() {
		exampleScript := `// 小驼峰命名转换示例
function toCamelCase(str) {
    return str.replace(/_([a-z])/g, (g) => g[1].toUpperCase())
              .replace(/^[A-Z]/, (g) => g.toLowerCase());
}

// 输入数据结构说明
// input = {
//   "tableName": "表名",
//   "fields": [
//     {
//       "name": "字段名", 
//       "tsType": "TypeScript类型",
//       "comment": "字段注释",
//       "isPrimary": false,
//       "isNullable": false
//     }
//   ]
// }

// 处理字段名并生成接口
let result = "export interface " + input.tableName + " {\n";
for (const field of input.fields) {
    const camelCaseName = toCamelCase(field.name);
    result += "  /** " + (field.comment || field.name) + " */\n";
    result += "  " + camelCaseName + (field.isNullable ? "?" : "") + ": " + field.tsType + ";\n\n";
}
result += "}\n";

// 设置输出
output = result;

// 调试信息（可选）
console.log("生成的接口代码:", result);`
		contentEntry.SetText(exampleScript)
		dialog.ShowInformation("复制成功", "完整的示例脚本已填充到内容区域", w)
	})

	// 创建复制内容按钮
	copyContentBtn := widget.NewButton("复制内容", func() {
		clipboard := w.Clipboard()
		clipboard.SetContent(contentEntry.Text)
		dialog.ShowInformation("复制成功", "脚本内容已复制到剪贴板", w)
	})

	// 创建按钮容器
	buttonContainer := container.NewHBox(
		copyExampleBtn,
		copyContentBtn,
		layout.NewSpacer(),
	)

	// 创建内容容器
	contentContainer := container.NewBorder(
		nil,
		buttonContainer,
		nil,
		nil,
		contentEntry,
	)

	// 创建表单
	form := widget.NewForm(
		widget.NewFormItem("脚本名称", nameEntry),
		widget.NewFormItem("脚本内容", contentContainer),
	)

	// 创建对话框
	dialog := dialog.NewCustomConfirm("脚本编辑器", "保存", "取消", form,
		func(confirmed bool) {
			if confirmed {
				newName := nameEntry.Text
				newContent := contentEntry.Text

				if newName == "" {
					dialog.ShowError(fmt.Errorf("脚本名称不能为空"), w)
					return
				}

				if newContent == "" {
					dialog.ShowError(fmt.Errorf("脚本内容不能为空"), w)
					return
				}

				// 保存脚本
				p.storage.SetScript(newName, newContent)
				p.scripts = p.storage.GetScripts()
				p.updateScriptNames()
				p.scriptList.Refresh()

				p.log.Infof("已保存脚本: %s", newName)
				dialog.ShowInformation("成功", "脚本保存成功", w)
			}
		}, w)

	dialog.Resize(fyne.NewSize(700, 550))
	dialog.Show()
}

// getScriptExamples 获取脚本示例
func (p *ScriptManagerPage) getScriptExamples() string {
	return `// 🎯 脚本编写指南：
// 1. 输入变量 input 是一个JSON对象，包含表结构信息
//    input.TableName: 表名
//    input.Fields: 字段数组，每个字段包含 Name, TsType, Comment 等属性
// 2. 必须设置输出变量 output 作为字符串，包含生成的TypeScript代码
// 3. 可以使用 console.log() 进行调试

// 📝 示例模板：

// 示例1: 基础类生成
let output = 'export class ' + input.TableName + ' {\\n';
for (const field of input.Fields) {
    output += '  ' + field.Name + ': ' + field.TsType + ';\\n';
}
output += '}';

// 示例2: 带注释的模型
let output = '/**\\n * ' + input.TableName + ' 实体类\\n';
if (input.Comment) {
    output += ' * ' + input.Comment + '\\n';
}
output += ' */\\nexport class ' + input.TableName + ' {\\n';
for (const field of input.Fields) {
    if (field.Comment) {
        output += '  /** ' + field.Comment + ' */\\n';
    }
    output += '  ' + field.Name + ': ' + field.TsType + ';\\n\\n';
}
output += '}';

// 示例3: 小驼峰命名转换
let output = 'export class ' + input.TableName + ' {\\n';
for (const field of input.Fields) {
    const camelCaseName = field.Name.replace(/_([a-z])/g, (g) => g[1].toUpperCase());
    output += '  ' + camelCaseName + ': ' + field.TsType + ';\\n';
}
output += '}';

// 示例4: 添加文件头注释
let output = '// 自动生成的TypeScript模型\\n';
output += '// 表名: ' + input.TableName + '\\n';
output += '// 生成时间: ' + new Date().toLocaleString() + '\\n\\n';
output += 'export class ' + input.TableName + ' {\\n';
for (const field of input.Fields) {
    output += '  ' + field.Name + ': ' + field.TsType + ';\\n';
}
output += '}';`
}
