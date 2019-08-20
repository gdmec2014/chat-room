package controllers

import (
	"chat-room/api/models"
)

// Operations about Api
type NeedLoginController struct {
	User models.User
	BaseController
}

func (this *NeedLoginController)Prepare()  {
	this.User = this.checkLogin()
}