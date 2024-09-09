package core

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/ons/nas-server/util"
)

const MODEL_NAME = "Core"

type Router struct {
	// 路由对象
	Object *gin.Engine
	// 端口号
	port string
	// 监听地址
	host string
	// 证书地址
	cert string
	// 公钥地址
	key string
}

// 构建路由器
func BuildRouter(release bool, port int, host, cert, key string, page embed.FS) *Router {
	if release {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.Discard
	router := &Router{
		Object: gin.Default(),
		cert:   cert,
		key:    key,
	}
	// 配置端口
	if port == 0 {
		if cert != "" && key != "" {
			router.port = "443"
		} else {
			router.port = "80"
		}
	} else {
		router.port = fmt.Sprint(port)
	}
	// 配置监听地址
	if host == "" {
		router.host = "0.0.0.0"
	} else {
		router.host = host
	}
	// 载入静态页面
	appPage, _ := fs.Sub(page, "page/dist")
	router.Object.StaticFS("/app", http.FS(appPage))

	// 挂载公共路由
	addPublicRoute(router.Object)
	// 兼容路由
	router.Object.NoRoute(func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/app"):
			getIndexFile(c, http.FS(appPage))
		default:
			c.Status(http.StatusNotFound)
		}
	})
	return router
}

// 获取首页
func getIndexFile(c *gin.Context, fileSystem http.FileSystem) {
	f, err := fileSystem.Open("index.html")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	var indexContent []byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		indexContent = append(indexContent, []byte(line)...)
		indexContent = append(indexContent, '\n')
	}
	if err := scanner.Err(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
}

// 挂载公共路由
func addPublicRoute(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		router.HandleContext(ctx)
	})
	control := CreateController()
	setting := CreateSettingServer()
	api := router.Group("/api")
	{
		api.GET("/setting", setting.Get)
		api.POST("/setting/pwd", setting.SetPassword)
		api.GET("/register", control.Register)
		api.GET("/conn/state", control.GetStatus)
		api.GET("/conn/open", control.Connect)
		api.GET("/conn/close", control.Disconnect)
	}
}

// 启动路由
func (r Router) Run() {
	util.OutLogf(MODEL_NAME, "starting from port "+r.port)
	// 启动服务
	go func() {
		var err error
		if r.cert == "" {
			err = r.Object.Run(":" + r.port)
		} else {
			err = r.Object.RunTLS(":"+r.port, r.cert, r.key)
		}
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				util.OutLogf(MODEL_NAME, "please add '--port=' after start command to change the port")
			}
			util.OutErr(MODEL_NAME, "router startup failed: %v", err)
		}
	}()

	util.OutLog(MODEL_NAME, "router is ready")
	// 等待中断信号
	r.wait()
}

// 等待关闭
func (r Router) wait() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	util.OutLog(MODEL_NAME, "router is offline")

	util.OutLog(MODEL_NAME, "server is stopped")
}
