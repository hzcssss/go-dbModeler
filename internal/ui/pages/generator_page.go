package pages

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go-DBmodeler/internal/config"
	"go-DBmodeler/internal/db/connector"
	"go-DBmodeler/internal/db/metadata"
	"go-DBmodeler/internal/generator"
	"go-DBmodeler/internal/ui/widgets"
	"go-DBmodeler/pkg/logger"
	"io"
	"strings"
)

// GeneratorPage 表示TS模型生成页面
type GeneratorPage struct {
	container       *fyne.Container
	log             *logger.Logger
	connections     []*ConnectionConfig
	processor       *metadata.Processor
	storage         *config.Storage
	templateManager *generator.TemplateManager

	// UI组件
	connectionSelect *widget.Select
	databaseSelect   *widget.Select
	tableSelect      *widget.Select
	scriptEditor     *widget.Entry
	scriptLoadBtn    *widget.Button
	generateBtn      *widget.Button
	copyBtn          *widget.Button
	saveBtn          *widget.Button
	tableView        *widgets.TableView
	codeContainer    *fyne.Container // 代码显示容器

	// 数据
	databases       []string
	tables          []string
	selectedTable   string
	generatedCode   string
	currentMetadata *connector.TableMetadata
}

// NewGeneratorPage 创建一个新的TS模型生成页面
func NewGeneratorPage(log *logger.Logger, connections []*ConnectionConfig, templateManager *generator.TemplateManager, storage *config.Storage) *fyne.Container {
	page := &GeneratorPage{
		log:             log,
		connections:     connections,
		templateManager: templateManager,
		storage:         storage,
	}

	// 构建UI并返回容器
	container := page.buildUI()

	return container
}

// buildUI 构建TS模型生成页面的UI
func (p *GeneratorPage) buildUI() *fyne.Container {
	// 创建连接选择器
	connectionNames := make([]string, 0, len(p.connections))
	for _, conn := range p.connections {
		connectionNames = append(connectionNames, conn.Name)
	}

	p.connectionSelect = widget.NewSelect(connectionNames, p.onConnectionSelected)
	p.connectionSelect.PlaceHolder = "选择数据库连接"

	// 创建数据库选择器
	p.databaseSelect = widget.NewSelect([]string{}, p.onDatabaseSelected)
	p.databaseSelect.PlaceHolder = "选择数据库"
	p.databaseSelect.Disable()

	// 创建表选择器
	p.tableSelect = widget.NewSelect([]string{}, p.onTableSelected)
	p.tableSelect.PlaceHolder = "选择表"
	p.tableSelect.Disable()

	// 使用默认模板，不再需要模板选择器

	// 创建代码容器 - 使用更大的div块来展示代码
	p.codeContainer = container.NewVBox()

	// 创建JavaScript脚本编辑器 - 增大显示区域
	p.scriptEditor = widget.NewMultiLineEntry()
	p.scriptEditor.SetPlaceHolder("输入JavaScript脚本来处理表结构数据\n\n" +
		"// 脚本编写指南：\n" +
		"// 1. 输入变量 input 是一个JSON对象，包含表结构信息\n" +
		"//    input.TableName: 表名\n" +
		"//    input.Fields: 字段数组，每个字段包含 Name, TsType, Comment 等属性\n" +
		"// 2. 必须设置输出变量 output 作为字符串，包含生成的TypeScript代码\n" +
		"// 3. 可以使用 console.log() 进行调试\n" +
		"// 4. 示例：\n" +
		"//    let result = `export class ${input.TableName} {}`;\n" +
		"//    for (const field of input.Fields) {\n" +
		"//      result += `\\n  ${field.Name}: ${field.TsType};`;\n" +
		"//    }\n" +
		"//    output = result + '\\n}';\n")
	p.scriptEditor.Wrapping = fyne.TextWrapOff
	p.scriptEditor.SetMinRowsVisible(10) // 设置最小行数，确保足够的编辑空间

	// 创建脚本加载按钮
	p.scriptLoadBtn = widget.NewButton("加载脚本文件", p.onScriptLoadClicked)

	// 获取所有脚本名称
	scriptNames := make([]string, 0, len(p.storage.GetScripts()))
	for name := range p.storage.GetScripts() {
		scriptNames = append(scriptNames, name)
	}

	// 创建选项卡容器
	tabs := container.NewAppTabs(
		container.NewTabItem("自定义脚本", container.NewBorder(
			container.NewHBox(
				widget.NewLabel("JavaScript脚本:"),
				p.scriptLoadBtn,
			),
			nil,
			nil,
			nil,
			p.scriptEditor,
		)),
		container.NewTabItem("常用脚本", container.NewVBox(
			widget.NewLabel("选择预设脚本:"),
			widget.NewSelect(scriptNames, func(selected string) {
				p.loadCommonScript(selected)
			}),
			widget.NewLabel("\n脚本说明:"),
			widget.NewLabel("• 小驼峰命名转换: 将字段名转换为小驼峰命名\n• 添加文件头注释: 自动添加文件头注释\n• 格式化代码: 格式化生成的TypeScript代码\n• 添加类型导入: 自动添加必要的类型导入"),
		)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// 创建表视图
	p.tableView = widgets.NewTableView(&connector.TableMetadata{
		Name:   "",
		Fields: []connector.FieldInfo{},
	})

	// 创建按钮
	p.generateBtn = widget.NewButton("生成TS模型", p.onGenerateClicked)
	p.copyBtn = widget.NewButton("复制到剪贴板", p.onCopyClicked)
	p.saveBtn = widget.NewButton("保存为文件", p.onSaveClicked)

	// 禁用按钮
	p.generateBtn.Disable()
	p.copyBtn.Disable()
	p.saveBtn.Disable()

	// 创建左侧面板 - 简化版
	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("连接配置", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("连接:"),
			nil,
			p.connectionSelect,
		),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("数据库:"),
			nil,
			p.databaseSelect,
		),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("表:"),
			nil,
			p.tableSelect,
		),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabel(""),
			p.generateBtn,
		),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("脚本处理", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		tabs,
	)

	// 创建代码显示区域 - 优化布局，占满窗口
	codeHeader := container.NewHBox(
		widget.NewLabelWithStyle("📋 TypeScript 模型代码", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		p.copyBtn,
		p.saveBtn,
	)

	// 创建代码容器 - 使用更大的div块来展示代码
	p.codeContainer = container.NewVBox()

	// 创建右侧面板 - 只显示代码，移除结构选项卡
	rightPanel := container.NewBorder(
		codeHeader,
		nil,
		nil,
		nil,
		p.codeContainer, // 直接使用代码容器，不需要额外的滚动容器
	)

	// 创建分割布局 - 使用固定大小800x600
	split := container.NewHSplit(
		container.NewVScroll(
			container.NewPadded(leftPanel),
		),
		container.NewBorder(
			nil,
			nil,
			nil,
			nil,
			rightPanel,
		),
	)
	split.Offset = 0.15 // 左侧占15%，右侧占85%，适合固定大小显示

	// 创建主容器 - 使用带边距的布局
	p.container = container.NewPadded(split)

	return p.container
}

// onConnectionSelected 处理连接选择事件
func (p *GeneratorPage) onConnectionSelected(connName string) {
	// 查找选中的连接配置
	var selectedConn *ConnectionConfig
	for _, conn := range p.connections {
		if conn.Name == connName {
			selectedConn = conn
			break
		}
	}

	if selectedConn == nil {
		return
	}

	// 关闭之前的连接
	if p.processor != nil {
		p.processor.Close()
	}

	// 创建连接器
	conn, err := connector.NewConnector(&connector.ConnectionConfig{
		Type:     selectedConn.Type,
		Host:     selectedConn.Host,
		Port:     selectedConn.Port,
		Username: selectedConn.Username,
		Password: selectedConn.Password,
		Database: selectedConn.Database,
	})

	if err != nil {
		p.log.Errorf("创建连接器失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 连接数据库
	_, err = conn.Connect()
	if err != nil {
		p.log.Errorf("连接数据库失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 创建元数据处理器
	p.processor = metadata.NewProcessor(conn)

	// 获取数据库列表
	databases, err := p.processor.GetDatabases()
	if err != nil {
		p.log.Errorf("获取数据库列表失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 更新数据库选择器
	p.databases = databases
	p.databaseSelect.Options = databases
	p.databaseSelect.Enable()
	p.databaseSelect.Refresh()

	// 清空表选择器
	p.tables = []string{}
	p.tableSelect.Options = []string{}
	p.tableSelect.Disable()
	p.tableSelect.Refresh()

	// 禁用生成按钮
	p.generateBtn.Disable()
}

// onDatabaseSelected 处理数据库选择事件
func (p *GeneratorPage) onDatabaseSelected(dbName string) {
	if p.processor == nil {
		return
	}

	// 获取表列表
	tables, err := p.processor.GetTables(dbName)
	if err != nil {
		p.log.Errorf("获取表列表失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 更新表选择器
	p.tables = tables
	p.tableSelect.Options = tables
	p.tableSelect.Enable()
	p.tableSelect.Refresh()

	// 清空选择
	p.tableSelect.SetSelected("")
	p.selectedTable = ""
	p.generateBtn.Disable()
}

// onTableSelected 处理表选择事件
func (p *GeneratorPage) onTableSelected(tableName string) {
	p.selectedTable = tableName
	if tableName != "" {
		p.generateBtn.Enable()
	} else {
		p.generateBtn.Disable()
	}
}

// 删除不再需要的函数

// onGenerateClicked 处理生成按钮点击事件
func (p *GeneratorPage) onGenerateClicked() {
	p.generateCode()
}

// generateCode 生成TypeScript代码
func (p *GeneratorPage) generateCode() {
	if p.processor == nil || p.selectedTable == "" || p.databaseSelect.Selected == "" {
		return
	}

	// 获取表元数据
	metadata, err := p.processor.GetTableMetadata(p.databaseSelect.Selected, p.selectedTable)
	if err != nil {
		p.log.Errorf("获取表元数据失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 保存元数据用于表视图
	p.currentMetadata = metadata

	// 更新表视图
	p.tableView.SetMetadata(metadata)

	// 始终使用默认模板
	templateStr := generator.DefaultTemplate()

	// 创建生成器
	gen, err := generator.NewGenerator(p.connectionSelect.Selected, templateStr, p.log)
	if err != nil {
		p.log.Errorf("创建生成器失败: %v", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 设置JavaScript脚本（每次生成时都重新设置，确保实时更新）
	if p.scriptEditor.Text != "" {
		gen.SetScript(p.scriptEditor.Text)
	}

	// 生成代码
	code, err := gen.Generate(metadata)
	if err != nil {
		p.log.Errorf("生成代码失败: %v", err)
		// 检查是否是JavaScript脚本执行错误
		if strings.Contains(err.Error(), "JavaScript执行错误") {
			dialog.ShowError(fmt.Errorf("JavaScript脚本执行失败: %v\n请检查脚本语法是否正确", err), fyne.CurrentApp().Driver().AllWindows()[0])
		} else {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		}
		return
	}

	// 保存生成的代码
	p.generatedCode = code

	// 清空代码容器并添加语法高亮显示
	p.codeContainer.Objects = nil

	// 创建语法高亮器并添加代码
	highlighter := widgets.NewSyntaxHighlighter()
	highlightedCode := highlighter.HighlightTypeScript(code)
	p.codeContainer.Add(highlightedCode)
	p.codeContainer.Refresh()

	// 启用复制和保存按钮
	p.copyBtn.Enable()
	p.saveBtn.Enable()
}

// onCopyClicked 处理复制按钮点击事件
func (p *GeneratorPage) onCopyClicked() {
	if p.generatedCode == "" {
		return
	}

	// 复制到剪贴板
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	w.Clipboard().SetContent(p.generatedCode)

	dialog.ShowInformation("复制成功", "TypeScript模型已复制到剪贴板", w)
}

// onScriptLoadClicked 处理脚本加载按钮点击事件
func (p *GeneratorPage) onScriptLoadClicked() {
	// 创建文件打开对话框
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			p.log.Errorf("文件打开对话框错误: %v", err)
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return // 用户取消了操作
		}

		defer reader.Close()

		// 读取文件内容
		content, err := io.ReadAll(reader)
		if err != nil {
			p.log.Errorf("读取脚本文件失败: %v", err)
			dialog.ShowError(err, w)
			return
		}

		// 设置脚本编辑器内容
		p.scriptEditor.SetText(string(content))
		p.log.Infof("已加载脚本文件: %s", reader.URI().Path())
		dialog.ShowInformation("成功", "脚本文件加载成功", w)
	}, w)
}

// onSaveClicked 处理保存按钮点击事件
func (p *GeneratorPage) onSaveClicked() {
	if p.generatedCode == "" {
		return
	}

	// 创建保存对话框
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if writer == nil {
			return
		}

		// 写入文件
		_, err = writer.Write([]byte(p.generatedCode))
		writer.Close()

		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("保存成功", "TypeScript模型已保存到文件", w)
	}, w)

	// 设置默认文件名
	saveDialog.SetFileName(p.selectedTable + ".ts")

	// 注意：在某些Fyne版本中，SetLocation需要ListableURI，而不是普通URI
	// 这里我们只设置文件名，不设置位置

	// 显示对话框
	saveDialog.Show()
}

// loadCommonScript 加载常用脚本
func (p *GeneratorPage) loadCommonScript(scriptName string) {
	if scriptName == "" {
		return
	}

	// 从存储中获取脚本内容
	scripts := p.storage.GetScripts()
	scriptContent, exists := scripts[scriptName]
	if !exists {
		p.log.Errorf("脚本不存在: %s", scriptName)
		dialog.ShowError(fmt.Errorf("脚本 '%s' 不存在", scriptName), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	// 设置脚本编辑器内容
	p.scriptEditor.SetText(scriptContent)

	p.log.Infof("已加载常用脚本: %s", scriptName)
	dialog.ShowInformation("成功", fmt.Sprintf("已加载脚本: %s", scriptName), fyne.CurrentApp().Driver().AllWindows()[0])
}
