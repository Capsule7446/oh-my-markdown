# oh-my-markdown 代码地图

**最后更新**：2026-04-19  
**项目语言**：Go 1.25.0+  
**CLI 框架**：Cobra 1.10.2

## 概览

本文档提供了 oh-my-markdown 项目的完整代码地图，包括模块关系、函数流程、依赖关系等。

## 快速导航

| 模块 | 文件 | 入口 | 责任 |
|------|------|------|------|
| **CLI 层** | `cmd/root.go` | `Execute()` | 命令路由、日志初始化 |
| **前置元数据** | `cmd/frontmatter.go` | `runFrontMatter()` | front-matter 命令实现 |
| **业务逻辑** | `internal/frontmatter/parser.go` | `ReadDir()`, `ParseFile()` | YAML 解析、目录遍历 |
| **日志系统** | `internal/logger/logger.go` | `Init()` | 日志初始化、轮转配置 |

---

## 模块详解

### 1. 命令层 (cmd/)

#### cmd/root.go

**导出接口**：
```go
var rootCmd *cobra.Command        // Cobra 根命令
func Execute() error              // 应用入口（在 main.go 中调用）
func initLogger(cmd, args) error  // 日志初始化钩子
```

**调用链**：
```
main.go
└─> Execute()
    └─> rootCmd.Execute() (Cobra)
        ├─> PersistentPreRunE: initLogger()
        │   └─> logger.Init()
        └─> 子命令路由
```

**关键代码**：
```go
var rootCmd = &cobra.Command{
    Use:              "oh-my-markdown",
    Short:            "批次處理 Markdown 的工具",
    SilenceErrors:    true,        // 禁用自动错误输出
    SilenceUsage:     true,        // 隐藏使用帮助
    PersistentPreRunE: initLogger, // 每个命令执行前都运行
}

func Execute() error {
    return rootCmd.Execute()
}
```

**测试状态**：❌ 无测试文件（E2E 测试在计划中）

#### cmd/frontmatter.go

**导出接口**：
```go
var frontMatterCmd *cobra.Command       // front-matter 子命令
func runFrontMatter(cmd, args) error    // 命令处理函数
```

**参数解析**：

| 位置参数 | 说明 |
|---------|------|
| `<directory>` | 必需，要扫描的目录 |

| 标志 | 简写 | 类型 | 默认 | 说明 |
|------|------|------|------|------|
| `--output` | `-o` | string | 空 | 输出文件路径（stdout） |
| `--recursive` | `-r` | bool | true | 递归遍历子目录 |

**执行流程**：
```
runFrontMatter(cmd *cobra.Command, args []string)
│
├─ 1. 获取位置参数
│     └─> dir := args[0]
│
├─ 2. 解析标志
│     ├─> outputFile, _ := cmd.Flags().GetString("output")
│     └─> recursive, _ := cmd.Flags().GetBool("recursive")
│
├─ 3. 调用业务逻辑
│     └─> result, err := frontmatter.ReadDir(dir, recursive)
│
├─ 4. 处理错误
│     ├─> 文件级错误：记录到 slog.Warn()
│     └─> 致命错误：返回包装的错误
│
├─ 5. 序列化为 JSON
│     └─> jsonData, _ := json.MarshalIndent(result.Results)
│
└─ 6. 输出结果
      ├─> if outputFile != "": os.WriteFile()
      └─> else: fmt.Fprintln(stdout)
```

**关键代码**：
```go
var frontMatterCmd = &cobra.Command{
    Use:   "front-matter <directory>",
    Short: "读取目录下所有 Markdown 文件的 front matter 并输出 JSON",
    Args:  cobra.ExactArgs(1),
    RunE:  runFrontMatter,
}

func init() {
    rootCmd.AddCommand(frontMatterCmd)
    frontMatterCmd.Flags().StringP("output", "o", "", "输出文件路径")
    frontMatterCmd.Flags().BoolP("recursive", "r", true, "递归遍历子目录")
}

func runFrontMatter(cmd *cobra.Command, args []string) error {
    dir := args[0]
    outputFile, _ := cmd.Flags().GetString("output")
    recursive, _ := cmd.Flags().GetBool("recursive")
    
    slog.Info("开始读取 front matter", "dir", dir)
    
    result, err := frontmatter.ReadDir(dir, recursive)
    if err != nil {
        return fmt.Errorf("failed to read front matter: %w", err)
    }
    
    jsonData, _ := json.MarshalIndent(result.Results, "", "  ")
    
    if outputFile != "" {
        os.WriteFile(outputFile, jsonData, 0644)
    } else {
        fmt.Fprintln(cmd.OutOrStdout(), string(jsonData))
    }
    return nil
}
```

**测试状态**：❌ 无测试文件

---

### 2. 业务逻辑层 (internal/frontmatter/)

#### internal/frontmatter/parser.go

**导出类型**：
```go
type FrontMatter map[string]any              // Markdown 前置元数据

type ReadDirResult struct {
    Results []FrontMatter      // 解析成功的前置元数据列表
    Errors  []string           // 文件级错误信息
}
```

**导出函数**：
```go
func ParseFile(path string) (FrontMatter, error)
func ReadDir(dir string, recursive bool) (*ReadDirResult, error)
```

#### ParseFile(path string) (FrontMatter, error)

**职责**：解析单个 Markdown 文件的 front matter

**参数**：
- `path`：文件路径（绝对路径或相对路径）

**返回值**：
- `FrontMatter`：解析的前置元数据映射（无 front matter 返回 nil）
- `error`：文件读取或 YAML 解析失败时返回包装的错误

**算法**（伪代码）：

```
1. 读取文件内容到内存
2. 检查文件是否为空或不以 "---" 开头
   如果是：返回 (nil, nil)
3. 去掉开头的 "---\n" 或 "---\r\n"
4. 按 \n 分行
5. 逐行扫描，寻找行首的 "---"
   如果找到：
     - 提取之前的所有行为 front matter 内容
     - 使用 yaml.Unmarshal() 解析 YAML
     - 返回 (FrontMatter, nil)
   如果未找到：
     - 返回 (nil, nil)
```

**特殊处理**：
- ✅ Windows 风格换行符 (`\r\n`)
- ✅ Unix 风格换行符 (`\n`)
- ✅ 无 front matter 的文件（不视为错误）
- ✅ 无效的 YAML 格式（返回错误，由调用者处理）

**关键代码**：
```go
func ParseFile(path string) (FrontMatter, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }

    // 检查前置元数据块
    if len(content) == 0 || !bytes.HasPrefix(content, []byte("---")) {
        return nil, nil
    }

    // 去掉开头的 "---\n"
    rest := content[3:]
    if len(rest) > 0 && rest[0] == '\n' {
        rest = rest[1:]
    } else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
        rest = rest[2:]
    }

    // 查找结束符
    lines := bytes.Split(rest, []byte("\n"))
    var frontmatterContent []byte
    var foundEnd bool

    for i, line := range lines {
        trimmedLine := bytes.TrimRight(line, "\r")
        if bytes.Equal(trimmedLine, []byte("---")) {
            if i > 0 {
                frontmatterContent = bytes.Join(lines[:i], []byte("\n"))
            }
            foundEnd = true
            break
        }
    }

    if !foundEnd {
        return nil, nil
    }

    // 解析 YAML
    var data map[string]any
    if err := yaml.Unmarshal(frontmatterContent, &data); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return FrontMatter(data), nil
}
```

**性能特点**：
- 时间复杂度：O(n)（n = 文件大小）
- 空间复杂度：O(n)（完整读取到内存）
- 适用场景：<1MB 的文件

**测试覆盖**：✅ 表驱动测试，覆盖率 >80%

#### ReadDir(dir string, recursive bool) (*ReadDirResult, error)

**职责**：遍历目录，提取所有 Markdown 文件的 front matter

**参数**：
- `dir`：目录路径
- `recursive`：是否递归遍历子目录

**返回值**：
- `*ReadDirResult`：包含解析成功的结果列表和错误列表
- `error`：目录访问失败时返回包装的错误

**算法**（伪代码）：

```
1. 初始化 results 切片和 errors 切片
2. 使用 filepath.WalkDir() 遍历目录树：
   对于每个条目 (path, DirEntry, error):
     a. 如果是目录：
        - 如果非递归且不是根目录：SkipDir
        - 否则：继续
     b. 如果不是 .md 文件：跳过
     c. 调用 ParseFile(path)：
        - 如果成功且不为 nil：append 到 results
        - 如果出错：append 错误信息到 errors，继续
3. 返回 ReadDirResult{results, errors}
```

**特殊处理**：
- 非致命错误：单个文件的解析错误不中断整个过程
- 目录遍历：支持递归和非递归两种模式
- 相对路径：记录相对于根目录的路径到错误信息

**关键代码**：
```go
func ReadDir(dir string, recursive bool) (*ReadDirResult, error) {
    var results []FrontMatter
    var errors []string

    walkFunc := func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            if !recursive && path != dir {
                return filepath.SkipDir
            }
            return nil
        }

        if !strings.HasSuffix(path, ".md") {
            return nil
        }

        fm, err := ParseFile(path)
        if err != nil {
            relPath, _ := filepath.Rel(dir, path)
            errors = append(errors, fmt.Sprintf("%s: %v", relPath, err))
            return nil
        }

        if fm != nil {
            results = append(results, fm)
        }

        return nil
    }

    if err := filepath.WalkDir(dir, walkFunc); err != nil {
        return nil, fmt.Errorf("failed to walk directory: %w", err)
    }

    return &ReadDirResult{Results: results, Errors: errors}, nil
}
```

**性能特点**：
- 时间复杂度：O(n)（n = 文件数）
- 空间复杂度：O(m)（m = 有 front matter 的文件数）
- 适用场景：<10,000 个文件

**测试覆盖**：✅ 表驱动测试，覆盖率 >80%

---

### 3. 日志系统层 (internal/logger/)

#### internal/logger/logger.go

**导出函数**：
```go
func Init() error         // 初始化日志系统
func Dir() string         // 返回日志目录路径（供测试使用）
```

#### Init() error

**职责**：初始化全局日志系统

**初始化步骤**：

```
1. 构建日志目录路径
   └─> {TempDir}/oh-my-markdown
   
2. 创建目录（如不存在）
   └─> os.MkdirAll(logDir, 0755)
   
3. 配置日志轮转器 (lumberjack.Logger)
   ├─> Filename: {logDir}/oh-my-markdown.log
   ├─> MaxAge: 5 天
   ├─> Compress: true (gzip)
   └─> MaxBackups: 0 (由 MaxAge 管理)
   
4. 创建 slog JSON 处理器
   ├─> 输出格式：JSON
   ├─> 日志级别：DEBUG
   └─> 输出目标：日志轮转器
   
5. 设置为全局默认 logger
   └─> slog.SetDefault(slog.New(handler))
   
6. 写入初始化日志
   └─> slog.Info("Logger initialized", "path", logPath)
```

**关键配置**：

| 配置项 | 值 | 说明 |
|--------|-----|------|
| 日志位置 | `{TempDir}/oh-my-markdown/oh-my-markdown.log` | Unix: `/tmp/oh-my-markdown/oh-my-markdown.log`; Windows: `C:\...\Temp\oh-my-markdown\oh-my-markdown.log` |
| 日志格式 | JSON | 便于后续解析和分析 |
| 日志级别 | DEBUG | 记录所有调试信息 |
| 轮转周期 | 5 天 | 超过 5 天的日志自动删除 |
| 压缩 | 启用 | 使用 gzip 压缩旧日志 |
| 备份数 | 0 | 无限制（由 MaxAge 控制） |

**关键代码**：
```go
func Init() error {
    logDir = filepath.Join(os.TempDir(), "oh-my-markdown")
    os.MkdirAll(logDir, 0755)
    logPath := filepath.Join(logDir, "oh-my-markdown.log")

    lj := &lumberjack.Logger{
        Filename:   logPath,
        MaxAge:     5,
        Compress:   true,
        MaxBackups: 0,
    }

    handler := slog.NewJSONHandler(lj, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })

    slog.SetDefault(slog.New(handler))
    slog.Info("Logger initialized", "path", logPath)
    return nil
}
```

**日志使用示例**：
```go
slog.Info("开始处理", "count", 100)
slog.Warn("可能的问题", "detail", "value")
slog.Error("处理失败", "error", err)
slog.Debug("详细信息", "variable", value)
```

**测试覆盖**：✅ 测试初始化和目录创建

#### Dir() string

**职责**：返回日志目录路径（供测试和调试使用）

**实现**：
```go
func Dir() string {
    if logDir == "" {
        return filepath.Join(os.TempDir(), "oh-my-markdown")
    }
    return logDir
}
```

---

## 依赖关系图

```
oh-my-markdown/
│
├─> main.go
│   └─> cmd.Execute()
│       └─> cmd/root.go
│           ├─> logger.Init()
│           │   └─> internal/logger/logger.go
│           │       ├─> slog (标准库)
│           │       └─> lumberjack
│           │
│           └─> cmd/frontmatter.go
│               └─> frontmatter.ReadDir()
│                   └─> internal/frontmatter/parser.go
│                       ├─> bytes (标准库)
│                       ├─> os (标准库)
│                       ├─> filepath (标准库)
│                       ├─> strings (标准库)
│                       └─> yaml.v3
│
└─> go.mod
    ├─> github.com/spf13/cobra v1.10.2
    ├─> gopkg.in/yaml.v3 v3.0.1
    └─> gopkg.in/natefinch/lumberjack.v2 v2.2.1
```

---

## 调用链路分析

### 场景：用户执行 `oh-my-markdown front-matter ./docs -o output.json`

```
1. 启动
   └─> main.main()
       └─> cmd.Execute()

2. 初始化
   └─> rootCmd.Execute() (Cobra)
       └─> rootCmd.PersistentPreRunE()
           └─> initLogger()
               └─> logger.Init()
                   └─> 创建日志目录和配置

3. 命令路由
   └─> Cobra 解析标志和参数
       ├─> 目录：./docs
       ├─> 标志：-o output.json
       └─> 路由到 frontMatterCmd

4. 命令执行
   └─> frontmatter.go:runFrontMatter()
       │
       ├─> 获取参数
       │   ├─> dir = "./docs"
       │   ├─> outputFile = "output.json"
       │   └─> recursive = true (默认)
       │
       ├─> 调用业务逻辑
       │   └─> frontmatter.ReadDir("./docs", true)
       │       │
       │       ├─> filepath.WalkDir() 遍历目录
       │       │   └─ 对每个 .md 文件：
       │       │      └─> ParseFile(path)
       │       │          ├─> 读取文件
       │       │          ├─> 查找 --- 分隔符
       │       │          └─> yaml.Unmarshal()
       │       │
       │       └─> 返回 ReadDirResult{Results, Errors}
       │
       ├─> 处理错误
       │   └─> 记录到 slog.Warn()
       │
       ├─> JSON 序列化
       │   └─> json.MarshalIndent(Results)
       │
       └─> 输出
           └─> os.WriteFile("output.json", jsonData, 0644)

5. 结果
   └─> output.json 包含所有前置元数据的 JSON 数组
```

---

## 数据结构详解

### FrontMatter

```go
type FrontMatter map[string]any
```

**示例**：
```go
FrontMatter{
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "title": "测试文章",
    "date": "2024-01-01",
    "tags": []any{"go", "markdown"},
    "author": map[string]any{
        "name": "张三",
        "email": "zhangsan@example.com",
    },
}
```

### ReadDirResult

```go
type ReadDirResult struct {
    Results []FrontMatter  // 成功解析的前置元数据
    Errors  []string       // 失败的文件和错误信息
}
```

**JSON 序列化示例**：
```json
[
  {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Article 1",
    "date": "2024-01-01"
  },
  {
    "title": "Article 2",
    "author": "John"
  }
]
```

---

## 文件交互矩阵

| 源文件 | 目标文件 | 交互方式 | 数据流向 |
|--------|---------|---------|---------|
| `main.go` | `cmd/root.go` | 函数调用 | Execute() → 返回 error |
| `cmd/root.go` | `internal/logger/logger.go` | 函数调用 | Init() → 返回 error |
| `cmd/frontmatter.go` | `internal/frontmatter/parser.go` | 函数调用 | ReadDir() → 返回 ReadDirResult |
| `internal/frontmatter/parser.go` | `gopkg.in/yaml.v3` | 库调用 | yaml.Unmarshal() → 返回 map |
| 任何文件 | `internal/logger/logger.go` | 库调用 | slog.Info/Warn/Error() → 写入日志 |

---

## 扩展点

### 添加新命令

在 `cmd/` 中创建新文件，遵循现有模式：

```go
var newCmd = &cobra.Command{
    Use: "newcmd <args>",
    Short: "命令描述",
    Args: cobra.ExactArgs(1),
    RunE: runNewCommand,
}

func init() {
    rootCmd.AddCommand(newCmd)
}

func runNewCommand(cmd *cobra.Command, args []string) error {
    // 实现逻辑
    return nil
}
```

### 添加新的业务包

在 `internal/` 中创建新包，导出所需的函数和类型。

### 自定义输出格式

修改 `cmd/frontmatter.go` 中的 JSON 序列化部分，支持不同格式（CSV、XML 等）。

---

## 测试覆盖情况

| 模块 | 文件 | 覆盖率 | 状态 |
|------|------|--------|------|
| `cmd/` | `root.go` | - | ❌ 无测试 |
| `cmd/` | `frontmatter.go` | - | ❌ 无测试 |
| `internal/frontmatter/` | `parser.go` | >80% | ✅ 表驱动测试 |
| `internal/logger/` | `logger.go` | >80% | ✅ 单元测试 |

### 测试文件位置

```
internal/frontmatter/parser_test.go
internal/logger/logger_test.go
```

### 运行测试

```bash
go test -race -cover ./...
```

---

## 相关资源

- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构文档
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发者指南
- [README.md](../README.md) - 项目概述
