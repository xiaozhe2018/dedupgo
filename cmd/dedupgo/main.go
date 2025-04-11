package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xiaozhe/dedupgo/internal/config"
	"github.com/xiaozhe/dedupgo/internal/core"
)

var (
	configFile    string
	hashAlgorithm string
	minSize      string
	force        bool
	outputFormat string
	useTrash     bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	flag.StringVar(&hashAlgorithm, "hash", "md5", "哈希算法 (md5/sha256)")
	flag.StringVar(&minSize, "min-size", "0", "最小文件大小 (例如: 10MB)")
	flag.BoolVar(&force, "force", false, "强制删除重复文件")
	flag.StringVar(&outputFormat, "output", "txt", "输出格式 (txt/json)")
	flag.BoolVar(&useTrash, "trash", true, "使用回收站")
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 命令行参数覆盖配置文件
	if hashAlgorithm != "md5" {
		cfg.HashAlgorithm = hashAlgorithm
	}
	if minSize != "0" {
		cfg.MinSize = minSize
	}
	cfg.DryRun = !force
	if outputFormat != "txt" {
		cfg.OutputFormat = outputFormat
	}
	cfg.UseTrash = useTrash

	// 获取扫描目录
	dirs := flag.Args()
	if len(dirs) == 0 {
		fmt.Fprintln(os.Stderr, "错误: 请指定至少一个扫描目录")
		flag.Usage()
		os.Exit(1)
	}

	// 创建扫描器
	scanner := core.NewScanner(
		cfg.HashAlgorithm,
		0, // minSize 将在 Scanner 中解析
		cfg.IncludeTypes,
		cfg.ExcludePatterns,
	)

	// 执行扫描
	result, err := scanner.Scan(dirs...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "扫描失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	switch strings.ToLower(cfg.OutputFormat) {
	case "json":
		outputJSON(result)
	default:
		outputText(result, cfg.DryRun)
	}
}

func outputJSON(result *core.Result) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "JSON输出失败: %v\n", err)
		os.Exit(1)
	}
}

func outputText(result *core.Result, isDryRun bool) {
	fmt.Printf("扫描完成！\n")
	fmt.Printf("总文件数: %d\n", result.TotalFiles)
	fmt.Printf("总大小: %.2f MB\n", float64(result.TotalSize)/(1024*1024))
	fmt.Printf("可节省空间: %.2f MB\n\n", float64(result.SavedSize)/(1024*1024))

	if len(result.DuplicateGroups) == 0 {
		fmt.Println("未发现重复文件")
		return
	}

	fmt.Printf("发现 %d 组重复文件:\n\n", len(result.DuplicateGroups))
	for hash, files := range result.DuplicateGroups {
		fmt.Printf("哈希值: %s\n", hash)
		for i, file := range files {
			if i == 0 {
				fmt.Printf("  [保留] %s\n", file)
			} else {
				if isDryRun {
					fmt.Printf("  [待删除] %s\n", file)
				} else {
					fmt.Printf("  [已删除] %s\n", file)
				}
			}
		}
		fmt.Println()
	}

	if isDryRun {
		fmt.Println("提示: 这是预览模式。使用 --force 参数执行实际删除操作。")
	}
} 