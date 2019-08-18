package controllers

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"time"
)

// Operations about Api
type NeedLoginController struct {
	User models.User
	BaseController
}

func (this *NeedLoginController)Prepare()  {
	this.User = this.checkLogin()
}

func (this *NeedLoginController) checkLogin() (user models.User) {
	Token := this.Ctx.Input.Header("Authorization")
	if len(Token) < 1 {
		Token = this.GetString("token")
	}
	helper.Debug("Token :", Token)
	if len(Token) != 32 {
		this.SetReturnData(helper.TOKEN_ERROR, "no token", nil)
	}
	user = models.User{Token:Token}
	has, user , err := models.GetUserByToken(&user)

	helper.Debug(user)

	if helper.Error(err) {
		this.SetReturnData(helper.FAILED, err.Error(), err)
	}
	if !has {
		this.SetReturnData(helper.TOKEN_ERROR, "无效的token", nil)
	}
	//登陆有效期检测
	if time.Now().Unix()-user.LastLogin.Unix() > 604800 {
		this.SetReturnData(helper.LOGIN_EXPIRACTION, "登陆过期", nil)
	}

	return user
}