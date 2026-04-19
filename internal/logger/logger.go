package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var logDir string

// Init 初始化全局 slog logger，将日志写入 <TempDir>/oh-my-markdown/oh-my-markdown.log
// 日志配置：
// - 自动压缩（gzip）旧日志
// - 超过 5 天的日志自动删除
func Init() error {
	// 构建日志目录路径
	logDir = filepath.Join(os.TempDir(), "oh-my-markdown")

	// 创建目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "oh-my-markdown.log")

	// 配置 lumberjack 日志轮转
	lj := &lumberjack.Logger{
		Filename:   logPath,
		MaxAge:     5,    // 5 天后删除
		Compress:   true, // gzip 压缩
		MaxBackups: 0,    // 不限制备份文件数量，依赖 MaxAge 自动清理
	}

	// 创建 slog 处理器
	handler := slog.NewJSONHandler(lj, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// 设为默认 logger
	slog.SetDefault(slog.New(handler))

	// 写入初始日志确保文件被创建
	slog.Info("Logger initialized", "path", logPath)

	return nil
}

// Dir 返回日志目录路径（供测试和调试使用）
func Dir() string {
	if logDir == "" {
		return filepath.Join(os.TempDir(), "oh-my-markdown")
	}
	return logDir
}
