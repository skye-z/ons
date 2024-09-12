<template>
    <div class="card info-card mr-10">
        <div class="flex justify-between border-bottom">
            <div class="info-item pa-10 border-right">
                <div class="info-label">设备名称</div>
                <div class="info-value line1">{{ info.hostname }}</div>
            </div>
            <div class="info-item pa-10">
                <div class="info-label">NAT.ID</div>
                <div class="info-value">{{ info.natId ? info.natId : '请先向中控注册设备' }}</div>
            </div>
        </div>
        <div class="flex justify-between border-bottom">
            <div class="info-item pa-10 border-right">
                <div class="repass float-right" @click="buildPassword">
                    <n-icon>
                        <ArrowClockwise12Regular />
                    </n-icon>
                </div>
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
            <n-button class="info-auto full-width" v-if="info.auto" type="warning" @click="switchAuto">关闭自动启动</n-button>
            <n-button class="info-auto full-width" v-else type="primary" @click="switchAuto">开启自动启动</n-button>
        </template>
        <div v-if="wait" class="loading">
            <n-spin size="medium" />
        </div>
    </div>
</template>
<script>
import { ArrowClockwise12Regular } from '@vicons/fluent'
import { setting, device } from '../plugins/api'

export default {
    name: "InfoCard",
    components: { ArrowClockwise12Regular },
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
        register() {
            window.open('http://' + this.info.server + '/app/oauth2?uri=' + location.origin)
        },
        registerNext() {
            this.wait = true
            device.register().then(res => {
                this.wait = false
                if (res.state) {
                    this.info.natId = res.data
                    window.$message.success("设备注册成功");
                    location.href = location.origin
                } else window.$message.warning(res.message ? res.message : "注册失败");
            }).catch(() => {
                this.wait = false
                window.$message.error("发生意料之外的错误");
            })
        },
        buildPassword() {
            window.$dialog.warning({
                title: "操作确认",
                content: (this.info.password === '' ? '开启连接密码后将不再允许关闭' : '重新生成连接密码后会立即生效, 正在传输的不受影响, 但后续传输需使用新的密码') + ", 确认要继续吗?",
                positiveText: "确认",
                negativeText: "取消",
                onPositiveClick: () => {
                    setting.setPassword().then(res => {
                        if (res.state) {
                            this.init()
                            window.$message.success("连接密码已更新");
                        } else window.$message.warning(res.message ? res.message : "密码生成失败");
                    }).catch(() => {
                        window.$message.error("发生意料之外的错误");
                    })
                },
            });
        },
        switchAuto() {
            device.switchAuto().then(res => {
                if (res.state) {
                    this.info.auto = !this.info.auto
                } else window.$message.warning(res.message ? res.message : "切换失败");
            }).catch(() => {
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

.repass {
    cursor: pointer;
    font-size: 22px;
    line-height: 22px;
    margin: -6px -6px 0 0;
}

.repass:hover {
    color: #999999;
}
</style>