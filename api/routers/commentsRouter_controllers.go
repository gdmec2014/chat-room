package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["chat-room/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["chat-room/api/controllers:AuthController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chat-room/api/controllers:AuthController"] = append(beego.GlobalControllerRouter["chat-room/api/controllers:AuthController"],
        beego.ControllerComments{
            Method: "Register",
            Router: `/register`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chat-room/api/controllers:UserController"] = append(beego.GlobalControllerRouter["chat-room/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "UpdateAvatar",
            Router: `/update_avatar`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
