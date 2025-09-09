package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 创建应用
	a := app.New()
	w := a.NewWindow("中文测试")

	// 创建标签
	label := widget.NewLabel("这是一个测试中文显示的程序")

	// 创建按钮
	button := widget.NewButton("点击我", func() {
		label.SetText("按钮被点击了！")
	})

	// 创建布局
	content := container.NewVBox(
		label,
		button,
	)

	// 设置窗口内容
	w.SetContent(content)

	// 显示窗口
	w.ShowAndRun()
}
