package cmd

import (
	"encoding/json"
	"fmt"
	"oh-my-markdown/internal/frontmatter"
	"os"

	"github.com/spf13/cobra"
)

var frontMatterCmd = &cobra.Command{
	Use:   "front-matter <directory>",
	Short: "读取目录下所有 Markdown 文件的 front matter 并输出 JSON",
	Args:  cobra.ExactArgs(1),
	RunE:  runFrontMatter,
}

func init() {
	rootCmd.AddCommand(frontMatterCmd)
	frontMatterCmd.Flags().StringP("output", "o", "", "输出文件路径（默认 stdout）")
	frontMatterCmd.Flags().BoolP("recursive", "r", true, "递归遍历子目录")
}

func runFrontMatter(cmd *cobra.Command, args []string) error {
	dir := args[0]

	// 从 cobra flags 读取参数值，避免包级别可变状态的竞态条件
	outputFile, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}

	recursive, err := cmd.Flags().GetBool("recursive")
	if err != nil {
		return fmt.Errorf("failed to get recursive flag: %w", err)
	}

	// 读取 front matter
	result, err := frontmatter.ReadDir(dir, recursive)
	if err != nil {
		return fmt.Errorf("failed to read front matter: %w", err)
	}

	// 如果有文件级别的错误，输出到 stderr
	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "警告：%d 个文件解析失败\n", len(result.Errors))
		for _, errMsg := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", errMsg)
		}
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(result.Results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 输出结果
	if outputFile != "" {
		// 写入文件
		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "已写入到文件：%s\n", outputFile); err != nil {
			return fmt.Errorf("failed to write output message: %w", err)
		}
	} else {
		// 输出到 stdout
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jsonData)); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}
