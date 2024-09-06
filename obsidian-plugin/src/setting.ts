import NSPlugin from 'main';
import { App, PluginSettingTab, Setting } from 'obsidian';

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
			.setName('同步模式')
			.setDesc('在互联网中访问可以访问你的 NAS 的地址')
			.addDropdown(dropdown => dropdown
				.addOption('auto', '自动同步')
				.addOption('manual', '手动同步')
				.setValue(this.plugin.settings.model)
				.onChange(async (value) => {
					this.plugin.settings.model = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('穿透服务器')
			.setDesc('无法从互联网访问你的 NAS 时, 将由此服务器提供穿透服务')
			.addText(text => text
				.setPlaceholder('ons.betax.dev')
				.setValue(this.plugin.settings.server)
				.onChange(async (value) => {
					this.plugin.settings.server = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('穿透标识')
			.setDesc('NAS 中的唯一标识, 让你能在穿透服务器中找到你的 NAS')
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
					this.plugin.initializePeerManager();
				}));
	}
}
