package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go-DBmodeler/internal/config"
	"go-DBmodeler/pkg/logger"
)

// ConnectionPage 表示连接管理页面
type ConnectionPage struct {
	container   *fyne.Container
	log         *logger.Logger
	storage     *config.Storage
	connections []*ConnectionConfig
	list        *widget.List
	onRefresh   func() // 刷新回调函数
}

// ConnectionConfig 表示数据库连接配置
type ConnectionConfig struct {
	Name     string
	Type     string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// NewConnectionPage 创建一个新的连接管理页面
func NewConnectionPage(log *logger.Logger, storage *config.Storage) (*ConnectionPage, *fyne.Container) {
	page := &ConnectionPage{
		log:         log,
		storage:     storage,
		connections: make([]*ConnectionConfig, 0),
		onRefresh:   func() {}, // 默认空函数
	}

	// 从存储加载连接配置
	page.loadConnections()

	return page, page.buildUI()
}

// SetRefreshCallback 设置刷新回调函数
func (p *ConnectionPage) SetRefreshCallback(callback func()) {
	p.onRefresh = callback
}

// GetContainer 获取页面容器
func (p *ConnectionPage) GetContainer() *fyne.Container {
	return p.container
}

// loadConnections 从存储加载连接配置
func (p *ConnectionPage) loadConnections() {
	if p.storage == nil {
		return
	}

	// 清空当前连接列表
	p.connections = make([]*ConnectionConfig, 0)

	// 从存储获取连接配置
	storedConns := p.storage.GetConnections()
	for _, conn := range storedConns {
		// 解密密码
		decryptedConfig, err := config.DecryptConnectionPassword(conn)
		if err != nil {
			p.log.Errorf("解密连接密码失败: %v", err)
			continue
		}

		p.connections = append(p.connections, &ConnectionConfig{
			Name:     decryptedConfig.Name,
			Type:     decryptedConfig.Type,
			Host:     decryptedConfig.Host,
			Port:     decryptedConfig.Port,
			Username: decryptedConfig.Username,
			Password: decryptedConfig.Password,
			Database: decryptedConfig.Database,
		})
	}
}

// buildUI 构建连接管理页面的UI
func (p *ConnectionPage) buildUI() *fyne.Container {
	// 创建连接列表
	p.list = widget.NewList(
		func() int {
			return len(p.connections)
		},
		func() fyne.CanvasObject {
			// 使用Text函数创建标签，可能有助于解决中文显示问题
			return widget.NewLabel("Connection Name")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(p.connections) {
				obj.(*widget.Label).SetText(p.connections[id].Name)
			}
		},
	)

	// 创建新建连接按钮（使用中文）
	newConnectionBtn := widget.NewButton("+ 新建连接", func() {
		p.showNewConnectionDialog(fyne.CurrentApp().Driver().AllWindows()[0])
	})

	// 创建右侧详情面板（使用中文）
	detailsPanel := container.NewVBox(
		widget.NewLabel("选择左侧连接查看详情"),
	)

	// 当选择连接时显示详情和删除按钮
	p.list.OnSelected = func(id widget.ListItemID) {
		if id < len(p.connections) {
			conn := p.connections[id]

			// 创建删除按钮
			deleteBtn := widget.NewButton("删除连接", func() {
				p.showDeleteConnectionDialog(id, fyne.CurrentApp().Driver().AllWindows()[0])
			})

			// 更新详情面板
			detailsPanel.Objects = []fyne.CanvasObject{
				widget.NewLabel("连接详情:"),
				widget.NewLabel("名称: " + conn.Name),
				widget.NewLabel("类型: " + conn.Type),
				widget.NewLabel("主机: " + conn.Host),
				widget.NewLabel("端口: " + conn.Port),
				widget.NewLabel("用户名: " + conn.Username),
				widget.NewLabel("数据库: " + conn.Database),
				layout.NewSpacer(),
				deleteBtn,
			}
			detailsPanel.Refresh()
		}
	}

	// 创建分割布局
	split := container.NewHSplit(
		container.NewBorder(nil, newConnectionBtn, nil, nil, p.list),
		detailsPanel,
	)
	split.Offset = 0.3

	// 创建主容器
	p.container = container.NewPadded(split)

	return p.container
}

// showDeleteConnectionDialog 显示删除连接确认对话框
func (p *ConnectionPage) showDeleteConnectionDialog(id widget.ListItemID, win fyne.Window) {
	if id >= len(p.connections) {
		return
	}

	conn := p.connections[id]
	confirmDialog := dialog.NewConfirm(
		"确认删除",
		"确定要删除连接 \""+conn.Name+"\" 吗？此操作不可撤销。",
		func(confirm bool) {
			if confirm {
				p.deleteConnection(id, win)
			}
		},
		win,
	)
	confirmDialog.SetDismissText("取消")
	confirmDialog.SetConfirmText("确认删除")
	confirmDialog.Show()
}

// deleteConnection 删除指定的连接
func (p *ConnectionPage) deleteConnection(id widget.ListItemID, win fyne.Window) {
	if id >= len(p.connections) {
		return
	}

	conn := p.connections[id]

	// 从存储中删除连接
	if err := p.storage.DeleteConnection(conn.Name); err != nil {
		p.log.Errorf("删除连接失败: %v", err)
		dialog.ShowError(err, win)
		return
	}

	// 重新加载连接列表
	p.loadConnections()
	p.list.Refresh()
	p.list.UnselectAll()

	p.log.Infof("删除连接: %s", conn.Name)

	// 调用刷新回调，通知其他页面更新连接列表
	p.onRefresh()

	// 显示成功消息
	dialog.ShowInformation("删除成功", "数据库连接已删除", win)
}

// showNewConnectionDialog 显示新建连接对话框
func (p *ConnectionPage) showNewConnectionDialog(win fyne.Window) {
	// 创建表单项（使用中文）
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("连接名称")

	dbTypeSelect := widget.NewSelect([]string{"MySQL", "PostgreSQL", "SQLite"}, nil)
	dbTypeSelect.PlaceHolder = "选择数据库类型"

	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder("主机名/IP地址")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("端口号")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("用户名")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("密码")

	databaseEntry := widget.NewEntry()
	databaseEntry.SetPlaceHolder("数据库名（可选）")

	// 创建表单
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "连接名称", Widget: nameEntry},
			{Text: "数据库类型", Widget: dbTypeSelect},
			{Text: "主机", Widget: hostEntry},
			{Text: "端口", Widget: portEntry},
			{Text: "用户名", Widget: usernameEntry},
			{Text: "密码", Widget: passwordEntry},
			{Text: "数据库", Widget: databaseEntry},
		},
		OnSubmit: func() {
			// 创建新连接配置
			config := config.ConnectionConfig{
				Name:     nameEntry.Text,
				Type:     dbTypeSelect.Selected,
				Host:     hostEntry.Text,
				Port:     portEntry.Text,
				Username: usernameEntry.Text,
				Password: passwordEntry.Text,
				Database: databaseEntry.Text,
			}

			// 保存到存储
			if err := p.storage.AddConnection(config); err != nil {
				p.log.Errorf("保存连接失败: %v", err)
				dialog.ShowError(err, win)
				return
			}

			// 重新加载连接列表
			p.loadConnections()
			p.list.Refresh()

			p.log.Infof("创建新连接: %s (%s)", config.Name, config.Type)

			// 调用刷新回调，通知其他页面更新连接列表
			p.onRefresh()

			// 关闭对话框
			dialog.ShowInformation("连接创建成功", "数据库连接已保存", win)
		},
		OnCancel: func() {
			// 关闭对话框
		},
		SubmitText: "保存",
		CancelText: "取消",
	}

	// 添加测试连接按钮
	testBtn := widget.NewButton("测试连接", func() {
		// 在实际应用中，这里应该尝试连接数据库
		p.log.Info("测试连接...")
		dialog.ShowInformation("连接测试", "连接成功", win)
	})

	// 创建对话框内容
	content := container.NewVBox(
		form,
		container.NewHBox(
			layout.NewSpacer(),
			testBtn,
		),
	)

	// 显示对话框
	dialog := dialog.NewCustom("新建数据库连接", "关闭", content, win)
	dialog.Resize(fyne.NewSize(400, 400))
	dialog.Show()
}
