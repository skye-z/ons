// 插件设置
export interface NSPluginSettings {
	server: string;
	devId: string;
	pwd: string;
	lastSync: number;
	stunMain: string;
	stunBackup: string;
}

// 插件设置默认值
export const NSDefaultSettings: NSPluginSettings = {
	server: 'ons.betax.dev',
	devId: '',
	pwd: '',
	lastSync: 0,
	stunMain: 'stun:stun.l.google.com:19302',
	stunBackup: 'stun:stun.nextcloud.com:443'
}