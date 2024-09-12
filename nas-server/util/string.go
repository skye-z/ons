/*
字符串工具

BetaX Blog
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateRandomString(length int) string {
	return GenerateRandom("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", length)
}

func GenerateRandomNumber(length int) string {
	return GenerateRandom("0123456789", length)
}

func GenerateRandom(charset string, length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func CalculateMD5(input string) string {
	data := []byte(input)
	hash := md5.Sum(data)
	hashInHex := hex.EncodeToString(hash[:])
	return hashInHex
}

func GetPostBool(ctx *gin.Context, name string, defaultValue bool) bool {
	cache := ctx.PostForm(name)
	if cache == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(cache)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetPostInt(ctx *gin.Context, name string, defaultValue int) int {
	cache := ctx.PostForm(name)
	if cache == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(cache)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetPostInt64(ctx *gin.Context, name string, defaultValue int64) int64 {
	cache := ctx.PostForm(name)
	if cache == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(cache, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}
