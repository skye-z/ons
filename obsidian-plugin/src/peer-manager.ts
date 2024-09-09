import { arrayBufferToBase64, Notice, TFile, TFolder } from 'obsidian';

// 消息模型
interface Message {
  event: string;
  data: any;
  to?: string;
  from?: string;
  pass?: string;
}

interface SyncMessage {
  type: 'text' | 'binary' | 'directory' | undefined; // 消息类型
  operate: 'create' | 'delete' | 'update' | 'rename' | undefined; // 操作类型
  path: string | undefined; // 所在路径
  name: string | undefined; // 对象名称
  data: string | undefined; // 实际数据
}

export class PeerManager {
  // 服务器A, 作为信令和中转服务器
  private nsa: WebSocket;
  // NSA 地址(wss://....)
  private nsaPath: string;
  // NAS编号
  private nabId: string;
  // NAS连接密码
  private pass: string;
  // 点对点连接
  private p2pCon: RTCPeerConnection;
  // 数据通道
  private channel: RTCDataChannel;
  // 添加一个候选队列
  private iceCandidateQueue: RTCIceCandidateInit[] = [];
  // 构造函数
  constructor(url: string, nabId: string, pass: string) {
    this.nsaPath = 'wss://' + url + '/nat';
    this.nabId = nabId;
    this.pass = pass;
    // 创建点对点连接
    this.p2pCon = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: 'stun:stun.nextcloud.com:443' }
      ],
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
    this.channel.onclose = () => console.log('数据通道已关闭');
    this.channel.onerror = (error) => {
      console.error('数据通道错误:', error)
    };
    // 监听网络节点变动
    this.p2pCon.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('网络节点信息', event.candidate);
        if (this.p2pCon.localDescription && this.p2pCon.remoteDescription) {
          const candidateMsg: Message = {
            event: 'p2p-node',
            to: this.nabId,
            from: 'NSC',
            data: event.candidate,
            pass: this.pass
          };
          this.sendMessage(candidateMsg);
        } else {
          console.log('等待描述设置完成再发送候选');
        }
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

    this.p2pCon.ondatachannel = (event) => {
      const dataChannel = event.channel;
      dataChannel.onopen = () => {
        new Notice("NAS 已连接");
        console.log('数据通道开启');
        this.syncFiles();
      };
      dataChannel.onmessage = (event) => {
        console.log('收到数据:', event.data);
      };
    };
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
      } else if (message.event === 'p2p-error' || message.event === 'error') {
        this.outError(message.data);
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
    const msg: Message = {
      event: 'p2p-exchange',
      to: this.nabId,
      from: 'NSC',
      data: this.p2pCon.localDescription,
      pass: this.pass
    };
    this.sendMessage(msg);
  }

  // 第五步 设置远程连接信息
  private async setRemoteInfo(data: any) {
    if (data.sdp) {
      const sdp = new RTCSessionDescription(data.sdp);
      await this.p2pCon.setRemoteDescription(sdp);

      // 设置远程描述成功后，处理 ICE 候选队列
      this.iceCandidateQueue.forEach(candidate => this.setNodeInfo(candidate));
      this.iceCandidateQueue = []; // 清空队列
    } else if (data.candidate) {
      // 如果还没设置远程描述，将候选缓存到队列
      if (!this.p2pCon.remoteDescription) {
        this.iceCandidateQueue.push(data);
      } else {
        await this.setNodeInfo(data);
      }
    }
  }

  // 第五步 设置节点信息
  private async setNodeInfo(data: RTCIceCandidateInit) {
    try {
      const candidate = new RTCIceCandidate(data);
      await this.p2pCon.addIceCandidate(candidate);
    } catch (error) {
      console.error('添加 ICE 候选失败:', error);
    }
  }

  private async outError(data: any) {
    console.log(data)
    let msg = data
    switch (data) {
      case 'password error':
        msg = '连接密码错误'
        break
      case 10001:
        msg = '协议无法对齐'
        break
      case 10002:
        msg = '不支持的接入类型'
        break
      case 10003:
        msg = '不支持的消息类型'
        break
      case 10004:
        msg = '设备不存在'
        break
      case 10005:
        msg = '不支持的指令'
        break
      case 10006:
        msg = '未能读取到消息内容'
        break
      case 10007:
        msg = '不支持的消息格式'
        break
      case 10008:
        msg = 'NAS 已离线'
        break
      default:
        msg = '不支持的消息格式'
        break
    }
    new Notice(msg);
  }

  // [工具] 发送信息给服务器A
  private sendMessage(message: Message) {
    if (this.nsa && this.nsa.readyState === WebSocket.OPEN) {
      this.nsa.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not open. Ready state: ', this.nsa.readyState);
    }
  }

  syncFiles() {
    console.log('文件同步正在进行...');
    // 1. 检查是否需要更新
    // 2. nas新则拉取
    // 3. 本地新则推送
    this.channel.send('Hello from JavaScript!')
  }

  sendOperate(operate: 'create' | 'delete' | 'update' | 'rename', file: TFile | TFolder, old: string | undefined) {
    const blockSize = 40 * 1024;
    // 发送文本消息
    let msg: SyncMessage = {
      path: file.parent?.path,
      name: file.name,
      type: undefined,
      data: undefined,
      operate
    };
    // 删除操作, 直接发送
    if (operate === 'delete') {
      this.channel.send(JSON.stringify(msg));
      return
    }
    // 重命名操作, 直接发送
    if (operate === 'rename') {
      msg.data = old;
      this.channel.send(JSON.stringify(msg));
      return
    }
    // 目标为文件夹, 直接发送
    if (file instanceof TFolder) {
      msg.type = 'directory'
      this.channel.send(JSON.stringify(msg));
      return
    }
    // 判断文件类型
    if (file.extension === 'md') {
      msg.type = 'text'
      file.vault.cachedRead(file).then(data => {
        const encoder = new TextEncoder();
        const encodedText = encoder.encode(data);
        msg.data = btoa(String.fromCharCode(...new Uint8Array(encodedText)))
        this.channel.send(JSON.stringify(msg));
      })
    } else {
      msg.type = 'binary'
      file.vault.readBinary(file).then(data => {
        let index = 1;
        const chunks = this.splitData(data, blockSize);
        console.log('文件分块: ' + chunks.length + '块')
        // 发送每个分块
        chunks.forEach(chunk => {
          msg.data = index + ':' + chunks.length + ':' + arrayBufferToBase64(chunk)
          index++
          this.channel.send(JSON.stringify(msg));
        });
      })
    }
  }

  private splitData(data: ArrayBuffer, blockSize: number): ArrayBuffer[] {
    const chunks: ArrayBuffer[] = [];
    const totalSize = data.byteLength;
    const numChunks = Math.ceil(totalSize / blockSize);

    for (let i = 0; i < totalSize; i += blockSize) {
      const end = Math.min(i + blockSize, totalSize);
      const chunk = data.slice(i, end);
      const chunkInfo = new TextEncoder().encode(`${i}:${totalSize}:${numChunks}:`);
      const combinedChunk = new Uint8Array(chunkInfo.byteLength + chunk.byteLength);
      combinedChunk.set(new Uint8Array(chunkInfo), 0);
      combinedChunk.set(new Uint8Array(chunk), chunkInfo.byteLength);
      chunks.push(combinedChunk.buffer);
    }

    return chunks;
  }

  close() {
    this.channel.close()
    this.p2pCon.close()
    this.nsa.close()
  }
}