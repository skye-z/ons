import { defineConfig } from 'vitepress'

export default defineConfig({
  base: "/ons/",
  title: "BetaX ONS",
  description: "BetaX Obsidian NAS Sync Server",
  head: [['link', { rel: 'icon', href: '/icon/icon-light@1x.png' }]],
  themeConfig: {
    logo: '/icon/icon-light@2x.png',
    nav: [
      { text: '指南', link: '/guide/about' },
      { text: '插件', link: '/guide/about' },
      {
        text: '服务', items: [
          { text: 'NAS服务', link: '/guide/about' },
          { text: '中控服务', link: '/guide/about' }
        ]
      },
    ],
    sidebar: [
      {
        text: '指南',
        items: [
          { text: '关于 ONS', link: '/guide/about' },
          { text: '快速开始', link: '/guide/start' },
          { text: '兼容性', link: '/guide/compatible' },
          { text: '初始化', link: '/guide/reinit' },
          { text: '最小化', link: '/guide/min' }
        ]
      },
      {
        text: 'Obsidian',
        items: [
          { text: '基本配置', link: '/guide/about' },
          { text: '常见问题', link: '/guide/about' },
        ]
      },
      {
        text: 'NAS 服务',
        items: [
          { text: '注册设备', link: '/guide/about' },
          { text: '同步控制', link: '/guide/about' },
          { text: '连接密码', link: '/guide/about' },
        ]
      },
      {
        text: '中控服务',
        items: [
          { text: '注册登录', link: '/guide/about' },
          { text: '设备管理', link: '/guide/about' },
          { text: '自行部署', link: '/guide/about' },
        ]
      }
    ],
    outline: {
      level: [2, 6]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/skye-z/ons' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present Skye Zhang'
    },
    search: {
      provider: 'local'
    },
  },
  markdown: {
    image: {
      lazyLoading: true
    }
  }
})
