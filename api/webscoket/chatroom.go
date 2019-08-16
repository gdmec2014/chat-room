package webscoket

import (
	"chat-room/api/models"
	"container/list"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

//TODO webscoket 房间事件处理

type Subscription struct {
	Archive []Event      // 所有事件
	New     <-chan Event // 新的事件
}

type Room struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Type UserType `json:"type"`
}

//创建新的房间
func newEvent(ep EventType, user models.User, msg, roomId string) Event {
	return Event{ep, user, int(time.Now().Unix()), msg, roomId}
}

//加入房间
func Join(user models.User, ws *websocket.Conn, roomId string, userType UserType) {
	subscribe <- Subscriber{User: user, Conn: ws, RoomId: roomId, Type: userType}
}

//销毁房间
func Leave(id int64) {
	unsubscribe <- id
}

type Subscriber struct {
	RoomId string          `json:"room_id"` //房间ID
	User   models.User     `json:"user"`    //user
	Conn   *websocket.Conn `json:"conn"`    // 房间链接
	Type   UserType        `json:"type"`
}

var (
	// 加入新房间的阻塞通道
	subscribe = make(chan Subscriber, 10)
	// 退出房间的阻塞通道
	unsubscribe = make(chan int64, 10)
	// 推送消息的阻塞通道
	publish = make(chan Event, 10)
	// 长轮询列表
	waitingList = list.New()
	subscribers = list.New()
)

//此函数处理所有传入的阻塞消息。
func chatRoom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.User.Id) {
				subscribers.PushBack(sub) // 将用户添加到列表末尾。
				// 推送加入事件
				publish <- newEvent(EVENT_JOIN, sub.User, "", sub.RoomId)
				beego.Info("New user:", sub.User.Id, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.User.Id, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:

			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			NewArchive(event)

			if event.Type == EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).User.Id == unsub {
					subscribers.Remove(sub)
					// 关闭链接
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(EVENT_LEAVE, sub.Value.(Subscriber).User, "", sub.Value.(Subscriber).RoomId) // 推送离开事件
					break
				}
			}
		}
	}
}

func init() {
	//死循环都要开线程去搞他
	go chatRoom()
}

//判断是否已经加入了链接
func isUserExist(subscribers *list.List, id int64) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).User.Id == id {
			return true
		}
	}
	return false
}
