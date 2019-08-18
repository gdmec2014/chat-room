package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
)

//TODO webscoket api

// @Title 握手
// @Description 握手
// @Param body body Room "啊啊啊"
// @Success helper.SUCCESS {object} helper.RestfulReturn "加入成功"
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn "加入失败"
// @router /hand [get]
func (this *WebSocketController) Hand() {
	if models.DBOk {
		helper.Debug("数据库还没连接呢！")
		this.join()
	}
	this.SetReturnData(helper.SUCCESS, "握手成功", nil, false)
}

// @Title 获取房间成员
// @Description 获取房间成员
// @Param body body Room "房间ID"
// @Success helper.SUCCESS {object} helper.RestfulReturn
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn
// @router /get_room_member [get]
func (this *WebSocketController) GetRoomMember() {
	room := getRoom(this.GetString("room_id"))
	this.SetReturnData(helper.SUCCESS, "获取成功", getMemberByRoom(room), false)
}

// @Title 获取全部房间
// @Description 获取全部房间
// @Success helper.SUCCESS {object} helper.RestfulReturn
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn
// @router /get_all_room [get]
func (this *WebSocketController) GetAllRoom() {
	rooms := getAllRooms()
	if len(rooms) < 1 {
		this.SetReturnData(helper.SUCCESS, "还没有记录喔", []string{}, false)
	} else {
		this.SetReturnData(helper.SUCCESS, "获取成功", rooms, false)
	}
}
