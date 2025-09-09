package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// SyntaxHighlighter 提供基本的语法高亮功能
type SyntaxHighlighter struct {
	container *fyne.Container
}

// NewSyntaxHighlighter 创建一个新的语法高亮器
func NewSyntaxHighlighter() *SyntaxHighlighter {
	return &SyntaxHighlighter{
		container: container.NewVBox(),
	}
}

// HighlightTypeScript 对TypeScript代码进行语法高亮
func (sh *SyntaxHighlighter) HighlightTypeScript(code string) fyne.CanvasObject {
	lines := strings.Split(code, "\n")
	sh.container.Objects = nil

	// 创建一个大的代码块容器
	codeBlock := container.NewVBox()

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			codeBlock.Add(widget.NewLabel(" "))
			continue
		}

		// 简单的语法高亮规则
		highlightedLine := sh.highlightLine(line)
		codeBlock.Add(highlightedLine)
	}

	// 设置代码块的样式
	codeBlockContainer := container.NewPadded(codeBlock)

	// 将代码块添加到主容器
	sh.container.Add(codeBlockContainer)

	// 返回带滚动条的容器
	scroll := container.NewVScroll(sh.container)
	scroll.SetMinSize(fyne.NewSize(800, 600)) // 设置固定大小800x600

	return scroll
}

// highlightLine 高亮单行代码
func (sh *SyntaxHighlighter) highlightLine(line string) fyne.CanvasObject {
	// 检测不同的语法元素
	if sh.isComment(line) {
		return sh.createCommentLabel(line)
	} else if sh.isInterface(line) {
		return sh.createInterfaceLabel(line)
	} else if sh.isField(line) {
		return sh.createFieldLabel(line)
	} else if sh.isBrace(line) {
		return sh.createBraceLabel(line)
	} else {
		return sh.createNormalLabel(line)
	}
}

// isComment 判断是否为注释行
func (sh *SyntaxHighlighter) isComment(line string) bool {
	return strings.Contains(line, "//") || strings.Contains(line, "/*") || strings.Contains(line, "*/")
}

// isInterface 判断是否为接口定义行
func (sh *SyntaxHighlighter) isInterface(line string) bool {
	return strings.Contains(line, "interface") || strings.Contains(line, "export interface")
}

// isField 判断是否为字段定义行
func (sh *SyntaxHighlighter) isField(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.Contains(trimmed, ":") && !strings.Contains(trimmed, "//") &&
		!strings.Contains(trimmed, "interface") && !strings.Contains(trimmed, "export")
}

// isBrace 判断是否为括号行
func (sh *SyntaxHighlighter) isBrace(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "{" || trimmed == "}" || trimmed == "};"
}

// createCommentLabel 创建注释标签
func (sh *SyntaxHighlighter) createCommentLabel(line string) fyne.CanvasObject {
	label := widget.NewLabel(line)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	return label
}

// createInterfaceLabel 创建接口标签
func (sh *SyntaxHighlighter) createInterfaceLabel(line string) fyne.CanvasObject {
	label := widget.NewLabel(line)
	label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	return label
}

// createFieldLabel 创建字段标签
func (sh *SyntaxHighlighter) createFieldLabel(line string) fyne.CanvasObject {
	label := widget.NewLabel(line)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	return label
}

// createBraceLabel 创建括号标签
func (sh *SyntaxHighlighter) createBraceLabel(line string) fyne.CanvasObject {
	label := widget.NewLabel(line)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	return label
}

// createNormalLabel 创建普通标签
func (sh *SyntaxHighlighter) createNormalLabel(line string) fyne.CanvasObject {
	label := widget.NewLabel(line)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	return label
}

// SimpleCodeDisplay 简单的代码显示组件
type SimpleCodeDisplay struct {
	*widget.Label
}

// NewSimpleCodeDisplay 创建简单的代码显示组件
func NewSimpleCodeDisplay(text string) *SimpleCodeDisplay {
	display := &SimpleCodeDisplay{
		Label: widget.NewLabel(text),
	}
	display.TextStyle = fyne.TextStyle{Monospace: true}
	return display
}

// SetText 设置文本内容
func (scd *SimpleCodeDisplay) SetText(text string) {
	scd.Label.SetText(text)
}
