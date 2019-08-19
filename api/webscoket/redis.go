package webscoket

import (
	"chat-room/api/helper"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var RidesClient *redis.Client

func init() {
	helper.Debug("init RidesClient")
	RidesClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := RidesClient.Ping().Result()
	if helper.Error(err) {
		helper.Error("init redis failed")
	}
	beego.Info(pong)
}

func SetMap(name string, data map[string]interface{}) {
	RidesClient.HMSet(name, data)
}

func SetSAdd(name string, data ...interface{}) (err error) {
	r := RidesClient.SAdd(name, data...)
	err = r.Err()
	helper.Error(err)
	return
}

func GetSet(name string, data interface{}) {
	r := RidesClient.SMembers(name)
	switch data.(type) {
	case *[]string:
		d := r.Val()
		bd, _ := json.Marshal(&d)
		json.Unmarshal(bd, &data)
		break
	case *Room:

		break
	}
}

func IsInSet(name string, data interface{}) bool {
	isMember, err := RidesClient.SIsMember(name, data).Result()
	helper.Error(err)
	return isMember
}

func DelSet(name string, data interface{}) {
	RidesClient.SRem(name, data)
}
