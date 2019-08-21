package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "GetAllRoom",
			Router: `/get_all_room`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "GetRoomMember",
			Router: `/get_room_member`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"] = append(beego.GlobalControllerRouter["chat-room/api/webscoket:WebSocketController"],
		beego.ControllerComments{
			Method: "Hand",
			Router: `/hand`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
