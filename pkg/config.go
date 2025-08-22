package pkg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	OpenAI OpenAIConfig `yaml:"openai"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果没有指定配置文件路径，使用默认路径
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 解析YAML
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// 验证必要配置
	if config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("openai.api_key is required in config file")
	}

	// 设置默认值
	if config.OpenAI.Model == "" {
		config.OpenAI.Model = "gpt-4o"
	}

	return &config, nil
}

// GetDefaultConfig 返回默认配置（用于生成示例配置文件）
func GetDefaultConfig() *Config {
	return &Config{
		OpenAI: OpenAIConfig{
			APIKey:  "your-openai-api-key",
			BaseURL: "https://api.openai.com/v1/", // 或者 "https://api.shubiaobiao.cn/v1/"
			Model:   "gpt-4o",
		},
	}
}

// SaveExampleConfig 保存示例配置文件
func SaveExampleConfig(configPath string) error {
	config := GetDefaultConfig()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
