/*
配置工具

BetaX Blog
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/spf13/viper"
)

const Version = "0.2.0"

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefault()
		} else {
			OutErr("Config", "read failed: %v", err)
		}
	}
}

func Reload() {
	viper.ReadInConfig()
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
	viper.WriteConfig()
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func createDefault() {
	// 安装状态
	viper.SetDefault("basic.install", "0")
	// SSL配置
	viper.SetDefault("basic.sslCert", "")
	viper.SetDefault("basic.sslKey", "")
	// 是否开放注册
	viper.SetDefault("connect.auto", "false")
	viper.SetDefault("connect.server", "192.168.1.160:9891")
	viper.SetDefault("connect.natId", "")
	viper.SetDefault("connect.password", "")
	// 令牌密钥
	secret, err := generateSecret()
	if err != nil {
		panic(err)
	}
	viper.SetDefault("token.secret", secret)
	// 令牌有效期/小时
	viper.SetDefault("token.exp", 24)
	viper.SafeWriteConfig()
}

func generateSecret() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
