package main

import (
	"embed"
	"flag"

	"github.com/skye-z/nas-sync/cloud-server/core"
	"github.com/skye-z/nas-sync/cloud-server/util"
)

//go:embed page/dist/*
var page embed.FS

func main() {
	// 初始化系统配置
	util.InitConfig()
	// 初始化数据库
	engine := util.InitDB()
	go util.InitDBTable(engine)
	// 定义一个命令行参数
	debug := flag.Bool("debug", false, "output debug logs")
	// 定义一个命令行参数
	port := flag.Int("port", 9891, "the port to listen on")
	// 解析命令行参数
	flag.Parse()
	// 初始化路由器
	router := core.BuildRouter(!*debug, *port, "0.0.0.0", util.GetString("basic.sslCert"), util.GetString("basic.sslKey"), engine, page)
	router.Run(engine)
}
