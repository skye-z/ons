package core

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/cloud-server/model"
	"github.com/skye-z/cloud-server/util"
	"xorm.io/xorm"
)

type DeviceService struct {
	Data *model.DeviceModel
}

func CreateDeviceService(engine *xorm.Engine) *DeviceService {
	data := &model.DeviceModel{
		DB: engine,
	}
	return &DeviceService{
		Data: data,
	}
}

// 注册 NAS 设备
func (ds DeviceService) Register(ctx *gin.Context) {
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	name := ctx.PostForm("name")
	if len(name) == 0 {
		util.ReturnMessage(ctx, false, "设备名称不能为空")
		return
	}
	list, err := ds.Data.GetDeviceList(uid, 1, 5)
	if err == nil && len(list) >= 3 {
		util.ReturnMessage(ctx, false, "当前账户注册设备已达上限")
		return
	} else if err != nil {
		util.ReturnMessage(ctx, false, "设备服务异常")
		return
	}
	nat := util.GenerateRandomNumber(6)
	if ds.Data.CheckNATId(name) {
		nat = util.GenerateRandomNumber(6)
		if ds.Data.CheckNATId(name) {
			nat = util.GenerateRandomNumber(6)
		}
	}

	if ds.Data.AddDevice(uid, name, nat) {
		util.ReturnMessage(ctx, true, nat)
	} else {
		util.ReturnMessage(ctx, false, "设备注册失败")
	}
}

// 重命名设备
func (ds DeviceService) ReName(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	name := ctx.PostForm("name")
	if len(name) == 0 {
		util.ReturnMessage(ctx, false, "设备名称不能为空")
		return
	}
	info, err := ds.Data.GetDevice(id)
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	if err != nil {
		util.ReturnError(ctx, util.Errors.UnexpectedError)
		return
	} else if info == nil {
		util.ReturnMessage(ctx, false, "设备不存在")
		return
	} else if uid != 1 && uid != info.UId {
		util.ReturnMessage(ctx, false, "非法访问")
		return
	}
	if ds.Data.UpdateName(id, name) {
		util.ReturnMessage(ctx, false, "重命名失败")
	} else {
		util.ReturnMessage(ctx, true, "重命名成功")
	}
}

// 获取设备列表
func (ds DeviceService) GetList(ctx *gin.Context) {
	page := "1"
	num := "5"
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	if uid == 1 {
		if ctx.Query("page") != "" {
			page = ctx.Query("page")
		}
		if ctx.Query("number") != "" {
			num = ctx.Query("number")
		}
	}
	iPage, err1 := strconv.Atoi(page)
	iNum, err2 := strconv.Atoi(num)
	if err1 != nil || err2 != nil {
		util.ReturnError(ctx, util.Errors.ParamIllegalError)
		return
	}
	list, err := ds.Data.GetDeviceList(uid, iPage, iNum)
	if err != nil {
		log.Println(err)
		util.ReturnError(ctx, util.Errors.UnexpectedError)
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 获取设备信息
func (ds DeviceService) GetInfo(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	info, err := ds.Data.GetDevice(id)
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	if err != nil {
		util.ReturnError(ctx, util.Errors.UnexpectedError)
	} else if uid != 1 && uid != info.UId {
		util.ReturnMessage(ctx, false, "非法访问")
	} else {
		util.ReturnData(ctx, true, info)
	}
}

// 删除设备
func (ds DeviceService) Del(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	info, err := ds.Data.GetDevice(id)
	uid, _ := strconv.ParseInt(ctx.GetString("uid"), 10, 64)
	if err != nil {
		util.ReturnError(ctx, util.Errors.UnexpectedError)
		return
	} else if info == nil {
		util.ReturnMessage(ctx, false, "设备不存在")
		return
	} else if uid != 1 && uid != info.UId {
		util.ReturnMessage(ctx, false, "非法访问")
		return
	}
	if ds.Data.DelDevice(id) {
		util.ReturnData(ctx, true, "删除成功")
	} else {
		util.ReturnData(ctx, false, "删除失败")
	}
}
