package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skye-z/cloud-server/model"
	"github.com/skye-z/cloud-server/util"
	"xorm.io/xorm"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	peers         = make(map[string]*websocket.Conn)
	clients       = make(map[*websocket.Conn]string)
	clientPeers   = make(map[string]*websocket.Conn)
	clientMapping = make(map[string]string)
	mu            sync.Mutex
)

type P2PService struct {
	Data *model.DeviceModel
}

func CreateP2PService(engine *xorm.Engine) *P2PService {
	data := &model.DeviceModel{
		DB: engine,
	}
	return &P2PService{
		Data: data,
	}
}

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
	To    string          `json:"to,omitempty"`
	From  string          `json:"from,omitempty"`
	Pass  string          `json:"pass,omitempty"`
}

func (ps P2PService) Assess(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ps.sendError(ws, "10001")
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("[P2P] connection closure failed: %v", err)
		}
	}()
	_, message, err := ws.ReadMessage()
	if err != nil {
		ps.sendError(ws, "10002")
		return
	}
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		ps.sendError(ws, "10003")
		return
	}
	if msg.Event == "register" {
		// NSB注册
		info, err := ps.Data.NATGetDevice(strings.Trim(string(msg.Data), `"`))
		if err != nil {
			ps.sendError(ws, "10010")
		} else if info == nil {
			ps.sendError(ws, "10004")
		} else {
			// 更新在线时间
			ps.Data.UpdateOnlineTime(info.Id)
			// 注册设备
			ps.register(ws, info.NatId)
		}
	} else if msg.Event == "connect" {
		// NSC接入
		ps.connet(ws, message)
	} else {
		ps.sendError(ws, "10005")
		ws.Close()
	}
}

// 注册设备
func (ps P2PService) register(ws *websocket.Conn, natId string) {
	log.Printf("[P2P] register NSB: %s", natId)
	mu.Lock()
	peers[natId] = ws
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(peers, natId)
		mu.Unlock()
	}()

	ps.sendMessage(ws, Message{
		Event: "online",
		From:  "NSA",
	})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			ps.sendError(ws, "10006")
			break
		}

		// 转发信令消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg.To != "" {
				ps.relay(ws, msg)
			}
		} else {
			ps.sendError(ws, "10007")
		}
	}
}

// 连接设备
func (ps P2PService) connet(ws *websocket.Conn, firstMessage []byte) {
	clientID := util.GenerateRandomNumber(6)
	log.Printf("[P2P] NSC #%s is connected", clientID)

	defer func() {
		mu.Lock()
		if clientID != "" {
			delete(clients, ws)
			delete(clientPeers, clientMapping[clientID])
			delete(clientMapping, clientID)
		}
		mu.Unlock()
		log.Printf("[P2P] NSC #%s disconnected", clientID)
	}()

	// 处理客户端的初始连接消息
	ps.handleClientMessage(ws, clientID, firstMessage)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			ps.sendError(ws, "10006")
			break
		}

		ps.handleClientMessage(ws, clientID, message)
	}
}

// 处理客户端消息
func (ps P2PService) handleClientMessage(ws *websocket.Conn, clientID string, message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		ps.sendError(ws, "10007")
		return
	}

	// 客户端请求连接到 Peer
	if msg.Event == "connect" && msg.To != "" {
		mu.Lock()
		clients[ws] = msg.To
		clientPeers[msg.To] = ws
		clientMapping[clientID] = msg.To
		mu.Unlock()

		// 发送确认消息给客户端
		msg := Message{
			Event: "connect",
			Data:  json.RawMessage(`"准许连接 #` + msg.To + ` NSB"`),
			To:    msg.To,
			From:  "NSA",
		}
		if err := ps.sendMessage(ws, msg); err != nil {
			log.Printf("[P2P] mssage sending failed %s: %v", msg.To, err)
		} else {
			// 更新连接时间
			ps.Data.UpdateConnectTime(msg.To)
			log.Printf("[P2P] NSC applies to connect #%s NSB", msg.To)
		}
	} else if msg.Event == "p2p-error" || msg.Event == "p2p-exchange" || msg.Event == "p2p-node" {
		ps.relay(ws, msg)
	}
}

// 转发消息
func (ps P2PService) relay(now *websocket.Conn, msg Message) {
	var ws *websocket.Conn
	var peerExists bool
	if msg.From == "NSC" {
		mu.Lock()
		ws, peerExists = peers[msg.To]
		mu.Unlock()
	} else if msg.From == "NSB" {
		mu.Lock()
		ws, peerExists = clientPeers[msg.To]
		mu.Unlock()
	}

	if peerExists {
		if err := ps.sendMessage(ws, msg); err != nil {
			log.Printf("[P2P] mssage sending failed %s: %v", msg.To, err)
		}
	} else {
		ps.sendError(now, "10008")
	}
}

// 发送消息
func (ps P2PService) sendError(ws *websocket.Conn, msg string) error {
	log.Printf("[P2P] error: %s", msg)
	if ws == nil {
		return nil
	}
	msgBytes, _ := json.Marshal(Message{
		Event: "error",
		Data:  json.RawMessage(msg),
		From:  "NSA",
	})
	return ws.WriteMessage(websocket.TextMessage, msgBytes)
}

// 发送消息
func (ps P2PService) sendMessage(ws *websocket.Conn, msg Message) error {
	msgBytes, _ := json.Marshal(msg)
	return ws.WriteMessage(websocket.TextMessage, msgBytes)
}

// 检查连接状态
func (ps P2PService) CheckOnline(ctx *gin.Context) {
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	list, err := ps.Data.GetDeviceList(uid, 1, 100)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取状态失败")
		return
	}

	date := make(map[string]map[string]bool)
	online := make(map[string]bool)
	connect := make(map[string]bool)
	for _, device := range list {
		id := device.NatId
		_, oe := peers[id]
		online[id] = oe
		_, ce := clientPeers[id]
		connect[id] = ce
	}
	date["online"] = online
	date["connect"] = connect
	util.ReturnData(ctx, true, date)
}
