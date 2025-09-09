package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 检查是否安装了 rsvg-convert 工具
	_, err := exec.LookPath("rsvg-convert")
	if err != nil {
		fmt.Println("未找到 rsvg-convert 工具，请安装：brew install librsvg")
		return
	}

	// 使用 rsvg-convert 将 SVG 转换为 PNG
	cmd := exec.Command("rsvg-convert", "-h", "1024", "-w", "1024", "-o", "temp_icons/app.png", "resources/icons/app.svg")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	fmt.Println("成功将 SVG 转换为 PNG: temp_icons/app.png")
}
