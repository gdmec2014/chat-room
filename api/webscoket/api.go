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
		this.CheckLogin()
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

// @Title 重新開始遊戲
// @Description 重新開始遊戲
// @Success helper.SUCCESS {object} helper.RestfulReturn
// @Failure helper.SQL_ERROR {object} helper.RestfulReturn
// @router /restart [get]
func (this *WebSocketController) ReStart() {
	user := this.CheckLogin()
	room, has := getRoomByMember(user.Id)

	if !has {
		this.SetReturnData(helper.SUCCESS, "还没有房間记录喔~不能開始遊戲", nil, false)
		return
	}

	//超过两个人就可以开始游戏
	if len(room.Member) < 2 {
		helper.Debug("人數不夠~不能開始")
		this.SetReturnData(helper.FAILED, "人數不夠~不能開始", nil, false)
		return
	}

	if isGameStart(room) {
		this.SetReturnData(helper.FAILED, "游戏已经开始了喔~", nil, false)
		return
	}

	//helper.DebugStructToString(room)

	for i, _ := range room.Member {
		room.Member[i].UserType = VIEWER
	}

	room.Times = 0
	room.Mark = make(map[int][]Mark, 0)
	updateRooms(room)
	newWS(user, room, EVENT_GAME_RE_START)
	this.SetReturnData(helper.SUCCESS, "成功", nil, false)
}
