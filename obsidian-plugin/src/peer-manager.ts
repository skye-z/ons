import { arrayBufferToBase64, base64ToArrayBuffer, Notice, TAbstractFile, TFile, TFolder, Vault } from 'obsidian';
import NSPlugin from 'main';
import t from './i18n/locale';

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
  operate: 'create' | 'delete' | 'update' | 'rename' | 'check' | 'tree' | 'tree-none' | undefined; // 操作类型
  path: string | undefined; // 所在路径
  name: string | undefined; // 对象名称
  data: string | undefined | null; // 实际数据
}

export class PeerManager {
  private isSync: boolean;
  private syncStateTimer: NodeJS.Timeout;
  private reConnectNumber: number;
  private reConnectTimer: NodeJS.Timeout;
  private chunkCache: Map<string, string[]> = new Map();
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
  constructor(app: NSPlugin) {
    this.isSync = false;
    this.reConnectNumber = 0;
    this.nsaPath = 'wss://' + app.settings.server + '/nat';
    this.nabId = app.settings.devId;
    this.pass = app.settings.pwd;
    // 创建点对点连接
    this.p2pCon = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: 'stun:stun.nextcloud.com:443' }
      ],
    });
    // 第一步 生成本地描述信息
    this.settingLocalInfo(app, app.app.vault)
    // 第二步 连接 NSA
    this.nsa = this.connectnsa(app);
  }

  // 第一步 生成本地连接信息
  private async settingLocalInfo(app: NSPlugin, vault: Vault) {
    // 创建数据通道(必须在最前面)
    this.channel = this.p2pCon.createDataChannel('NSChanel')
    // this.channel.onclose = () => console.log('数据通道已关闭');
    // this.channel.onerror = (error) => {
    //   console.error('数据通道错误:', error)
    // };
    // 监听网络节点变动
    this.p2pCon.onicecandidate = (event) => {
      if (event.candidate) {
        // console.log('网络节点信息', event.candidate);
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
          // console.log('等待描述设置完成再发送候选');
        }
      }
    };
    this.p2pCon.oniceconnectionstatechange = () => {
      // console.log('连接状态更新:', this.p2pCon.iceConnectionState);
      if (this.p2pCon.iceConnectionState === 'disconnected') {
        new Notice("⛓️‍💥 NAS " + t("TIP_CONNECTION_DROPPED"));
        app.status.setText('🟡 NAS ' + t("STATUS_CONNECTION_DROPPED"));
        this.reConnect(app)
      }
    };

    // 创建本地连接信息
    const offer = await this.p2pCon.createOffer();
    // 设置本地连接信息
    await this.p2pCon.setLocalDescription(offer);

    this.p2pCon.ondatachannel = (event) => {
      const dataChannel = event.channel;
      dataChannel.onopen = () => {
        clearTimeout(this.reConnectTimer)
        this.reConnectNumber = 0
        app.status.setText('🟢 NAS 已连接');
        new Notice("🚀 NAS 已连接");
      };
      dataChannel.onmessage = (event) => {
        let msg: SyncMessage = JSON.parse(event.data)
        // console.log('收到数据:', msg);

        if (msg.operate === 'tree') this.handleTree(app, vault, msg)
        else if (msg.operate === 'tree-none') {
          new Notice("😆 同步结束, 数据已是最新");
          this.syncOver();
        }
        else if (msg.operate === 'create') this.handleCreate(app, vault, msg)
        else if (msg.operate === 'delete') this.handleDelete(app, vault, msg)
        else if (msg.operate === 'update') this.handleUpdate(app, vault, msg)
      };
    };
  }

  private reConnect(app: NSPlugin){
    if (this.reConnectNumber < 3) {
      clearTimeout(this.reConnectTimer)
      this.reConnectTimer = setTimeout(() => {
        this.reConnectNumber++
        new Notice("第"+this.reConnectNumber+"次尝试重新连接...");
        let iceServers = [];
        if (app.settings.stunMain != '') iceServers.push({urls:app.settings.stunMain})
        if (app.settings.stunBackup != '') iceServers.push({urls:app.settings.stunBackup})
        this.p2pCon = new RTCPeerConnection({
          iceServers
        });
        this.settingLocalInfo(app, app.app.vault)
        this.nsa = this.connectnsa(app);
      }, 3000)
    }
  }

  private async handleCreate(app: NSPlugin, vault: Vault, msg: SyncMessage) {
    this.isSync = true;
    if (msg.path === undefined) return
    if ((msg.path === '.' || msg.path === '/' || msg.path.startsWith('.')) && msg.name === '') return
    if (vault.getFolderByPath(msg.path) == null) await vault.createFolder(msg.path);
    vault.create(msg.path + '/' + msg.name, "")
    this.updateSyncTime(app)
    this.syncOver();
  }

  private handleDelete(app: NSPlugin, vault: Vault, msg: SyncMessage) {
    this.isSync = true;
    let path
    if (msg.name === '') path = msg.path
    else path = msg.path === '.' ? (msg.name) : (msg.path + '/' + msg.name)
    if (path == undefined) return
    let file = vault.getAbstractFileByPath(path)
    if (file == null) return
    app.app.fileManager.trashFile(file)
    this.updateSyncTime(app)
    this.syncOver();
  }

  private handleUpdate(app: NSPlugin, vault: Vault, msg: SyncMessage) {
    this.isSync = true;
    let path = msg.path === '.' ? (msg.name) : (msg.path + '/' + msg.name)
    if (path == undefined) return
    let file = vault.getAbstractFileByPath(path)
    if (file == null) return
    if (msg.type === 'text') {
      if (file instanceof TFile) {
        if (msg.data == null || msg.data == undefined) return;
        const decoder = new TextDecoder("utf-8");
        vault.modify(file, decoder.decode(new Uint8Array(Array.from(atob(msg.data), c => c.charCodeAt(0)))))
      }
    } else if (msg.type === 'binary') {
      if (file instanceof TFile) {
        if (msg.data == null || msg.data == undefined) return;
        this.handleBinaryChunk(vault, file, msg.data);
      }
    }
    this.updateSyncTime(app);
    this.syncOver();
  }

  private syncOver() {
    clearTimeout(this.syncStateTimer)
    this.syncStateTimer = setTimeout(() => {
      this.isSync = false
    }, 2000)
  }

  private handleBinaryChunk(vault: Vault, file: TFile, msg: string) {
    const parts = msg.split(':');
    if (parts.length !== 3) {
      console.error('Invalid chunk data format:', msg);
      return;
    }

    const chunkIdx = parseInt(parts[0], 10);
    const totalChunks = parseInt(parts[1], 10);
    const chunkBase64 = parts[2];

    const key = `${file.path}:${totalChunks}`;
    let chunks = this.chunkCache.get(key) || [];
    chunks[chunkIdx - 1] = chunkBase64;
    this.chunkCache.set(key, chunks);

    if (chunks.length === totalChunks) {
      const code = chunks.join('');
      vault.modifyBinary(file, base64ToArrayBuffer(code));
      this.chunkCache.delete(key);
    }
  }

  private handleTree(app: NSPlugin, vault: Vault, msg: SyncMessage) {
    let list = vault.getAllLoadedFiles()
    this.isSync = true;
    if (msg.data === undefined || msg.data === null || msg.data === "") {
      // 发送本地文件树
      let tree = [];
      for (let i in list) {
        let item = list[i]
        if (item instanceof TFile) {
          tree.push({
            name: item.name,
            path: item.path,
            mtime: Math.floor(item.stat.mtime / 1000),
            size: item.stat.size
          })
        } else if (item instanceof TFolder) {
          tree.push({
            name: null,
            path: item.path,
            mtime: null,
            size: null
          })
        }
      }
      let msg: SyncMessage = {
        path: "",
        name: "",
        type: undefined,
        data: JSON.stringify(tree),
        operate: 'tree'
      };
      this.channel.send(JSON.stringify(msg));
      return false
    }
    let data = JSON.parse(msg.data)
    // 本地有云端没有
    for (let i in list) {
      let cloud;
      let item = list[i]
      if (item.path === '.' || item.path === '/' || item.path.startsWith('.')) continue;
      let exist = false
      for (let x in data) {
        if (item.path === data[x].path) {
          exist = true;
          cloud = data[x]
          break;
        }
      }
      if (!exist) {
        // 新建文件
        this.sendOperate(app, "create", item, undefined, true)
        // 间隔一段时间后发送文件内容
        setTimeout(() => {
          this.sendOperate(app, "update", item, undefined, true)
        }, 2000);
      } else if (item instanceof TFile && item.stat.size !== cloud.size && item.stat.mtime - cloud.mtime > 3) {
        this.sendOperate(app, "update", item, undefined, true)
      }
    }
    // 云端有本地没有
    for (let x in data) {
      if (data[x].path === '.' || data[x].path === '/' || data[x].path.startsWith('.')) continue;
      let exist = false
      for (let i in list) {
        if (list[i].path === data[x].path) {
          exist = true;
          break;
        }
      }
      if (!exist) {
        let msg: SyncMessage = {
          path: data[x].path,
          name: data[x].name,
          type: undefined,
          data: undefined,
          operate: 'delete'
        };
        this.channel.send(JSON.stringify(msg));
      }
    }
  }

  // 第二步 连接 NSA
  private connectnsa(app: NSPlugin): WebSocket {
    const nsa = new WebSocket(`${this.nsaPath}`);
    nsa.onopen = () => {
      // console.log('与 NSA 的通信端口已打开');
      // 第三步 在 NSA 上注册连接
      this.register();
    };
    nsa.onmessage = (event: MessageEvent) => {
      const message: Message = JSON.parse(event.data);
      // console.log('收到 NSA 消息:', message);
      // 连接注册响应
      if (message.event === 'connect') {
        // 第四步 发送本地连接信息
        this.sendLocalInfo()
      } else if (message.event === 'p2p-exchange') {
        this.setRemoteInfo(message.data);
      } else if (message.event === 'p2p-node') {
        this.setNodeInfo(message.data);
      } else if (message.event === 'p2p-error' || message.event === 'error') {
        this.outError(app, message.data);
      }
    };
    // nsa.onerror = (error) => { console.error('与 NSA 连接出错:', error) };
    // nsa.onclose = (event) => { console.log('与 NSA 的连接已关闭:', event) };
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

  private async outError(app: NSPlugin, data: any) {
    let msg = data
    switch (data) {
      case 'password error':
        app.status.setText('连接密码错误');
        msg = '连接密码错误'
        break
      case 10001:
        app.status.setText('连接失败');
        msg = '协议无法对齐'
        break
      case 10002:
        app.status.setText('连接失败');
        msg = '不支持的接入类型'
        break
      case 10003:
        msg = '不支持的消息类型'
        break
      case 10004:
        app.status.setText('设备不存在');
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
        app.status.setText('🔴 NAS 已离线');
        msg = 'NAS 已离线'
        this.reConnect(app)
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

  syncFiles(lastSync: number) {
    if (this.channel.readyState != 'open') {
      new Notice("⚠️ 未连接到 NAS, 请重新连接后再试");
      return false
    }
    // console.log('已请求文件同步');
    let msg: SyncMessage = {
      path: './',
      name: '.synclog',
      type: 'text',
      data: lastSync + '',
      operate: 'check'
    };
    this.channel.send(JSON.stringify(msg));
  }

  async getSyncCheckTime(vault: Vault) {
    let logPath = "/.synclog";
    let checkFile = vault.getFileByPath(logPath)
    if (checkFile == null) {
      vault.create(".synclog", '0')
      return 0
    } else return parseInt(await vault.read(checkFile))
  }

  setSyncCheckTime(vault: Vault, time: number) {
    let logPath = "/.synclog";
    let checkFile = vault.getFileByPath(logPath)
    if (checkFile == null) {
      vault.create(".synclog", time + '')
    } else {
      vault.modify(checkFile, time + '')
    }
  }

  sendOperate(app: NSPlugin, operate: 'create' | 'delete' | 'update' | 'rename', file: TFile | TFolder | TAbstractFile, old: string | undefined, force: boolean) {
    if (this.isSync && !force) return false
    const blockSize = 40 * 1024;
    // 发送文本消息
    let msg: SyncMessage = {
      path: file.parent?.path,
      name: file.name,
      type: undefined,
      data: undefined,
      operate
    };
    if (this.channel.readyState != 'open') return
    // 删除操作, 直接发送
    if (operate === 'delete') {
      msg.path = file.path
      this.channel.send(JSON.stringify(msg));
      this.updateSyncTime(app)
      return
    }
    // 重命名操作, 直接发送
    if (operate === 'rename') {
      msg.data = old;
      this.channel.send(JSON.stringify(msg));
      this.updateSyncTime(app)
      return
    }
    // 目标为文件夹, 直接发送
    if (file instanceof TFolder) {
      msg.type = 'directory'
      this.channel.send(JSON.stringify(msg));
      this.updateSyncTime(app)
      return
    } else if (file instanceof TFile) {
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
          const base64Data = arrayBufferToBase64(data);
          const chunks = this.splitBase64Data(base64Data, blockSize);
          chunks.forEach(chunk => {
            msg.data = index + ':' + chunks.length + ':' + chunk;
            index++;
            this.channel.send(JSON.stringify(msg));
          });
        });
      }
      this.updateSyncTime(app)
    }
  }

  private splitBase64Data(data: string, blockSize: number): string[] {
    const chunks = [];
    for (let i = 0; i < data.length; i += blockSize) {
      const end = Math.min(i + blockSize, data.length);
      chunks.push(data.substring(i, end));
    }
    return chunks;
  }

  private updateSyncTime(app: NSPlugin) {
    app.settings.lastSync = Math.trunc(Date.now() / 1000);
    app.saveData(app.settings)
  }

  close() {
    this.channel.close()
    this.p2pCon.close()
    this.nsa.close()
  }
}