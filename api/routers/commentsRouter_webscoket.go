package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/create`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "GetRoomMember",
			Router: `/get_room_member`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "Join",
			Router: `/join`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
