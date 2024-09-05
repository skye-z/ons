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
        <router-view />
        <foot-bar />
      </div>
    </n-scrollbar>
  </n-config-provider>
</template>

<script>
import { zhCN, dateZhCN, darkTheme, lightTheme } from 'naive-ui'
import GlobalApi from './components/globalApi.vue'
import HeadBar from './components/headBar.vue'
import FootBar from './components/footBar.vue'
import { useThemeStore } from './plugins/store'
import theme from './theme.json'

export default {
  name: "App",
  components: { GlobalApi, HeadBar },
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
    },
  },
  mounted() {
    document.body.classList.toggle('dark', this.isDark);
    this.init()
  }
};
</script>