package controllers

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"strconv"
	"time"
)

// 登录注册
type AuthController struct {
	BaseController
}

// @Title 注册
// @Description 用户注册
// @Param body body models.User "用户注册"
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn 注册失败
// @Failure helper.EXIST_FAILED {object} helper.RestfulReturn 用户名已存在
// @Failure helper.REPASSWORD_FAIELD {object} helper.RestfulReturn 两次密码不一致
// @router /register [post]
func (this *AuthController) Register() {
	user := models.User{}
	this.GetPostDataNotStop(&user)
	this.NeedPostData(user.Name, user.Password, user.RePassword)
	if user.Password != user.RePassword {
		this.SetReturnData(helper.REPASSWORD_FAIELD, "两次密码不一致", nil)
	}
	has, _, err := models.GetUserByName(user.Name)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
	}
	if has {
		this.SetReturnData(helper.EXIST_FAILED, "用户名已存在", nil)
	}
	user.Token = helper.Md5(time.Now().String() + user.Name)
	user.LastLogin = time.Now()
	err = models.Insert(&user)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
	}

	user.Password = helper.Md5(user.Password + models.SaltCode + strconv.Itoa(int(user.Id)))

	var whereData []interface{}

	err = models.Update(&user, "", whereData, "password")
	if helper.Error(err) {
		go models.Delete(&user, "", whereData)
		this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
	}
	user.Password = ""
	this.SetReturnData(helper.SUCCESS, "注册成功", user)
}

// @Title 登录
// @Description 用户登录
// @Param body body models.User "用户登录参数"
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn 登录失败
// @Failure helper.NOT_EXIST_FAILED {object} helper.RestfulReturn 用户不存在
// @Failure helper.PASSWORD_ERROR {object} helper.RestfulReturn 密码错误
// @router /login [post]
func (this *AuthController) Login() {
	postData := models.User{}
	this.GetPostDataNotStop(&postData)
	this.NeedPostData(postData.Name, postData.Password)
	has, user, err := models.GetUserByName(postData.Name)
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "登录失败", err.Error())
	}

	if !has {
		this.SetReturnData(helper.NOT_EXIST_FAILED, "用户不存在", nil)
	}

	md5 := helper.Md5(postData.Password + models.SaltCode + strconv.Itoa(int(user.Id)))

	if md5 != user.Password {
		this.SetReturnData(helper.PASSWORD_ERROR, "密码错误", nil)
	}

	user.Token = helper.Md5(time.Now().String() + user.Name)

	var whereData []interface{}

	err = models.Update(&user, "", whereData, "token")
	if helper.Error(err) {
		this.SetReturnData(helper.SQL_ERROR, "登录失败", err.Error())
	}
	user.Password = ""
	this.SetReturnData(helper.SUCCESS, "登录成功", user)
}
