<template>
    <div></div>
</template>
<script setup>
import { useMessage, useDialog, useNotification, useLoadingBar } from 'naive-ui'

window.$notify = useNotification()
window.$message = useMessage()
window.$dialog = useDialog()
window.$loading = useLoadingBar()
</script>