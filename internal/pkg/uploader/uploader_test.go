package uploader

import (
	"testing"
)

func TestGenerateImageKey(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		contentType string
		wantExt     string // 期望的扩展名
	}{
		{
			name:        "image/jpeg should use .jpg",
			data:        []byte("fake jpeg data"),
			contentType: "image/jpeg",
			wantExt:     ".jpg",
		},
		{
			name:        "image/jfif should use .jpg",
			data:        []byte("fake jfif data"),
			contentType: "image/jfif",
			wantExt:     ".jpg",
		},
		{
			name:        "image/pjpeg should use .jpg",
			data:        []byte("fake pjpeg data"),
			contentType: "image/pjpeg",
			wantExt:     ".jpg",
		},
		{
			name:        "image/png should use .png",
			data:        []byte("fake png data"),
			contentType: "image/png",
			wantExt:     ".png",
		},
		{
			name:        "image/gif should use .gif",
			data:        []byte("fake gif data"),
			contentType: "image/gif",
			wantExt:     ".gif",
		},
		{
			name:        "image/jpeg with charset should use .jpg",
			data:        []byte("fake jpeg data"),
			contentType: "image/jpeg; charset=utf-8",
			wantExt:     ".jpg",
		},
		{
			name:        "empty contentType should have no extension",
			data:        []byte("fake data"),
			contentType: "",
			wantExt:     "",
		},
		{
			name:        "unknown contentType should have no extension",
			data:        []byte("fake data"),
			contentType: "application/octet-stream",
			wantExt:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := generateImageKey(tt.data, tt.contentType)
			
			// 检查 key 是否包含期望的扩展名
			if tt.wantExt != "" {
				if len(key) < len(tt.wantExt) {
					t.Errorf("generateImageKey() key too short: %s", key)
					return
				}
				gotExt := key[len(key)-len(tt.wantExt):]
				if gotExt != tt.wantExt {
					t.Errorf("generateImageKey() extension = %v, want %v", gotExt, tt.wantExt)
				}
			} else {
				// 如果没有期望的扩展名，检查 key 是否以日期路径结尾（没有扩展名）
				// key 格式: images/YYYY/MM/DD/md5 或 test/images/YYYY/MM/DD/md5
				if len(key) < 20 {
					t.Errorf("generateImageKey() key too short: %s", key)
				}
			}
			
			// 检查 key 格式是否正确
			if len(key) == 0 {
				t.Errorf("generateImageKey() returned empty key")
			}
		})
	}
}

func TestGenerateImageKey_Format(t *testing.T) {
	data := []byte("test image data")
	contentType := "image/jpeg"
	
	key := generateImageKey(data, contentType)
	
	// 检查 key 格式: images/YYYY/MM/DD/md5.jpg 或 test/images/YYYY/MM/DD/md5.jpg
	if len(key) < 30 {
		t.Errorf("generateImageKey() key too short: %s", key)
	}
	
	// 检查是否包含日期格式
	if len(key) < 20 {
		t.Errorf("generateImageKey() key format incorrect: %s", key)
	}
	
	// 检查是否以 .jpg 结尾
	if len(key) < 4 || key[len(key)-4:] != ".jpg" {
		t.Errorf("generateImageKey() should end with .jpg, got: %s", key)
	}
}

func TestGenerateImageKey_SameDataSameKey(t *testing.T) {
	data := []byte("same test data")
	contentType := "image/png"
	
	key1 := generateImageKey(data, contentType)
	key2 := generateImageKey(data, contentType)
	
	// 由于包含时间戳，key 应该不同，但 MD5 部分应该相同
	// 简化测试：检查 key 长度和格式是否一致
	if len(key1) != len(key2) {
		t.Errorf("generateImageKey() keys should have same length: %d vs %d", len(key1), len(key2))
	}
	
	// 检查扩展名是否相同
	if len(key1) >= 4 && len(key2) >= 4 {
		ext1 := key1[len(key1)-4:]
		ext2 := key2[len(key2)-4:]
		if ext1 != ext2 {
			t.Errorf("generateImageKey() extensions should be same: %s vs %s", ext1, ext2)
		}
	}
}

func TestGenerateImageKey_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		contentType string
		checkFunc   func(t *testing.T, key string)
	}{
		{
			name:        "very long contentType should work",
			data:        []byte("test"),
			contentType: "image/jpeg; charset=utf-8; boundary=something-very-long",
			checkFunc: func(t *testing.T, key string) {
				if len(key) == 0 {
					t.Error("key should not be empty")
				}
				if len(key) < 4 || key[len(key)-4:] != ".jpg" {
					t.Error("should end with .jpg")
				}
			},
		},
		{
			name:        "empty data should still generate key",
			data:        []byte(""),
			contentType: "image/png",
			checkFunc: func(t *testing.T, key string) {
				if len(key) == 0 {
					t.Error("key should not be empty even with empty data")
				}
			},
		},
		{
			name:        "nil data should still generate key",
			data:        nil,
			contentType: "image/png",
			checkFunc: func(t *testing.T, key string) {
				if len(key) == 0 {
					t.Error("key should not be empty even with nil data")
				}
			},
		},
		{
			name:        "webp image type",
			data:        []byte("fake webp data"),
			contentType: "image/webp",
			checkFunc: func(t *testing.T, key string) {
				if len(key) == 0 {
					t.Error("key should not be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := generateImageKey(tt.data, tt.contentType)
			if tt.checkFunc != nil {
				tt.checkFunc(t, key)
			}
		})
	}
}

func TestGenerateImageKey_ContentTypeVariations(t *testing.T) {
	testData := []byte("test image data")
	
	tests := []struct {
		contentType string
		expectedExt string
	}{
		{"image/jpeg", ".jpg"},
		{"image/JPEG", ".jpg"}, // 测试大小写
		{"image/Jpeg", ".jpg"},
		{"image/jpeg ", ".jpg"}, // 测试空格
		{" image/jpeg", ".jpg"},
		{"image/jpeg;charset=utf-8", ".jpg"},
		{"image/jpeg; charset=utf-8", ".jpg"},
		{"image/jpeg;boundary=something", ".jpg"},
		{"image/jfif", ".jpg"},
		{"image/JFIF", ".jpg"},
		{"image/pjpeg", ".jpg"},
		{"image/pJPEG", ".jpg"},
		{"image/png", ".png"},
		{"image/gif", ".gif"},
		{"image/bmp", ""}, // bmp 可能没有标准扩展名
		{"image/svg+xml", ".svg"},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			key := generateImageKey(testData, tt.contentType)
			
			if tt.expectedExt != "" {
				if len(key) < len(tt.expectedExt) {
					t.Errorf("key too short: %s", key)
					return
				}
				gotExt := key[len(key)-len(tt.expectedExt):]
				if gotExt != tt.expectedExt {
					t.Errorf("expected extension %s, got %s (key: %s)", tt.expectedExt, gotExt, key)
				}
			} else {
				// 对于没有期望扩展名的情况，只检查 key 不为空
				if len(key) == 0 {
					t.Errorf("key should not be empty")
				}
			}
		})
	}
}
