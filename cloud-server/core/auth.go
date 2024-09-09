package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/skye-z/ons/cloud-server/model"
	"github.com/skye-z/ons/cloud-server/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"xorm.io/xorm"
)

const IssuerName = "BetaX ONS"
const tokenKey = "token.secret"

type AuthService struct {
	Config *oauth2.Config
	UM     *model.UserModel
}

// 创建鉴权服务
func NewAuthService(engine *xorm.Engine) *AuthService {
	as := new(AuthService)
	as.Config = BuildAuthConfig()
	as.UM = &model.UserModel{
		DB: engine,
	}
	return as
}

func BuildAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     util.GetString("github.clientId"),
		ClientSecret: util.GetString("github.clientSecret"),
		RedirectURL:  util.GetString("github.redirectUrl"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// 发起 OAuth2 登陆
func (as AuthService) Login(c *gin.Context) {
	if as.Config.ClientID == "" {
		util.Reload()
		as.Config = BuildAuthConfig()
	}
	url := as.Config.AuthCodeURL(util.GenerateRandomString(6))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 处理 OAuth2 回调
func (as AuthService) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	res, err := as.Config.Exchange(ctx, code)
	if err != nil {
		// 授权服务不可用
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=1")
		return
	}
	// 换取授权信息
	token := res.AccessToken
	info, err := as.QueryUserInfo(token)
	if err != nil || info == nil {
		// 授权信息无效
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=2")
		return
	}
	// 查询用户是否存在
	user, err := as.UM.GetOAuthUser(info.Id)
	if err != nil {
		// 账户查询错误
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=3")
		return
	} else if user == nil {
		if util.GetBool("basic.register") {
			user = &model.User{
				GId:      info.Id,
				Avatar:   info.Avatar,
				Nickname: info.Nickname,
				Username: info.Username,
				Email:    info.Email,
			}
			if !as.UM.AddUser(user) {
				// 自动注册失败
				ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=4")
				return
			}
		} else {
			// 账户不存在
			ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=5")
			return
		}
	}
	// 签发令牌
	token, exp, err := GenerateToken(user.Id)
	if err != nil {
		// 令牌签发失败
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/auth?state=6")
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/app/auth?state=9&code=%s&exp=%v", token, exp))
}

type User struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Version  string `json:"version"`
	Email    string `json:"email"`
}

// 查询 OAuth2 用户信息
func (as AuthService) QueryUserInfo(token string) (*User, error) {
	// 获取用户信息
	result, err := as.getUserInfo(token)
	if err != nil {
		return nil, err
	}

	user := User{}
	for key, value := range result {
		if key == "id" {
			if strVal, ok := value.(string); ok {
				user.Id = strVal
			} else if strVal, ok := value.(int16); ok {
				user.Id = strconv.FormatInt(int64(strVal), 10)
			} else if strVal, ok := value.(float64); ok {
				user.Id = strconv.FormatFloat(float64(strVal), 'f', -1, 64)
			}
		}
		if strVal, ok := value.(string); ok {
			if key == "login" {
				user.Username = strVal
			} else if key == "name" {
				user.Nickname = strVal
			} else if key == "avatar_url" {
				user.Avatar = strVal
			} else if key == "email" {
				user.Email = strVal
			}
		}
	}
	return &user, err
}

// 获取 OAuth2 用户信息
func (as AuthService) getUserInfo(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

// 生成令牌
func GenerateToken(id int64) (string, int64, error) {
	secret := util.GetString(tokenKey)
	expTime := util.GetInt("token.exp")
	exp := time.Now().Add(time.Hour * time.Duration(expTime)).Unix()
	tc := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": exp,
			"iss": IssuerName,
			"sub": fmt.Sprintf("%v", id),
		},
	)
	key, _ := base64.StdEncoding.DecodeString(secret)
	token, err := tc.SignedString(key)
	return token, exp, err
}

func AuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Request.Header.Get("Authorization")
		if code == "" {
			util.ReturnError(ctx, util.Errors.NotLoginError)
			return
		}
		if strings.Contains(code, " ") {
			code = code[strings.Index(code, " ")+1:]
		}
		info := jwt.MapClaims{}
		secret := util.GetString(tokenKey)
		token, err := jwt.ParseWithClaims(code, &info, func(token *jwt.Token) (interface{}, error) {
			key, err := base64.StdEncoding.DecodeString(secret)
			return key, err
		})
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		if !token.Valid {
			util.ReturnError(ctx, util.Errors.TokenInvalidError)
			return
		}
		iss, err := info.GetIssuer()
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		if iss != IssuerName {
			util.ReturnError(ctx, util.Errors.TokenIllegalError)
			return
		}
		sub, err := info.GetSubject()
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		ctx.Set("uid", sub)
	}
}
