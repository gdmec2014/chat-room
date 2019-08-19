package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"container/list"
	"encoding/json"
	"time"
)

//TODO webscoket 数据处理

type EventType int

const (
	EVENT_HAND    = 10 //握手事件
	EVENT_CREATE  = 11 //创房事件
	EVENT_JOIN    = 12 //加房事件
	EVENT_LEAVE   = 13 //离线事件
	EVENT_MESSAGE = 14 //消息事件
	EVENT_INVAILD = 15 //无效事件
)

type UserType int

const (
	VIEWER = 100
	PLAYER = 101
)

type youPerformIGuess struct {
	MaxNumber int `json:"max_number"`
}

var (
	allUser                            = list.New()            // 全部用户和链接
	allRooms                           = list.New()            // 全部房间
	YouPerformIGuess *youPerformIGuess = new(youPerformIGuess) // 配置
)

//初始化
func init() {
	YouPerformIGuess.MaxNumber = models.GetAppConfInt("youPerformIGuess::maxNumber")
	if YouPerformIGuess.MaxNumber < 1 {
		YouPerformIGuess.MaxNumber = 6
	}
}

func addRooms(room Room) {
	go updateRedisRooms(room)
	allRooms.PushBack(room)
}

func addUser(user models.User) {
	allUser.PushBack(user)
}

//更新房间成员
func updateRoomsMember(roomId, roomName string, member Member) Room {

	go updateRedisRoomsMember(roomId, roomName, member)

	newMember := make([]Member, 0)
	newMember = append(newMember, member)

	newRoom := Room{}

	has := false

	for r := allRooms.Front(); r != nil; r = r.Next() {
		helper.Debug(r.Value.(Room))
		if r.Value.(Room).Id == roomId {
			helper.Debug("存在房间，更新成员")
			for _, m := range r.Value.(Room).Member {
				if m.UserId != member.UserId {
					has = true
					newMember = append(newMember, m)
				}
			}
			newRoom = Room{
				Id:     r.Value.(Room).Id,
				Name:   r.Value.(Room).Name,
				Member: newMember}
			allRooms.Remove(r)
			addRooms(newRoom)
		}
	}

	//不存在房间，就加进去
	if !has {
		helper.Debug("不存在房间，新加")
		newRoom = Room{
			Id:     roomId,
			Name:   roomName,
			Member: newMember}
		addRooms(newRoom)
	}

	return newRoom
}

//更新用户连接
func updateUserConn(user models.User) {
	has := false
	for r := allUser.Front(); r != nil; r = r.Next() {
		if r.Value.(models.User).Id == user.Id {
			//关闭旧的连接
			//r.Value.(models.User).Conn.Close()
			helper.Debug("存在用户，更新")
			has = true
			newUser := r.Value.(models.User)
			newUser.Conn = user.Conn
			allUser.Remove(r)
			addUser(newUser)
			break
		}
	}
	//不存在用户，就加进去
	if !has {
		helper.Debug("不存在用户，新加")
		addUser(user)
	}
}

//根据房间ID获取房间
func getRoom(roomId string) (room Room) {
	for r := allRooms.Front(); r != nil; r = r.Next() {
		if r.Value.(Room).Id == roomId {
			room = r.Value.(Room)
			break
		}
	}
	return
}

//获取全部房间
func getAllRooms() (rooms []Room) {
	for r := allRooms.Front(); r != nil; r = r.Next() {
		rooms = append(rooms, r.Value.(Room))
	}
	return
}

//获取房间的用户
func getMemberByRoom(room Room) (user []models.User) {

	for r := allUser.Front(); r != nil; r = r.Next() {
		for _, m := range room.Member {
			if r.Value.(models.User).Id == m.UserId {
				if !helper.Contain(user, r.Value.(models.User)) {
					user = append(user, r.Value.(models.User))
				}
			}
		}
	}

	return
}

//获取全部的用户
func getAllMember() (user []models.User) {

	for r := allUser.Front(); r != nil; r = r.Next() {
		if !helper.Contain(user, r.Value.(models.User)) {
			user = append(user, r.Value.(models.User))
		}
	}

	return
}

//检查用户是否存在
func hasMember(id int64) (has bool, user models.User) {

	has = false

	for r := allUser.Front(); r != nil; r = r.Next() {
		if r.Value.(models.User).Id == id {
			has = true
			user = r.Value.(models.User)
			break
		}
	}

	return
}

//use redis
func updateRedisRooms(room Room) {
	helper.Debug("updateRedisRooms --- ", room)
	//保存在总的房间
	if !IsInSet("room", room.Id) {
		helper.Debug("设置 哦哦哦哦哦")
		SetSAdd("room", room.Id)
	}
	var oldRoom Room
	GetSet("room", &oldRoom)
	//房间信息
	roomData := make(map[string]interface{})
	roomData["Name"] = room.Name
	roomData["TimeUnix"] = time.Now().Unix()
	SetMap(room.Id, roomData)
	//成员
	for _, rm := range room.Member {
		//保证唯一
		bm, err := json.Marshal(&rm)
		if helper.Error(err) {
			continue
		}
		if IsInSet(room.Id+"_member", string(bm)) {
			DelSet(room.Id+"_member", string(bm))
		}
		SetSAdd(room.Id+"_member", string(bm))
	}
}

//更新房间成员
func updateRedisRoomsMember(roomId, roomName string, member Member) (room Room) {
	if IsInSet(roomId+"_member", member) {
		DelSet(room.Id+"_member", member)
	}
	SetSAdd(roomId+"_member", member)
	var oldMember Member
	GetSet(roomId+"_member", &oldMember)
	return
}
