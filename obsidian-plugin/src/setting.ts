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
			.setName('中控服务器')
			.setDesc('中控负责登记 NAS 设备, 提供鉴权并促成连接的建立')
			.addText(text => text
				.setPlaceholder('ons.betax.dev')
				.setValue(this.plugin.settings.server)
				.onChange(async (value) => {
					this.plugin.settings.server = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('唯一标识')
			.setDesc('让 Obsidian 能在中控服务器中找到你的 NAS 的唯一标识')
			.addText(text => text
				.setPlaceholder('000000')
				.setValue(this.plugin.settings.devId)
				.onChange(async (value) => {
					this.plugin.settings.devId = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('连接密码')
			.setDesc('NAS 中的设置的密码, NAS 会判断密码来决定是否建立连接')
			.addText(text => text
				.setPlaceholder('8-24位')
				.setValue(this.plugin.settings.pwd)
				.onChange(async (value) => {
					this.plugin.settings.pwd = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('连接测试')
			.setDesc('配置完成后可使用此功能测试配置是否正确')
			.addButton(text => text
				.setButtonText("开始测试")
				.onClick(async () => {
					new Notice("开始测试, 请留意右上角提示与右下角的状态")
					this.plugin.initPeerManager();
				}));
	}
}
