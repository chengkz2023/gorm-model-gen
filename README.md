# GORM结构体生成工具

基于gorm/gen开发的结构体生成工具，可以连接数据库来生成gorm结构体。

## 功能特点

- 支持读取配置文件生成gorm结构体
- 支持命令行参数指定生成参数
- 支持生成单个表或全部表的结构体
- 支持生成JSON标签
- 支持自定义输出目录和包名
- 支持移除表前缀

## 安装

```bash
go get github.com/yourusername/gen
```

## 使用方法

### 配置文件方式

创建配置文件`config.yaml`：

```yaml
# 数据库配置
database:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: test
  charset: utf8mb4

# 生成配置
generate:
  # 输出目录
  output_dir: ./models
  # 包名
  package: models
  # 要生成的表，为空表示所有表
  tables: []
  # 表前缀，会从生成的结构体名称中移除
  table_prefix: ""
  # 是否生成JSON标签
  with_json_tag: true
  # 是否生成数据库查询函数
  with_query_functions: false
```

然后运行：

```bash
gen generate -c config.yaml
```

### 命令行参数方式

```bash
# 生成指定表的结构体
gen generate -t user

# 生成全部表的结构体
gen generate

# 生成结构体到指定目录
gen generate -o ./models

# 生成带JSON标签的结构体
gen generate --json

# 移除表前缀
gen generate --prefix tbl_

# 指定包名
gen generate -p models

# 显示帮助信息
gen -h
gen generate -h

# 显示版本信息
gen version
gen -v
```

## 配置项说明

### 数据库配置

- `host`: 数据库服务器地址
- `port`: 数据库端口
- `user`: 数据库用户名
- `password`: 数据库密码
- `dbname`: 数据库名称
- `charset`: 字符集

### 生成配置

- `output_dir`: 生成的代码输出目录
- `package`: 生成的代码包名
- `tables`: 要生成的表名列表，为空表示所有表
- `table_prefix`: 表前缀，会从生成的结构体名称中移除
- `with_json_tag`: 是否生成JSON标签
- `with_query_functions`: 是否生成数据库查询函数 