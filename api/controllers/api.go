package controllers

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"time"
)

// Operations about Api
type ApiController struct {
	BaseController
}

func (this *ApiController) GetRunTime() {
	this.SetReturnData(helper.SUCCESS, "love you", time.Now().Unix() - models.App.RunTime)
}

func (this *ApiController) GetUserByToken() {
	user := this.CheckLogin()
	this.SetReturnData(helper.SUCCESS, "love you", user)
}

func (this *ApiController) GetUserByWXId() {
	user := this.CheckLogin()
	this.SetReturnData(helper.SUCCESS, "love you", user)
}