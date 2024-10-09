package download

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// 参数:
//   - source: 可以是磁力链接或磁力文件的路径 (必填)
//   - output: 输出路径 (可选)
func AnimeDownloader(source string, output ...string) error {
	// 检查参数
	if source == "" {
		return fmt.Errorf("source is required")
	}

	//确定outputDir
	downloadDir := viper.GetString("download_dir")
	if len(output) > 0 && output[0] != "" {
		downloadDir = output[0]
	}

	// 检查输出路径是否存在且可写入
	if err := ensureWritableDirectory(downloadDir); err != nil {
		return fmt.Errorf("download path is not writable or cannot be created: %w", err)
	}

	client, err := NewClient(downloadDir)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	if isMagnetLink(source) {
		return client.DownloadFromMagnet(source)
	}

	if _, err := os.Stat(source); err == nil {
		return client.DownloadFromFile(source)
	}

	return fmt.Errorf("invalid source: %s", source)
}

// isMagnetLink 检查给定的字符串是否为磁力链接
func isMagnetLink(s string) bool {
	return len(s) > 8 && s[:8] == "magnet:?"
}

// ensureWritableDirectory checks if the given path is a writable directory
// If it doesn't exist, it tries to create it
func ensureWritableDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	testFile := filepath.Join(absPath, ".write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %w", err)
	}
	file.Close()
	os.Remove(testFile)

	return nil
}
