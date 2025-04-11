package gui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/fonts/SourceHanSansCN-Normal.ttf
var sourceHanSans []byte

// DedupTheme 自定义主题
type DedupTheme struct{}

var _ fyne.Theme = (*DedupTheme)(nil)

// 颜色定义
var (
	primaryColor   = color.NRGBA{R: 66, G: 133, B: 244, A: 255} // Google Blue
	secondaryColor = color.NRGBA{R: 52, G: 168, B: 83, A: 255}  // Google Green
	errorColor     = color.NRGBA{R: 234, G: 67, B: 53, A: 255}  // Google Red
	warningColor   = color.NRGBA{R: 251, G: 188, B: 4, A: 255}  // Google Yellow
)

func NewDedupTheme() *DedupTheme {
	return &DedupTheme{}
}

func (t *DedupTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantLight {
		switch name {
		case theme.ColorNamePrimary:
			return primaryColor
		case theme.ColorNameSuccess:
			return secondaryColor
		case theme.ColorNameError:
			return errorColor
		case theme.ColorNameWarning:
			return warningColor
		case theme.ColorNameBackground:
			return color.White
		case theme.ColorNameForeground:
			return color.Black
		case theme.ColorNameButton:
			return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 180, G: 180, B: 180, A: 255}
		}
	}

	// 深色主题
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 33, G: 33, B: 33, A: 255}
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameButton:
		return color.NRGBA{R: 51, G: 51, B: 51, A: 255}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *DedupTheme) Font(style fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "SourceHanSansCN-Normal.ttf",
		StaticContent: sourceHanSans,
	}
}

func (t *DedupTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *DedupTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	default:
		return theme.DefaultTheme().Size(name)
	}
} 