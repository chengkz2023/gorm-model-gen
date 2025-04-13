package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gen/config"
	"github.com/yourusername/gen/generator"
)

var (
	// 输出目录
	outputDir string
	// 包名
	packageName string
	// 生成JSON标签
	withJSONTag bool
	// 表前缀
	tablePrefix string
)

// generateCmd 表示生成代码的命令
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "生成gorm结构体",
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		// 载入配置
		var cfg *config.Config
		var err error

		if cfgFile != "" {
			// 加载配置文件
			cfg, err = config.LoadConfig(cfgFile)
			if err != nil {
				fmt.Printf("加载配置文件失败: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 使用默认配置
			cfg = config.DefaultConfig()
		}

		// 命令行参数覆盖配置文件
		if outputDir != "" {
			cfg.Generate.OutputDir = outputDir
		}
		if packageName != "" {
			cfg.Generate.Package = packageName
		}
		if tablePrefix != "" {
			cfg.Generate.TablePrefix = tablePrefix
		}
		if cmd.Flags().Changed("json") {
			cfg.Generate.WithJSONTag = withJSONTag
		}
		if tableName != "" {
			cfg.Generate.Tables = []string{tableName}
		}

		// 创建生成器
		gen, err := generator.NewGenerator(cfg)
		if err != nil {
			fmt.Printf("创建生成器失败: %v\n", err)
			os.Exit(1)
		}

		// 生成代码
		if tableName != "" {
			err = gen.GenerateTable(tableName)
		} else {
			err = gen.Generate()
		}

		if err != nil {
			fmt.Printf("生成代码失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// 添加命令行标志
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "输出目录")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "", "包名")
	generateCmd.Flags().BoolVar(&withJSONTag, "json", true, "生成JSON标签")
	generateCmd.Flags().StringVar(&tablePrefix, "prefix", "", "表前缀")
}
