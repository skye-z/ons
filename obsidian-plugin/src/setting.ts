import NSPlugin from 'main';
import { App, Notice, PluginSettingTab, Setting } from 'obsidian';

// 设置面板
export class NSSettingTab extends PluginSettingTab {
	plugin: NSPlugin;

	constructor(app: App, plugin: NSPlugin) {
		super(app, plugin);
		this.plugin = plugin;
	}

	display(): void {
		const { containerEl } = this;
		containerEl.empty();

		new Setting(containerEl)
			.setName('Central control server')
			.setDesc('中控服务器')
			.addText(text => text
				.setPlaceholder('ons.betax.dev')
				.setValue(this.plugin.settings.server)
				.onChange(async (value) => {
					this.plugin.settings.server = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('NAT.ID')
			.setDesc('唯一标识')
			.addText(text => text
				.setPlaceholder('000000')
				.setValue(this.plugin.settings.devId)
				.onChange(async (value) => {
					this.plugin.settings.devId = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('Connection password')
			.setDesc('连接密码')
			.addText(text => text
				.setPlaceholder('8-24位')
				.setValue(this.plugin.settings.pwd)
				.onChange(async (value) => {
					this.plugin.settings.pwd = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('Main stun server')
			.setDesc('信令服务器')
			.addText(text => text
				.setPlaceholder('stun:domain.com:443')
				.setValue(this.plugin.settings.stunMain)
				.onChange(async (value) => {
					this.plugin.settings.stunMain = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('Backup stun server')
			.setDesc('备用信令服务器')
			.addText(text => text
				.setPlaceholder('stun:domain.com:443')
				.setValue(this.plugin.settings.stunBackup)
				.onChange(async (value) => {
					this.plugin.settings.stunBackup = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('Connection test')
			.setDesc('连接测试')
			.addButton(text => text
				.setButtonText("Start")
				.onClick(async () => {
					new Notice("开始测试, 请留意右上角提示与右下角的状态")
					this.plugin.initPeerManager();
				}));
	}
}
