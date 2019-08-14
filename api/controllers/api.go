package controllers

import (
	"chat-room/go/helper"
	"chat-room/go/models"
	"time"
)

// Operations about Api
type ApiController struct {
	BaseController
}

func (this *ApiController) GetRunTime() {
	this.SetReturnData(helper.SUCCESS, "love you", time.Now().Unix() - models.App.RunTime)
}
