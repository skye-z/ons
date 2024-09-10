package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/skye-z/ons/nas-server/util"
)

const vaultPath = "./vault"

var (
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

// 存储库操作
func VaultOperate(channel *webrtc.DataChannel, data []byte) {
	var syncMsg SyncMessage
	if err := json.Unmarshal(data, &syncMsg); err != nil {
		log.Printf("[Vault] failed to unmarshal message: %v", err)
		return
	}
	log.Printf("[Vault] received operate: %s", syncMsg.Operate)

	// 根据操作类型执行对应的操作
	switch syncMsg.Operate {
	case "check":
		handleCheck(channel, syncMsg)
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

// 读取.synclog文件中的时间戳
func getSyncCheckTime() (int64, error) {
	syncLogPath := filepath.Join(vaultPath, ".synclog")
	data, err := os.ReadFile(syncLogPath)
	if os.IsNotExist(err) {
		// 文件不存在，初始化时间为当前时间
		// saveSyncLog(time.Now())
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	return timestamp, nil
}

// 保存操作日志
func saveSyncLog() {
	logPath := filepath.Join(vaultPath, ".synclog")
	if err := os.WriteFile(logPath, []byte(fmt.Sprint(time.Now().Unix()-1)), 0644); err != nil {
		log.Printf("[Vault] error writing sync log: %v", err)
	}
}

// 处理新旧检查任务
func handleCheck(channel *webrtc.DataChannel, msg SyncMessage) {
	// 获取客户端同步时间
	clientDate, err := strconv.ParseInt(msg.Data, 10, 64)
	if err != nil {
		log.Printf("[Vault] error parsing client date: %v", err)
		return
	}

	// 读取服务端.synclog中的时间
	serverDate, err := getSyncCheckTime()
	if err != nil {
		log.Printf("[Vault] error getting server sync check time: %v", err)
		return
	}

	log.Printf("%v ~ %v", clientDate, serverDate)
	// 比对时间, 如果客户端新则发送服务端.synclog中的时间, 如果服务端新则直接发送服务端文件给客户端
	if clientDate-serverDate <= 3 && clientDate-serverDate >= -3 {
		log.Println("无需同步")
	} else if clientDate < serverDate {
		log.Println("服务端新, 要求客户端发来文件树, 服务端比对后返回变更操作")
		// 如果客户端的时间戳较新，则发送服务端发送文件树
		msgBytes, _ := json.Marshal(SyncMessage{
			Type:    "text",
			Operate: "tree",
			Path:    ".",
			Name:    "",
			Data:    "",
		})
		channel.SendText(string(msgBytes))
	} else {
		log.Println("客户端新, 服务器主动发送文件树, 客户端比对后发回变更操作")
		scan, err := util.ScanDirectory(vaultPath)
		if err != nil {
			log.Printf("[Vault] scan directory error: %v", err)
			return
		}
		scanBytes, _ := json.Marshal(scan)
		// 如果服务端的时间戳较新，则要求客户端发送文件树
		msgBytes, _ := json.Marshal(SyncMessage{
			Type:    "text",
			Operate: "tree",
			Path:    ".",
			Name:    "",
			Data:    string(scanBytes),
		})
		channel.SendText(string(msgBytes))
	}
}

// 处理创建任务
func handleCreate(msg SyncMessage) {
	msg.Path = filepath.Join(vaultPath, msg.Path)

	if msg.Type == "directory" {
		if err := os.MkdirAll(msg.Path, os.ModePerm); err != nil {
			log.Printf("[Vault] error creating directory: %v", err)
		} else {
			saveSyncLog()
		}
	} else {
		handleChunkedDataIfBinary(msg)
	}
}

// 处理删除任务
func handleDelete(path string) {
	log.Printf("rename: %s", path)
	path = filepath.Join(vaultPath, path)
	if err := os.Remove(path); err != nil {
		log.Printf("[Vault] error removing file or directory: %v", err)
	}
	saveSyncLog()
}

// 处理更新任务
func handleUpdate(msg SyncMessage) {
	msg.Path = filepath.Join(vaultPath, msg.Path)
	handleChunkedDataIfBinary(msg)
}

// 处理重命名任务
func handleRename(path, name, oldName string) {
	path = filepath.Join(vaultPath, path)

	// 确保路径存在
	if err := util.EnsureDirExists(path); err != nil {
		log.Printf("[Vault] error ensuring directory exists: %v", err)
		return
	}

	if err := os.Rename(filepath.Join(vaultPath, oldName), filepath.Join(path, name)); err != nil {
		log.Printf("[Vault] error renaming file or directory: %v", err)
	}
	saveSyncLog()
}

// 处理数据分块合并任务
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
		base64ChunkData := parts[2]

		// 获取文件名的完整路径
		filePath := filepath.Join(msg.Path, msg.Name)

		// 确保路径存在
		dirPath := filepath.Dir(filePath)
		if err := util.EnsureDirExists(dirPath); err != nil {
			log.Printf("[Vault] error ensuring directory exists: %v", err)
			return
		}

		// 初始化文件的分块映射表
		if fileChunks[filePath] == nil {
			fileChunks[filePath] = make(map[int]string)
		}

		// 将数据追加到对应文件的分块列表中
		fileChunks[filePath][currentChunk] = base64ChunkData

		// 检查是否所有分块都已经接收完毕
		if len(fileChunks[filePath]) == totalChunks {
			// 合并所有分块
			mergedData := util.MergeChunks(fileChunks[filePath])
			delete(fileChunks, filePath)

			// 将合并后的数据写入文件
			if err := os.WriteFile(filePath, mergedData, 0644); err != nil {
				log.Printf("[Vault] error writing merged data to file: %v", err)
			} else {
				saveSyncLog()
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
		if err := util.EnsureDirExists(dirPath); err != nil {
			log.Printf("[Vault] error ensuring directory exists: %v", err)
			return
		}

		// 将解码后的数据写入文件
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			log.Printf("[Vault] error writing file: %v", err)
		} else {
			saveSyncLog()
		}
	}
}
