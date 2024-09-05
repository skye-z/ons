package core

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/cloud-server/model"
	"github.com/skye-z/cloud-server/util"
	"xorm.io/xorm"
)

type UserService struct {
	Data *model.UserModel
}

func CreateUserService(engine *xorm.Engine) *UserService {
	data := &model.UserModel{
		DB: engine,
	}
	return &UserService{
		Data: data,
	}
}

func (us UserService) GetLoginUser(ctx *gin.Context) {
	uid := int64(ctx.GetInt("uid"))
	info, err := us.Data.GetUser(uid)
	if err == nil {
		util.ReturnData(ctx, true, info)
	} else {
		util.ReturnMessage(ctx, false, "获取用户信息失败")
	}
}
