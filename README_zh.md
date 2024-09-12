# BetaX Obsidian NAS Sync

[English](README.md)

这是一个`Obsidian Vault`同步服务, 它能让你的`Obsidian`库与家中内网的 NAS 进行同步

适用于家中有一台 NAS, 并且将其作为存储中枢的用户

## 前提条件

1. 一台安装了 [Nas Server](nas-server) 的 NAS
2. 一个最新版本的 `Obsidian` 客户端
3. 一个 Github 账户

## 使用说明

1. 你需要先在中控服务器上注册你的 NAS, 获得`NAT.ID`
2. 在 `NAS Server` 上开启连接密码(安全起见最好开启)
3. 安装并启用本插件
4. 在设置中填写 `唯一标识` 和 `连接密码`
5. 点击 `开始测试` 测试是否可以正常连接

### 如何同步

> 下面的信息非常重要, 请一定要认真阅读!!!!

你是否首次使用本服务

- 是, 我从未使用过 `BetaX Obsidian NAS Sync` 的任何服务
  - 首先请确保 `NAS Server` 中不存在 `Vault`, 否则可能导致数据被覆盖, 检查完成后建立连接, 连接成功即可开始使用啦
- 否, 我使用过 `NAS Server` 或本插件
  - 我要将 `NAS Server` 中的 `Vault`拉取到本地
    - 首先请确保本地 `Obsidian Vault` 是一个全新的 `Vault`, 然后建立连接, 打开命令面板选择`BetaX NAS Sync: 手动更新`即可
  - 我要将本地的 `Vault` 推送到 `NAS Server`
    - 推送到新的 `NAS Server`
      - 首先请确保 `NAS Server` 中不存在 `Vault`, 然后建立连接, 打开命令面板选择`BetaX NAS Sync: 手动更新`即可
    - 推送到旧的 `NAS Server`
      - 这是一个危险操作, 必须保证两端的 `Vault` 一致才能继续使用, 如果强行使用, 会使得本地 `Obsidian Vault` 被覆盖

## 无法连接

如果你在配置好 `NAS Server` 和本插件后无法连接, 请更换网络环境后再试

本服务使用现代化的P2P打洞技术, 但仍有一些特殊网络无法穿透, 如果多次更换网络仍然无法连接, 那可能意味着你的网络环境无法使用本服务
