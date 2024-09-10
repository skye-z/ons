package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const vaultPath = "./vault"
const syncLogFileName = ".synclog"
const idleTimeout = 10 * time.Second

var (
	lastOperationTimestamp time.Time
	mu                     sync.Mutex
	ticker                 *time.Ticker
	stop                   chan struct{}

	fileChunks map[string]map[int]string // 存储文件名和对应的分块数据
	fileMutex  sync.Mutex                // 保护文件分块数据的互斥锁
)

type SyncMessage struct {
	Type    string `json:"type"`
	Operate string `json:"operate"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Data    string `json:"data"`
}

func init() {
	lastOperationTimestamp = time.Now()
	stop = make(chan struct{})
	fileChunks = make(map[string]map[int]string)
	ticker = time.NewTicker(idleTimeout / 2) // 使用一半的间隔来检查是否需要保存
	go checkAndSaveSyncLog()
}

func VaultOperate(data []byte) {
	var syncMsg SyncMessage
	if err := json.Unmarshal(data, &syncMsg); err != nil {
		log.Printf("[Vault] failed to unmarshal message: %v", err)
		return
	}
	log.Printf("[Vault] received operate: %s", syncMsg.Operate)

	updateLastOperationTime()

	// 根据操作类型执行对应的操作
	switch syncMsg.Operate {
	case "create":
		handleCreate(syncMsg)
	case "delete":
		handleDelete(syncMsg.Path)
	case "update":
		handleUpdate(syncMsg)
	case "rename":
		handleRename(syncMsg.Path, syncMsg.Name, syncMsg.Data)
	default:
		log.Println("[Vault] unknown operation:", syncMsg.Operate)
	}
}

func updateLastOperationTime() {
	mu.Lock()
	defer mu.Unlock()
	lastOperationTimestamp = time.Now()
}

func checkAndSaveSyncLog() {
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			now := time.Now()
			if now.Sub(lastOperationTimestamp) >= idleTimeout {
				saveSyncLog(lastOperationTimestamp)
				mu.Unlock()
				return
			}
			mu.Unlock()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func saveSyncLog(timestamp time.Time) {
	logPath := filepath.Join(vaultPath, syncLogFileName)
	data := timestamp.Format(time.RFC3339)
	if err := os.WriteFile(logPath, []byte(data), 0644); err != nil {
		log.Printf("[Vault] error writing sync log: %v", err)
	}
}

func handleCreate(msg SyncMessage) {
	msg.Path = filepath.Join(vaultPath, msg.Path)

	if msg.Type == "directory" {
		if err := os.MkdirAll(msg.Path, os.ModePerm); err != nil {
			log.Printf("[Vault] error creating directory: %v", err)
		}
	} else {
		handleChunkedDataIfBinary(msg)
	}
}

func handleDelete(path string) {
	log.Printf("rename: %s", path)
	path = filepath.Join(vaultPath, path)
	if err := os.Remove(path); err != nil {
		log.Printf("[Vault] error removing file or directory: %v", err)
	}
}

func handleUpdate(msg SyncMessage) {
	msg.Path = filepath.Join(vaultPath, msg.Path)
	handleChunkedDataIfBinary(msg)
}

func handleRename(path, name, oldName string) {
	path = filepath.Join(vaultPath, path)

	// 确保路径存在
	if err := ensureDirExists(path); err != nil {
		log.Printf("[Vault] error ensuring directory exists: %v", err)
		return
	}

	if err := os.Rename(filepath.Join(vaultPath, oldName), filepath.Join(path, name)); err != nil {
		log.Printf("[Vault] error renaming file or directory: %v", err)
	}
}

func handleChunkedDataIfBinary(msg SyncMessage) {
	if msg.Type == "binary" {
		fileMutex.Lock()
		defer fileMutex.Unlock()

		// 解析分块数据
		parts := strings.Split(msg.Data, ":")
		if len(parts) < 3 {
			log.Printf("[Vault] invalid chunk data format: %s", msg.Data)
			return
		}

		currentChunk, _ := strconv.Atoi(parts[0])
		totalChunks, _ := strconv.Atoi(parts[1])
		chunkData := parts[2]

		// 获取文件名的完整路径
		filePath := filepath.Join(msg.Path, msg.Name)

		// 确保路径存在
		dirPath := filepath.Dir(filePath)
		if err := ensureDirExists(dirPath); err != nil {
			log.Printf("[Vault] error ensuring directory exists: %v", err)
			return
		}

		// 初始化文件的分块映射表
		if fileChunks[filePath] == nil {
			fileChunks[filePath] = make(map[int]string)
		}

		// 将数据追加到对应文件的分块列表中
		fileChunks[filePath][currentChunk] = chunkData

		// 检查是否所有分块都已经接收完毕
		if len(fileChunks[filePath]) == totalChunks {
			// 合并所有分块
			mergedData := mergeChunks(fileChunks[filePath])
			delete(fileChunks, filePath)

			// 将合并后的数据写入文件
			if err := os.WriteFile(filePath, []byte(mergedData), 0644); err != nil {
				log.Printf("[Vault] error writing merged data to file: %v", err)
			}
		}
	} else {
		// 处理非二进制数据
		data, err := base64.StdEncoding.DecodeString(msg.Data)
		if err != nil {
			log.Printf("[Vault] error decoding base64 data: %v", err)
			return
		}
		// 获取文件名的完整路径
		filePath := filepath.Join(msg.Path, msg.Name)

		// 确保路径存在
		dirPath := filepath.Dir(filePath)
		if err := ensureDirExists(dirPath); err != nil {
			log.Printf("[Vault] error ensuring directory exists: %v", err)
			return
		}
		if err := os.WriteFile(filepath.Join(msg.Path, msg.Name), data, 0644); err != nil {
			log.Printf("[Vault] error writing file: %v", err)
		}
	}
}

func mergeChunks(chunks map[int]string) []byte {
	var mergedData []byte
	var keys []int

	// 将keys排序
	for k := range chunks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// 按照序号合并数据
	for _, k := range keys {
		decodedChunk, err := base64.StdEncoding.DecodeString(chunks[k])
		if err != nil {
			log.Printf("[Vault] error decoding base64 chunk data: %v", err)
			continue
		}
		mergedData = append(mergedData, decodedChunk...)
	}
	return mergedData
}

// 确保目录存在
func ensureDirExists(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	return nil
}
