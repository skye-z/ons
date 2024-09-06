<template>
    <div class="head-bar pa-10 flex align-center justify-between">
        <div class="flex align-center">
            <a href="/"><img class="logo mr-10" src="../assets/logo.png" /></a>
            <div class="title">
                <div style="font-size: 12px;line-height: 12px;">Obsidian</div>
                <div style="font-size: 22px;line-height: 22px;margin-bottom: 5px;">NAS Sync</div>
            </div>
        </div>
        <div class="right-bar flex align-center justify-end">
            <n-button quaternary circle @click="toggleTheme">
                <template #icon>
                    <n-icon>
                        <Lightbulb24Filled v-if="isDark" />
                        <LightbulbFilament24Filled v-else />
                    </n-icon>
                </template>
            </n-button>
        </div>
    </div>
</template>

<script>
import { Lightbulb24Filled, LightbulbFilament24Filled } from '@vicons/fluent'
import { Server, Github } from '@vicons/fa'
import { useThemeStore } from '../plugins/store'

export default {
    name: "HeadBar",
    components: { Server, Github, Lightbulb24Filled, LightbulbFilament24Filled },
    computed: {
        isDark() {
            const themeStore = useThemeStore();
            return themeStore.isDark;
        }
    },
    methods: {
        toggleTheme() {
            const themeStore = useThemeStore();
            themeStore.toggleTheme();
        }
    }
};
</script>

<style scoped>
.logo {
    width: 48px;
    border-radius: 25px;
}
</style>
