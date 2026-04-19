# oh-my-markdown 快速开始指南

**最后更新**：2026-04-19

## 5 分钟快速开始

### 安装

```bash
# 克隆项目
git clone https://github.com/yourusername/oh-my-markdown.git
cd oh-my-markdown

# 构建（Windows）
go build -o oh-my-markdown.exe .

# 验证安装
oh-my-markdown.exe --help
```

### 基本用法

```bash
# 提取目录中的所有 front matter，输出到终端
oh-my-markdown.exe front-matter ./docs

# 保存结果到文件
oh-my-markdown.exe front-matter ./docs -o result.json

# 仅处理顶级文件（不递归）
oh-my-markdown.exe front-matter ./docs --recursive=false
```

---

## 常用命令

### front-matter 命令

**提取 Markdown 文件的 YAML 前置元数据并输出为 JSON**

#### 基本语法
```bash
oh-my-markdown front-matter <directory> [flags]
```

#### 快速示例

| 需求 | 命令 |
|------|------|
| 提取当前目录下的所有 front matter | `oh-my-markdown.exe front-matter .` |
| 提取 docs 文件夹的 front matter | `oh-my-markdown.exe front-matter ./docs` |
| 保存到指定文件 | `oh-my-markdown.exe front-matter ./docs -o output.json` |
| 仅处理顶级文件 | `oh-my-markdown.exe front-matter ./docs -r false` |
| 查看帮助信息 | `oh-my-markdown.exe front-matter --help` |

#### 标志详解

**`--output` / `-o`**
```bash
# 输出到文件
oh-my-markdown.exe front-matter ./docs -o result.json

# 输出到标准输出（默认）
oh-my-markdown.exe front-matter ./docs

# 也可以使用长形式
oh-my-markdown.exe front-matter ./docs --output=result.json
```

**`--recursive` / `-r`**
```bash
# 递归遍历所有子目录（默认）
oh-my-markdown.exe front-matter ./docs

# 等价于
oh-my-markdown.exe front-matter ./docs -r true

# 仅处理顶级目录
oh-my-markdown.exe front-matter ./docs -r false

# 或
oh-my-markdown.exe front-matter ./docs --recursive=false
```

#### 输出示例

**输入文件** (`docs/article.md`):
```markdown
---
uuid: "550e8400-e29b-41d4-a716-446655440000"
title: "Go 语言入门"
date: "2024-01-15"
author: "张三"
tags:
  - go
  - tutorial
---

# Go 语言入门

这是文章内容...
```

**命令**:
```bash
oh-my-markdown.exe front-matter ./docs
```

**输出** (JSON):
```json
[
  {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Go 语言入门",
    "date": "2024-01-15",
    "author": "张三",
    "tags": [
      "go",
      "tutorial"
    ]
  }
]
```

---

## 工作流示例

### 示例 1：批量提取博客 front matter

```bash
# 目录结构
blog/
├── 2024-01/
│   ├── article1.md
│   └── article2.md
└── 2024-02/
    └── article3.md

# 提取所有 front matter
oh-my-markdown.exe front-matter ./blog -o frontmatters.json

# 查看结果
type frontmatters.json
```

### 示例 2：处理特定目录

```bash
# 仅处理 blog 目录下的一级文件（不进入子文件夹）
oh-my-markdown.exe front-matter ./blog --recursive=false

# 输出：仅包含 blog/ 目录下的 .md 文件的 front matter
```

### 示例 3：与其他工具结合

```bash
# 提取并使用 jq 过滤特定字段（需要安装 jq）
oh-my-markdown.exe front-matter ./docs | jq '.[].title'

# 输出：
# "Go 语言入门"
# "Python 快速入门"
# ...

# 统计 front matter 数量
oh-my-markdown.exe front-matter ./docs | jq 'length'

# 输出：
# 42
```

---

## 开发工作流

### 第一次设置

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/oh-my-markdown.git
cd oh-my-markdown

# 2. 下载依赖
go mod download

# 3. 验证测试通过
go test -race ./...

# 4. 构建项目
go build -o oh-my-markdown.exe .
```

### 日常开发

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -run TestParseFile ./internal/frontmatter

# 检查代码格式
go fmt ./...

# 静态分析
go vet ./...

# 运行项目
go run main.go front-matter ./.test_data
```

### 快速测试流程

```bash
# 一行命令完整检查
go fmt ./... && go vet ./... && go test -race ./... && go build -o oh-my-markdown.exe .
```

---

## 常见问题速查

### Q: 如何查看日志？

**A:** 日志存储在 `{TempDir}/oh-my-markdown/oh-my-markdown.log`

**Windows 日志位置**：`C:\Users\<YourUsername>\AppData\Local\Temp\oh-my-markdown\oh-my-markdown.log`

使用记事本或任何文本编辑器打开此文件即可查看。

### Q: 如何处理无 front matter 的文件？

**A:** 这些文件会被自动跳过，不会出现在输出中。这不是错误，是正常行为。

### Q: 输出 JSON 为空表示什么？

**A:** 表示目录中没有文件含有有效的 front matter。

### Q: 如何处理包含中文的 front matter？

**A:** 完全支持 UTF-8 编码，无需特殊配置。

```bash
# 这个命令可以正确处理中文
./oh-my-markdown front-matter ./chinese-docs -o output.json
```

### Q: 文件很多时需要多久？

**A:** 取决于文件数量和大小。一般来说：
- 1,000 个小文件：< 1 秒
- 10,000 个小文件：1-5 秒
- 100,000 个小文件：10-30 秒

### Q: 如何只处理特定目录（不包含子目录）？

**A:** 使用 `--recursive=false` 标志

```bash
./oh-my-markdown front-matter ./docs --recursive=false
```

### Q: 输出文件已存在，会覆盖吗？

**A:** 是的，会直接覆盖。如需保留，请先备份或使用不同的文件名。

```bash
# 安全做法：使用时间戳
./oh-my-markdown front-matter ./docs -o "output_$(date +%Y%m%d_%H%M%S).json"
```

---

## 快速命令参考

### 构建和运行

```bash
go build -o oh-my-markdown.exe .                          # 构建（Windows）
oh-my-markdown.exe front-matter ./docs                    # 运行
go run main.go front-matter ./docs             # 直接运行（无需构建）
```

### 测试

```bash
go test ./...                                  # 运行所有测试
go test -race ./...                           # 检测数据竞争
go test -cover ./...                          # 显示覆盖率
go test -v ./...                              # 详细输出
go test -run TestParseFile ./internal/...     # 运行特定测试
```

### 代码质量

```bash
go fmt ./...                                   # 格式化代码
goimports -w .                                 # 修复 import
go vet ./...                                   # 静态检查
gosec ./...                                    # 安全扫描
golangci-lint run ./...                       # 综合检查
```

### 依赖管理

```bash
go mod download                                # 下载依赖
go mod tidy                                    # 清理依赖
go get -u ./...                               # 更新依赖
go mod verify                                 # 验证依赖
```

---

## 故障排除

### 问题：`go build` 失败

**解决**：
```bash
# 清理缓存
go clean -cache

# 重新下载依赖
go mod download

# 重新构建
go build -o oh-my-markdown .
```

### 问题：找不到文件

**解决**：
```bash
# 确保路径正确
./oh-my-markdown front-matter ./docs          # 相对路径
./oh-my-markdown front-matter /absolute/path  # 绝对路径

# 检查文件是否存在
ls -la ./docs
```

### 问题：输出为空数组 `[]`

**原因**：
1. 目录中没有 `.md` 文件
2. 文件中没有 front matter
3. 路径不存在或无访问权限

**调试**：
```bash
# 检查目录内容
dir /b ./docs

# 检查是否有 .md 文件
where /R ./docs *.md

# 运行工具时添加调试信息
oh-my-markdown.exe front-matter ./docs
```

---

## 下一步

- 📖 [完整 README](../README.md) - 详细的项目文档
- 🏗️ [架构文档](ARCHITECTURE.md) - 系统设计和模块关系
- 👨‍💻 [开发指南](DEVELOPMENT.md) - 如何为项目做贡献
- 🗺️ [代码地图](CODEMAPS.md) - 深入的代码结构分析

---

## 获得帮助

```bash
# 查看基本帮助
oh-my-markdown.exe --help

# 查看 front-matter 命令的帮助
oh-my-markdown.exe front-matter --help

# 查看特定标志的帮助
oh-my-markdown.exe front-matter -h
```

---

**提示**：将 `oh-my-markdown.exe` 添加到 Windows PATH，使其可从任何目录运行

1. 将 `oh-my-markdown.exe` 复制到某个目录，如 `C:\tools\`
2. 添加该目录到 Windows PATH 环境变量
3. 然后在任何目录下可以直接运行 `oh-my-markdown.exe`
