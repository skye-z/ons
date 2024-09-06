import { Notice, Plugin, TFile, TFolder } from 'obsidian';
import { NSPluginSettings, NSDefaultSettings } from './src/model';
import { PeerManager } from './src/peer-manager';
import { NSSettingTab } from './src/setting';

// 插件主体
export default class NSPlugin extends Plugin {
	settings: NSPluginSettings;
	status: HTMLElement;
	private peerManager: PeerManager | null = null;
	private isSyncing: boolean = false;
	// 插件加载
	async onload() {
		// 加载设置
		await this.loadSettings();
		// 创建状态栏显示区
		this.status = this.addStatusBarItem();
		this.status.setText(this.settings.model === 'auto' ? '自动模式' : '手动模式');
		// 添加更新命令
		this.addCommand({
			id: 'nas-sync',
			name: '手动更新',
			callback: () => this.syncFilesManually()
		});
		// 创建设置选项卡
		this.addSettingTab(new NSSettingTab(this.app, this));
		// 初始化 PeerManager
		this.initPeerManager();
		// 初始化监听器
		this.initListener();
		// 注册自动同步计时器
		this.registerInterval(window.setInterval(
			() => this.syncFilesAutomatically(), 5 * 60 * 1000
		));
	}
	// 插件卸载
	onunload() {
		if (this.peerManager) {
			this.peerManager.close();
			this.peerManager = null;
		}
	}
	// 加载设置数据
	async loadSettings() {
		this.settings = Object.assign({}, NSDefaultSettings, await this.loadData());
	}
	// 保存设置数据
	async saveSettings() {
		await this.saveData(this.settings);
	}
	// 初始化 PeerManager
	initPeerManager() {
		if (this.peerManager) this.peerManager.close()
		if (this.settings.server)
			this.peerManager = new PeerManager(this.settings.server, this.settings.devId, this.settings.pwd);
	}
	// 初始化监听器
	initListener() {
		const { vault } = this.app;
		vault.on('create', (file) => {
			this.syncWork('create', file.name, file.path)
		})
		vault.on('delete', (file) => {
			this.syncWork('delete', file.name, file.path)
		})
		vault.on('rename', (file) => {
			this.syncWork('rename', file.name, file.path)
		})
		vault.on('modify', (file) => {
			this.syncWork('modify', file.name, file.path)
		})
	}
	syncWork(type: string, name: string, path: string) {
		var stat = this.app.vault.getAbstractFileByPath(path)
		if (stat instanceof TFile) {
			stat.vault.cachedRead(stat).then(res => {
				console.log('file', type, name, path, stat, res)
			})
		} else if (stat instanceof TFolder) {
			console.log('folder', type, name, path, stat)
		}
	}
	// 执行手动同步
	private syncFilesManually() {
		if (this.isSyncing) {
			new Notice('同步正在进行，请稍后再试。');
			return;
		}

		if (this.peerManager) {
			this.isSyncing = true; // 设置同步状态为进行中
			try {
				// 在这里调用同步文件的逻辑
				this.syncFilesAutomatically();
				new Notice('手动同步已触发');
			} finally {
				this.isSyncing = false; // 同步完成后重置状态
			}
		}
	}
	// 执行自动同步
	private syncFilesAutomatically() {
		if (this.peerManager) {
			// 在这里调用同步文件的逻辑
			console.log('自动同步正在进行...');
			// 实际的同步逻辑应该在这里实现
			this.peerManager.syncFiles()
		}
	}
}
