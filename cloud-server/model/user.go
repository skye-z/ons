package model

import (
	"xorm.io/xorm"
)

type User struct {
	Id       int64  `json:"id"`       // 用户编号 - 管理员固定为1
	GId      string `json:"gid"`      // Github编号
	Avatar   string `json:"avatar"`   // 头像地址
	Nickname string `json:"nickname"` // 昵称
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱地址
}

type UserModel struct {
	DB *xorm.Engine
}

// 新增用户
func (model UserModel) AddUser(user *User) bool {
	_, err := model.DB.Insert(user)
	return err == nil
}

// 更新用户
func (model UserModel) EditUser(user *User) bool {
	if user.Id == 0 {
		return false
	}
	_, err := model.DB.ID(user.Id).Update(user)
	return err == nil
}

// 删除用户
func (model UserModel) DelUser(user *User) bool {
	if user.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(user)
	return err == nil
}

// 获取用户列表
func (model UserModel) GetUserList(keyword string, page int, num int) ([]User, error) {
	var users []User
	var err error
	if len(keyword) == 0 {
		err = model.DB.Limit(page*num, (page-1)*num).Find(&users)
	} else {
		keyword = "%" + keyword + "%"
		err = model.DB.Where("nickname like ? or username like ? or email like ?", keyword, keyword, keyword).Limit(page*num, (page-1)*num).Find(&users)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 获取授权用户
func (model UserModel) GetOAuthUser(oauthId string) (*User, error) {
	user := &User{
		GId: oauthId,
	}
	has, err := model.DB.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

// 获取用户信息
func (model UserModel) GetUser(id int64) (*User, error) {
	user := &User{
		Id: id,
	}
	has, err := model.DB.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}
