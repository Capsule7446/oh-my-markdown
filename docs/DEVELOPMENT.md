# oh-my-markdown 开发者指南

**最后更新**：2026-04-18

## 目录

1. [开发环境设置](#开发环境设置)
2. [项目结构导航](#项目结构导航)
3. [常用开发任务](#常用开发任务)
4. [编码规范](#编码规范)
5. [测试指南](#测试指南)
6. [调试技巧](#调试技巧)
7. [添加新功能](#添加新功能)

## 开发环境设置

### 前置要求

- **Go 版本**：1.25.0 或更高
  ```bash
  go version
  ```
- **Git**：用于版本控制
- **编辑器/IDE**：推荐 VS Code 或 GoLand

### 初始化开发环境

```bash
# 克隆项目
git clone https://github.com/yourusername/oh-my-markdown.git
cd oh-my-markdown

# 下载依赖
go mod download

# 构建项目（Windows）
go build -o oh-my-markdown.exe .

# 验证构建成功
oh-my-markdown.exe --help
```

### IDE 配置（VS Code）

推荐安装以下扩展：

1. **Go** - Official extension (golang.go)
2. **Makefile Tools** - (charliermarsh.ruff) - 可选

推荐的 `.vscode/settings.json` 配置：

```json
{
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.useLanguageServer": true,
  "go.coverageInViewGutter": true
}
```

## 项目结构导航

### 快速参考

| 目录 | 用途 | 修改频率 |
|------|------|---------|
| `cmd/` | CLI 命令定义 | 高（添加新命令） |
| `internal/frontmatter/` | 核心业务逻辑 | 中（算法改进） |
| `internal/logger/` | 日志系统 | 低（配置调整） |
| `testdata/` | 测试文件 | 中（添加测试用例） |
| `docs/` | 文档 | 中（维护文档） |

### 模块简介

#### `cmd/` - 命令层

**何时编辑**：
- 添加新的子命令
- 修改命令行标志
- 改进命令输出格式

**关键文件**：
- `root.go`：根命令，通常不需要修改
- `frontmatter.go`：front-matter 命令的实现

#### `internal/frontmatter/` - Front Matter 解析

**何时编辑**：
- 改进 YAML 解析逻辑
- 添加新的前置元数据字段验证
- 优化目录遍历算法

**关键函数**：
- `ParseFile()`：解析单个文件
- `ReadDir()`：遍历目录

#### `internal/logger/` - 日志系统

**何时编辑**：
- 调整日志级别
- 更改日志输出目录
- 修改日志轮转策略

**关键函数**：
- `Init()`：初始化日志系统

## 常用开发任务

### 构建和运行

```bash
# 构建项目（Windows）
go build -o oh-my-markdown.exe .

# 运行项目
oh-my-markdown.exe front-matter ./.test_data

# 构建并运行
go run main.go front-matter ./.test_data

# 构建发布版本（Windows 64-bit，优化）
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o oh-my-markdown.exe .
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行所有测试（检测数据竞争）
go test -race ./...

# 运行特定包的测试
go test -race ./internal/frontmatter

# 运行特定测试函数
go test -run TestParseFile ./internal/frontmatter

# 显示测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 详细测试输出
go test -v ./...

# 运行单个测试并显示日志
go test -v -run TestParseFile/正常_YAML_front_matter_文件 ./internal/frontmatter
```

### 代码质量检查

```bash
# 代码格式化（必需）
go fmt ./...

# 自动修复 import 语句
goimports -w .

# 静态分析
go vet ./...

# 安全扫描
gosec ./...

# 完整的代码检查流程
go fmt ./... && \
goimports -w . && \
go vet ./... && \
go test -race ./...
```

### 依赖管理

```bash
# 查看依赖树
go mod graph

# 添加新依赖
go get github.com/some/package@latest

# 更新依赖
go get -u ./...

# 清理未使用的依赖
go mod tidy

# 验证依赖
go mod verify

# 下载依赖到本地
go mod download
```

## 编码规范

### 错误处理

**必须遵守**：始终包装错误并提供上下文

```go
// ❌ 错误做法
return err

// ✅ 正确做法
return fmt.Errorf("failed to process file %s: %w", path, err)

// ✅ 多层包装
if err != nil {
    return fmt.Errorf("failed to read directory: %w", err)
}
```

**错误处理位置**：
- 文件 I/O：总是检查错误
- YAML 解析：总是检查错误
- 标志解析：总是检查错误

### 变量声明

```go
// ❌ 避免包级别的可变状态
var globalState map[string]string

// ✅ 使用局部变量或参数
func processData(state map[string]string) error {
    return nil
}
```

### 函数签名

```go
// ✅ 接受接口，返回结构体
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

// ❌ 不要接受结构体参数（如果该结构体有接口）
func NewService(repo *SqlRepository) *Service {
    return &Service{repo: repo}
}
```

### 命名约定

| 类型 | 约定 | 示例 |
|------|------|------|
| 包名 | 小写单词 | `frontmatter`, `logger` |
| 导出函数 | PascalCase | `ParseFile`, `ReadDir` |
| 导出变量 | PascalCase | `FrontMatter`, `ReadDirResult` |
| 私有函数 | camelCase | `parseYaml`, `readFile` |
| 常量 | UPPER_SNAKE_CASE | `DefaultMaxAge` |
| 接收器 | 单字母或简写 | `func (f FrontMatter) String()` |

### 注释规范

```go
// ✅ 导出函数必须有注释
// ParseFile 读取单个 .md 文件，提取 front matter
// 无 front matter 返回 nil, nil
func ParseFile(path string) (FrontMatter, error) {
    // ...
}

// ✅ 复杂逻辑需要说明意图
// 检查是否以 --- 开头（三个连字符）
if bytes.HasPrefix(content, []byte("---")) {
    // ...
}

// ❌ 避免无价值的注释
// 读取文件
content, err := os.ReadFile(path)
```

## 测试指南

### 测试文件位置

```
internal/frontmatter/
├── parser.go
└── parser_test.go           # 必须在同一包中

internal/logger/
├── logger.go
└── logger_test.go           # 必须在同一包中
```

### 表驱动测试模式

本项目使用表驱动测试。参考现有的 `parser_test.go`：

```go
func TestParseFile(t *testing.T) {
    tests := []struct {
        name    string      // 测试用例名称
        path    string      // 输入
        want    FrontMatter // 期望输出
        wantErr bool        // 是否期望错误
    }{
        {
            name: "正常 YAML front matter 文件",
            path: ".test_data/valid.md",
            want: FrontMatter{
                "uuid": "550e8400-e29b-41d4-a716-446655440000",
                "title": "测试文章",
            },
            wantErr: false,
        },
        {
            name: "无 front matter 的文件",
            path: ".test_data/no-frontmatter.md",
            want: nil,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseFile(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseFile() got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 运行特定测试

```bash
# 运行 TestParseFile 的所有子测试
go test -run TestParseFile ./internal/frontmatter

# 运行特定子测试
go test -run TestParseFile/正常_YAML_front_matter_文件 ./internal/frontmatter

# 显示详细输出
go test -v -run TestParseFile ./internal/frontmatter
```

### 添加新测试

1. **确定测试类型**：
   - 单元测试：测试单个函数
   - 集成测试：测试多个模块交互

2. **遵循 AAA 模式**：
   ```go
   t.Run("description", func(t *testing.T) {
       // Arrange - 准备测试数据
       input := "test input"
       
       // Act - 执行被测函数
       result, err := FunctionUnderTest(input)
       
       // Assert - 验证结果
       if err != nil {
           t.Fatalf("expected no error, got %v", err)
       }
       if result != expected {
           t.Errorf("got %v, want %v", result, expected)
       }
   })
   ```

3. **生成测试数据**：
   - 添加到 `testdata/` 目录
   - 使用有意义的文件名
   - 更新 `README` 说明测试文件用途

### 覆盖率检查

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看 HTML 报告
go tool cover -html=coverage.out

# 查看命令行覆盖率
go tool cover -func=coverage.out

# 检查特定包的覆盖率
go test -cover ./internal/frontmatter
```

**目标覆盖率**：>80%（新增代码必须）

## 调试技巧

### 使用日志调试

应用使用 `slog` 结构化日志。查看日志：

```bash
# Windows：日志位置
# C:\Users\<User>\AppData\Local\Temp\oh-my-markdown\oh-my-markdown.log

# 使用记事本或文本编辑器打开日志文件查看
```

### 添加临时调试日志

```go
import "log/slog"

// 在代码中添加
slog.Debug("processing file", "path", filePath, "size", fileSize)
slog.Warn("unexpected value", "value", value)
slog.Error("operation failed", "error", err)
```

### 使用 `go run` 调试

```bash
# 带参数运行，易于重复测试
go run main.go front-matter ./.test_data -o /tmp/output.json

# 带竞争条件检测
go run -race main.go front-matter ./.test_data
```

### IDE 调试器（VS Code）

在 `.vscode/launch.json` 中添加：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Process",
      "type": "go",
      "mode": "local",
      "request": "launch",
      "program": "${workspaceFolder}",
      "args": ["front-matter", "./.test_data"]
    }
  ]
}
```

## 添加新功能

### 添加新命令

**步骤 1**：在 `cmd/` 中创建新文件

```go
// cmd/newcommand.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "newcommand <args>",
    Short: "命令描述",
    Args:  cobra.ExactArgs(1),
    RunE:  runNewCommand,
}

func init() {
    rootCmd.AddCommand(newCmd)
    newCmd.Flags().StringP("output", "o", "", "输出文件路径")
}

func runNewCommand(cmd *cobra.Command, args []string) error {
    // 实现逻辑
    return nil
}
```

**步骤 2**：添加测试

```go
// cmd/newcommand_test.go
package cmd

import (
    "testing"
)

func TestNewCommand(t *testing.T) {
    // 实现测试
}
```

**步骤 3**：更新文档

- 更新 README.md 中的使用说明
- 在 ARCHITECTURE.md 中记录新命令
- 为命令添加帮助文本

### 修改现有命令

1. **在 `cmd/` 中的对应文件中修改**
2. **更新命令的 Short 和 Long 描述**
3. **添加或修改标志**
4. **编写测试验证新功能**
5. **更新文档**

### 添加新的包

**命名约定**：
- 在 `internal/` 下创建新包（不导出给外部）
- 包名小写，使用有意义的单词

**结构**：
```
internal/newpackage/
├── main.go          # 主逻辑
├── types.go         # 数据类型定义
├── helpers.go       # 辅助函数
└── main_test.go     # 测试文件
```

## 代码审查检查清单

在提交 PR 前，确保满足以下条件：

- [ ] 代码通过 `go fmt ./...`
- [ ] 通过 `go vet ./...`
- [ ] 新增代码有合理的注释
- [ ] 添加了必要的测试
- [ ] 测试覆盖率 ≥80%
- [ ] 所有测试通过：`go test -race ./...`
- [ ] 没有硬编码的敏感信息
- [ ] 函数长度 <50 行（复杂逻辑应拆分）
- [ ] 没有使用全局可变状态
- [ ] 错误都被妥善处理
- [ ] README/文档已更新

## 常见问题

### Q：如何快速验证代码修改？

**A**：使用以下命令链：
```bash
go fmt ./... && \
goimports -w . && \
go vet ./... && \
go test -race ./... && \
go test -cover ./...
```

### Q：如何查看特定函数的实现？

**A**：使用 `grep` 搜索：
```bash
grep -r "func ParseFile" internal/
```

### Q：如何调试失败的测试？

**A**：使用 `-v` 标志运行详细输出，或在测试中添加 `slog` 日志。

### Q：如何添加外部依赖？

**A**：
```bash
go get github.com/some/package@latest
go mod tidy
```

然后在代码中导入使用。确保在 PR 中说明为什么需要新依赖。

## 相关文档

- [README.md](../README.md) - 项目概述
- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构
- [CODEMAPS.md](CODEMAPS.md) - 代码地图
