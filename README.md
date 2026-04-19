# oh-my-markdown

一个高效的 Go CLI 工具，用于批量处理 Markdown 文件。该工具采用 Cobra 框架构建，提供强大的文件处理能力，特别是对 Markdown 文件的 front matter 提取和处理。

## 项目特点

- **高性能**：使用 Go 编写，具有优秀的并发处理能力
- **灵活的命令行界面**：基于 Cobra 框架，提供直观的命令和标志
- **Front Matter 提取**：支持读取 Markdown 文件中的 YAML front matter，并输出为 JSON
- **递归目录支持**：支持递归遍历子目录，灵活处理嵌套文件结构
- **错误隐断**：完善的错误处理和日志记录机制
- **高测试覆盖率**：代码覆盖率达到 80% 以上

## 系统要求

- **Go 版本**：1.25.0 或更高
- **操作系统**：Windows（64-bit）

## 安装

### 从源码构建

```bash
# 克隆项目
git clone https://github.com/yourusername/oh-my-markdown.git
cd oh-my-markdown

# 构建项目
go build -o oh-my-markdown.exe .

# 验证安装
./oh-my-markdown.exe --help
```

## 快速开始

### 基本用法

```bash
# 查看帮助信息
oh-my-markdown.exe --help

# 运行具体命令
oh-my-markdown.exe front-matter ./path/to/markdown/files
```

### Front Matter 命令

`front-matter` 命令用于提取 Markdown 文件的 front matter 并输出为 JSON 格式。

#### 基本语法

```bash
oh-my-markdown front-matter <directory> [flags]
```

#### 参数

- `<directory>` **必需**：要扫描的目录路径

#### 标志（Flags）

| 标志 | 简写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--output` | `-o` | string | 无（stdout） | 输出文件路径。如果指定，结果将写入该文件；否则输出到标准输出 |
| `--recursive` | `-r` | bool | `true` | 是否递归遍历子目录。默认为 `true`，设置为 `false` 则仅处理顶级目录中的文件 |

#### 示例

**示例 1：读取目录中的所有 front matter 并输出到标准输出**

```bash
oh-my-markdown.exe front-matter ./docs
```

输出（JSON 格式）：
```json
[
  {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "title": "测试文章",
    "date": "2024-01-01",
    "tags": [
      "go",
      "markdown"
    ]
  }
]
```

**示例 2：读取目录中的 front matter 并保存到文件**

```bash
oh-my-markdown.exe front-matter ./docs -o ./output.json
```

**示例 3：仅读取顶级目录中的文件（不递归）**

```bash
oh-my-markdown.exe front-matter ./docs --recursive=false
```

**示例 4：组合多个标志**

```bash
oh-my-markdown.exe front-matter ./docs -r=false -o ./frontmatters.json
```

## 开发指南

### 项目结构

```
oh-my-markdown/
├── main.go                          # 应用入口点
├── go.mod                           # Go 模块定义
├── go.sum                           # 依赖锁定文件
├── cmd/                             # 命令定义（Cobra 子命令）
│   ├── root.go                      # 根命令和全局设置
│   └── frontmatter.go               # front-matter 子命令
├── internal/                        # 内部包（不导出）
│   ├── frontmatter/                 # Front matter 解析逻辑
│   │   ├── parser.go                # 核心 YAML 解析和目录遍历
│   │   └── parser_test.go           # 测试
│   └── logger/                      # 日志系统
│       ├── logger.go                # 日志初始化和配置
│       └── logger_test.go           # 测试
├── testdata/                        # 测试数据文件
├── docs/                            # 文档（开发者指南、架构等）
└── README.md                        # 本文件
```

### 核心模块

#### `cmd/` - 命令层

- **`root.go`**：定义根命令，初始化全局日志系统
  - 所有子命令继承根命令的配置
  - `initLogger` 函数在命令执行前初始化日志

- **`frontmatter.go`**：实现 `front-matter` 子命令
  - 解析命令行标志（输出文件、递归选项）
  - 调用 `frontmatter` 包的核心逻辑
  - 处理 JSON 序列化和输出

#### `internal/frontmatter/` - Front Matter 解析

- **`parser.go`**：核心业务逻辑
  - `ParseFile(path string)`：解析单个 Markdown 文件的 front matter
    - 支持标准 YAML front matter 格式（`---` 分隔符）
    - 处理 Windows (`\r\n`) 和 Unix (`\n`) 风格换行符
    - 如果文件无 front matter，返回 `nil, nil`
  - `ReadDir(dir string, recursive bool)`：遍历目录
    - 支持递归和非递归两种模式
    - 收集成功解析的结果和出错信息
    - 继续处理其他文件（非致命错误处理）

#### `internal/logger/` - 日志系统

- **`logger.go`**：日志初始化和管理
  - 使用 `slog` 作为日志系统
  - 日志存储在 `{TempDir}/oh-my-markdown/oh-my-markdown.log`
  - 自动轮转：5 天后删除、gzip 压缩
  - `Init()` 函数在应用启动时调用
  - `Dir()` 函数返回日志目录路径

### 构建和测试

#### 构建项目

```bash
go build -o oh-my-markdown .
```

#### 运行所有测试

```bash
go test ./...
```

#### 运行测试并显示覆盖率

```bash
go test -cover ./...
```

#### 运行测试（检测数据竞争）

```bash
go test -race ./...
```

#### 运行单个测试

```bash
go test -run TestParseFile ./internal/frontmatter
```

#### 详细测试输出

```bash
go test -v ./...
```

### 代码质量检查

```bash
# 代码格式化
go fmt ./...

# 导入语句优化
goimports -w .

# 静态分析
go vet ./...

# 安全扫描
gosec ./...

# 扩展静态检查（如已配置）
golangci-lint run ./...
```

## 编码规范

本项目遵循以下编码规范：

### 错误处理

始终使用上下文信息包装错误：

```go
if err != nil {
    return fmt.Errorf("failed to process markdown: %w", err)
}
```

### 接口设计

- 保持接口简洁（1-3 个方法）
- 接受接口，返回结构体
- 在使用接口的地方定义接口，而不是在实现处

### 不可变性

- 优先使用不可变数据结构
- 避免在导出函数中使用可变全局状态
- 使用值接收器而非指针接收器进行不可变操作

### 测试

- 使用表驱动测试（Go 标准做法）
- 新代码目标覆盖率 ≥80%
- 始终使用 `-race` 标志运行测试以检测并发问题

### 格式化

使用 `gofmt` 和 `goimports` 进行强制格式化（无风格辩论）。

## 依赖管理

本项目使用最小化的依赖集：

| 包 | 版本 | 用途 |
|----|------|------|
| `github.com/spf13/cobra` | v1.10.2 | CLI 框架 |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML 解析 |
| `gopkg.in/natefinch/lumberjack.v2` | v2.2.1 | 日志轮转 |

添加新依赖时，请确保：
- 包的维护状态良好
- 有充分的测试
- 与项目目标相符

## 日志配置

应用的日志配置位于 `internal/logger/logger.go`：

- **日志位置**：`{TempDir}/oh-my-markdown/oh-my-markdown.log`
- **日志级别**：DEBUG
- **日志格式**：JSON
- **轮转策略**：
  - 超过 5 天的日志自动删除
  - 旧日志自动 gzip 压缩
  - 备份文件数量无限制（由 MaxAge 自动清理）

### 访问日志

```bash
# Windows
# 日志位置：C:\Users\<YourUsername>\AppData\Local\Temp\oh-my-markdown\oh-my-markdown.log
```

## 常见问题

### Q：如何获取特定格式的输出？
**A：** 使用 `front-matter` 命令的 `-o` 标志将输出保存为文件，然后可以用其他工具处理 JSON 数据。

### Q：递归选项的默认值是什么？
**A：** 默认为 `true`（递归遍历子目录）。使用 `--recursive=false` 或 `-r=false` 禁用递归。

### Q：如何处理无 front matter 的文件？
**A：** 这些文件会被自动跳过，不会出现在输出结果中。

### Q：错误发生时会发生什么？
**A：** 
- **文件级错误**：记录到日志，处理继续进行
- **致命错误**：中断处理并返回错误信息

## 开发与构建

### 构建发布版本

```bash
# 仅限 Windows 64-bit 平台
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o oh-my-markdown.exe .
```

## 贡献指南

欢迎提交问题和拉取请求！请确保：

1. 遵循项目的编码规范
2. 为新功能添加测试
3. 保持测试覆盖率 ≥80%
4. 运行 `go fmt ./...` 格式化代码
5. 提交前通过所有测试：`go test -race ./...`

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 联系方式

如有任何问题或建议，请提交 Issue 或 PR。