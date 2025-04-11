package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed resources/SourceHanSansCN-Normal.ttf
var sourceHanSans []byte

var defaultTheme = theme.DefaultTheme()

type myTheme struct{}

func newMyTheme() fyne.Theme {
	return &myTheme{}
}

func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}
	return defaultTheme.Color(name, variant)
}

func (m *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "SourceHanSansCN-Normal.ttf",
		StaticContent: sourceHanSans,
	}
}

func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return defaultTheme.Icon(name)
}

func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // 增大默认文字大小
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	default:
		return defaultTheme.Size(name)
	}
} 