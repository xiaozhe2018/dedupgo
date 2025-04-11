package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	HashAlgorithm    string   `yaml:"hash_algorithm"`
	MinSize          string   `yaml:"min_size"`
	ExcludePatterns  []string `yaml:"exclude_patterns"`
	IncludeTypes     []string `yaml:"include_types"`
	DryRun          bool     `yaml:"dry_run"`
	OutputFormat     string   `yaml:"output_format"`
	UseTrash        bool     `yaml:"use_trash"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		HashAlgorithm: "md5",
		MinSize:      "0",
		ExcludePatterns: []string{
			"*.tmp",
			"*.temp",
			"node_modules",
			".git",
		},
		DryRun:      true,
		OutputFormat: "txt",
		UseTrash:    true,
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	if path == "" {
		// 尝试从默认位置加载
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config, nil
		}
		path = filepath.Join(homeDir, ".config", "dedupgo", "config.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(homeDir, ".config", "dedupgo")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
		path = filepath.Join(configDir, "config.yaml")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
} 