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

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"` // 修改类型为 json.RawMessage
	To    string          `json:"to,omitempty"`
	From  string          `json:"from,omitempty"`
}

type WebRTCServer struct {
	peerID    string
	signaling *websocket.Conn
}

func NewWebRTCServer(peerID string, signalingURL string) *WebRTCServer {
	// Connect to signaling server
	u := url.URL{Scheme: "ws", Host: signalingURL, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	server := &WebRTCServer{
		peerID:    peerID,
		signaling: c,
	}

	go server.handleSignalingMessages()
	server.register()

	return server
}

// Register peer to signaling server
func (s *WebRTCServer) register() {
	message := Message{
		Event: "register",
		Data:  json.RawMessage(`"` + s.peerID + `"`), // 确保发送的数据是字符串形式
	}
	s.sendMessage(message)
}

// 处理来自信令服务器的消息
func (s *WebRTCServer) handleSignalingMessages() {
	for {
		_, msgBytes, err := s.signaling.ReadMessage()
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
		case "webrtc-signal":
			var signalData map[string]interface{} // 使用 map[string]interface{} 解析信令数据
			if err := json.Unmarshal(msg.Data, &signalData); err != nil {
				log.Printf("Failed to unmarshal signal data: %v", err)
				continue
			}

			if sdpMap, ok := signalData["sdp"].(map[string]interface{}); ok {
				if sdp, ok := sdpMap["sdp"].(string); ok {
					if sdpType, ok := sdpMap["type"].(string); ok && sdpType == "offer" {
						log.Printf("Received SDP offer: %v", sdp)
						s.handleSDPOffer(sdp, msg.From) // 添加 msg.From 参数以用于回答
					}
				}
			}

		default:
			log.Printf("Unknown message event: %s", msg.Event)
		}
	}
}

// Send signaling message to signaling server
func (s *WebRTCServer) sendMessage(message Message) {
	msgBytes, _ := json.Marshal(message)
	s.signaling.WriteMessage(websocket.TextMessage, msgBytes)
}

func (s *WebRTCServer) handleSDPOffer(offerSDP string, sender string) {
	// 使用 pion/webrtc 创建一个新的 PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.nextcloud.com:443"},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create PeerConnection: %v", err)
	}

	// 设置远端描述 (offer)
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	}
	log.Printf("Received SDP offer: %v", offerSDP)
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Fatalf("Failed to set remote description: %v", err)
	}

	// 创建 Answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Fatalf("Failed to create answer: %v", err)
	}

	// 设置本地描述 (answer)
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		log.Fatalf("Failed to set local description: %v", err)
	}

	// 等待 ICE Gathering 完成
	<-gatherComplete

	// 发送 answer 回信令服务器
	answerData := map[string]interface{}{
		"sdp": map[string]interface{}{
			"type": "answer",
			"sdp":  peerConnection.LocalDescription().SDP,
		},
	}

	// 将 answerData 编码为 JSON
	answerDataBytes, err := json.Marshal(answerData)
	if err != nil {
		log.Fatalf("Failed to marshal answer data: %v", err)
	}

	answerMsg := Message{
		Event: "webrtc-signal",
		Data:  json.RawMessage(answerDataBytes),
		To:    sender, // 回复给原始 offer 发送者
		From:  s.peerID,
	}
	s.sendMessage(answerMsg)
}

func main() {
	server := NewWebRTCServer("749601", "192.168.1.160:8080")

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	server.signaling.Close()
}
