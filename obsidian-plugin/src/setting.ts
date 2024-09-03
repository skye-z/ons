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
			.setName('NAS 地址')
			.setDesc('在互联网中访问可以访问你的 NAS 的地址')
			.addText(text => text
				.setPlaceholder('https://your-nas/sync')
				.setValue(this.plugin.settings.nas)
				.onChange(async (value) => {
					this.plugin.settings.nas = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('允许穿透')
			.setDesc('如果无法从互联网访问你的 NAS, 必须开启穿透才有可能正常同步')
			.addToggle(toggle => toggle
				.setValue(this.plugin.settings.nat)
				.onChange(async (value) => {
					this.plugin.settings.nat = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('穿透服务器')
			.setDesc('无法从互联网访问你的 NAS 时, 将由此服务器提供穿透服务')
			.addText(text => text
				.setPlaceholder('https://signal.betax.dev')
				.setValue(this.plugin.settings.signal)
				.onChange(async (value) => {
					this.plugin.settings.signal = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName('穿透标识')
			.setDesc('NAS 中的唯一标识, 让你能在穿透服务器中找到你的 NAS')
			.addText(text => text
				.setPlaceholder('#000000')
				.setValue(this.plugin.settings.devId)
				.onChange(async (value) => {
					this.plugin.settings.devId = value;
					await this.plugin.saveSettings();
				}));
	}
}
