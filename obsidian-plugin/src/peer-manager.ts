// 消息模型
interface Message {
  event: string;
  data: any;
  to?: string;
  from?: string;
}

export class PeerManager {
  // 服务器A, 作为信令和中转服务器
  private nsa: WebSocket;
  // NSA 地址(ws://....)
  private nsaPath: string;
  // NAS编号
  private nabId: string;
  // 点对点连接
  private p2pCon: RTCPeerConnection;
  // 数据通道
  private channel: RTCDataChannel;
  // 构造函数
  constructor(url: string, nabId: string) {
    this.nsaPath = url;
    this.nabId = nabId;
    // 创建点对点连接
    this.p2pCon = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: 'stun:stun.nextcloud.com:443' }
      ]
    });
    // 第一步 生成本地描述信息
    this.settingLocalInfo()
    // 第二步 连接 NSA
    this.nsa = this.connectnsa();
  }

  // 第一步 生成本地连接信息
  private async settingLocalInfo() {
    // 创建数据通道(必须在最前面)
    this.channel = this.p2pCon.createDataChannel('NSChanel')
    this.channel.onopen = () => {
      console.log('数据通道开启');
      this.syncFiles();  // 确保通道打开后立即发送同步文件的请求
    };
    this.channel.onclose = () => console.log('数据通道已关闭');
    this.channel.onerror = (error) => console.error('数据通道错误:', error);
    // 监听网络节点变动
    this.p2pCon.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('网络节点信息', event.candidate);
        const candidateMsg: Message = {
          event: 'p2p-node',
          to: this.nabId,
          from: 'NSC',
          data: event.candidate,
        };
        this.sendMessage(candidateMsg);
      }
    };
    
    this.p2pCon.oniceconnectionstatechange = () => {
      console.log('连接状态更新:', this.p2pCon.iceConnectionState);
      if (this.p2pCon.iceConnectionState === 'connected') {
        console.log('对等连接已建立');
      }
    };

    // 创建本地连接信息
    const offer = await this.p2pCon.createOffer();
    // 设置本地连接信息
    await this.p2pCon.setLocalDescription(offer);
  }

  // 第二步 连接 NSA
  private connectnsa(): WebSocket {
    const nsa = new WebSocket(`${this.nsaPath}`);
    nsa.onopen = () => {
      console.log('与 NSA 的通信端口已打开');
      // 第三步 在 NSA 上注册连接
      this.register();
    };
    nsa.onmessage = (event: MessageEvent) => {
      const message: Message = JSON.parse(event.data);
      console.log('收到 NSA 消息:', message);
      // 连接注册响应
      if (message.event === 'connect') {
        // 第四步 发送本地连接信息
        this.sendLocalInfo()
      } else if (message.event === 'p2p-exchange') {
        this.setRemoteInfo(message.data);
      } else if (message.event === 'p2p-node') {
        this.setNodeInfo(message.data);
      }
    };
    nsa.onerror = (error) => { console.error('与 NSA 连接出错:', error) };
    nsa.onclose = (event) => { console.log('与 NSA 的连接已关闭:', event) };
    return nsa;
  }

  // 第三步 在 NSA 上注册连接
  private register() {
    const connectMsg: Message = { event: 'connect', to: this.nabId, from: 'NSC', data: '' };
    this.sendMessage(connectMsg);
  }

  // 第四步 发送本地连接信息
  public async sendLocalInfo() {
    const msg: Message = { event: 'p2p-exchange', to: this.nabId, from: 'NSC', data: this.p2pCon.localDescription };
    this.sendMessage(msg);
  }

  // 第五步 设置远程连接信息
  private async setRemoteInfo(data: any) {
    if (data.sdp) {
      const sdp = new RTCSessionDescription(data.sdp);
      await this.p2pCon.setRemoteDescription(sdp);
    } else if (data.candidate) {
      this.setNodeInfo(data);
    }
  }

  // 第五步 设置节点信息
  private async setNodeInfo(data: RTCIceCandidateInit) {
    const candidate = new RTCIceCandidate(data);
    await this.p2pCon.addIceCandidate(candidate);
  }

  // [工具] 发送信息给服务器A
  private sendMessage(message: Message) {
    if (this.nsa && this.nsa.readyState === WebSocket.OPEN) {
      this.nsa.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not open. Ready state: ', this.nsa.readyState);
    }
  }

  syncFiles(){
    console.log('文件同步正在进行...');
    this.channel.send('文件数据')
  }

}