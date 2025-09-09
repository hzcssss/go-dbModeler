package main

import (
	"fmt"
	"go-DBmodeler/internal/app"
	"go-DBmodeler/pkg/logger"
	"os"
	"runtime"
)

func main() {
	// 设置Fyne的环境变量，可能有助于解决中文显示问题
	os.Setenv("FYNE_SCALE", "1.0")   // 设置UI缩放比例
	os.Setenv("FYNE_THEME", "light") // 设置默认主题

	// 根据不同操作系统设置不同的字体路径
	switch runtime.GOOS {
	case "darwin": // macOS
		// 尝试几种常见的中文字体
		fontPaths := []string{
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/STHeiti Light.ttc",
			"/System/Library/Fonts/STHeiti Medium.ttc",
			"/Library/Fonts/Arial Unicode.ttf",
		}
		for _, path := range fontPaths {
			if _, err := os.Stat(path); err == nil {
				os.Setenv("FYNE_FONT", path)
				break
			}
		}
	case "windows": // Windows
		os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\msyh.ttc") // 微软雅黑
	case "linux": // Linux
		// 尝试常见的中文字体路径
		fontPaths := []string{
			"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
			"/usr/share/fonts/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc",
		}
		for _, path := range fontPaths {
			if _, err := os.Stat(path); err == nil {
				os.Setenv("FYNE_FONT", path)
				break
			}
		}
	}

	// 初始化日志
	log := logger.New()
	log.Info("GoDBModeler 启动中...")

	// 打印欢迎信息（使用中文）
	fmt.Println("欢迎使用 GoDBModeler - 数据库建模工具")
	fmt.Println("版本: v0.1.0")
	fmt.Println("作者: GoDBModeler Team")
	fmt.Println("-----------------------------------")

	// 检查工作目录
	wd, err := os.Getwd()
	if err == nil {
		log.Info("工作目录: " + wd)
	}

	// 启动应用
	application := app.New()
	application.Run()
}
