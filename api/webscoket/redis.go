package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
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
		models.Redis.Enable = false
	} else {
		//从redis读取旧的房间数据
		beego.Info("从redis读取旧的房间数据")
		InitData()
	}
	beego.Info(pong)
}

func SetMap(name string, data map[string]interface{}) {
	if !models.Redis.Enable {
		return
	}
	RidesClient.HMSet(name, data)
}

func GetMap(name string) (data map[string]string) {
	if models.Redis.Enable {
		return
	}
	var err error
	r := RidesClient.HGetAll(name)
	data, err = r.Result()
	helper.Error(err)
	helper.Debug("GetMap:", data)
	return data
}

func SetSAdd(name string, data ...interface{}) (err error) {
	if models.Redis.Enable {
		return
	}
	dataStrs := make([]interface{}, 0)
	for _, str := range data {
		dataStrs = append(dataStrs, helper.GetString(str))
	}
	r := RidesClient.SAdd(name, dataStrs...)
	err = r.Err()
	helper.Error(err)
	return
}

func GetSet(name string, data interface{}) {
	if models.Redis.Enable {
		return
	}
	r := RidesClient.SMembers(name)
	values, err := r.Result()
	if helper.Error(err) {
		helper.Error("GetSet error")
		return
	}

	helper.Debug("GetSet : ", values)

	switch data.(type) {
	case *[]Member:
		var ms []Member
		for _, d := range values {
			var m = Member{}
			err := json.Unmarshal([]byte(d), &m)
			if helper.Error(err) {
				continue
			}
			ms = append(ms, m)
		}
		b := helper.GetByte(ms)
		err := json.Unmarshal(b, &data)
		helper.Error(err)
		break
	case *[]string:
		b := helper.GetByte(values)
		err := json.Unmarshal(b, &data)
		helper.Error(err)
	}
}

func IsInSet(name string, data interface{}) (isMember bool) {
	if models.Redis.Enable {
		return
	}
	var err error
	dataStr := helper.GetString(data)
	isMember, err = RidesClient.SIsMember(name, dataStr).Result()
	helper.Error(err)
	return isMember
}

func DelSet(name string, data interface{}) {
	if models.Redis.Enable {
		return
	}
	RidesClient.SRem(name, data)
}
