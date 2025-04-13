package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示程序的配置
type Config struct {
	// 数据库配置
	Database DatabaseConfig `yaml:"database"`
	// 生成配置
	Generate GenerateConfig `yaml:"generate"`
}

// DatabaseConfig 表示数据库连接配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

// GenerateConfig 表示代码生成配置
type GenerateConfig struct {
	// 输出目录
	OutputDir string `yaml:"output_dir"`
	// 包名
	Package string `yaml:"package"`
	// 要生成的表，空表示所有表
	Tables []string `yaml:"tables"`
	// 表名前缀，会被移除
	TablePrefix string `yaml:"table_prefix"`
	// 是否生成JSON标签
	WithJSONTag bool `yaml:"with_json_tag"`
	// 是否生成数据库查询函数
	WithQueryFunctions bool `yaml:"with_query_functions"`
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	return config, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			DBName:   "test",
			Charset:  "utf8mb4",
		},
		Generate: GenerateConfig{
			OutputDir:          "./models",
			Package:            "models",
			Tables:             []string{},
			TablePrefix:        "",
			WithJSONTag:        true,
			WithQueryFunctions: false,
		},
	}
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("生成YAML配置失败: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
