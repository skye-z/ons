<template>
    <div class="tips-box no-select">
        <n-result :status="state == 9 ? 'success' : 'warning'" :title="getTitle()" :description="getDescription()">
            <template v-if="state != 9" #footer>
                <n-button @click="back">返回首页</n-button>
            </template>
        </n-result>
    </div>
</template>

<script>

export default {
    name: "Auth",
    data: () => ({
        state: 0
    }),
    methods: {
        init() {
            console.log('[Auth] login...')
            let query = this.$route.query;
            this.state = query.state;
            if (query.state == 9) {
                localStorage.setItem('access:token', query.code)
                setTimeout(() => {
                    this.$router.push('/')
                }, 1500);
            }
        },
        getTitle() {
            if (this.state == 0) return '请先登录'
            else if (this.state == 1) return '授权服务不可用'
            else if (this.state == 2) return '授权信息无效'
            else if (this.state == 3) return '账户不存在'
            else if (this.state == 4) return '令牌签发失败'
            else if (this.state == 9) return '登录成功'
        },
        getDescription() {
            if (this.state == 0) return '禁止访问受限资源, 请登录账户后再试'
            else if (this.state == 1) return '无法与 Github OAuth2 服务建立连接, 请检查网络与授权配置'
            else if (this.state == 2) return '无法解析回调数据, 请检查 OAuth2 授权配置'
            else if (this.state == 3) return '当前授权的 Github 账户与系统绑定账户不一致'
            else if (this.state == 4) return '生成令牌时发生错误, 请检查密钥配置是否正确'
            else if (this.state == 9) return '欢迎回来, 正在跳转控制台...'
        },
        back() {
            this.$router.push('/')
        }
    },
    mounted() {
        this.init()
    },
};
</script>

<style scoped>
.tips-box {
    min-height: calc(100vh - 114px);
    padding-top: 10vh;
}
</style>
