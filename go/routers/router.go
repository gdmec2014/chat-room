// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"chat-room/go/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"chat-room/go/controllers"
)

func init() {
	// API
	beego.AutoRouter(&controllers.ApiController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("*",
			//Options用于跨域复杂请求预检
			beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.Debug("models.Domain", models.Domain)
	beego.AddNamespace(ns)


	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}
