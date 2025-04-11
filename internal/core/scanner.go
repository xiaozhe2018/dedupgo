package core

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Scanner 文件扫描器
type Scanner struct {
	HashAlgorithm string
	MinSize       int64
	FileTypes     []string
	ExcludePatterns []string
	concurrent    int
}

// FileInfo 存储文件信息
type FileInfo struct {
	Path     string
	Size     int64
	Hash     string
	FileType string
}

// Result 扫描结果
type Result struct {
	DuplicateGroups map[string][]string
	TotalFiles      int
	TotalSize       int64
	SavedSize       int64
}

// NewScanner 创建新的扫描器实例
func NewScanner(hashAlgo string, minSize int64, fileTypes []string, excludePatterns []string) *Scanner {
	return &Scanner{
		HashAlgorithm:    hashAlgo,
		MinSize:          minSize,
		FileTypes:        fileTypes,
		ExcludePatterns: excludePatterns,
		concurrent:      5, // 默认并发数
	}
}

// getHasher 根据配置返回相应的哈希函数
func (s *Scanner) getHasher() hash.Hash {
	switch s.HashAlgorithm {
	case "sha256":
		return sha256.New()
	default:
		return md5.New()
	}
}

// calculateFileHash 计算文件哈希值
func (s *Scanner) calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := s.getHasher()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// Scan 执行扫描操作
func (s *Scanner) Scan(paths ...string) (*Result, error) {
	fileMap := make(map[string][]string)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.concurrent)

	for _, root := range paths {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() || info.Size() < s.MinSize {
				return nil
			}

			// 检查是否匹配排除模式
			for _, pattern := range s.ExcludePatterns {
				if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
					return nil
				}
			}

			wg.Add(1)
			go func(filePath string, size int64) {
				defer wg.Done()
				semaphore <- struct{}{} // 获取信号量
				defer func() { <-semaphore }() // 释放信号量

				hash, err := s.calculateFileHash(filePath)
				if err != nil {
					return
				}

				mutex.Lock()
				fileMap[hash] = append(fileMap[hash], filePath)
				mutex.Unlock()
			}(path, info.Size())

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	wg.Wait()

	// 处理结果
	result := &Result{
		DuplicateGroups: make(map[string][]string),
	}

	for hash, files := range fileMap {
		if len(files) > 1 {
			result.DuplicateGroups[hash] = files
			fileInfo, _ := os.Stat(files[0])
			result.SavedSize += fileInfo.Size() * int64(len(files)-1)
		}
		result.TotalFiles += len(files)
		fileInfo, _ := os.Stat(files[0])
		result.TotalSize += fileInfo.Size() * int64(len(files))
	}

	return result, nil
}

// MoveToTrash 将文件移动到系统回收站
func MoveToTrash(filePath string) error {
	// 在 macOS 上使用 osascript 将文件移动到回收站
	if strings.HasPrefix(runtime.GOOS, "darwin") {
		script := fmt.Sprintf(`tell app "Finder" to delete POSIX file "%s"`, filePath)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()
	}

	// 在 Windows 上使用 PowerShell 将文件移动到回收站
	if runtime.GOOS == "windows" {
		script := fmt.Sprintf(`Add-Type -AssemblyName Microsoft.VisualBasic
[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('%s','OnlyErrorDialogs','SendToRecycleBin')`, filePath)
		cmd := exec.Command("powershell", "-Command", script)
		return cmd.Run()
	}

	// 在 Linux 上使用 gio 将文件移动到回收站
	if runtime.GOOS == "linux" {
		cmd := exec.Command("gio", "trash", filePath)
		return cmd.Run()
	}

	// 如果以上方法都不适用，则返回错误
	return fmt.Errorf("不支持在当前操作系统(%s)上使用回收站功能", runtime.GOOS)
} 