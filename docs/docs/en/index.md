---
layout: home

hero:
  name: "BetaX"
  text: "Obsidian\nNAS Sync Server"
  tagline: Keep Obsidian synced with your home NAS
  image:
    light: /icon/icon-dark@4x.png
    dark: /icon/icon-dark@4x.png
    alt: VitePress
  actions:
    - theme: brand
      text: Quick Start
      link: /en/guide/start
    - theme: alt
      text: Github
      link: https://github.com/skye-z/ons

features:
  - icon: 💰
    title: Zero Cost
    details: P2P tech, self-deploy cloud server
  - icon: 🚀
    title: Fast Transfer
    details: Direct P2P sync, max bandwidth use
  - icon: 🧠
    title: User-Friendly
    details: Simple UI, easy steps, full docs
  - icon: ⚡️
    title: Low Overhead
    details: Idle ~3% CPU, ~20MB RAM. Transfer 3-8% CPU, ~50MB RAM
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
