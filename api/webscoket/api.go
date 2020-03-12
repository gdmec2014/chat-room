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

	if len(room.Member) < YouPerformIGuess.MaxNumber {
		helper.Debug("人數不夠~不能開始")
		this.SetReturnData(helper.FAILED, "人數不夠~不能開始", nil, false)
		return
	}
	noMasterNum := 0
	for _, m := range room.Member {
		if m.UserType != NO_MASTER {
			//存在出題人，則遊戲沒有結束
			helper.Debug("存在出題人，則遊戲沒有結束2")
			this.SetReturnData(helper.FAILED, "遊戲沒有結束,不能開始新的", nil, false)
			return
		} else {
			noMasterNum++
		}
	}

	if noMasterNum > 0 {
		if noMasterNum != len(room.Member) {
			helper.Debug("存在出題人，則遊戲沒有結束1")
			this.SetReturnData(helper.FAILED, "遊戲沒有結束,不能開始新的", nil, false)
			return
		}
	}

	helper.DebugStructToString(room)

	for i, _ := range room.Member {
		room.Member[i].UserType = VIEWER
	}

	newWS(user, room, EVENT_GAME_RE_START)
	this.SetReturnData(helper.SUCCESS, "成功", nil, false)
}
