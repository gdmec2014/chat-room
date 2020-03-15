package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

//TODO webscoket 房间事件处理

type Room struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Member        []Member       `json:"member"`
	TimeUnix      int64          `json:"time_unix"`      //创建时间
	KeyWord       string         `json:"key_word"`       //当前游戏正确答案
	CorrectNumber int            `json:"correct_number"` //当前轮游戏回答正确人数
	MaxMember     int            `json:"max_member"`     //游戏最大人数
	Times         int            `json:"times"`          //经历了几轮游戏
	Mark          map[int][]Mark `json:"mark"`           //每一轮的分数
	MasterId      int64          `json:"master_id"`      //谁发送的消息
	IsOver        bool           `json:"is_over"`        //是否结束本轮
}

type Mark struct {
	Id        int64 `json:"id"`
	Point     int   `json:"point"`
	HasAnswer bool  `json:"has_answer"`
}

type Member struct {
	UserType UserType `json:"user_type"` // 用户类型跟房间，因为不是每一个房间的身份都一样
	UserId   int64    `json:"user_id"`   // 用户ID
	UserName string   `json:"user_name"` // 用户名
	NickName string   `json:"nick_name"` // 微信用到
	Avatar   string   `json:"avatar"`    // 用户头像
}

type Event struct {
	EventType EventType   `json:"event_type"` // 消息类型
	Room      Room        `json:"room"`       // 房间       //前端發來的字段
	Msg       string      `json:"msg"`        // 消息
	TimeUnix  int64       `json:"time_unix"`  // 消息时间戳
	Data      interface{} `json:"data"`       // 附带数据    //返回後端的字段
	Uid       int64       `json:"uid"`        // 用戶id 房間人滿的時候用到
	Mutex     sync.Mutex  `json:"mutex"`      //锁
}

var (
	// 推送消息的阻塞通道
	publish = make(chan Event, 10)
)

//加入房间
func Join(user models.User, roomId string) {
	room := getRoom(roomId)
	//helper.Debug("Join", room)
	newWS(user, room, EVENT_JOIN)
}

//创建房间
func Create(user models.User, roomId, roomName string) {
	room := Room{
		Id:            roomId,
		Name:          roomName,
		TimeUnix:      time.Now().Unix(),
		KeyWord:       "",
		CorrectNumber: 0,
		MaxMember:     YouPerformIGuess.MaxNumber,
	}
	room.Mark = make(map[int][]Mark, 0)
	newWS(user, room, EVENT_CREATE)
}

//封装消息
func newWS(user models.User, room Room, eventType EventType) {

	var newRoom Room

	newRoom.TimeUnix = time.Now().Unix()
	msg := "加入成功"

	member := Member{
		NickName: user.NickName,
		UserType: VIEWER,
		UserId:   user.Id,
		UserName: user.Name,
		Avatar:   user.Avatar,
	}

	//房间nil是握手的房间
	if room.Id != "nil" {
		//helper.Debug("更新房间成员")
		//更新房间成员
		var code EventType
		if EVENT_GAME_RE_START != eventType {
			newRoom, code = updateRoomsMember(room, member)
		} else {
			//helper.Debug("重新開始遊戲")
			code = EVENT_GAME_CAN_START
			newRoom = room
		}
		msg = user.Name + " 加入了房间 " + newRoom.Name

		if eventType != EVENT_CREATE {
			//helper.Debug("code:", code)
			eventType = code
			switch code {
			case EVENT_GAME_CAN_START:
				//开始游戏
				msg = "准备好了么~要开始了喔~"
				//helper.Debug("准备好了么~要开始了喔~")
				//赋予玩家身份
				for i, _ := range newRoom.Member {
					if i == 0 {
						newRoom.Member[i].UserType = MASTER
					} else {
						newRoom.Member[i].UserType = PLAYER
					}
				}
				//随机获取问题答案
				vocabulary := models.Vocabulary
				lenVocabulary := len(vocabulary)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				key := r.Intn(lenVocabulary)
				newRoom.KeyWord = vocabulary[key]
				newRoom.MasterId = user.Id

				updateRooms(newRoom)
				go startGame(newRoom)
				break
			case EVENT_GAME_NO_START:
				//人数还不够，不可以开始喔
				msg = "人数还不够，不可以开始喔~"
				//helper.Debug("人数还不够，不可以开始喔~")
				break
			case EVENT_NO_PLACE:
				//房间不能加人了
				msg = "房间已经满人了喔！"
				//helper.Debug("房间已经满人了喔！")
				break
			}
		}

	} else {
		eventType = EVENT_INVAILD
		msg = "无效事件"
		//helper.Debug("无效事件~")
		members := make([]Member, 0)
		members = append(members, member)
		newRoom = Room{
			Id:     "nil",
			Name:   "无效事件",
			Member: members,
		}
	}

	event := Event{
		Data:      user,
		TimeUnix:  time.Now().Unix(),
		Msg:       msg,
		EventType: eventType,
		Room:      newRoom}

	////helper.DebugStructToString(event)

	event.Uid = user.Id

	//推送数据到浏览器
	publish <- event
}

//此函数处理所有传入的阻塞消息。
func chatRoom() {
	for {
		select {
		case event := <-publish:
			switch event.EventType {
			case EVENT_HAND:
				//helper.Debug("握手")
				event.Msg = "握手成功"
				break
			case EVENT_JOIN, EVENT_GAME_NO_START:
				//加入房间
				//helper.Debug("加入房间")
				break
			case EVENT_CREATE:
				//创建房间
				//helper.Debug("创建房间")
				event.Room.Id = helper.GetRandomString(16)
				break
			case EVENT_MESSAGE:
				//helper.Debug("发送信息")
				break
			case EVENT_DRAW:
				event.Msg = "啊！请接收绘图"
				//helper.Debug("绘图")
				break
			case EVENT_BREAK_DRAW:
				event.Msg = "啊！请中断绘图"
				//helper.Debug("中断绘图")
				break
			case EVENT_GAME_CAN_START:
				event.Msg = "还不可以开始喔~"
				//helper.Debug("还不可以开始喔~")
				break
			case EVENT_NO_PLACE:
				event.Msg = "人已经满了喔~"
				//helper.Debug("人已经满了喔~")
				break
			case EVENT_GAME_ANSWER:
				//helper.Debug("回答问题")
				break
			case EVENT_GAME_BONUS:
				//helper.Debug("加分事件")
				break
			case EVENT_NEW_DRAW:
				//helper.Debug("新的绘图事件")
				break
			case EVENT_SYSTEM_MESSAGE:
				//helper.Debug("系统消息")
				break
			default:
				//握手时候，没有房间号
				//helper.Debug("假装握手")
				event.EventType = EVENT_HAND
				event.Msg = "握手成功"
			}
			//helper.DebugStructToString(event)
			broadcastWebSocket(event)
		}
	}
}

// 广播消息
func broadcastWebSocket(event Event) {

	event.TimeUnix = time.Now().Unix()

	data, err := json.Marshal(event)
	if helper.Error(err) {
		helper.Error("broadcastWebSocket 发生错误，不能发送消息")
		return
	}

	var room Room
	var member []models.User

	////helper.Debug("event.EventType -- ",event.EventType)

	switch event.EventType {
	case EVENT_CREATE:
		//创建了房间，要通知所有人更新房间列表
		member = getAllMember()
		break
	case EVENT_HAND:
		//helper.Debug("新建房间，握手，需要发送数据给自己 -> ", event.Room.Member[0].UserId)
		//新建房间，握手，需要发送数据给自己
		if len(event.Room.Member) > 0 {
			var has bool
			has, u := hasMember(event.Room.Member[0].UserId)
			if has {
				member = append(member, u)
			} else {
				helper.Error("broadcastWebSocket 发生错误，用户数据丢失")
				return
			}
		}
		break
	case EVENT_NO_PLACE:
		event.Msg = "人已经满了喔~"
		//helper.Debug("人已经满了喔~")
		room = getRoom(event.Room.Id)
		has, user := hasMember(event.Uid)
		if has {
			member = make([]models.User, 0)
			member = append(member, user)
		}
		break
	case EVENT_GAME_IS_START:
		event.Msg = "遊戲正在進行中喔~"
		//helper.Debug("遊戲正在進行中喔~")
		room = getRoom(event.Room.Id)
		has, user := hasMember(event.Uid)
		if has {
			member = make([]models.User, 0)
			member = append(member, user)
		}
		break
	case EVENT_GAME_MEMBER_NOT_ENOUGH:
		event.Msg = "人數沒满喔~不可以開始"
		//helper.Debug("人數沒满喔~不可以開始")
		room = getRoom(event.Room.Id)
		has, user := hasMember(event.Uid)
		//helper.Debug("has", has, event.Uid)
		if has {
			member = make([]models.User, 0)
			member = append(member, user)
		}
		break
	case EVENT_NEW_DRAW:
		//不用发给自己了
		room = getRoom(event.Room.Id)
		if len(room.Member) > 0 {
			member = getMemberByRoom(room)
			newMember := make([]models.User, 0)

			//helper.Debug("event.Room.MasterId -- ", event.Room.MasterId)

			for _, m := range member {
				if event.Room.MasterId != m.Id {
					newMember = append(newMember, m)
				}
			}

			member = newMember
			////helper.DebugStructToString(member)
		} else {
			//已经没有用户了，应该销毁他
			//helper.Debug("已经没有用户了，应该销毁他")
			return
		}
		break
	default:
		room = getRoom(event.Room.Id)
		if len(room.Member) > 0 {
			member = getMemberByRoom(room)
			////helper.DebugStructToString(member)
		} else {
			//已经没有用户了，应该销毁他
			//helper.Debug("已经没有用户了，应该销毁他")
			return
		}
	}

	if event.EventType == EVENT_GAME_OVER {
		//helper.Debug("已經~晚了")
	}

	event.Mutex.Lock()
	for _, m := range member {
		ws := m.Conn
		if ws != nil {
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				//发生错误，这里应该作重连机制
				beego.Error(err)
				ws.Close()
			}
		}
	}
	defer event.Mutex.Unlock()
}

//开始游戏
func startGame(room Room) bool {

	t := 0
	t1 := time.NewTimer(time.Second * 1)

	room.IsOver = false

	event := Event{
		EventType: 0,
		Room:      room,
		Msg:       "嘿嘿嘿",
		TimeUnix:  0,
		Data:      nil,
	}

loop:
	for {
		select {
		case <-t1.C:
			t++
			if t > YouPerformIGuess.TimeOver {
				//开始新的一轮游戏啦
				t = 0
				//转换身份事件
				isOver := 0 //是否应该结束游戏呢? =最大人數 => true
				hasMaster := false
				for i, m := range room.Member {
					switch m.UserType {
					case MASTER:
						//随机获取问题答案
						vocabulary := models.Vocabulary
						lenVocabulary := len(vocabulary)
						r := rand.New(rand.NewSource(time.Now().UnixNano()))
						key := r.Intn(lenVocabulary)
						event.Room.KeyWord = vocabulary[key]
						event.Room.Member[i].UserType = NO_MASTER
						isOver++
						updateRooms(event.Room)
						break
					case NO_MASTER:
						isOver++
						break
					case PLAYER:
						if !hasMaster {
							hasMaster = true
							event.Room.Member[i].UserType = MASTER
						}
						updateRooms(event.Room)
						break
					}
				}
				//helper.Debug("isOver:", isOver)

				if isOver == len(event.Room.Member) {
					room.Times++
					if YouPerformIGuess.Times < room.Times {
						//中断游戏
						//helper.Debug("游戏已經結束啦")
						event.EventType = EVENT_GAME_OVER
						broadcastWebSocket(event)
					} else {
						helper.Debug("第", room.Times, "轮游戏")

						//新的一轮交换顺序
						newMember := make([]Member, 0)
						for i := len(room.Member) - 1; i > -1; i-- {
							room.Member[i].UserType = VIEWER
							newMember = append(newMember, room.Member[i])
						}
						room.Member = newMember
						user := models.User{}
						newWS(user, room, EVENT_GAME_RE_START)
					}
					newRoom := getRoom(room.Id)
					updateRooms(newRoom)
					break loop
				} else {
					//更新身份
					event.EventType = EVENT_GIVE_IDENTITY
					broadcastWebSocket(event)
				}
			} else {
				newRoom := getRoom(room.Id)
				if newRoom.IsOver {
					room = newRoom
					updateRooms(newRoom)
					helper.Debug("全部人已经回答完了")
					//全部人已经回答完了
					t = YouPerformIGuess.TimeOver + 1
				} else {
					//推送本轮游戏剩余时间
					event.Room = newRoom
					event.EventType = EVENT_GAME_TIME
					event.Data = YouPerformIGuess.TimeOver - t
					broadcastWebSocket(event)
				}
			}
			t1.Reset(time.Second * 1)
		}
	}

	return true
}

//初始化函数
func init() {
	//死循环都要开线程去搞他
	go chatRoom()
}
