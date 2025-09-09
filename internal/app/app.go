package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"go-DBmodeler/internal/config"
	"go-DBmodeler/internal/generator"
	"go-DBmodeler/internal/ui/pages"
	apptheme "go-DBmodeler/internal/ui/theme"
	"go-DBmodeler/pkg/logger"
	"path/filepath"
)

// Application 表示GoDBModeler应用
type Application struct {
	fyneApp         fyne.App
	mainWindow      fyne.Window
	log             *logger.Logger
	storage         *config.Storage
	templateManager *generator.TemplateManager

	// 数据
	connections []config.ConnectionConfig
}

// New 创建一个新的应用实例
func New() *Application {
	// 创建日志实例
	log := logger.New()

	// 创建Fyne应用
	fyneApp := app.New()

	return &Application{
		fyneApp:     fyneApp,
		log:         log,
		connections: make([]config.ConnectionConfig, 0),
	}
}

// Run 启动应用
func (a *Application) Run() {
	// 初始化配置
	a.initConfig()

	// 设置主题
	a.setTheme()

	// 创建主窗口
	a.mainWindow = a.fyneApp.NewWindow("GoDBModeler - 数据库建模工具")
	a.mainWindow.Resize(fyne.NewSize(1024, 768))

	// 创建主界面
	a.setupUI()

	// 显示窗口并运行应用
	a.mainWindow.ShowAndRun()
}

// initConfig 初始化配置
func (a *Application) initConfig() {
	// 创建配置存储
	storage, err := config.NewStorage()
	if err != nil {
		a.log.Errorf("创建配置存储失败: %v", err)
		return
	}

	a.storage = storage

	// 加载连接配置
	a.connections = storage.GetConnections()

	// 创建模板管理器
	templateDir := filepath.Join(".", "templates", "imported")
	a.templateManager = generator.NewTemplateManager(a.log, templateDir)

	// 初始化默认模板
	if err := a.templateManager.InitializeDefaultTemplates(); err != nil {
		a.log.Warnf("初始化默认模板失败: %v", err)
	}
}

// setTheme 设置主题
func (a *Application) setTheme() {
	// 始终使用暗色主题
	a.fyneApp.Settings().SetTheme(apptheme.NewDarkTheme())
}

// setupUI 设置应用界面
func (a *Application) setupUI() {
	// 创建连接页面
	connectionPage, connectionContainer := pages.NewConnectionPage(a.log, a.storage)

	// 创建生成器页面
	generatorPage := pages.NewGeneratorPage(a.log, toConnectionConfigArray(a.connections), a.templateManager, a.storage)

	// 创建脚本管理页面
	scriptManagerPage := pages.NewScriptManagerPage(a.log, a.storage)

	// 设置连接页面的刷新回调
	connectionPage.SetRefreshCallback(func() {
		// 重新加载连接配置
		a.connections = a.storage.GetConnections()
		// 重新创建生成器页面以更新连接列表
		generatorPage = pages.NewGeneratorPage(a.log, toConnectionConfigArray(a.connections), a.templateManager, a.storage)

		// 更新标签页内容
		if tabs := a.mainWindow.Content().(*container.AppTabs); tabs != nil {
			tabs.Items[1].Content = generatorPage
			tabs.Refresh()
		}
	})

	// 创建标签页，使用中文
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("连接管理", theme.ComputerIcon(), connectionContainer),
		container.NewTabItemWithIcon("TS模型生成", theme.DocumentCreateIcon(), generatorPage),
		container.NewTabItemWithIcon("脚本管理", theme.DocumentIcon(), scriptManagerPage),
	)

	// 添加标签页切换事件处理
	tabs.OnChanged = func(tab *container.TabItem) {
		// 当切换到TS模型生成页面时，更新模板列表
		if tab.Text == "TS模型生成" {
			// 重新创建生成器页面以更新模板列表
			tab.Content = pages.NewGeneratorPage(a.log, toConnectionConfigArray(a.connections), a.templateManager, a.storage)
			tabs.Refresh()
		}
	}

	tabs.SetTabLocation(container.TabLocationLeading)

	// 设置主窗口内容
	a.mainWindow.SetContent(tabs)
}

// toConnectionConfigArray 将config.ConnectionConfig数组转换为pages.ConnectionConfig数组
func toConnectionConfigArray(configs []config.ConnectionConfig) []*pages.ConnectionConfig {
	result := make([]*pages.ConnectionConfig, 0, len(configs))

	for _, cfg := range configs {
		// 解密密码
		decryptedConfig, err := config.DecryptConnectionPassword(cfg)
		if err != nil {
			continue
		}

		result = append(result, &pages.ConnectionConfig{
			Name:     decryptedConfig.Name,
			Type:     decryptedConfig.Type,
			Host:     decryptedConfig.Host,
			Port:     decryptedConfig.Port,
			Username: decryptedConfig.Username,
			Password: decryptedConfig.Password,
			Database: decryptedConfig.Database,
		})
	}

	return result
}
