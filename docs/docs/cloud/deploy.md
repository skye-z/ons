# 自行部署

## Docker 容器部署

1. 拉取`skyezhang/ons-server`
2. 开放`80`和`443`端口
3. 启动容器

## 初始化设置

当你安装完成中控服务后需要编辑`config.ini`:

1. 将你的`Github OAuth2`信息填写到`github`中
2. 访问中控服务, 第一个访问者将成为管理员
3. 修改`register`, 决定是否开启注册
