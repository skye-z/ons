package core

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/ons/nas-server/util"
)

type Setting struct {
	Hostname string `json:"hostname"`
	Auto     bool   `json:"auto"`
	Server   string `json:"server"`
	NatId    string `json:"natId"`
	Password string `json:"password"`
}

type SettingServer struct {
}

func CreateSettingServer() *SettingServer {
	return &SettingServer{}
}

func (ss SettingServer) Get(ctx *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "匿名主机"
	}
	setting := &Setting{
		Hostname: hostname,
		Auto:     util.GetBool("connect.auto"),
		Server:   util.GetString("connect.server"),
		NatId:    util.GetString("connect.natId"),
		Password: util.GetString("connect.password"),
	}
	util.ReturnData(ctx, true, setting)
}

func (ss SettingServer) SetPassword(ctx *gin.Context) {
	util.Set("connect.password", util.GenerateRandomString(8))
	util.ReturnMessage(ctx, true, "密码已更新")
}
