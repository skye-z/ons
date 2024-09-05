package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/nas-server/util"
)

type Controller struct {
	Server *P2PServer
}

func CreateController() *Controller {
	control := &Controller{}
	if util.GetBool("connect.auto") {
		control.open()
	}
	return control
}

func (c *Controller) GetStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"connected": c.Server.connect != nil && !c.Server.isClosed(),
	})
}

func (c *Controller) Connect(ctx *gin.Context) {
	c.open()
	ctx.JSON(http.StatusOK, gin.H{"message": "连接已开启"})
}

func (c *Controller) open() {
	if c.Server != nil {
		c.Server.connect = nil
		c.Server.ticker.Stop()
	}
	c.Server = NewP2PServer(util.GetString("connect.natId"), util.GetString("connect.server"))
}

func (c *Controller) Disconnect(ctx *gin.Context) {
	if c.Server != nil && c.Server.connect != nil {
		c.Server.connect.Close()
		c.Server.connect = nil
		c.Server.ticker.Stop()
		c.Server.ticker = time.NewTicker(5 * time.Minute) // 重启定时器
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "连接已关闭"})
}
