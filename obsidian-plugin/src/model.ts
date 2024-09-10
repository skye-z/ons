// 插件设置
export interface NSPluginSettings {
	model: string;
	server: string;
	devId: string;
	pwd: string;
	lastSync: number;
}

// 插件设置默认值
export const NSDefaultSettings: NSPluginSettings = {
	model: 'auto',
	server: 'ons.betax.dev',
	devId: '',
	pwd: '',
	lastSync: 0
}