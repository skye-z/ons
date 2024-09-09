<template>
    <div class="app-content no-select">
        <div class="gradient"></div>
        <div class="hello">
            <n-gradient-text class="logo" type="info">BetaX</n-gradient-text>
            <div class="name">Obsidian</div>
            <div class="name">NAS Sync Server</div>
            <div class="slogan">让 Obsidian 与你家中内网的 NAS 时刻保持同步</div>
            <div class="mt-10 flex align-center">
                <a class="flex" href="https://skye-z.github.io/ons/" target="_blank" style="text-decoration: none;">
                    <n-button round type="primary" class="mr-10">使用帮助</n-button>
                </a>
                <a class="flex" href="https://github.com/skye-z/ons/stargazers" target="_blank">
                    <img class="stars" :src="start" />
                </a>
            </div>
        </div>
    </div>
</template>

<script>

export default {
    name: "Home",
    data: () => ({
        start: 'https://img.shields.io/github/stars/skye-z/ons?style=for-the-badge&logo=github&labelColor=%2318a058&color=%2350cf8c'
    }),
    methods: {
        init() {
        },
    },
    mounted() {
        this.init()
    },
};
</script>

<style scoped>
.hello {
    position: fixed;
    z-index: 99;
    top: 20vh;
}

.logo {
    font-size: 48px;
    line-height: 50px;
    font-weight: bold;
}

.name {
    font-size: 48px;
    line-height: 50px;
}

.slogan {
    font-size: 24px;
    margin-top: 10px;
}

.stars {
    border-radius: 18px;
    height: 34px;
}

@keyframes rotate {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}

.gradient {
    --size: 30vw;
    --speed: 50s;
    --easing: cubic-bezier(0.8, 0.2, 0.2, 0.8);
    transition: all 0.5s ease-out;

    width: var(--size);
    height: var(--size);
    filter: blur(calc(var(--size) / 5));
    background-image: linear-gradient(rgb(67, 119, 239), rgb(123, 239, 208));
    animation: rotate var(--speed) var(--easing) alternate infinite;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    position: fixed;
    right: 20vw;
    top: 18vh;
    z-index: 1;
}

@media only screen and (max-width:1000px) {
    .gradient {
        --size: 50vw;
        right: 20px;
    }
}
</style>
