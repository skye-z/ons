/*
数据库工具

BetaX Blog
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"github.com/skye-z/nas-sync/cloud-server/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func InitDB() *xorm.Engine {
	OutLogf("Data", "load engine")
	engine, err := xorm.NewEngine("sqlite", "./local.store")
	if err != nil {
		panic(err)
	}
	return engine
}

func InitDBTable(engine *xorm.Engine) {
	OutLogf("Data", "load table")
	err := engine.Sync2(new(model.User))
	if err != nil {
		panic(err)
	}
	err = engine.Sync2(new(model.Device))
	if err != nil {
		panic(err)
	}
}
