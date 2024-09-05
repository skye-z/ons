// stores/theme.js
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    isDark: false,
  }),
  actions: {
    toggleTheme() {
      this.isDark = !this.isDark;
      this.saveThemeToLocalStorage();
    },
    saveThemeToLocalStorage() {
      localStorage.setItem('cache:theme', this.isDark);
    },
    loadThemeFromLocalStorage() {
      const storedTheme = localStorage.getItem('cache:theme');
      if (storedTheme !== null) {
        this.isDark = storedTheme === 'true';
      }
    },
  },
})