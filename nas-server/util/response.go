/*
HTTP请求工具

BetaX Blog
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func ReturnError(ctx *gin.Context, err CustomError) {
	ctx.JSON(200, err)
	ctx.Abort()
}

type commonResponse struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    int64  `json:"time"`
}

func ReturnMessage(ctx *gin.Context, state bool, message string) {
	ctx.JSON(200, commonResponse{
		State:   state,
		Message: message,
		Time:    time.Now().Unix() * 1000,
	})
	ctx.Abort()
}

func ReturnData(ctx *gin.Context, state bool, obj any) {
	ctx.JSON(200, commonResponse{
		State: state,
		Data:  obj,
		Time:  time.Now().Unix() * 1000,
	})
	ctx.Abort()
}

func ReturnMessageData(ctx *gin.Context, state bool, message string, obj any) {
	ctx.JSON(200, commonResponse{
		State:   state,
		Message: message,
		Data:    obj,
		Time:    time.Now().Unix() * 1000,
	})
	ctx.Abort()
}

func ToReturnMessage(data []byte) *commonResponse {
	var msg commonResponse
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil
	}
	return &msg
}
