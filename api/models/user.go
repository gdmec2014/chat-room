package models

import (
	"chat-room/api/helper"
	"github.com/gorilla/websocket"
	"time"
)

type User struct {
	Id         int64           `xorm:"not null pk autoincr INT(64)" json:"uid,omitempty"`
	Name       string          `xorm:"not null comment('名字') unique VARCHAR(50)" json:"name,omitempty"`
	WxId       string          `xorm:"not null comment('微信ID') unique VARCHAR(255)" json:"wx_id,omitempty"`
	Avatar     string          `xorm:"null comment('头像') VARCHAR(255)" json:"avatar,omitempty"`
	Password   string          `xorm:"not null comment('密码') VARCHAR(255)" json:"password,omitempty"`
	Token      string          `xorm:"null comment('登录凭证') VARCHAR(255)" json:"token,omitempty"`
	LastLogin  time.Time       `xorm:"comment('最后登录时间') not null DATETIME" json:"last_login,omitempty"`
	DeleteTime time.Time       `xorm:"comment('删除时间') DATETIME" json:"delete_time,omitempty"`
	RePassword string          `xorm:"-" json:"re_password,omitempty"`
	Conn       *websocket.Conn `xorm:"-" json:"conn"` // 用户与后台握手
}

func AddUser(user *User) error {
	//先查是否存在数据，若存在则增加
	has, oldUser, err := GetUserByWxId(user.WxId)
	if helper.Error(err) {
		return err
	}
	if has {
		whereData := make([]interface{}, 0)
		oldUser.LastLogin = time.Now()
		whereData = append(whereData, oldUser.WxId)
		err = Update(oldUser, "wx_id=?", whereData, "last_login")
		user.LastLogin = oldUser.LastLogin
		return err
	}
	return Insert(user)
}

func GetUserByName(name string) (has bool, user User, err error) {
	has, err = db.Table("user").Where("name =?", name).Get(&user)
	return
}

func GetUserByWxId(wxId string) (has bool, user User, err error) {
	has, err = db.Table("user").Where("wx_id =?", wxId).Get(&user)
	return
}

func DeleteUser(user *User) error {
	var whereData []interface{}
	whereData = append(whereData, user.Id)
	return Delete(user, "id=?", whereData)
}

func GetUserByToken(user *User) (has bool, User User, err error) {
	has, err = db.Table("user").Where("token=?", user.Token).Get(&User)
	User.Password = ""
	if helper.Error(err) {
		return
	}
	return
}
