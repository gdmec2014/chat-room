package controllers

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/miniprogram"
	"strconv"
	"time"
)

var (
	wx  *wechat.Wechat
	wxa *miniprogram.MiniProgram
)

type WxData struct {
	EncryptedData string     `json:"encryptedData"`
	Iv            string     `json:"iv"`
	Signature     string     `json:"signature"`
	WxUserInfo    WxUserInfo `json:"userInfo"`
	JsCode        string     `json:"code"`
}

type WxUserInfo struct {
	AvatarUrl string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	NickName  string `json:"nickName"`
	Province  string `json:"province"`
}

// 微信类
type WeChatController struct {
	BaseController
}

func init() {
	config := &wechat.Config{
		AppID:     "wx217ad83857cab85a",
		AppSecret: "45a37f928fcf9873dd1729635d4fdb25",
	}
	wx = wechat.NewWechat(config)
	wxa = wx.GetMiniProgram()
}

// @Title 微信 login
// @Description 微信 login
// @Param body body WxData "WxData"
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn 登陆失败
// @router /login [post]
func (this *WeChatController) Login() {
	var wxData WxData
	this.GetPostDataNotStop(&wxData)
	this.NeedPostData(wxData.JsCode)
	r, err := wxa.Code2Session(wxData.JsCode)
	if helper.Error(err) {
		this.SetReturnData(helper.FAILED, "no", err.Error())
	}
	this.SetReturnData(helper.SUCCESS, "ok", r)
}

// @Title 微信注册
// @Description 微信注册
// @Param body body WxData "WxData"
// @Success helper.SUCCESS {object} models.User
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn 登陆失败
// @router /regist [post]
func (this *WeChatController) Regist() {
	var wxData WxData
	this.GetPostDataNotStop(&wxData)
	this.NeedPostData(wxData.Iv, wxData.EncryptedData, wxData.Signature,wxData.JsCode)
	r, err := wxa.Code2Session(wxData.JsCode)
	if helper.Error(err) {
		this.SetReturnData(helper.FAILED, "no", err.Error())
	}
	wxUser, err := wxa.Decrypt(r.SessionKey, wxData.EncryptedData, wxData.Iv)
	if helper.Error(err) {
		this.SetReturnData(helper.FAILED, "no", err.Error())
	}
	//检查数据库是否存在用户
	has, user, err := models.GetUserByWxId(r.OpenID)
	if helper.Error(err) {
		this.SetReturnData(helper.FAILED, "no", err.Error())
	}
	if !has {
		user.Name = helper.GetRandomString(10)
		user.WxId = r.OpenID
		user.Token = helper.Md5(time.Now().String() + user.WxId)
		user.LastLogin = time.Now()
		user.NickName = wxUser.NickName
		user.Avatar = wxUser.AvatarURL

		err = models.Insert(&user)
		if helper.Error(err) {
			this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
		}

		user.Password = helper.Md5(helper.GetString(10) + models.SaltCode + strconv.Itoa(int(user.Id)))

		var whereData []interface{}
		whereData = append(whereData, user.Name)

		err = models.Update(&user, "name=?", whereData, "password")
		if helper.Error(err) {
			go models.Delete(&user, "", whereData)
			this.SetReturnData(helper.SQL_ERROR, "注册失败", err.Error())
		}
		user.Password = ""
	}
	this.SetReturnData(helper.SUCCESS, "ok", user)
}
