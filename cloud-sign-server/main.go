package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	peers   = make(map[string]*websocket.Conn) // 用唯一编号注册的 Peer 连接
	clients = make(map[*websocket.Conn]string) // 客户端连接与其目标Peer的编号映射
	mu      sync.Mutex
)

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
	To    string          `json:"to,omitempty"`
	From  string          `json:"from,omitempty"`
}

// 处理 WebSocket 连接
func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("Failed to close WebSocket connection: %v", err)
		}
	}()

	// 读取第一条消息来判断是服务器B还是客户端
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Error reading first message: %v", err)
		return
	}

	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return
	}

	// 判断消息类型是服务器B的 'register' 还是客户端的 'connect'
	if msg.Event == "register" {
		peerID := string(msg.Data) // 服务器B的注册信息
		peerID = strings.Trim(peerID, `"`)
		if peerID != "" {
			registerPeer(ws, peerID)
		} else {
			log.Println("Peer ID is missing for a server connection")
		}
	} else if msg.Event == "connect" {
		registerClient(ws, message)
	} else {
		log.Println("Unknown first message event type")
		ws.Close()
	}
}

// 注册 Peer (服务器B)
func registerPeer(ws *websocket.Conn, peerID string) {
	log.Printf("Registering server (peer): %s", peerID)
	mu.Lock()
	peers[peerID] = ws
	log.Printf("Current peers map after registration: %v", peers) // Debug line
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(peers, peerID)
		log.Printf("Peer %s removed. Current peers map: %v", peerID, peers) // Debug line
		mu.Unlock()
		log.Printf("Server %s disconnected", peerID)
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Read error from server %s: %v", peerID, err)
			break
		}

		// 转发信令消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg.To != "" {
				forwardSignal(msg)
			}
		} else {
			log.Printf("Unmarshal error from server %s: %v", peerID, err) // Debug line
		}
	}
}

// 注册客户端
func registerClient(ws *websocket.Conn, firstMessage []byte) {
	log.Println("Registering client")
	clientID := "" // 在此获取或生成一个客户端ID

	defer func() {
		mu.Lock()
		if clientID != "" {
			delete(clients, ws)
		}
		mu.Unlock()
		log.Printf("Client %s disconnected", clientID)
	}()

	// 处理客户端的初始连接消息
	handleClientMessage(ws, firstMessage)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Read error from client: %v", err)
			break
		}

		handleClientMessage(ws, message)
	}
}

func handleClientMessage(ws *websocket.Conn, message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return
	}

	// 客户端请求连接到 Peer
	if msg.Event == "connect" && msg.To != "" {
		mu.Lock()
		clients[ws] = msg.To
		mu.Unlock()

		// 发送确认消息给客户端
		confirmationMsg := Message{
			Event: "connect-confirmation",
			Data:  json.RawMessage(`"Connected to peer ` + msg.To + `"`),
		}
		msgBytes, _ := json.Marshal(confirmationMsg)
		if err := ws.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Failed to send confirmation to client: %v", err)
		} else {
			log.Printf("Sent connect confirmation to client for peer %s", msg.To)
		}
	}

	// 转发信令消息
	if msg.Event == "webrtc-signal" {
		forwardSignal(msg)
	}
}

// 转发信令消息
func forwardSignal(msg Message) {
	log.Printf("Attempting to forward signal to peer: %s", msg.To) // Debug line

	mu.Lock()
	targetPeer, peerExists := peers[msg.To]
	log.Printf("Current peers map in forwardSignal: %v", peers) // Debug line
	mu.Unlock()

	if peerExists {
		// 将整个 Message 对象转换为 JSON 字符串
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			return
		}

		// 发送整个 JSON 对象
		if err := targetPeer.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Failed to forward signal to peer %s: %v", msg.To, err)
		} else {
			log.Printf("Forwarded signal from %s to peer %s", msg.From, msg.To)
		}
	} else {
		log.Printf("Peer %s not found", msg.To)
	}
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ws", handleConnections)

	log.Println("Server A is listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
