package webscoket

import (
	"chat-room/api/models"
	"container/list"
)

//TODO webscoket 数据处理

type EventType int

const (
	EVENT_JOIN = iota + 10
	EVENT_LEAVE
	EVENT_MESSAGE
)

type UserType int

const (
	VIEWER = 100
	PLAYER = 101
)

type Event struct {
	Type      EventType   `json:"type"`      // 消息类型
	User      models.User `json:"user"`      // 用户
	Timestamp int         `json:"timestamp"` // 时间
	Content   string      `json:"content"`   // 消息
	RoomId    string      `json:"room_id"`   // 房间id
}

type youPerformIGuess struct {
	MaxNumber int `json:"max_number"`
}

var (
	archive                            = list.New()            // 全部链接
	YouPerformIGuess *youPerformIGuess = new(youPerformIGuess) // 配置
)

//初始化
func init() {
	YouPerformIGuess.MaxNumber = models.GetAppConfInt("youPerformIGuess::maxNumber")
	if YouPerformIGuess.MaxNumber < 1 {
		YouPerformIGuess.MaxNumber = 6
	}
}

// 将新的事件保存到全局消息
func NewArchive(event Event) {
	archive.PushBack(event)
}

//根据房间ID获取组成员
func getRoomUsers(roomId string) (subscriber []Subscriber) {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).RoomId == roomId {
			ws := sub.Value.(Subscriber).Conn
			if ws != nil {
				//如果链接不正常，就抛弃他
				subscriber = append(subscriber, sub.Value.(Subscriber))
			}
		}
	}
	return
}
