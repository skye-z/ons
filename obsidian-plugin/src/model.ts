// 插件设置
export interface NSPluginSettings {
	model: string;
	nat: boolean;
	signal: string;
	devId: string;
	nas: string;
}

// 插件设置默认值
export const NSDefaultSettings: NSPluginSettings = {
	model: 'auto',
	nat: false,
	signal: 'ws://signal.betax.dev',
	devId: '',
	nas: ''
}