package theme

import (
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"os"
	"path/filepath"
)

//go:embed fonts/*.ttf
var fontData embed.FS

// CustomTheme 表示自定义主题
type CustomTheme struct {
	baseTheme fyne.Theme
	fontRes   fyne.Resource
}

// NewDarkTheme 创建一个新的暗色主题（默认主题）
func NewDarkTheme() fyne.Theme {
	// 尝试加载中文字体
	fontRes := loadChineseFont()

	return &CustomTheme{
		baseTheme: theme.DarkTheme(),
		fontRes:   fontRes,
	}
}

// loadChineseFont 加载中文字体
func loadChineseFont() fyne.Resource {
	// 首先尝试加载嵌入的字体文件
	fontFiles := []string{
		"fonts/Arial Unicode.ttf", // Arial Unicode 包含中文
	}

	for _, fontFile := range fontFiles {
		fontBytes, err := fontData.ReadFile(fontFile)
		if err == nil {
			return fyne.NewStaticResource(filepath.Base(fontFile), fontBytes)
		}
	}

	// 如果嵌入的字体加载失败，尝试从环境变量中获取字体路径
	if fontPath := os.Getenv("FYNE_FONT"); fontPath != "" {
		if fontBytes, err := os.ReadFile(fontPath); err == nil {
			return fyne.NewStaticResource("custom.ttf", fontBytes)
		}
	}

	// 尝试使用系统字体路径（macOS优先），排除 .ttc 文件
	systemFontPaths := []string{
		"/System/Library/Fonts/Supplemental/Arial Unicode.ttf", // Arial Unicode 包含中文
		"/Library/Fonts/Arial.ttf",                             // macOS Arial
		"/System/Library/Fonts/Geneva.ttf",                     // macOS Geneva
		"/System/Library/Fonts/Helvetica.ttf",                  // macOS Helvetica
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",      // Linux
		"C:/Windows/Fonts/msyh.ttf",                            // Windows 微软雅黑
		"C:/Windows/Fonts/simhei.ttf",                          // Windows 黑体
	}

	for _, fontPath := range systemFontPaths {
		if _, err := os.Stat(fontPath); err == nil {
			if fontBytes, err := os.ReadFile(fontPath); err == nil {
				return fyne.NewStaticResource(filepath.Base(fontPath), fontBytes)
			}
		}
	}

	// 所有方法都失败，返回nil，将使用默认字体
	return nil
}

// Color 返回指定主题颜色
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 自定义颜色
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 120, B: 212, A: 255} // 蓝色
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0, G: 180, B: 212, A: 255} // 亮蓝色
	}

	// 使用基础主题的颜色
	return t.baseTheme.Color(name, variant)
}

// Icon 返回指定主题图标
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.baseTheme.Icon(name)
}

// Font 返回指定主题字体
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 使用我们的中文字体
	if t.fontRes != nil {
		return t.fontRes
	}

	// 如果我们的字体加载失败，使用基础主题的字体
	return t.baseTheme.Font(style)
}

// Size 返回指定主题大小
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.baseTheme.Size(name)
}
