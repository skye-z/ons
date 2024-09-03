package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

// 消息模型
type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
	To    string          `json:"to,omitempty"`
	From  string          `json:"from,omitempty"`
}

type P2PServer struct {
	peerID  string
	connect *websocket.Conn
	p2p     *webrtc.PeerConnection
}

// 第一步 创建 P2P 服务
func NewP2PServer(peerID string, signalingURL string) *P2PServer {
	path := url.URL{Scheme: "ws", Host: signalingURL, Path: "/ws"}
	connect, _, err := websocket.DefaultDialer.Dial(path.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	server := &P2PServer{
		peerID:  peerID,
		connect: connect,
	}
	// 第二步 监听请求
	go server.handleMessages()
	// 第三步 注册 NSB
	server.register()

	return server
}

// 第二步 监听请求
func (s *P2PServer) handleMessages() {
	for {
		_, msgBytes, err := s.connect.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		switch msg.Event {
		case "p2p-exchange":
			signalData := webrtc.SessionDescription{}
			if err := json.Unmarshal(msg.Data, &signalData); err != nil {
				log.Printf("无法解析连接信息: %v", err)
				continue
			}
			log.Printf("收到 %s 发来的连接信息", msg.From)
			if signalData.Type == webrtc.SDPTypeOffer {
				s.setP2PInfo(signalData)
			}
		case "p2p-node":
			nodeData := webrtc.ICECandidateInit{}
			if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
				log.Printf("无法解析节点信息: %v", err)
				continue
			}
			s.setP2PNode(nodeData)
		default:
			log.Printf("Unknown message event: %s", msg.Event)
		}
	}
}

// 第三步 注册 NSB
func (s *P2PServer) register() {
	message := Message{
		Event: "register",
		Data:  json.RawMessage(`"` + s.peerID + `"`),
	}
	s.sendMessage(message)
}

// [工具] 发送消息
func (s *P2PServer) sendMessage(message Message) {
	msgBytes, _ := json.Marshal(message)
	s.connect.WriteMessage(websocket.TextMessage, msgBytes)
}

// 设置对等连接信息
func (s *P2PServer) setP2PInfo(data webrtc.SessionDescription) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
			{
				URLs: []string{"stun:stun.nextcloud.com:443"},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create PeerConnection: %v", err)
	}
	s.p2p = peerConnection
	// 设置 NSC 连接信息
	if err := peerConnection.SetRemoteDescription(data); err != nil {
		log.Printf("Failed to set remote description: %v", err)
	}
	log.Println("NSC 连接已设置")
	// 创建 NSB 本地连接信息
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Printf("Failed to create answer: %v", err)
	}
	// 设置本地连接信息
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		log.Fatalf("Failed to set local description: %v", err)
	}

	// 等待 ICE 采集完成
	<-gatherComplete

	// 监控节点更新
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		candidateJSON := candidate.ToJSON()
		jsonBytes, err := json.Marshal(candidateJSON)
		if err != nil {
			log.Println("JSON 序列化错误:", err)
			return
		}
		log.Printf("向 NSC 发送节点信息")
		answerMsg := Message{
			Event: "p2p-node",
			Data:  json.RawMessage(jsonBytes),
			To:    s.peerID,
			From:  "NSB",
		}
		s.sendMessage(answerMsg)
	})

	// 发送本地连接信息
	answerDataBytes, err := json.Marshal(map[string]interface{}{
		"sdp": map[string]interface{}{
			"type": "answer",
			"sdp":  peerConnection.LocalDescription().SDP,
		},
	})
	if err != nil {
		log.Fatalf("Failed to marshal answer data: %v", err)
	}

	answerMsg := Message{
		Event: "p2p-exchange",
		Data:  json.RawMessage(answerDataBytes),
		To:    s.peerID,
		From:  "NSB",
	}
	s.sendMessage(answerMsg)
}

// 设置节点信息
func (s *P2PServer) setP2PNode(data webrtc.ICECandidateInit) {
	log.Println("应用 NSC 节点信息")
	err := s.p2p.AddICECandidate(data)
	if err != nil {
		log.Fatalf("设置节点信息出错: %v", err)
	}
}

func main() {
	server := NewP2PServer("749601", "192.168.1.160:8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	server.connect.Close()
}
