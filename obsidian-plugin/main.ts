import { Notice, Plugin, TFile, TFolder, Vault } from 'obsidian';
import { NSPluginSettings, NSDefaultSettings } from './src/model';
import { PeerManager } from './src/peer-manager';
import { NSSettingTab } from './src/setting';

// 插件主体
export default class NSPlugin extends Plugin {
	settings: NSPluginSettings;
	status: HTMLElement;
	private peerManager: PeerManager | null = null;
	isSyncing: boolean = false;
	// 插件加载
	async onload() {
		// 加载设置
		await this.loadSettings();
		// 创建状态栏显示区
		this.status = this.addStatusBarItem();
		this.status.setText('连接中...');
		// 添加更新命令
		this.addCommand({
			id: 'nas-manual-update',
			name: '手动更新',
			callback: () => this.syncFilesManually()
		});
		this.addCommand({
			id: 'nas-reconnect',
			name: '重新连接',
			callback: () => this.initPeerManager()
		});
		// 创建设置选项卡
		this.addSettingTab(new NSSettingTab(this.app, this));
		// 初始化 PeerManager
		this.initPeerManager();
		// 初始化监听器
		this.initListener();
		// 注册自动同步计时器
		// this.registerInterval(window.setInterval(
		// 	() => this.syncFilesAuto(), 5 * 60 * 1000
		// ));
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
		if (this.settings.server) {
			this.status.setText('连接中...');
			this.peerManager = new PeerManager(this);
		}
	}
	// 初始化监听器
	initListener() {
		const { vault } = this.app;
		vault.on('create', (file) => {
			if (this.peerManager != null && !this.isSyncing)
				this.peerManager.sendOperate(this, 'create', file, undefined, false)
		})
		vault.on('delete', (file) => {
			if (this.peerManager != null && !this.isSyncing)
				this.peerManager.sendOperate(this, 'delete', file, undefined, false)
		})
		vault.on('modify', (file) => {
			if (this.peerManager != null && !this.isSyncing)
				this.peerManager.sendOperate(this, 'update', file, undefined, false)
		})
		vault.on('rename', (file, old) => {
			if (this.peerManager != null && !this.isSyncing)
				this.peerManager.sendOperate(this, 'rename', file, old, false)
		})
	}
	syncWork(type: string, name: string, path: string) {
		let stat = this.app.vault.getAbstractFileByPath(path)
		if (stat instanceof TFile) {
			stat.vault.cachedRead(stat);
		} 
		// ignore folder
		// else if (stat instanceof TFolder) {
		// 	console.log('folder', type, name, path, stat)
		// }
	}
	// 执行手动同步
	private syncFilesManually() {
		if (this.isSyncing) {
			new Notice('同步正在进行，请稍后再试。');
			return;
		}

		if (this.peerManager) {
			try {
				this.syncFiles();
			} finally {
				this.isSyncing = false; // 同步完成后重置状态
			}
		}
	}
	// 执行自动同步
	private syncFiles() {
		if (this.peerManager) {
			new Notice('正在同步中, 请勿编辑和操作');
			// 在这里调用同步文件的逻辑
			// console.log('文件同步准备中');
			// 实际的同步逻辑应该在这里实现
			this.peerManager.syncFiles(this.settings.lastSync)
		}
	}
}
