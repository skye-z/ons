<template>
    <div class="app-content no-select">
        <div v-if="state == 0" class="loading">
            <n-spin />
        </div>
        <div v-else-if="state == 1">
            {{ list }}
        </div>
        <n-result v-else class="tips" :status="state == 2 ? 'warning' : 'error'"
            :title="state == 2 ? '获取设备列表失败' : '获取设备列表出错'"
            :description="state == 2 ? '设备服务出现异常, 请稍后再试' : '无法与服务器建立连接'" />
    </div>
</template>

<script>
import { device } from '../plugins/api'

export default {
    name: "List",
    data: () => ({
        state: 0,
        list: []
    }),
    methods: {
        init() {
            this.getList();
        },
        getList() {
            device.getList().then(res => {
                this.state = res.state ? 1 : 2
                if (res.state) {
                    this.list = res.data
                }
            }).catch(err => {
                this.state = 3;
            })
        }
    },
    mounted() {
        this.init()
    },
};
</script>

<style scoped>
.tips {
    padding-top: 10vh;
}
</style>
