<template>
    <div class="card control-card mb-10">
        <div class="control-state border-bottom" :class="{ online: state }">
            <n-icon size="42">
                <LinkSquare24Filled />
            </n-icon>
            <div>{{ state ? '在线' : '离线' }}</div>
        </div>
        <n-button class="control-btn full-width" v-if="state" type="warning" @click="closeServer">断开同步网络</n-button>
        <n-button class="control-btn full-width" v-else type="primary" @click="openServer">连接同步网络</n-button>
    </div>
</template>
<script>
import { LinkSquare24Filled } from '@vicons/fluent'
import { device } from '../plugins/api'

export default {
    name: "ControlCard",
    components: { LinkSquare24Filled },
    data: () => ({
        state: false
    }),
    methods: {
        init() {
            device.getState().then(res => {
                this.state = res.state
            }).catch(() => {
                window.$message.error("发生意料之外的错误");
            })
        },
        openServer() {
            device.openServer().then(res => {
                if (res.state) this.state = true
                else window.$message.warning(res.message ? res.message : "同步网络连接失败");
            }).catch(() => {
                window.$message.error("发生意料之外的错误");
            })
        },
        closeServer() {
            device.closeServer().then(res => {
                if (res.state) this.state = false
                else window.$message.warning(res.message ? res.message : "同步网络断开失败");
            }).catch(() => {
                window.$message.error("发生意料之外的错误");
            })
        },
    },
    mounted() {
        this.init()
    },
};
</script>
<style scoped>
.control-card {
    height: 162px;
}

.control-state {
    text-align: center;
    padding: 28px 42px;
    font-weight: bold;
    color: #999;
}

.online {
    color: #58cb58;
}

.control-btn {
    border-radius: 0 0 8px 8px;
}
</style>