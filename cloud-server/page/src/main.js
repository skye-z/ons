import { createApp } from 'vue'
import router from './plugins/router'
import { createPinia } from 'pinia'
import { useThemeStore } from './plugins/store'
import './style.css'
import App from './App.vue'
// 导入等宽字体
import 'vfonts/FiraCode.css'

const app = createApp(App)
const pinia = createPinia()


app.use(router)
app.use(pinia)

const themeStore = useThemeStore()
themeStore.loadThemeFromLocalStorage()

themeStore.$subscribe((mutation, state) => {
  document.body.classList.toggle('dark', state.isDark);
  themeStore.saveThemeToLocalStorage();
})

app.mount('#app')