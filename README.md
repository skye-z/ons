# BetaX Obsidian NAS Sync

[中文](README_zh.md)

This is an `Obsidian Vault` synchronization service that allows you to synchronize your `Obsidian` library with your home intranet NAS.

For users who have a NAS at home and use it as a storage hub.

## Prerequisites

1. a NAS with [Nas Server](nas-server) installed
2. a recent version of the `Obsidian` client.

## Instructions for use

1. You need to register your NAS with the central server to get the `NAT.ID`. 2.
2. Enable the connection password on the `NAS Server` (it is better to enable it for security reasons)
3. Install and enable the plugin
4. Fill in the `唯一标识` and `连接密码` in the Settings.
5. Click `Start Test` to test if the connection works.

### How to synchronize

> The following information is very important, please read it carefully !!!!

Are you using this service for the first time?

- Yes, I have never used any of the `BetaX Obsidian NAS Sync` services.
  - First of all, please make sure that `Vault` does not exist in `NAS Server`, otherwise the data may be overwritten, after checking, establish the connection, and then start to use it if the connection is successful.
- No, I have used `NAS Server` or this plugin.
  - I want to pull the `Vault` from `NAS Server` to local.
    - First make sure the local `Obsidian Vault` is a brand new `Vault`, then establish a connection, open the command panel and select `BetaX NAS Sync: 手动更新`.
  - I want to push the local `Vault` to the `NAS Server`.
    - Push to the new `NAS Server`.
      - First make sure that the `Vault` does not exist in the `NAS Server`, then establish a connection, open the command panel and select `BetaX NAS Sync: 手动更新` to push to the old `NAS Server`.
    - Push to the old `NAS Server`.
      - This is a dangerous operation, you must ensure that the `Vault` on both ends is the same before you can continue to use it, if you force it, it will cause the local `Obsidian Vault` to be overwritten.

## Unable to connect

If you can't connect after configuring `NAS Server` and this plugin, please change your network environment and try again.

This service uses modern P2P hole-punching technology, but there are still some special networks that can't be penetrated, if you change networks several times and still can't connect, it may mean that your network environment can't use this service.
