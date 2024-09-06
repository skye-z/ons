<template>
    <div class="app-content no-select">
        <div v-if="state == 0" class="loading">
            <n-spin />
        </div>
        <template v-else-if="state == 1">
            <div class="card pa-10 flex align-center justify-between mb-10" v-for="item in list">
                <div>
                    <div class="nas-name">{{ item.name }}</div>
                    <div class="nas-id">NAT.ID {{ item.natId }}</div>
                </div>
                <div>
                    <div class="nas-time text-small text-right">
                        <n-time :time="item.lastOnline" :type="now - item.lastOnline > offset ? 'relative' : 'date'" />
                        <span class="ml-5">上线</span>
                    </div>
                    <div class="nas-time text-small text-right">
                        <n-time :time="item.lastConnect"
                            :type="now - item.lastConnect > offset ? 'relative' : 'date'" />
                        <span class="ml-5">连接</span>
                    </div>
                    <div class="nas-time mt-5 flex align-center text-small justify-end">
                        <template v-if="online[item.natId]">
                            <div class="dot dot-green"></div>
                            <div class="mr-5">在线</div>
                        </template>
                        <template v-else>
                            <div class="dot dot-red"></div>
                            <div class="mr-5">离线</div>
                        </template>
                        <template v-if="connect[item.natId]">
                            <div class="dot dot-green"></div>
                            <div>已连接</div>
                        </template>
                        <template v-else>
                            <div class="dot dot-red"></div>
                            <div>未连接</div>
                        </template>
                    </div>
                </div>
            </div>
        </template>
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
        list: [],
        connect: {},
        online: {},
        offset: 604800000,
        now: 0,
    }),
    methods: {
        init() {
            this.getList();
        },
        getList() {
            this.now = new Date().getTime()
            device.getList().then(res => {
                this.state = res.state ? 1 : 2
                if (res.state) {
                    this.list = res.data
                    if (this.list.length > 0) this.checkState()
                }
            }).catch(err => {
                this.state = 3;
            })
        },
        checkState() {
            device.getState().then(res => {
                if (res.state) {
                    this.connect = res.data.connect
                    this.online = res.data.online
                } else window.$message.warning("更新设备状态失败");
            }).catch(() => {
                window.$message.warning("更新设备状态出错");
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

.nas-name {
    font-size: 18px;
    font-weight: bold;
}

.nas-time{
    line-height: 14px;
}

.dot{
    border-radius: 8px;
    margin-right: 3px;
    height: 10px;
    width: 10px;
}

.dot.dot-red{
    background-color: #ff0000;
}

.dot.dot-green{
    background-color: #008000;
}
</style>
