package frontmatter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FrontMatter 代表单个 .md 文件的 front matter 内容
type FrontMatter map[string]any

// ParseFile 读取单个 .md 文件，提取 front matter；无 front matter 返回 nil, nil
func ParseFile(path string) (FrontMatter, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 检查是否为空或不以 --- 开头
	if len(content) == 0 || !bytes.HasPrefix(content, []byte("---")) {
		return nil, nil
	}

	// 去掉开头的 "---\n" 或 "---\r\n"
	rest := content[3:]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	// 按行查找结束符 "---"（只在行首匹配）
	lines := bytes.Split(rest, []byte("\n"))
	var frontmatterContent []byte
	var foundEnd bool

	for i, line := range lines {
		// 移除行尾的 \r（处理 Windows 风格的 \r\n）
		trimmedLine := bytes.TrimRight(line, "\r")

		if bytes.Equal(trimmedLine, []byte("---")) {
			// 找到了结束符，将前面的行合并为 front matter 内容
			if i > 0 {
				frontmatterContent = bytes.Join(lines[:i], []byte("\n"))
			}
			foundEnd = true
			break
		}
	}

	if !foundEnd {
		// 没有找到结束符，视为无 front matter
		return nil, nil
	}

	// 解析 YAML
	var data map[string]any
	if err := yaml.Unmarshal(frontmatterContent, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML front matter: %w", err)
	}

	return FrontMatter(data), nil
}

// ReadDirResult 包含成功的 FrontMatter 列表和解析过程中遇到的错误
type ReadDirResult struct {
	Results []FrontMatter
	Errors  []string // 包含出错文件的相对路径和错误信息
}

// ReadDir 遍历目录，返回所有 FrontMatter 列表（无 front matter 文件跳过）
// 返回的 ReadDirResult 包含成功结果和任何非致命的文件级别错误
func ReadDir(dir string, recursive bool) (*ReadDirResult, error) {
	var results []FrontMatter
	var errors []string

	walkFunc := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			// 如果非递归，跳过子目录
			if !recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 .md 文件
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// 解析文件
		fm, err := ParseFile(path)
		if err != nil {
			// 记录错误但继续处理其他文件
			relPath, _ := filepath.Rel(dir, path)
			errors = append(errors, fmt.Sprintf("%s: %v", relPath, err))
			return nil
		}

		// 跳过无 front matter 的文件
		if fm != nil {
			results = append(results, fm)
		}

		return nil
	}

	if err := filepath.WalkDir(dir, walkFunc); err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return &ReadDirResult{
		Results: results,
		Errors:  errors,
	}, nil
}
