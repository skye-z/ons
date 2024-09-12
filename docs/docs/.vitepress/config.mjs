import { defineConfig } from 'vitepress'

export default defineConfig({
  base: "/ons/",
  title: "BetaX ONS",
  description: "BetaX Obsidian NAS Sync Server",
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      link: '/en',
      themeConfig: {
        logo: '/icon/icon-light@2x.png',
        nav: [
          { text: 'Guide', link: '/en/guide/about' },
          { text: 'Plugin', link: '/en/obsidian/setting' },
          {
            text: 'Service', items: [
              { text: 'NAS Service', link: '/en/nas/register' },
              { text: 'Centr Control Server', link: '/en/cloud/register' }
            ]
          },
        ],
        sidebar: [
          {
            text: 'Guide',
            items: [
              { text: 'About ONS', link: '/en/guide/about' },
              { text: 'Quick Start', link: '/en/guide/start' },
              { text: 'Compatibility', link: '/en/guide/compatible' },
              { text: 'Initialization', link: '/en/guide/reinit' }
            ]
          },
          {
            text: 'Obsidian',
            items: [
              { text: 'Basic Config', link: '/en/obsidian/setting' },
              { text: 'Common Problem', link: '/en/obsidian/problem' },
            ]
          },
          {
            text: 'NAS Service',
            items: [
              { text: 'Register Device', link: '/en/nas/register' },
              { text: 'Sync Control', link: '/en/nas/connect' },
              { text: 'Connection Password', link: '/en/nas/pass' },
            ]
          },
          {
            text: 'Centr Control Server',
            items: [
              { text: 'Register & Login', link: '/en/cloud/register' },
              { text: 'Device Management', link: '/en/cloud/device' },
              { text: 'Self Deployment', link: '/en/cloud/deploy' },
            ]
          }
        ]
      },
    },
    zh: {
      label: '简体中文',
      lang: 'zh',
      link: '/zh',
      themeConfig: {
        logo: '/icon/icon-light@2x.png',
        nav: [
          { text: '指南', link: '/zh/guide/about' },
          { text: '插件', link: '/zh/obsidian/setting' },
          {
            text: '服务', items: [
              { text: 'NAS服务', link: '/zh/nas/register' },
              { text: '中控服务', link: '/zh/cloud/register' }
            ]
          },
        ],
        sidebar: [
          {
            text: '指南',
            items: [
              { text: '关于 ONS', link: '/zh/guide/about' },
              { text: '快速开始', link: '/zh/guide/start' },
              { text: '兼容性', link: '/zh/guide/compatible' },
              { text: '初始化', link: '/zh/guide/reinit' }
            ]
          },
          {
            text: 'Obsidian',
            items: [
              { text: '基本配置', link: '/zh/obsidian/setting' },
              { text: '常见问题', link: '/zh/obsidian/problem' },
            ]
          },
          {
            text: 'NAS 服务',
            items: [
              { text: '注册设备', link: '/zh/nas/register' },
              { text: '同步控制', link: '/zh/nas/connect' },
              { text: '连接密码', link: '/zh/nas/pass' },
            ]
          },
          {
            text: '中控服务',
            items: [
              { text: '注册登录', link: '/zh/cloud/register' },
              { text: '设备管理', link: '/zh/cloud/device' },
              { text: '自行部署', link: '/zh/cloud/deploy' },
            ]
          }
        ],
      },
    }
  },
  head: [['link', { rel: 'icon', href: '/ons/icon/icon-light@1x.png' }]],
  themeConfig: {
    search: {
      provider: 'local',
      options: {
        locales: {
          zh: {
            translations: {
              button: {
                buttonText: '搜索文档',
                buttonAriaLabel: '搜索文档'
              },
              modal: {
                noResultsText: '无法找到相关结果',
                resetButtonTitle: '清除查询条件',
                footer: {
                  selectText: '选择',
                  navigateText: '切换'
                }
              }
            }
          }
        }
      }
    },
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
  },
  markdown: {
    image: {
      lazyLoading: true
    }
  }
})
