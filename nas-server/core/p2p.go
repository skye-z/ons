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
	"github.com/skye-z/ons/nas-server/util"
)

// 消息模型
type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
	To    string          `json:"to,omitempty"`
	From  string          `json:"from,omitempty"`
	Pass  string          `json:"pass"`
}

type P2PServer struct {
	natId             string
	connect           *websocket.Conn
	p2p               *webrtc.PeerConnection
	ticker            *time.Ticker
	iceCandidateQueue []webrtc.ICECandidateInit
}

// 第一步 创建 P2P 服务
func NewP2PServer(natId string, host string) *P2PServer {
	path := url.URL{Scheme: "wss", Host: host, Path: "/nat"}
	log.Printf("[P2P] connect %s", path.String())
	connect, _, err := websocket.DefaultDialer.Dial(path.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	server := &P2PServer{
		natId:   natId,
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
			log.Println("[P2P] read:", err)
			return
		}
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("[P2P] unmarshal:", err)
			continue
		}

		switch msg.Event {
		case "p2p-exchange":
			if msg.Pass != util.GetString("connect.password") {
				log.Println("[P2P] NSC connection password error")
				s.sendMessage(Message{
					Event: "p2p-error",
					Data:  json.RawMessage(`"password error"`),
					To:    s.natId,
					From:  "NSB",
				})
				continue
			}
			signalData := webrtc.SessionDescription{}
			if err := json.Unmarshal(msg.Data, &signalData); err != nil {
				log.Printf("[P2P] unable to parse connection information: %v", err)
				continue
			}
			if signalData.Type == webrtc.SDPTypeOffer {
				s.setP2PInfo(signalData)
			}
		case "p2p-node":
			if msg.Pass != util.GetString("connect.password") {
				log.Println("[P2P] NSC connection password error")
				s.sendMessage(Message{
					Event: "p2p-error",
					Data:  json.RawMessage(`"password error"`),
					To:    s.natId,
					From:  "NSB",
				})
				continue
			}
			nodeData := webrtc.ICECandidateInit{}
			if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
				log.Printf("[P2P] unable to parse node information: %v", err)
				continue
			}
			s.setP2PNode(nodeData)
		case "online":
			log.Println("[P2P] connection successful")
		case "error":
			log.Printf("[P2P] connect failed: %v", string(msg.Data))
		default:
			log.Printf("[P2P] Unknown message event: %s", msg.Event)
		}
	}
}

// 第三步 注册 NSB
func (s *P2PServer) register() {
	log.Printf("[P2P] register device #%s", s.natId)
	message := Message{
		Event: "register",
		Data:  json.RawMessage(`"` + s.natId + `"`),
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
				log.Println("[P2P] attempt to reconnect failed:", err)
				continue
			}
			s.connect = connect
			log.Println("[P2P] successfully reconnected")
			go s.handleMessages() // 重启消息处理
		}
		log.Println("[P2P] check connection status:", t)
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
		log.Fatalf("[P2P] Failed to create PeerConnection: %v", err)
	}
	s.p2p = peerConnection
	// 设置 NSC 连接信息
	if err := s.p2p.SetRemoteDescription(data); err != nil {
		log.Printf("[P2P] Failed to set remote description: %v", err)
	}
	log.Println("[P2P] NSC connection has been set up")
	// 处理 ICE 候选队列
	for _, candidate := range s.iceCandidateQueue {
		err := s.p2p.AddICECandidate(candidate)
		if err != nil {
			log.Printf("[P2P] node addition failed: %v", err)
		}
	}
	// 清空候选队列
	s.iceCandidateQueue = nil
	// 创建 NSB 本地连接信息
	answer, err := s.p2p.CreateAnswer(nil)
	if err != nil {
		log.Printf("[P2P] Failed to create answer: %v", err)
	}
	// 设置本地连接信息
	if err := s.p2p.SetLocalDescription(answer); err != nil {
		log.Fatalf("[P2P] Failed to set local description: %v", err)
	}

	// 发送本地连接信息
	answerDataBytes, err := json.Marshal(map[string]interface{}{
		"sdp": map[string]interface{}{
			"type": "answer",
			"sdp":  s.p2p.LocalDescription().SDP,
		},
	})
	if err != nil {
		log.Fatalf("[P2P] Failed to marshal answer data: %v", err)
	}

	answerMsg := Message{
		Event: "p2p-exchange",
		Data:  json.RawMessage(answerDataBytes),
		To:    s.natId,
		From:  "NSB",
	}
	s.sendMessage(answerMsg)

	// 监控节点更新
	s.p2p.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		candidateJSON := candidate.ToJSON()
		jsonBytes, err := json.Marshal(candidateJSON)
		if err != nil {
			log.Println("JSON 序列化错误:", err)
			return
		}
		mgs := Message{
			Event: "p2p-node",
			Data:  json.RawMessage(jsonBytes),
			To:    s.natId,
			From:  "NSB",
		}
		s.sendMessage(mgs)
	})

	// 创建数据通道
	dataChannel, err := s.p2p.CreateDataChannel("NSChanel", nil)
	if err != nil {
		log.Fatalf("[P2P] Failed to create data channel: %v", err)
	}
	log.Println("[P2P] data channel created")

	s.p2p.OnDataChannel(func(channel *webrtc.DataChannel) {
		channel.OnOpen(func() {
			log.Println("[P2P] data channel open")
			err := dataChannel.SendText("Hello from Go!")
			if err != nil {
				log.Println("[P2P] Error sending initial message:", err)
			}
		})

		channel.OnClose(func() {
			log.Println("[P2P] data channel close")
		})

		channel.OnError(func(err error) {
			log.Printf("[P2P] data channel error: %s", err.Error())
		})
		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("[Channel] received message: %s", string(msg.Data))
		})
	})
}

// 设置节点信息
func (s *P2PServer) setP2PNode(data webrtc.ICECandidateInit) {
	if s.p2p == nil {
		// 如果 PeerConnection 还未准备好，先缓存候选
		s.iceCandidateQueue = append(s.iceCandidateQueue, data)
		return
	}
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
