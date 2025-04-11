package fileutil

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// MoveToTrash 将文件移动到回收站
func MoveToTrash(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return moveToTrashMacOS(path)
	case "windows":
		return moveToTrashWindows(path)
	default:
		return moveToTrashLinux(path)
	}
}

// moveToTrashMacOS 在 macOS 上将文件移动到回收站
func moveToTrashMacOS(path string) error {
	script := fmt.Sprintf(`tell app "Finder" to delete POSIX file "%s"`, path)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// moveToTrashWindows 在 Windows 上将文件移动到回收站
func moveToTrashWindows(path string) error {
	// Windows 上使用 PowerShell 的 RecycleBin
	script := fmt.Sprintf(`Add-Type -AssemblyName Microsoft.VisualBasic
[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('%s', 'OnlyErrorDialogs', 'SendToRecycleBin')`, path)
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// moveToTrashLinux 在 Linux 上将文件移动到回收站
func moveToTrashLinux(path string) error {
	// 在 Linux 上，我们创建自己的回收站目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(homeDir, ".local/share/Trash/files")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}

	// 生成唯一的文件名
	fileName := filepath.Base(path)
	trashPath := filepath.Join(trashDir, fileName)
	
	// 如果目标文件已存在，添加数字后缀
	for i := 1; ; i++ {
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(fileName)
		baseName := fileName[:len(fileName)-len(ext)]
		trashPath = filepath.Join(trashDir, fmt.Sprintf("%s_%d%s", baseName, i, ext))
	}

	return os.Rename(path, trashPath)
}

// GetFileType 获取文件类型
func GetFileType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件头部字节来判断文件类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// 使用 MIME 类型判断
	mimeType := http.DetectContentType(buffer)
	
	// 简化 MIME 类型
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image", nil
	case strings.HasPrefix(mimeType, "video/"):
		return "video", nil
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio", nil
	case strings.HasPrefix(mimeType, "text/"):
		return "text", nil
	case strings.HasPrefix(mimeType, "application/pdf"):
		return "pdf", nil
	case strings.HasPrefix(mimeType, "application/zip"),
		strings.HasPrefix(mimeType, "application/x-rar"),
		strings.HasPrefix(mimeType, "application/x-7z"):
		return "archive", nil
	default:
		return "other", nil
	}
}

// FormatFileSize 格式化文件大小显示
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// ParseFileSize 解析文件大小字符串
func ParseFileSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))
	if sizeStr == "" || sizeStr == "0" {
		return 0, nil
	}

	units := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	for suffix, multiplier := range units {
		if strings.HasSuffix(sizeStr, suffix) {
			numberStr := strings.TrimSuffix(sizeStr, suffix)
			value, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return 0, fmt.Errorf("无效的大小值: %s", sizeStr)
			}
			return int64(value * float64(multiplier)), nil
		}
	}

	// 如果没有单位，假设为字节
	value, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的大小值: %s", sizeStr)
	}
	return int64(value), nil
} 