---
layout: home

hero:
  name: "BetaX"
  text: "Obsidian\nNAS Sync Server"
  tagline: 让 Obsidian 与你家中内网的 NAS 时刻保持同步
  image:
    light: /icon/icon-dark@4x.png
    dark: /icon/icon-dark@4x.png
    alt: VitePress
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/start
    - theme: alt
      text: Github
      link: https://github.com/skye-z/ons

features:
  - icon: 💰
    title: 零成本
    details: 采用 P2P 技术实现 Obsidian 客户端与 NAS 端的点对点传输, 仅自行部署云端中控服务器存在成本
  - icon: 🚀
    title: 传输快
    details: 当通过中控服务器建立连接后, 点对点传输将得同步服务无需再进行中转, 可以最大化利用你的带宽
  - icon: 🧠
    title: 易上手
    details: 提供简单易用的可视化界面, 大幅精简必要的操作步骤, 以及完善的使用文档, 都可以帮助你快速上手
  - icon: ⚡️
    title: 开销低
    details: 闲置时仅 ~3% CPU 和 ~20MB 内存开销, 在传输 100MB 文件时也仅 3 ~ 8% CPU 和 ~50MB 的内存开销
---

<style>
  :root {
    --vp-button-brand-bg: #08BDC9;
    --vp-button-brand-hover-bg: #25D8E4;

    --vp-home-hero-name-color: transparent;
    --vp-home-hero-name-background: -webkit-linear-gradient(120deg, #bd34fe 30%, #08BDC9 70%);

    --vp-home-hero-image-filter: blur(68px);
    --vp-home-hero-image-background-image: linear-gradient(-45deg, #bd34fe 30%, #08BDC9 70%);
  }

  .VPHome{
    padding-bottom: 30px !important;
  }
</style>
