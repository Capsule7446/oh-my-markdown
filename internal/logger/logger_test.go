package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Init 成功初始化日志",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理之前的初始化（如有）
			_ = os.RemoveAll(Dir())

			err := Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 验证日志目录存在
				if info, err := os.Stat(Dir()); err != nil || !info.IsDir() {
					t.Errorf("Init() 日志目录不存在或不是目录: %s", Dir())
				}

				// 验证日志文件存在
				logFile := filepath.Join(Dir(), "oh-my-markdown.log")
				if _, err := os.Stat(logFile); err != nil {
					t.Errorf("Init() 日志文件不存在: %s", logFile)
				}

				// 验证 slog.Default() 不是 no-op handler
				// 通过写一条日志来间接验证
				slog.Info("test message")

				// 读取日志文件验证内容
				content, err := os.ReadFile(logFile)
				if err != nil {
					t.Errorf("Init() 无法读取日志文件: %v", err)
					return
				}
				if len(content) == 0 {
					t.Errorf("Init() 日志文件为空")
				}
			}

			// 清理测试数据
			_ = os.RemoveAll(Dir())
		})
	}
}

func TestDir(t *testing.T) {
	dir := Dir()

	// 验证路径包含 oh-my-markdown 子目录
	if !filepath.HasPrefix(dir, os.TempDir()) {
		t.Errorf("Dir() 路径不在临时目录中: %s", dir)
	}

	// 验证路径以 oh-my-markdown 结尾（目录名）
	if filepath.Base(dir) != "oh-my-markdown" {
		t.Errorf("Dir() 路径不以 oh-my-markdown 结尾: %s", dir)
	}
}
