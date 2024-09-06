<template>
    <div class="tips-box no-select">
        <n-result :status="state ? 'success' : 'warning'" :title="state ? '授权登录成功' : '授权失败'"
            :description="state ? '正在跳转, 请稍后...' : '请先登录后再试'">
        </n-result>
    </div>
</template>

<script>

export default {
    name: "OAuth",
    data: () => ({
        state: false
    }),
    methods: {
        init() {
            let token = localStorage.getItem('access:token');
            this.state = token != undefined
            setTimeout(() => {
                if (token) location.href = this.$route.query.uri + "?token=" + token
                else this.$router.push('/')
            }, 1500)
        },
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
