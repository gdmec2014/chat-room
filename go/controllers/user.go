package controllers

import (
	"chat-room/go/helper"
	"chat-room/go/models"
	"strconv"
	"time"
)

// Operations about Api
type UserController struct {
	BaseController
}

// @Title 注册
// @Description 用户注册
// @Param  models.User
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR 注册失败
// @Failure helper.EXIST_FAILED 用户名已存在
// @Failure helper.REPASSWORD_FAIELD 两次密码不一致
// @router /register
func (this *UserController) Register() {
	user := models.User{}
	this.GetPostDataNotStop(&user)
	this.NeedPostData(user.Name, user.Password, user.RePassword)
	if user.Password != user.RePassword {
		this.SetReturnData(helper.REPASSWORD_FAIELD, "两次密码不一致", nil)
		this.Finish()
		return
	}
	has, _, err := models.GetUserByName(user.Name)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
		this.Finish()
		return
	}
	if has {
		this.SetReturnData(helper.EXIST_FAILED, "用户名已存在", nil)
		this.Finish()
		return
	}
	user.Token = helper.Md5(time.Now().String() + user.Name)
	err = models.Insert(&user)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
		this.Finish()
		return
	}

	user.Password = helper.Md5(user.Password + models.SaltCode + strconv.Itoa(int(user.Id)))

	var whereData []interface{}

	err = models.Update(&user, "", whereData, "password")
	if helper.Error(err) {
		go models.Delete(&user, "", whereData)
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
		this.Finish()
		return
	}
	user.Password = ""
	this.SetReturnData(helper.SUCCESS, "注册成功", user)
}

// @Title 登录
// @Description 用户登录
// @Param  models.User
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR 登录失败
// @Failure helper.NOT_EXIST_FAILED 用户不存在
// @Failure helper.PASSWORD_ERROR 密码错误
// @router /login
func (this *UserController) Login() {
	postData := models.User{}
	this.GetPostDataNotStop(&postData)
	this.NeedPostData(postData.Name, postData.Password)
	has, user, err := models.GetUserByName(postData.Name)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "登录失败", err.Error())
		this.Finish()
		return
	}

	if !has {
		this.SetReturnData(helper.NOT_EXIST_FAILED, "用户不存在", nil)
		this.Finish()
		return
	}

	md5 := helper.Md5(postData.Password + models.SaltCode + strconv.Itoa(int(user.Id)))

	if md5 != user.Password {
		this.SetReturnData(helper.PASSWORD_ERROR, "密码错误", nil)
		this.Finish()
		return
	}

	user.Token = helper.Md5(time.Now().String() + user.Name)

	var whereData []interface{}

	err = models.Update(&user, "", whereData, "token")
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "登录失败", err.Error())
		this.Finish()
		return
	}
	user.Password = ""
	this.SetReturnData(helper.SUCCESS, "登录成功", user)
}