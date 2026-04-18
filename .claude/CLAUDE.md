# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在本仓库工作时提供指导。

## 项目概述

**oh-my-markdown** 是一个用 Go 编写的 CLI 工具，用于批量处理 Markdown 文件。该工具采用 Cobra 命令框架构建，旨在高效处理大规模
markdown 操作。

- **语言**：Go 1.25.0+
- **入口点**：`main.go` → `cmd.Execute()`
- **CLI 框架**：Cobra 1.10.2
- **二进制名称**：`omm`（oh-my-markdown）

## 项目规则

### 交流语言

**此项目内所有讨论、代码注释和文档应使用中文。** 代码本身保持 Go 语言的英文约定（函数名、变量名等），但所有非代码交流均使用中文。

## 快速开始

### 构建项目

```bash
go build -o omm .
```

### 运行工具

```bash
./omm <command> [flags]
```

### 运行测试（表驱动测试）

```bash
go test ./...           # 运行所有测试
go test -race ./...     # 运行时检测数据竞争
go test -cover ./...    # 显示覆盖率摘要
go test -v ./...        # 详细输出
```

### 运行单个测试

```bash
go test -run TestFunctionName ./...
```

### 代码质量检查

```bash
# 格式化代码（必需 - 使用 gofmt/goimports）
go fmt ./...

# 静态分析
go vet ./...

# 安全扫描
gosec ./...

# 扩展静态检查（如已配置）
golangci-lint run ./...
```

## 架构

### 目录结构

- **`main.go`** — 入口点，调用 `cmd.Execute()`
- **`cmd/root.go`** — Cobra 根命令定义和设置
- **`cmd/`** — 子命令目录（随功能扩展而增长）

### 执行流程

1. `main.go` 导入并调用 `cmd.Execute()`
2. `cmd.Execute()` 执行 Cobra 根命令
3. Cobra 解析命令行标志并路由到子命令
4. 错误处理会隐藏使用说明（参见 `SilenceUsage: true`）

### 添加新命令

添加子命令时的步骤：

1. 在 `cmd/` 中创建新文件（如 `cmd/process.go`）
2. 定义 Cobra 命令
3. 使用 `rootCmd.AddCommand()` 注册命令
4. 如果逻辑复杂，遵循函数式选项模式

## 编码规范

### 错误处理

始终用上下文信息包装错误：

```go
if err != nil {
return fmt.Errorf("failed to process markdown: %w", err)
}
```

### 测试

- 使用表驱动测试（Go 标准做法）
- 在同一包中使用 `_test.go` 文件
- 新代码的目标覆盖率 ≥80%
- 始终使用 `-race` 标志运行，以捕获并发问题

### 代码风格

- 使用 `gofmt` 和 `goimports`（强制）
- 保持接口简洁（1-3 个方法）
- 接受接口，返回结构体
- 导出函数中不使用可变状态（优先使用不可变性）

### 依赖管理

当前依赖最小化：

- `github.com/spf13/cobra` — CLI 框架
- `golang.org/x/sync` — 并发工具库

谨慎添加新依赖，并验证其维护状态良好。

## 配置与规则

项目遵循 `.claude/rules/` 中定义的规则：

- **common/** — 语言无关原则（编码风格、git 工作流、测试、安全）
- **golang/** — Go 特定约定（模式、钩子、安全）

强制执行的关键规则：

- TDD：先写测试，后写实现
- 禁止硬编码密钥和凭证
- 每个层级都要显式处理错误
- 最低测试覆盖率 80%
- 提交前进行代码审查

## 工作区设置

所有必要的规则、代理和配置都在 `.claude/` 中：

- `CLAUDE.md` — 本文件（Claude Code 指导）
- `rules/` — 编码标准和模式
- `agents.md` — 可用的专业代理（planner、code-reviewer、tdd-guide 等）

工作流程中的推荐做法：

1. **新功能开发** — 使用 **tdd-guide** 代理（先写测试）
2. **代码编写后** — 使用 **code-reviewer** 代理
3. **架构决策** — 使用 **architect** 代理
4. **遵循规则** — 参考 `rules/golang/` 中的 Go 编码风格规则

## 开发工作流

按以下顺序进行开发：

1. **研究与规划**
    - 使用 **planner** 代理制定实现计划
    - 在编码前生成规划文档

2. **TDD 方法**
    - 使用 **tdd-guide** 代理
    - 先写测试（RED）
    - 实现通过测试（GREEN）
    - 重构优化（IMPROVE）
    - 验证覆盖率 ≥80%

3. **代码审查**
    - 编写代码后立即使用 **code-reviewer** 代理
    - 解决所有关键和高优先级问题

4. **提交代码**
    - 编写详细的提交消息
    - 遵循约定式提交格式
    - 参考 `rules/common/git-workflow.md`