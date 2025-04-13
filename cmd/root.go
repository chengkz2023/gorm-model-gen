package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gen/config"
)

var (
	// 配置文件路径
	cfgFile string
	// 表名
	tableName string
	// 版本号
	version = "0.1.0"
)

// rootCmd 表示基础命令
var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成gorm结构体工具",
	Long: `基于gorm/gen开发的结构体生成工具，可以连接数据库来生成gorm结构体。
可以通过配置文件或命令行参数指定配置。`,
}

// versionCmd 表示显示版本的命令
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "显示版本信息",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gen version %s\n", version)
	},
}

// Execute 执行根命令
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 添加全局标志
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径 (默认为当前目录下的config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&tableName, "table", "t", "", "要生成的表名，不指定则生成全部表")

	// 添加子命令
	rootCmd.AddCommand(versionCmd)
}

// initConfig 初始化配置
func initConfig() {
	// 如果未指定配置文件，则尝试从当前目录加载config.yaml
	if cfgFile == "" {
		cfgFile = "config.yaml"
	}

	// 如果配置文件不存在，则使用默认配置
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在，将使用默认配置\n", cfgFile)
		return
	}

	// 加载配置文件
	_, err := config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}
}
