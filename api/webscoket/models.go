package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"container/list"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

//TODO webscoket 数据处理

type EventType int

const (
	EVENT_HAND                   = 10 //握手事件
	EVENT_CREATE                 = 11 //创房事件
	EVENT_JOIN                   = 12 //加房事件
	EVENT_LEAVE                  = 13 //离线事件
	EVENT_MESSAGE                = 14 //消息事件
	EVENT_INVAILD                = 15 //无效事件
	EVENT_DRAW                   = 16 //绘图事件
	EVENT_BREAK_DRAW             = 17 //中断绘画事件
	EVENT_GIVE_IDENTITY          = 18 //转换游戏身份事件
	EVENT_NO_PLACE               = 19 //房间满人不能加入事件
	EVENT_GAME_NO_START          = 20 //还不能开始游戏事件
	EVENT_GAME_CAN_START         = 21 //开始游戏事件
	EVENT_GAME_ANSWER            = 22 //回答问题事件
	EVENT_GAME_BONUS             = 23 //答对问题加分事件
	EVENT_GAME_TIME              = 24 //游戏计时事件
	EVENT_GAME_OVER              = 25 //游戏结束事件
	EVENT_GAME_IS_START          = 26 //遊戲正在進行，不能重複開始
	EVENT_GAME_MEMBER_NOT_ENOUGH = 27 //人數不夠，不能開始遊戲
	EVENT_GAME_RE_START          = 28 //重新開始遊戲
	EVENT_NEW_DRAW               = 29 //新的绘图事件
	EVENT_SYSTEM_MESSAGE         = 30 //系统消息
)

type UserType int

//刚刚进来就是观众身份，6个人到齐后就赋予身份
const (
	VIEWER    = 100 //观众
	PLAYER    = 101 //猜一方
	MASTER    = 102 //出题一方
	NO_MASTER = 103 //不可以再赋予出题一方，当全部玩家身份为103时候，游戏结束，结算分数
)

type youPerformIGuess struct {
	MaxNumber int `json:"max_number"`
	TimeOver  int `json:"time_over"`
	Times     int `json:"times"`
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
	YouPerformIGuess.TimeOver = models.GetAppConfInt("youPerformIGuess::timeOver")
	if YouPerformIGuess.TimeOver < 1 {
		YouPerformIGuess.TimeOver = 30
	}
	YouPerformIGuess.Times = models.GetAppConfInt("youPerformIGuess::times")
	if YouPerformIGuess.Times < 0 {
		YouPerformIGuess.Times = 0
	}
	beego.Info(YouPerformIGuess)
}

func InitData() {
	var redisRoomsId []string
	GetSet("room", &redisRoomsId)
	//helper.Debug("InitData redisRoomsId", redisRoomsId)
	for i, _ := range redisRoomsId {
		go func(i int) {
			var room Room
			redisRoomId := redisRoomsId[i]
			redisRoomId = strings.Replace(redisRoomId, "\"", "", -1)
			//helper.Debug("InitData redisRoomId", redisRoomId)
			if len(redisRoomId) < 1 {
				DelSet("room", redisRoomId)
				return
			}
			redisRoomMap := GetMap(redisRoomId)

			//helper.Debug("InitData redisRoomMap", redisRoomMap)

			//helper.Debug(redisRoomMap)
			_, ok := redisRoomMap["Name"]
			if ok {
				if len(redisRoomMap["Name"]) < 1 {
					DelSet("room", redisRoomId)
					return
				}
				room.Name = redisRoomMap["Name"]
			} else {
				DelSet("room", redisRoomId)
				return
			}
			//helper.Debug("InitData redisRoom name", room.Name)
			_, ok = redisRoomMap["TimeUnix"]
			if ok {
				timeUnix, err := strconv.ParseInt(redisRoomMap["TimeUnix"], 10, 64)
				if helper.Error(err) {
					timeUnix = 0
				}
				room.TimeUnix = timeUnix
			} else {
				DelSet("room", redisRoomId)
				return
			}
			//helper.Debug("InitData redisRoom TimeUnix", room.TimeUnix)
			var member []Member
			GetSet(redisRoomId+"_member", &member)
			if len(member) < 1 {
				//没有成员了，删除房间
				DelSet("room", redisRoomId)
				return
			}
			room.Member = member
			room.Id = redisRoomId
			//helper.Debug("InitData redisRoom Member", room.Member)
			allRooms.PushBack(room)
		}(i)
	}
}

func addRooms(room Room) {
	go updateRedisRooms(room)
	allRooms.PushBack(room)
}

func addUser(user models.User) {
	allUser.PushBack(user)
}

//更新分数
func updateRooms(room Room) {
	needUpdate := false
	for r := allRooms.Front(); r != nil; r = r.Next() {
		if r.Value.(Room).Id == room.Id {
			allRooms.Remove(r)
			needUpdate = true
		}
	}
	if needUpdate {
		addRooms(room)
	}

	//b, _ := json.MarshalIndent(room, "", " ")
	//_, files, line, ok := runtime.Caller(1)
	//if !ok {
	//	fmt.Println(fmt.Errorf("Error: Cant not print!"))
	//	return
	//}
	//fs := strings.Split(files, "/")
	//file := ""
	//file = fs[0]
	//if len(fs) > 2 {
	//	file = fs[len(fs)-2] + "/" + fs[len(fs)-1]
	//}
	//fileline := "[" + file + ":" + strconv.Itoa(line) + "]"
	//go beego.Debug(fileline, string(b), "\r\n")
	//
	//helper.DebugStructToString(getRoom(room.Id))
}

//更新房间成员
func updateRoomsMember(room Room, member Member) (newRoom Room, code EventType) {

	//判斷用戶是否在房間
	isInRoom := isMemberInRoom(room.Id, member.UserId)

	newMember := make([]Member, 0)

	if !isInRoom {
		newMember = append(newMember, member)
	}

	hasRoom := false

	for r := allRooms.Front(); r != nil; r = r.Next() {
		//helper.Debug(r.Value.(Room))
		if r.Value.(Room).Id == room.Id {
			hasRoom = true
			newRoom.MaxMember = YouPerformIGuess.MaxNumber
			//helper.Debug("存在房间，更新成员")
			for _, m := range r.Value.(Room).Member {
				if m.UserId != member.UserId {
					if len(r.Value.(Room).Member) < YouPerformIGuess.MaxNumber {
						newMember = append(newMember, m)
					} else {
						//人数已经满了，不可以再加了啦
						newRoom = Room{
							MaxMember: YouPerformIGuess.MaxNumber,
							Id:        r.Value.(Room).Id,
							Name:      r.Value.(Room).Name,
							Member:    r.Value.(Room).Member}
						code = EVENT_NO_PLACE
						return
					}
				}
			}
			newRoom = Room{
				MaxMember: YouPerformIGuess.MaxNumber,
				Id:        r.Value.(Room).Id,
				Name:      r.Value.(Room).Name,
				Member:    newMember}
			allRooms.Remove(r)
			addRooms(newRoom)
			if len(newMember) == YouPerformIGuess.MaxNumber {
				//helper.Debug("人数已经齐了，开打")
				code = EVENT_GAME_CAN_START
			} else {
				//人还没满呢
				code = EVENT_GAME_NO_START
			}
		}
	}

	//不存在房间，就加进去
	if !hasRoom {
		//helper.Debug("不存在房间，新加")
		newRoom = Room{
			MaxMember: YouPerformIGuess.MaxNumber,
			Id:        room.Id,
			Name:      room.Name,
			Member:    newMember}
		addRooms(newRoom)
	}

	go updateRedisRoomsMember(room.Id, room.Name, member)

	return
}

//更新用户连接
func updateUserConn(user models.User) {
	has := false
	for r := allUser.Front(); r != nil; r = r.Next() {
		if r.Value.(models.User).Id == user.Id {
			//关闭旧的连接
			//r.Value.(models.User).Conn.Close()
			//helper.Debug("存在用户，更新")
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
		//helper.Debug("不存在用户，新加")
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

//根據用戶獲取房間
func getRoomByMember(uid int64) (room Room, has bool) {
	rooms := getAllRooms()
	if len(rooms) < 1 {
		has = false
		return
	}

	for _, room = range rooms {
		for _, m := range room.Member {
			if uid == m.UserId {
				has = true
				return
			}
		}
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

//判斷用戶是否在房間
func isMemberInRoom(roomId string, uid int64) (has bool) {
	has = false
	room := getRoom(roomId)
	for _, m := range room.Member {
		if uid == m.UserId {
			has = true
			return
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
	//helper.Debug("updateRedisRooms --- ", room)
	//保存在总的房间
	if !IsInSet("room", room.Id) {
		//helper.Debug("updateRedisRooms set room id", room.Id)
		SetSAdd("room", room.Id)
	}

	if helper.IsDebug {
		var rs []string     //测试
		GetSet("room", &rs) //测试
		//helper.Debug("updateRedisRooms get rooms", rs) //测试
	}

	//房间信息
	roomData := make(map[string]interface{})
	roomData["Name"] = room.Name
	roomData["TimeUnix"] = time.Now().Unix()
	//helper.Debug("updateRedisRooms set room map ", roomData)

	SetMap(room.Id, roomData)

	if helper.IsDebug {
		GetMap(room.Id) //测试
	}
	//成员
	for _, rm := range room.Member {
		//保证唯一
		if IsInSet(room.Id+"_member", rm) {
			DelSet(room.Id+"_member", rm)
		}
		SetSAdd(room.Id+"_member", rm)
	}

	if helper.IsDebug {
		var member []Member                //测试
		GetSet(room.Id+"_member", &member) //测试
		//helper.Debug("updateRedisRooms get room  member :", member) //测试
	}
}

//更新房间成员
func updateRedisRoomsMember(roomId, roomName string, member Member) (room Room) {
	if IsInSet(roomId+"_member", member) {
		DelSet(room.Id+"_member", member)
	}
	SetSAdd(roomId+"_member", member)
	return
}

//判断游戏是不是在进行中
func isGameStart(room Room) bool {

	room = getRoom(room.Id)
	noMasterNum := 0

	for _, m := range room.Member {
		if m.UserType == PLAYER || m.UserType == MASTER {
			//存在出題人，則遊戲沒有結束
			//helper.Debug("存在出題人，則遊戲沒有結束2")
			return true
		}
		if m.UserType == NO_MASTER {
			noMasterNum++
		}
	}

	if noMasterNum > 0 {
		if noMasterNum != len(room.Member) {
			//helper.Debug("存在出題人，則遊戲沒有結束1")
			return true
		}
	}

	return false
}
