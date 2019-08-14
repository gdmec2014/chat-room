package models

import (
	"chat-room/go/helper"
	"time"
)

type User struct {
	Id         int64     `xorm:"not null pk autoincr INT(64)" json:"uid,omitempty"`
	Name       string    `xorm:"not null comment('名字') unique VARCHAR(255)" json:"name,omitempty"`
	WxId       string    `xorm:"not null comment('微信ID') unique VARCHAR(255)" json:"wx_id,omitempty"`
	LastLogin  time.Time `xorm:"comment('最后登录时间') not null DATETIME" json:"last_login,omitempty"`
	DeleteTime time.Time `xorm:"comment('删除时间') DATETIME" json:"delete_time,omitempty"`
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