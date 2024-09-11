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
      { text: '插件', link: '/obsidian/setting' },
      {
        text: '服务', items: [
          { text: 'NAS服务', link: '/nas/register' },
          { text: '中控服务', link: '/cloud/register' }
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
          { text: '初始化', link: '/guide/reinit' }
        ]
      },
      {
        text: 'Obsidian',
        items: [
          { text: '基本配置', link: '/obsidian/setting' },
          { text: '常见问题', link: '/obsidian/problem' },
        ]
      },
      {
        text: 'NAS 服务',
        items: [
          { text: '注册设备', link: '/nas/register' },
          { text: '同步控制', link: '/nas/connect' },
          { text: '连接密码', link: '/nas/pass' },
        ]
      },
      {
        text: '中控服务',
        items: [
          { text: '注册登录', link: '/cloud/register' },
          { text: '设备管理', link: '/cloud/device' },
          { text: '自行部署', link: '/cloud/deploy' },
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
