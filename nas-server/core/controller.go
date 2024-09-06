package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/nas-sync/nas-server/util"
)

type Controller struct {
	Server *P2PServer
}

func CreateController() *Controller {
	control := &Controller{}
	if util.GetBool("connect.auto") && util.GetString("connect.natId") != "" {
		log.Println("[P2P] auto connect")
		control.open()
	}
	return control
}

func (c *Controller) Register(ctx *gin.Context) {
	if util.GetString("connect.natId") != "" {
		util.ReturnMessage(ctx, false, "已完成设备注册, 请勿重复操作")
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "匿名主机"
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/nas/register", util.GetString("connect.server")), strings.NewReader(fmt.Sprintf("name=%s", hostname)))
	if err != nil {
		util.ReturnMessage(ctx, false, "中控服务器地址不可用")
		return
	}

	code := ctx.Request.Header.Get("Authorization")
	if code == "" {
		util.ReturnError(ctx, util.Errors.NotLoginError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", code)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		util.ReturnMessage(ctx, false, "中控服务器无响应")
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		util.ReturnMessage(ctx, false, "中控服务器故障")
		return
	}
	msg := util.ToReturnMessage(bodyBytes)
	if msg.State {
		util.Set("connect.natId", msg.Message)
		util.ReturnMessageData(ctx, true, "注册成功", msg.Message)
	} else {
		util.ReturnMessage(ctx, false, msg.Message)
	}
}

func (c *Controller) GetStatus(ctx *gin.Context) {
	if c.Server != nil && c.Server.connect != nil {
		util.ReturnMessage(ctx, true, "已连接")
	} else {
		util.ReturnMessage(ctx, false, "未连接")
	}
}

func (c *Controller) Connect(ctx *gin.Context) {
	if util.GetString("connect.natId") == "" {
		util.ReturnMessage(ctx, false, "请先完成设备注册")
		return
	}
	go c.open()
	util.ReturnMessage(ctx, true, "开始连接")
}

func (c *Controller) open() {
	if c.Server != nil {
		c.Server.connect = nil
		c.Server.ticker.Stop()
	}
	c.Server = NewP2PServer(util.GetString("connect.natId"), util.GetString("connect.server"))
}

func (c *Controller) Disconnect(ctx *gin.Context) {
	if util.GetString("connect.natId") == "" {
		util.ReturnMessage(ctx, false, "请先完成设备注册")
		return
	}
	if c.Server != nil && c.Server.connect != nil {
		c.Server.connect.Close()
		c.Server.connect = nil
		c.Server.ticker.Stop()
		c.Server.ticker = time.NewTicker(5 * time.Minute)
	}
	log.Println("[P2P] connection close")
	util.ReturnMessage(ctx, true, "连接关闭")
}
