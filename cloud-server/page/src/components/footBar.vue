<template>
    <div class="flex align-center justify-center full-width text-gray pa-10">
        <div class="mr-5">Copyright &copy; {{ year }}</div>
        <a href="https://github.com/skye-z" rel="noopener" target="_blank">
            <n-button text class="text-gray">Skye Zhang</n-button>
        </a>
        <div class="ml-10">&hearts;</div>
        <div class="ml-10 mr-5">Built on <a href="https://github.com/skye-z/nas-sync" rel="noopener"
                target="_blank"><n-button text class="text-gray">BetaX NAS Sync</n-button></a></div>
    </div>
</template>
<script>
export default {
    name: "FootBar",
    data: () => ({
        login: false,
        year: '',
    }),
    mounted() {
        this.year = new Date().getFullYear();
    },
};
</script>
<style scoped></style>