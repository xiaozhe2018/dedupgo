package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSize 解析文件大小字符串（如 "10MB"）为字节数
func ParseSize(size string) (int64, error) {
	size = strings.TrimSpace(strings.ToUpper(size))
	if size == "" || size == "0" {
		return 0, nil
	}

	units := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	var value float64

	for suffix, multiplier := range units {
		if strings.HasSuffix(size, suffix) {
			numberStr := strings.TrimSuffix(size, suffix)
			var err error
			value, err = strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return 0, fmt.Errorf("无效的大小值: %s", size)
			}
			return int64(value * float64(multiplier)), nil
		}
	}

	// 如果没有单位，假设为字节
	value, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的大小值: %s", size)
	}
	return int64(value), nil
}

// FormatSize 将字节数格式化为人类可读的字符串
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
} 