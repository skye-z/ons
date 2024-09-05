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
            <n-button class="mr-10" strong secondary v-if="info.id" @click="openList">
                <template #icon>
                    <n-icon>
                        <Server />
                    </n-icon>
                </template>
                我的设备
            </n-button>
            <a href="/login" rel="console" v-else>
                <n-button class="mr-10" strong secondary>
                    <template #icon>
                        <n-icon>
                            <Github />
                        </n-icon>
                    </template>
                    Github 登录
                </n-button>
            </a>
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
import { user } from '../plugins/api'

export default {
    name: "HeadBar",
    components: { Server, Github, Lightbulb24Filled, LightbulbFilament24Filled },
    data: () => ({
        info: {}
    }),
    computed: {
        isDark() {
            const themeStore = useThemeStore();
            return themeStore.isDark;
        }
    },
    methods: {
        init() {
            if(localStorage.getItem('access:token') == undefined) return false
            user.now().then(res => {
                if (res.state) {
                    localStorage.setItem('cache:user', JSON.stringify(res.data))
                    this.info = res.data;
                } else this.info = {}
            }).catch(() => {
                this.info = {}
            })
        },
        openList() {
            this.$router.push('/list')
        },
        toggleTheme() {
            const themeStore = useThemeStore();
            themeStore.toggleTheme();
        }
    },
    mounted() {
        this.init()
    },
    watch: {
        '$route': function (to, from) {
            if (from.name == 'Auth' && to.name != 'Auth') {
                this.init()
            }
        },
    },
};
</script>

<style scoped>
.logo {
    width: 48px;
    border-radius: 25px;
}
</style>
