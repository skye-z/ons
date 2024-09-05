package core

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

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
	ticker  *time.Ticker
}

// 第一步 创建 P2P 服务
func NewP2PServer(natId string, host string) *P2PServer {
	path := url.URL{Scheme: "ws", Host: host, Path: "/nat"}
	log.Printf("[P2P] connect %s", path.String())
	connect, _, err := websocket.DefaultDialer.Dial(path.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	server := &P2PServer{
		peerID:  natId,
		connect: connect,
		ticker:  time.NewTicker(5 * time.Minute), // 每5分钟检查一次
	}
	// 第二步 监听请求
	go server.handleMessages()
	// 第三步 注册 NSB
	server.register()
	// 第四步 启动自动重连
	go server.reconnect(host)
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
			log.Println(len(msgBytes))
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
	log.Printf("[P2P] register device")
	message := Message{
		Event: "register",
		Data:  json.RawMessage(`"` + s.peerID + `"`),
	}
	s.sendMessage(message)
}

// 第四步 启动自动重连
func (s *P2PServer) reconnect(host string) {
	for t := range s.ticker.C {
		if s.connect == nil || s.isClosed() {
			path := url.URL{Scheme: "ws", Host: host, Path: "/nat"}
			connect, _, err := websocket.DefaultDialer.Dial(path.String(), nil)
			if err != nil {
				log.Println("尝试重新连接失败:", err)
				continue
			}
			s.connect = connect
			log.Println("成功重新连接")
			go s.handleMessages() // 重启消息处理
		}
		log.Println("检查连接状态:", t)
	}
}

// [工具] 发送消息
func (s *P2PServer) sendMessage(message Message) {
	msgBytes, _ := json.Marshal(message)
	s.connect.WriteMessage(websocket.TextMessage, msgBytes)
}

// [工具] 检查连接是否已关闭
func (s *P2PServer) isClosed() bool {
	_, _, err := s.connect.ReadMessage()
	return err != nil
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

	// 创建数据通道
	dataChannel, err := s.p2p.CreateDataChannel("NSChanel", nil)
	if err != nil {
		log.Fatalf("Failed to create data channel: %v", err)
	}
	log.Println("Data channel created")

	// 设置数据通道的打开回调
	dataChannel.OnOpen(func() {
		log.Println("Data channel is now open, sending message...")
		dataChannel.SendText("Hello from Go!")
	})

	// 设置数据通道的消息接收回调
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Received message: %s\n", msg.Data)
	})
}

// 设置节点信息
func (s *P2PServer) setP2PNode(data webrtc.ICECandidateInit) {
	err := s.p2p.AddICECandidate(data)
	if err != nil {
		log.Fatalf("设置节点信息出错: %v", err)
	}
}

func (s *P2PServer) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	s.connect.Close()
	s.ticker.Stop()
}
