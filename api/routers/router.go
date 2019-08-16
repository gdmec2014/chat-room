// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"chat-room/api/controllers"
	"chat-room/api/webscoket"

	"github.com/astaxie/beego"
)

func init() {
	// API
	beego.AutoRouter(&controllers.ApiController{})

	// 二级api
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/room",
			beego.NSInclude(
				&webscoket.WebSocketController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
