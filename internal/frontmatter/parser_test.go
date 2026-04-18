package frontmatter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		wantErr   bool
		wantNil   bool
		wantUUID  string
		wantTitle string
	}{
		{
			name:      "正常 YAML front matter 文件",
			filePath:  ".test_data/valid.md",
			wantErr:   false,
			wantNil:   false,
			wantUUID:  "550e8400-e29b-41d4-a716-446655440000",
			wantTitle: "测试文章",
		},
		{
			name:     "无 front matter 的文件",
			filePath: ".test_data/no-frontmatter.md",
			wantErr:  false,
			wantNil:  true,
		},
		{
			name:     "空文件",
			filePath: ".test_data/empty.md",
			wantErr:  false,
			wantNil:  true,
		},
		{
			name:     "无效 YAML 格式",
			filePath: ".test_data/invalid-yaml.md",
			wantErr:  true,
			wantNil:  true,
		},
		{
			name:      "Windows 风格的换行符（\\r\\n）",
			filePath:  ".test_data/crlf.md",
			wantErr:   false,
			wantNil:   false,
			wantUUID:  "770e8400-e29b-41d4-a716-446655440002",
			wantTitle: "CRLF 文章",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调整路径使其相对于项目根目录
			absPath := filepath.Join("..", "..", tt.filePath)

			// 检查文件是否存在
			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				t.Fatalf("测试文件不存在：%s", absPath)
			}

			got, err := ParseFile(absPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("ParseFile() got = %v, wantNil %v", got, tt.wantNil)
				return
			}

			if !tt.wantNil && got != nil {
				if got["uuid"] != tt.wantUUID {
					t.Errorf("ParseFile() uuid = %v, want %v", got["uuid"], tt.wantUUID)
				}
				if got["title"] != tt.wantTitle {
					t.Errorf("ParseFile() title = %v, want %v", got["title"], tt.wantTitle)
				}
			}
		})
	}
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		recursive      bool
		wantErr        bool
		wantCount      int
		wantUUIDs      map[string]bool
		wantErrorCount int
	}{
		{
			name:           "读取 .test_data 目录（非递归）",
			dir:            ".test_data",
			recursive:      false,
			wantErr:        false,
			wantCount:      2, // valid.md 和 crlf.md 有 front matter
			wantErrorCount: 1, // invalid-yaml.md 解析失败
			wantUUIDs: map[string]bool{
				"550e8400-e29b-41d4-a716-446655440000": true,
				"770e8400-e29b-41d4-a716-446655440002": true,
			},
		},
		{
			name:           "读取 .test_data 目录（递归）",
			dir:            ".test_data",
			recursive:      true,
			wantErr:        false,
			wantCount:      3, // valid.md、crlf.md 和 subdir/nested.md
			wantErrorCount: 1, // invalid-yaml.md 解析失败
			wantUUIDs: map[string]bool{
				"550e8400-e29b-41d4-a716-446655440000": true,
				"660e8400-e29b-41d4-a716-446655440001": true,
				"770e8400-e29b-41d4-a716-446655440002": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调整路径相对于项目根目录
			absPath := filepath.Join("..", "..", tt.dir)

			got, err := ReadDir(absPath, tt.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Results) != tt.wantCount {
				t.Errorf("ReadDir() got %d results, want %d", len(got.Results), tt.wantCount)
				return
			}

			if len(got.Errors) != tt.wantErrorCount {
				t.Errorf("ReadDir() got %d errors, want %d", len(got.Errors), tt.wantErrorCount)
			}

			// 验证所有 uuid 都在期望的列表中
			for _, fm := range got.Results {
				uuid, ok := fm["uuid"]
				if !ok {
					t.Errorf("ReadDir() result missing uuid field")
					continue
				}
				uuidStr, ok := uuid.(string)
				if !ok {
					t.Errorf("ReadDir() uuid field is not a string, got %T: %v", uuid, uuid)
					continue
				}
				if !tt.wantUUIDs[uuidStr] {
					t.Errorf("ReadDir() unexpected uuid = %v", uuidStr)
				}
			}
		})
	}
}
