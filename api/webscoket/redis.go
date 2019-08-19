package webscoket

import (
	"chat-room/api/helper"
	"fmt"
	"github.com/go-redis/redis"
)

var RidesClient *redis.Client

func init() {
	RidesClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := RidesClient.Ping().Result()
	fmt.Println(pong, err)
	GetSet("room")
}

func SetMap(name string, data map[string]interface{}) {
	RidesClient.HMSet(name, data)
}

func SetSAdd(name string, data ...interface{}) {
	RidesClient.SAdd(name, data)
}

func GetSet(name string) {
	helper.Debug("GetSet : ",RidesClient.SMembers(name))
}

func IsInSet(name string, data interface{}) bool {
	isMember, err := RidesClient.SIsMember(name, data).Result()
	helper.Error(err)
	return isMember
}

func DelSet(name string, data interface{}) {
	RidesClient.SRem(name, data)
}
