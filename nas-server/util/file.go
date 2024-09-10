package util

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Ctime string `json:"ctime"`
	Mtime string `json:"mtime"`
	Size  int64  `json:"size"`
}

// 扫描目录下的内容
func ScanDirectory(vaultPath string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 相对于根目录的路径
		relativePath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return err
		}

		// 如果是文件夹，则Name为空
		name := ""
		if !info.IsDir() {
			name = info.Name()
		}
		if name == ".synclog" {
			return nil
		}

		// 创建时间和最后修改时间
		ctime := info.Sys().(*syscall.Stat_t).Ctimespec.Nsec
		mtime := info.ModTime().Unix()

		// 添加文件信息到列表
		files = append(files, FileInfo{
			Name:  name,
			Path:  relativePath,
			Ctime: fmt.Sprintf("%d", ctime),
			Mtime: fmt.Sprintf("%d", mtime),
			Size:  info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// 检查目录是否存在
func EnsureDirExists(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	return nil
}

// 合并分块数据
func MergeChunks(chunks map[int]string) []byte {
	var mergedData string
	var keys []int

	// 将keys排序
	for k := range chunks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// 按照序号合并数据
	for _, k := range keys {
		mergedData += chunks[k]
	}

	// 解码Base64编码的数据
	decodedChunkData, err := base64.StdEncoding.DecodeString(mergedData)
	if err != nil {
		log.Printf("[Vault] error decoding base64 chunk data: %v", err)
		return nil
	}
	return decodedChunkData
}
