package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/xiaozhe/dedupgo/internal/core"
)

func init() {
	// 在应用启动前设置字体
	switch runtime.GOOS {
	case "darwin":
		// macOS 字体设置
		fonts := []string{
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/STHeiti Light.ttc",
			"/System/Library/Fonts/STHeiti Medium.ttc",
			"/Library/Fonts/Arial Unicode.ttf",
		}
		for _, font := range fonts {
			if _, err := os.Stat(font); err == nil {
				os.Setenv("FYNE_FONT", font)
				break
			}
		}
	case "linux":
		// Linux 字体设置
		fonts := []string{
			"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		}
		for _, font := range fonts {
			if _, err := os.Stat(font); err == nil {
				os.Setenv("FYNE_FONT", font)
				break
			}
		}
	case "windows":
		// Windows 字体设置
		fonts := []string{
			"C:\\Windows\\Fonts\\msyh.ttc",
			"C:\\Windows\\Fonts\\simsun.ttc",
			"C:\\Windows\\Fonts\\simhei.ttf",
		}
		for _, font := range fonts {
			if _, err := os.Stat(font); err == nil {
				os.Setenv("FYNE_FONT", font)
				break
			}
		}
	}
}

// customTheme 自定义主题以支持中文
type customTheme struct {
	fyne.Theme
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func main() {
	// 设置默认编码为 UTF-8
	os.Setenv("LANG", "zh_CN.UTF-8")
	
	myApp := app.New()
	myApp.Settings().SetTheme(newMyTheme())
	
	myWindow := myApp.NewWindow("DedupGo - 文件去重工具")

	// 创建标题和副标题，使用更现代的样式
	title := canvas.NewText("DedupGo", nil)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Color = theme.PrimaryColor()

	subtitle := widget.NewLabelWithStyle(
		"文件去重工具",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	description := widget.NewLabelWithStyle(
		"轻量、安全、高效的文件去重解决方案",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// 创建主要控件
	pathList := widget.NewMultiLineEntry()
	pathList.SetText("等待添加扫描目录...\n")
	pathList.Disable()
	pathList.TextStyle = fyne.TextStyle{Monospace: true}
	pathList.Wrapping = fyne.TextWrapBreak

	// 使用 MultiLineEntry 替代 TextGrid
	resultArea := widget.NewMultiLineEntry()
	resultArea.SetText("\n  📊扫描结果将在这里显示\n\n")
	resultArea.Disable()
	resultArea.TextStyle = fyne.TextStyle{Monospace: true}
	resultArea.Wrapping = fyne.TextWrapBreak

	// 添加删除按钮（初始隐藏）
	deleteButton := widget.NewButtonWithIcon("删除重复文件", theme.DeleteIcon(), nil)
	deleteButton.Hide()
	deleteButton.Importance = widget.DangerImportance

	// 创建滚动容器，设置相同的最小高度和宽度
	containerSize := fyne.NewSize(500, 500)  // 设置固定的宽度和高度
	pathListScroll := container.NewScroll(pathList)
	pathListScroll.SetMinSize(containerSize)
	resultScroll := container.NewScroll(resultArea)
	resultScroll.SetMinSize(containerSize)

	// 创建带有边框和标题的容器
	pathListBorder := container.NewBorder(
		container.NewHBox(
			widget.NewIcon(theme.FolderIcon()),
			widget.NewLabelWithStyle("已选择的目录", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		),
		nil, nil, nil,
		container.NewPadded(pathListScroll),
	)

	resultBorder := container.NewBorder(
		container.NewHBox(
			widget.NewIcon(theme.DocumentIcon()),
			widget.NewLabelWithStyle("扫描结果", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		),
		nil, nil, nil,
		container.NewPadded(resultScroll),
	)

	// 创建主要布局容器，使用固定的分割比例
	mainContent := container.NewHSplit(
		container.NewPadded(pathListBorder),
		container.NewVBox(
			container.NewPadded(resultBorder),
			container.NewPadded(container.NewHBox(
				layout.NewSpacer(),
				deleteButton,
				layout.NewSpacer(),
			)),
		),
	)

	var selectedPaths []string
	
	// 添加目录按钮
	addButton := widget.NewButtonWithIcon("添加目录", theme.FolderOpenIcon(), nil)
	addButton.Importance = widget.HighImportance

	// 扫描选项样式优化
	hashAlgo := widget.NewSelect([]string{"md5", "sha256"}, nil)
	hashAlgo.SetSelected("md5")
	hashAlgo.PlaceHolder = "选择哈希算法"
	
	minSizeEntry := widget.NewEntry()
	minSizeEntry.SetPlaceHolder("最小文件大小（如：1MB）")
	minSizeEntry.Resize(fyne.NewSize(150, minSizeEntry.MinSize().Height))

	// 状态标签样式优化
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel.Hide()

	// 扫描按钮样式优化
	scanButton := widget.NewButtonWithIcon("开始扫描", theme.SearchIcon(), nil)
	scanButton.Importance = widget.HighImportance

	// 清除按钮样式优化
	clearButton := widget.NewButtonWithIcon("清除", theme.DeleteIcon(), nil)

	// 优化按钮布局
	buttons := container.NewHBox(
		container.NewPadded(addButton),
		widget.NewSeparator(),
		container.NewPadded(scanButton),
		widget.NewSeparator(),
		container.NewPadded(clearButton),
	)

	// 优化选项布局
	options := container.NewHBox(
		widget.NewLabelWithStyle("哈希算法", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(hashAlgo),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("最小大小", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(minSizeEntry),
	)

	// 优化头部布局
	header := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		container.NewCenter(description),
		widget.NewSeparator(),
		container.NewPadded(buttons),
		widget.NewSeparator(),
		container.NewPadded(options),
		widget.NewSeparator(),
	)

	// 设置主布局
	content := container.NewBorder(
		header,
		container.NewPadded(statusLabel),
		nil, nil,
		container.NewPadded(mainContent),
	)

	var currentResult *core.Result
	
	// 添加目录按钮的事件处理
	addButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uri == nil {
				return
			}
			path := uri.Path()
			selectedPaths = append(selectedPaths, path)
			currentText := pathList.Text
			if currentText == "等待添加扫描目录...\n" {
				currentText = ""
			}
			pathList.SetText(currentText + "📁 " + path + "\n")
		}, myWindow)
	}

	// 扫描按钮的事件处理
	scanButton.OnTapped = func() {
		if len(selectedPaths) == 0 {
			dialog.ShowInformation("提示", "请先添加要扫描的目录", myWindow)
			return
		}

		deleteButton.Hide()
		statusLabel.SetText("🔍 正在扫描文件...")
		statusLabel.Show()
		resultArea.SetText("\n  正在扫描中，请稍候...\n  这可能需要一些时间，具体取决于文件数量\n")
		scanButton.Disable()
		addButton.Disable()

		scanner := core.NewScanner(
			hashAlgo.Selected,
			0,
			nil,
			nil,
		)

		go func() {
			result, err := scanner.Scan(selectedPaths...)
			if err != nil {
				dialog.ShowError(err, myWindow)
				scanButton.Enable()
				addButton.Enable()
				statusLabel.Hide()
				return
			}

			currentResult = result

			var sb strings.Builder
			// 使用表格样式展示统计信息
			sb.WriteString("\n")  // 添加顶部间距
			sb.WriteString("  扫描结果统计:\n")
			sb.WriteString("  ┌─────────────────┬─────────────┐\n")
			sb.WriteString(fmt.Sprintf("  │ 📁 总文件数    │ %9d │\n", result.TotalFiles))
			sb.WriteString(fmt.Sprintf("  │ 💾 总大小      │ %8.1f MB│\n", float64(result.TotalSize)/(1024*1024)))
			sb.WriteString(fmt.Sprintf("  │ 🗑️ 可节省空间  │ %8.1f MB│\n", float64(result.SavedSize)/(1024*1024)))
			sb.WriteString(fmt.Sprintf("  │ 🔍 重复文件组  │ %9d │\n", len(result.DuplicateGroups)))
			sb.WriteString("  └─────────────────┴─────────────┘\n\n")

			if len(result.DuplicateGroups) == 0 {
				sb.WriteString("  ✨ 恭喜！未发现重复文件\n")
			} else {
				sb.WriteString("  📑 重复文件列表:\n\n")
				groupNum := 1
				for _, files := range result.DuplicateGroups {
					fileInfo, _ := os.Stat(files[0])
					fileSize := fileInfo.Size()
					savedSpace := float64(fileSize * int64(len(files)-1)) / (1024 * 1024)

					sb.WriteString("  ┌───────────────────────────────┐\n")
					sb.WriteString(fmt.Sprintf("  │ 📌 第 %d 组                   │\n", groupNum))
					sb.WriteString("  ├───────────────────────────────┤\n")
					sb.WriteString(fmt.Sprintf("  │ 📦 文件数: %-3d               │\n", len(files)))
					sb.WriteString(fmt.Sprintf("  │ 📏 大小: %-6.1f MB           │\n", float64(fileSize)/(1024*1024)))
					sb.WriteString(fmt.Sprintf("  │ 💾 节省: %-6.1f MB           │\n", savedSpace))
					sb.WriteString("  ├───────────────────────────────┤\n")
					
					for i, file := range files {
						fInfo, _ := os.Stat(file)
						if i == 0 {
							sb.WriteString("  │ 🟢 原始文件                   │\n")
						} else {
							sb.WriteString("  │ 🔴 重复文件                   │\n")
						}
						sb.WriteString(fmt.Sprintf("  │   %s\n", file))
						sb.WriteString(fmt.Sprintf("  │   修改于: %s   │\n", fInfo.ModTime().Format("2006-01-02 15:04")))
						sb.WriteString("  │                               │\n")
					}
					sb.WriteString("  └───────────────────────────────┘\n\n")
					groupNum++
				}
				
				if len(result.DuplicateGroups) > 0 {
					deleteButton.Show()
				}
			}

			resultArea.SetText(sb.String())
			statusLabel.Hide()
			scanButton.Enable()
			addButton.Enable()
		}()
	}

	// 设置删除按钮的动作
	deleteButton.OnTapped = func() {
		if currentResult == nil || len(currentResult.DuplicateGroups) == 0 {
			return
		}

		// 计算要删除的文件数量和总大小
		var totalFiles int
		var totalSize int64
		for _, files := range currentResult.DuplicateGroups {
			totalFiles += len(files) - 1 // 减去每组保留的文件
			for i, file := range files {
				if i > 0 { // 跳过每组的第一个文件（保留文件）
					if info, err := os.Stat(file); err == nil {
						totalSize += info.Size()
					}
				}
			}
		}

		// 显示确认对话框
		dialog.ShowConfirm(
			"确认删除",
			fmt.Sprintf("确定要删除 %d 个重复文件吗？\n总计可释放 %.2f MB 空间\n\n注意：删除的文件将被移动到回收站", 
				totalFiles, 
				float64(totalSize)/(1024*1024)),
			func(confirm bool) {
				if !confirm {
					return
				}

				// 执行删除操作
				statusLabel.SetText("🗑️ 正在删除文件...")
				statusLabel.Show()
				deleteButton.Disable()
				scanButton.Disable()
				
				go func() {
					var deletedCount int
					var errorCount int
					
					for _, files := range currentResult.DuplicateGroups {
						for i, file := range files {
							if i > 0 { // 跳过每组的第一个文件（保留文件）
								if err := core.MoveToTrash(file); err != nil {
									errorCount++
								} else {
									deletedCount++
								}
							}
						}
					}

					// 更新结果显示
					resultArea.SetText(fmt.Sprintf(
						"🗑️ 删除操作完成！\n\n"+
							"✅ 成功删除: %d 个文件\n"+
							"❌ 删除失败: %d 个文件\n\n"+
							"提示：删除的文件已移动到回收站，可以随时恢复。",
						deletedCount,
						errorCount,
					))

					statusLabel.Hide()
					deleteButton.Hide() // 隐藏删除按钮
					scanButton.Enable()
					currentResult = nil // 清除当前结果
				}()
			},
			myWindow,
		)
	}

	// 清除按钮的事件处理
	clearButton.OnTapped = func() {
		selectedPaths = nil
		pathList.SetText("等待添加扫描目录...\n")
		resultArea.SetText("\n  📊 欢迎使用 DedupGo\n\n  扫描结果将在这里显示\n  请先添加要扫描的目录...\n")
		statusLabel.Hide()
		deleteButton.Hide()
		scanButton.Enable()
		addButton.Enable()
	}

	// 设置窗口
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1200, 800))
	mainContent.SetOffset(0.35)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
} 