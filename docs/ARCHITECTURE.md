# oh-my-markdown 架构文档

**最后更新**：2026-04-18

## 概述

`oh-my-markdown` 是一个用 Go 编写的 CLI 工具，专门用于批量处理 Markdown 文件。该工具使用 Cobra 框架实现命令行界面，采用模块化架构设计，具有清晰的职责分离。

## 系统架构

```
┌─────────────────────────────────────────────────────┐
│                   CLI 用户界面                        │
│              (Cobra - 命令行框架)                     │
└────────────────────┬────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
   ┌────▼─────┐ ┌──▼────┐  ┌───▼──────┐
   │ root.go  │ │front- │  │  其他命令  │
   │          │ │matter.│  │  (扩展)   │
   │(全局配置)│ │go     │  └───────────┘
   └────┬─────┘ └──┬────┘
        │          │
   ┌───▼──────────▼────────────────┐
   │   内部业务逻辑层 (internal/)   │
   └───────────┬─────────────────┬─┘
               │                 │
        ┌──────▼──────┐   ┌──────▼──────┐
        │ frontmatter │   │   logger    │
        │   package   │   │   package   │
        └──────┬──────┘   └──────┬──────┘
               │                 │
        ┌──────▼──────┐   ┌──────▼──────────┐
        │ parser.go   │   │ logger.go       │
        │             │   │                 │
        │ · ParseFile │   │ · Init()        │
        │ · ReadDir   │   │ · Dir()         │
        └─────────────┘   └─────────────────┘
               │
        ┌──────▼──────────────────┐
        │  标准库 & 第三方依赖    │
        │  · encoding/json        │
        │  · gopkg.in/yaml.v3     │
        │  · lumberjack (日志)    │
        └─────────────────────────┘
```

## 目录结构

```
oh-my-markdown/
├── main.go                          # 应用入口点
├── README.md                        # 项目概述
├── go.mod                           # Go 模块定义
├── go.sum                           # 依赖锁定文件
│
├── cmd/                             # 命令层（Cobra 命令定义）
│   ├── root.go                      # 根命令及全局初始化
│   └── frontmatter.go               # front-matter 子命令
│
├── internal/                        # 内部包（未导出）
│   ├── frontmatter/                 # Front matter 解析模块
│   │   ├── parser.go                # 核心 YAML 解析和目录遍历
│   │   ├── parser_test.go           # 表驱动单元测试
│   │   └── testdata/                # 测试数据文件
│   │
│   └── logger/                      # 日志系统模块
│       ├── logger.go                # 日志初始化和管理
│       └── logger_test.go           # 日志系统测试
│
├── testdata/                        # 项目级测试数据
│   ├── valid.md                     # 有效的 front matter 文件
│   ├── empty.md                     # 空文件
│   ├── no-frontmatter.md            # 无 front matter 的文件
│   ├── invalid-yaml.md              # 无效的 YAML
│   ├── crlf.md                      # Windows 换行符测试
│   └── subdir/
│       └── nested.md                # 嵌套目录测试文件
│
└── docs/                            # 文档目录
    ├── ARCHITECTURE.md              # 本文件（架构设计）
    ├── DEVELOPMENT.md               # 开发者指南
    └── CODEMAPS.md                  # 代码地图
```

## 执行流程

### 应用启动流程

```
1. main.go
   └─> cmd.Execute()
       └─> rootCmd.Execute() (Cobra 框架)
           │
           ├─> PersistentPreRunE: initLogger()
           │   └─> logger.Init()
           │       └─> 初始化 slog，配置日志轮转
           │
           └─> 路由到子命令
               └─> front-matter 命令
                   └─> runFrontMatter()
```

### Front Matter 处理流程

```
runFrontMatter(cmd, args)
│
├─ 1. 解析命令行参数
│     ├─ directory: 从位置参数获取
│     ├─ output: 从 -o 标志获取
│     └─ recursive: 从 -r 标志获取
│
├─ 2. 调用核心业务逻辑
│     └─ frontmatter.ReadDir(dir, recursive)
│        │
│        ├─ 初始化 results 和 errors 切片
│        │
│        ├─ filepath.WalkDir() 遍历目录
│        │  │
│        │  └─ 对于每个文件：
│        │     ├─ 跳过目录和非 .md 文件
│        │     ├─ 调用 ParseFile(path)
│        │     │  │
│        │     │  ├─ 读取文件内容
│        │     │  ├─ 检查 --- 开头
│        │     │  ├─ 查找结束符 ---
│        │     │  └─ 使用 yaml.Unmarshal() 解析 YAML
│        │     │
│        │     ├─ 如果成功：append 到 results
│        │     └─ 如果失败：append 错误信息
│        │
│        └─ 返回 ReadDirResult{Results, Errors}
│
├─ 3. JSON 序列化
│     └─ json.MarshalIndent(results)
│
└─ 4. 输出结果
      ├─ 如果指定 -o：写入文件 os.WriteFile()
      └─ 否则：输出到 stdout fmt.Fprintln()
```

## 核心模块详解

### 1. 命令层 (`cmd/`)

#### `root.go` - 根命令

**职责**：
- 定义根命令的基本属性（名称、描述）
- 初始化全局日志系统
- 作为所有子命令的父命令

**关键函数**：
- `initLogger(cmd, args)` - 在任何命令执行前初始化日志
- `Execute()` - 执行根命令

**设计考虑**：
```go
var rootCmd = &cobra.Command{
    Use:              "omm",
    Short:            "批次處理 Markdown 的工具",
    SilenceErrors:    true,     // 禁用自动错误输出
    SilenceUsage:     true,     // 隐藏使用说明
    PersistentPreRunE: initLogger, // 在任何命令前初始化
}
```

#### `frontmatter.go` - Front Matter 命令

**职责**：
- 定义 `front-matter` 子命令
- 解析命令行标志和参数
- 调用业务逻辑模块
- 处理输出

**关键函数**：
- `runFrontMatter(cmd, args)` - 命令执行函数

**标志设计**：
- `--output/-o`：输出文件路径（可选）
- `--recursive/-r`：是否递归（默认 true）

**错误处理策略**：
- 参数错误：直接返回
- 文件级错误：记录后继续
- JSON 序列化错误：返回错误信息

### 2. 业务逻辑层 (`internal/`)

#### `frontmatter/parser.go` - Front Matter 解析器

**数据结构**：

```go
// FrontMatter 代表单个 Markdown 文件的前置元数据
type FrontMatter map[string]any

// ReadDirResult 包含解析结果和错误信息
type ReadDirResult struct {
    Results []FrontMatter  // 成功解析的前置元数据列表
    Errors  []string       // 文件级错误信息
}
```

**核心函数**：

1. **ParseFile(path string) (FrontMatter, error)**
   - 职责：解析单个 Markdown 文件的 front matter
   - 算法：
     1. 读取文件内容
     2. 检查是否以 `---` 开头
     3. 寻找结束符 `---`（仅在行首匹配）
     4. 使用 `yaml.Unmarshal()` 解析 YAML
   - 特殊处理：
     - Windows 风格换行符（`\r\n`）
     - 无 front matter 文件（返回 nil, nil）
     - YAML 解析错误（返回包装的错误）

2. **ReadDir(dir string, recursive bool) (*ReadDirResult, error)**
   - 职责：递归遍历目录，收集所有 Markdown 文件的 front matter
   - 算法：使用 `filepath.WalkDir()`，对每个文件调用 `ParseFile()`
   - 特殊处理：
     - 非递归模式：跳过子目录
     - 文件级错误：记录但继续处理
     - 无 front matter 文件：跳过（不记录为错误）

**错误处理设计**：
- **致命错误**（返回错误）：无法访问目录、文件读取失败
- **非致命错误**（记录为警告）：单个文件的 YAML 解析失败
- **可恢复情况**（忽略）：无 front matter 的文件

#### `logger/logger.go` - 日志系统

**初始化流程**：

```go
func Init() error {
    1. 创建日志目录 {TempDir}/omm
    2. 配置 lumberjack 日志轮转器
       - MaxAge: 5 天后删除
       - Compress: 启用 gzip 压缩
       - MaxBackups: 0（由 MaxAge 控制）
    3. 创建 slog 处理器（JSON 格式，DEBUG 级别）
    4. 设置为全局默认 logger
    5. 写入初始化日志
}
```

**日志配置**：
- **输出格式**：JSON
- **日志级别**：DEBUG
- **日志文件**：`{TempDir}/omm/omm.log`
- **自动轮转**：5 天
- **压缩**：启用 gzip

## 数据流

### Front Matter 提取数据流

```
输入文件：
---
title: "My Article"
date: "2024-01-01"
tags: ["a", "b"]
---
Content here...

↓ ParseFile()
  - 读取文件 → bytes
  - 查找分隔符 → 提取 YAML 块
  - yaml.Unmarshal → map[string]any

↓ ReadDir()
  - 遍历目录
  - 调用 ParseFile() 对每个 .md 文件
  - 收集结果

↓ json.MarshalIndent()

输出：
[
  {
    "title": "My Article",
    "date": "2024-01-01",
    "tags": ["a", "b"]
  }
]

↓ 输出到 stdout 或文件
```

## 关键设计决策

### 1. 模块化架构
- **分离关注点**：命令层、业务逻辑层、日志系统分离
- **优势**：易于测试、易于扩展、易于维护

### 2. 错误处理策略
- **单文件错误不中断**：ReadDir() 收集错误但继续处理
- **致命错误直接返回**：目录访问失败则中止
- **错误日志**：文件级错误记录到日志，用户看到日志消息

### 3. 无状态设计
- **cmd 层无全局状态**：从 Cobra flags 读取参数
- **业务逻辑纯函数**：ReadDir() 和 ParseFile() 无副作用
- **优势**：支持并发调用、易于测试

### 4. 日志系统集中管理
- **单一全局实例**：slog 全局 logger
- **自动轮转**：无需手动管理旧日志
- **JSON 格式**：便于后续解析和分析

## 依赖关系

```
oh-my-markdown
│
├─> github.com/spf13/cobra v1.10.2
│   └─> CLI 框架，提供命令定义和标志解析
│
├─> gopkg.in/yaml.v3 v3.0.1
│   └─> YAML 解析，用于前置元数据解析
│
└─> gopkg.in/natefinch/lumberjack.v2 v2.2.1
    └─> 日志轮转，自动压缩和删除旧日志

标准库依赖：
├─> encoding/json - JSON 序列化
├─> os - 文件操作
├─> path/filepath - 路径处理
├─> bytes - 字节操作
├─> log/slog - 结构化日志
└─> fmt - 字符串格式化
```

## 扩展点

### 添加新命令

```go
// 1. 在 cmd/ 中创建新文件，如 cmd/process.go
package cmd

var processCmd = &cobra.Command{
    Use: "process <directory>",
    Short: "处理 Markdown 文件",
    Args: cobra.ExactArgs(1),
    RunE: runProcess,
}

func init() {
    rootCmd.AddCommand(processCmd)
    processCmd.Flags().StringP("output", "o", "", "输出文件路径")
}

func runProcess(cmd *cobra.Command, args []string) error {
    // 实现逻辑
    return nil
}
```

### 扩展 Front Matter 处理

当前系统支持的扩展点：

1. **新的解析器**：在 `internal/frontmatter/` 中添加新文件
   - 例如：`markdown_processor.go`、`yaml_processor.go`

2. **新的输出格式**：修改 `frontmatter.go` 中的序列化逻辑
   - 例如：CSV、XML、Protocol Buffers

3. **新的日志目标**：配置 `logger.go` 中的 slog 处理器
   - 例如：远程日志服务、数据库

## 测试架构

### 测试组织

```
internal/frontmatter/
├── parser.go
└── parser_test.go          # 表驱动测试

internal/logger/
├── logger.go
└── logger_test.go          # 初始化测试

testdata/                   # 共享测试数据
├── valid.md
├── empty.md
├── no-frontmatter.md
├── invalid-yaml.md
├── crlf.md
└── subdir/nested.md
```

### 测试覆盖范围

| 模块 | 测试覆盖率 | 测试类型 |
|------|-----------|---------|
| `parser.go` | >80% | 单元测试（表驱动） |
| `logger.go` | >80% | 单元测试 |
| `cmd/` | 无测试文件 | 集成测试（使用 e2e） |

## 性能考虑

### 当前设计

- **目录遍历**：使用 `filepath.WalkDir()`，O(n) 复杂度
- **文件读取**：全量读取到内存，适合中等大小文件
- **YAML 解析**：使用第三方库 `gopkg.in/yaml.v3`，经过优化

### 适用场景

- 小到中等规模的 Markdown 文件集（< 10,000 文件）
- 单个文件大小 < 1 MB
- 目录深度 < 20 级

### 潜在优化方向

1. **并发处理**：使用 goroutines 并行处理多个文件
2. **流式处理**：对于大型目录，使用流式处理而非一次全量加载
3. **缓存**：缓存解析结果，减少重复解析

## 安全考虑

### 当前措施

1. **路径安全**：使用 `filepath` 包，自动处理路径规范化
2. **错误信息**：不暴露敏感路径信息
3. **文件权限**：日志文件权限为 0644（读写所有者，读其他）

### 需要关注的方面

1. **符号链接**：当前不处理符号链接，可能进入循环
2. **大文件**：无文件大小限制，大文件可能导致内存溢出
3. **YAML 注入**：YAML 解析可能执行代码（如 Go 标签）

## 部署考虑

### 依赖项

- Go 1.25.0+（仅编译时）
- 无运行时依赖（单个可执行文件）

### 跨平台支持

- Linux：完全支持
- macOS：完全支持
- Windows：完全支持（已处理 CRLF 换行符）

### 日志目录权限

- Linux/macOS：`/tmp/omm/`（或 `$TMPDIR/omm/`）
- Windows：`C:\Users\<User>\AppData\Local\Temp\omm\`

## 相关文档

- [README.md](../README.md) - 项目概述和使用指南
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发者指南
- [CODEMAPS.md](CODEMAPS.md) - 代码地图和模块关系