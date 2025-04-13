package generator

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yourusername/gen/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Generator 表示代码生成器
type Generator struct {
	Config *config.Config
	DB     *gorm.DB
}

// NewGenerator 创建一个新的生成器
func NewGenerator(cfg *config.Config) (*Generator, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &Generator{
		Config: cfg,
		DB:     db,
	}, nil
}

// Generate 生成代码
func (g *Generator) Generate() error {
	// 确保输出目录存在
	outputDir := g.Config.Generate.OutputDir
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建生成器
	genConfig := gen.Config{
		OutPath:      outputDir, // 输出路径
		Mode:         gen.WithDefaultQuery | gen.WithoutContext,
		ModelPkgPath: g.Config.Generate.Package, // 包名
	}

	// 配置JSON标签
	if g.Config.Generate.WithJSONTag {
		// 使用WithJSONTagNameStrategy配置JSON标签
		genConfig.WithJSONTagNameStrategy(func(columnName string) string {
			return columnName
		})
	}

	generator := gen.NewGenerator(genConfig)

	// 设置目标数据库
	generator.UseDB(g.DB)

	// 获取表
	tables := g.Config.Generate.Tables
	tablePrefix := g.Config.Generate.TablePrefix

	// 如果没有指定表，则获取所有表
	if len(tables) == 0 {
		var dbTables []string
		if err := g.DB.Raw("SHOW TABLES").Scan(&dbTables).Error; err != nil {
			return fmt.Errorf("获取表列表失败: %w", err)
		}
		tables = dbTables
	}

	// 为每个表创建模型
	for _, tableName := range tables {
		structName := tableName
		// 移除前缀
		if tablePrefix != "" && strings.HasPrefix(tableName, tablePrefix) {
			structName = strings.TrimPrefix(tableName, tablePrefix)
		}
		// 驼峰命名
		structName = toCamelCase(structName)

		// 生成模型选项
		opts := []gen.ModelOpt{}

		// 生成结构体
		generator.GenerateModel(tableName, opts...)
	}

	// 执行代码生成
	generator.Execute()

	log.Printf("已成功生成代码到 %s 目录", outputDir)
	return nil
}

// GenerateTable 生成单个表的结构体
func (g *Generator) GenerateTable(tableName string) error {
	// 确保输出目录存在
	outputDir := g.Config.Generate.OutputDir
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建生成器
	genConfig := gen.Config{
		OutPath:      outputDir, // 输出路径
		Mode:         gen.WithDefaultQuery | gen.WithoutContext,
		ModelPkgPath: g.Config.Generate.Package, // 包名
	}

	// 配置JSON标签
	if g.Config.Generate.WithJSONTag {
		// 使用WithJSONTagNameStrategy配置JSON标签
		genConfig.WithJSONTagNameStrategy(func(columnName string) string {
			return columnName
		})
	}

	generator := gen.NewGenerator(genConfig)

	// 设置目标数据库
	generator.UseDB(g.DB)

	structName := tableName
	// 移除前缀
	tablePrefix := g.Config.Generate.TablePrefix
	if tablePrefix != "" && strings.HasPrefix(tableName, tablePrefix) {
		structName = strings.TrimPrefix(tableName, tablePrefix)
	}
	// 驼峰命名
	structName = toCamelCase(structName)

	// 生成模型选项
	opts := []gen.ModelOpt{}

	// 检查表是否存在
	var count int64
	result := g.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&count)
	if result.Error != nil {
		return fmt.Errorf("检查表是否存在失败: %w", result.Error)
	}
	if count == 0 {
		return fmt.Errorf("表 %s 不存在", tableName)
	}

	// 生成结构体
	generator.GenerateModel(tableName, opts...)

	// 执行代码生成
	generator.Execute()

	log.Printf("已成功生成表 %s 的结构体到 %s 目录", tableName, outputDir)
	return nil
}

// toCamelCase 转换表名为驼峰命名
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

// GetAllTables 获取数据库中的所有表
func (g *Generator) GetAllTables() ([]string, error) {
	var tables []string
	if err := g.DB.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return nil, fmt.Errorf("获取表列表失败: %w", err)
	}
	return tables, nil
}
