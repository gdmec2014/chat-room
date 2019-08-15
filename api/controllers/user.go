package controllers

import (
	"chat-room/api/helper"
	"chat-room/api/models"
)

// 用户类
type UserController struct {
	NeedLoginController
}

// @Title 修改头像
// @Description 修改头像
// @Param body body models.User "头像地址"
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn 修改失败
// @router /update_avatar [post]
func (this *UserController) UpdateAvatar() {
	user := models.User{}
	this.GetPostDataNotStop(&user)
	this.NeedPostData(user.Avatar)

	var whereData []interface{}

	this.User.Avatar = user.Avatar

	err := models.Update(&this.User, "", whereData, "avatar")
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "修改失败", err.Error())
	}
	this.SetReturnData(helper.SUCCESS, "修改成功", user)
}
