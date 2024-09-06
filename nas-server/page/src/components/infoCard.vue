<template>
    <div class="card info-card mr-10">
        <div class="flex justify-between border-bottom">
            <div class="info-item pa-10 border-right">
                <div class="info-label">设备名称</div>
                <div class="info-value line1">{{ info.hostname }}</div>
            </div>
            <div class="info-item pa-10">
                <div class="info-label">NAT.ID</div>
                <div class="info-value">{{ info.natId ? info.natId:'请先向中控注册设备' }}</div>
            </div>
        </div>
        <div class="flex justify-between border-bottom">
            <div class="info-item pa-10 border-right">
                <div class="info-label">连接密码</div>
                <div class="info-value line1">{{ info.password ? info.password : '未设置密码' }}</div>
            </div>
            <div class="info-item pa-10">
                <div class="info-label">中控服务</div>
                <div class="info-value line1">{{ info.server }}</div>
            </div>
        </div>
        <n-button class="info-auto full-width" v-if="info.natId == ''" type="primary" @click="register">注册设备</n-button>
        <template v-else>
        <n-button class="info-auto full-width" v-if="info.auto" type="warning">关闭自动启动</n-button>
        <n-button class="info-auto full-width" v-else type="primary">开启自动启动</n-button>
        </template>
        <div v-if="wait" class="loading">
            <n-spin size="medium" />
        </div>
    </div>
</template>
<script>
import { setting, device } from '../plugins/api'

export default {
    name: "InfoCard",
    data: () => ({
        info: {
            auto: false,
            hostname: "",
            natId: "",
            password: "",
            server: "",
        },
        wait: false
    }),
    methods: {
        init() {
            setting.all().then(res => {
                if (res.state) {
                    this.info = res.data
                }
            }).catch(() => {
                window.$message.error("发生意料之外的错误");
            })
        },
        register(){
            window.open('http://' + this.info.server + '/oauth2?uri=' + location.origin)
        },
        registerNext(){
            this.wait = true
            device.register().then(res => {
                this.wait = false
                if (res.state) {
                    this.info.natId = res.data
                    window.$message.success("设备注册成功");
                } else window.$message.warning(res.message ? res.message : "注册失败");
            }).catch(() => {
                this.wait = false
                window.$message.error("发生意料之外的错误");
            })
        }
    },
    mounted() {
        this.init()
    },
};
</script>
<style scoped>
.info-card {
    height: 162px;
}

.info-item {
    min-width: 200px;
    height: 63px;
}

.info-label {
    font-size: 12px;
    line-height: 16px;
    margin-bottom: 5px;
}

.info-value {
    font-size: 18px;
    line-height: 22px;
}

.info-auto {
    border-radius: 0 0 8px 8px;
}
</style>