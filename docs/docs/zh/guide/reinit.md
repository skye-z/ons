# 初始化

## 初始配置

### NAS 服务

当你安装完成`NAS`服务端后需要:

1. 向中控注册设备, 注册后`NAS`会获得一个`NAT.ID`
2. 开启“连接密码”, 系统会为你生成一个安全的随机密码
3. 开启“开启自动启动”

### Obsidian 插件

当你安装完成`Obsidian`插件后需要在设置中:

1. 填写“唯一标识”, 即向中控注册`NAS`时获得的`NAT.ID`
2. 填写“连接密码”

## 存储库配置(重要)

请根据实际情况选择对应章节进行

- 我从未使用过本服务
  - 我在 `NAS` 上存储有 `Obsidian Vault`
    - G4 方案
  - 我只在本地有 `Obsidian Vault`
    - G1 方案
- 我经使用过本服务中任意组件
  - 我在 `NAS Server` 中已有 `Vault`
    - G2 方案
  - 我在本地 `Obsidian` 中已有同步过的 `Vault`
    - G3 方案

### G1 方案

你无需额外操作, 按照前面的步骤完成连接后正常使用即可

### G2 方案

由于你已有 `NAS Vault` 所以你需要先确认

- 本地 `Obsidian` 是一个全新的 `Vault`
- 本地 `Obsidian Vault` 中的文件与数据允许被覆盖

> Q: 如果我有一个本地 `Obsidian Vault` 而且数据重要不能覆盖怎么办呢? <br/> A: 请创建一个新的 `Obsidian Vault`

确认无误后按照前面的步骤完成连接, 连接成功后在 `Obsidian` 中打开命令面板选择`BetaX NAS Sync: 手动更新`即可

### G3 方案

现在你已有本地 `Obsidian Vault`, 如果你没有 `NAS Vault`, 那么你可以跳转到[G1 方案](#G1 方案)来将数据推送到新的 `NAS Server` 中

如果你有包含 `Vault` 数据的 `NAS Server` , 那就要判断此本地 `Vault` 和 `NAS Vault` 是否一致, 如果一致, 那么你可以跳转到[G1 方案](#G1 方案), 二者根据时间先后同步差异即可

如果不一致, 那么你需要清除无需保留的那一端, 最好是整个文件夹删除后重新创建

清除完成后按照前面的步骤完成连接, 连接成功后在 `Obsidian` 中打开命令面板选择`BetaX NAS Sync: 手动更新`即可

### G4 方案

看来你曾经使用 `NAS` 简单的文件共享存储过 `Obsidian Vault`

你需要先部署完成 `NAS Server` 后, 将你 `NAS` 上的 `Obsidian Vault` 映射到容器的 `/app/vault` 目录

然后在 `Obsidian Vault` 目录中创建 `.synclog` 文件, 文件没有任何后缀

创建文件后用记事本或其他编辑器打开它, 在里面输入以秒为单位的10位时间戳数字, 最后保存修改

上述完成后按照最前面的步骤完成连接, 连接成功后在 `Obsidian` 中打开命令面板选择`BetaX NAS Sync: 手动更新`即可自动拉取 `NAS` 上的 `Obsidian Vault` 到本地
