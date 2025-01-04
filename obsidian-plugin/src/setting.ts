import NSPlugin from 'main';
import { App, Notice, PluginSettingTab, Setting } from 'obsidian';
import t from './i18n/locale';

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
			.setName(t("SETTING_CONTROL_SERVER"))
			.setDesc(t("SETTING_CONTROL_SERVER_DESC"))
			.addText(text => text
				.setPlaceholder('ons.betax.dev')
				.setValue(this.plugin.settings.server)
				.onChange(async (value) => {
					this.plugin.settings.server = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName(t("SETTING_NAT_ID"))
			.setDesc(t("SETTING_NAT_ID_DESC"))
			.addText(text => text
				.setPlaceholder('000000')
				.setValue(this.plugin.settings.devId)
				.onChange(async (value) => {
					this.plugin.settings.devId = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName(t("SETTING_PASSWORD"))
			.setDesc(t("SETTING_PASSWORD_DESC"))
			.addText(text => text
				.setPlaceholder('8-24位')
				.setValue(this.plugin.settings.pwd)
				.onChange(async (value) => {
					this.plugin.settings.pwd = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName(t("SETTING_MAIN_STUN"))
			.setDesc(t("SETTING_MAIN_STUN_DESC"))
			.addText(text => text
				.setPlaceholder('stun:domain.com:443')
				.setValue(this.plugin.settings.stunMain)
				.onChange(async (value) => {
					this.plugin.settings.stunMain = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName(t("SETTING_BACKUP_STUN"))
			.setDesc(t("SETTING_BACKUP_STUN_DESC"))
			.addText(text => text
				.setPlaceholder('stun:domain.com:443')
				.setValue(this.plugin.settings.stunBackup)
				.onChange(async (value) => {
					this.plugin.settings.stunBackup = value;
					await this.plugin.saveSettings();
				}));
		new Setting(containerEl)
			.setName(t("SETTING_TEST"))
			.setDesc(t("SETTING_TEST_DESC"))
			.addButton(text => text
				.setButtonText("Start")
				.onClick(async () => {
					new Notice("开始测试, 请留意右上角提示与右下角的状态")
					this.plugin.initPeerManager();
				}));
	}
}
