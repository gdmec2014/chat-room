package webscoket

import (
	"chat-room/api/helper"
)

//TODO webscoket api

// @Title 加入房间
// @Description 加入房间
// @Param body body Room "房间ID"
// @Success helper.SUCCESS {object} helper.RestfulReturn "加入成功"
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn "加入失败"
// @router /join [post]
func (this *WebSocketController) Join() {

	roomData := Room{}
	this.GetPostDataNotStop(&roomData)
	this.NeedPostData(roomData.Id, roomData.Type)

	this.join(roomData.Id, roomData.Type)

	this.SetReturnData(helper.SUCCESS, "加入成功", nil, false)
}

// @Title 创建房间
// @Description 创建房间
// @Param body body Room "房间名字"
// @Success helper.SUCCESS {object} helper.RestfulReturn
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn
// @router /create [post]
func (this *WebSocketController) Create() {
	roomData := Room{}
	this.GetPostDataNotStop(&roomData)
	this.NeedPostData(roomData.Name)
	roomData.Id = helper.GetRandomString(16)

	this.join(roomData.Id, roomData.Type)

	this.SetReturnData(helper.SUCCESS, "创建成功", roomData, false)
}

// @Title 获取房间成员
// @Description 获取房间成员
// @Param body body Room "房间ID"
// @Success helper.SUCCESS {object} helper.RestfulReturn
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn
// @router /get_room_member [post]
func (this *WebSocketController) GetRoomMember() {

	roomData := Room{}
	this.GetPostDataNotStop(&roomData)
	this.NeedPostData(roomData.Id)

	this.SetReturnData(helper.SUCCESS, "获取成功", getRoomUsers(roomData.Id), false)
}
