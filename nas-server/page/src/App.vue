<template>
  <n-config-provider :locale="i18n.main" :date-locale="i18n.date" :theme-overrides="theme">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-message-provider>
          <n-notification-provider>
            <global-api />
          </n-notification-provider>
        </n-message-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
    <n-scrollbar style="height: 100vh">
      <div id="app-center">
        <head-bar />
        <div class="app-content">
          <div class="flex justify-center">
            <info-card ref="infoCard" />
            <control-card />
          </div>
          <div class="text-center mt-10 text-gray">警告: 本服务可视化控制面板不受密码保护, 切勿暴露到公网</div>
        </div>
        <foot-bar />
      </div>
    </n-scrollbar>
  </n-config-provider>
</template>

<script>
import { zhCN, dateZhCN, darkTheme, lightTheme } from 'naive-ui'
import GlobalApi from './components/globalApi.vue'
import { useThemeStore } from './plugins/store'

export default {
  name: "App",
  components: { GlobalApi },
  data: () => ({
    i18n: {
      main: zhCN,
      date: dateZhCN
    }
  }),
  computed: {
    isDark() {
      const themeStore = useThemeStore();
      return themeStore.isDark;
    },
    theme() {
      const themeStore = useThemeStore();
      return themeStore.isDark ? darkTheme : lightTheme;
    }
  },
  methods: {
    init() {
      let urlParams = new URLSearchParams(window.location.search);
      if (urlParams.get('token')) {
        localStorage.setItem("access:token", urlParams.get('token'));
        this.$refs.infoCard.registerNext()
      }
    },
  },
  mounted() {
    document.body.classList.toggle('dark', this.isDark);
    this.init()
  }
};
</script>