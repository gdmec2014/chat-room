<template>
  <div>
    <div>
      <div v-show="in_creat">
        <van-cell-group>
          <van-field
            @change="SetRoomName"
            required
            maxlength="16"
            label="房间名"
            v-model="ws_data.room.name"
            placeholder="请输入房间名"
          />

          <van-radio-group :value="ws_data.room.max_member">
            <van-cell-group>
              <van-cell
                title="2 人房间"
                value-class="value-class"
                clickable
                data-name="2"
                @click="SetRoomMaxMember('2')"
              >
                <van-radio name="2"/>
              </van-cell>
              <van-cell
                title="4 人房间"
                value-class="value-class"
                clickable
                data-name="4"
                @click="SetRoomMaxMember('4')"
              >
                <van-radio name="4"/>
              </van-cell>
            </van-cell-group>
          </van-radio-group>
        </van-cell-group>

        <div class="bnt">
          <van-button @click="createRoom" class="bnt" size="normal" type="primary">创建</van-button>
        </div>
      </div>
    </div>
    <div>
      <van-card desc="描述信息" title="商品标题"/>
    </div>
  </div>
</template> 

<script>
import { GetAllRoom } from "../../api/room";

export default {
  name: "room",
  data() {
    return {
      all_room: [],
      in_creat: true,
      room_msg: [],
      in_room: false,
      open_room_data: {
        id: "",
        name: "",
        user_type: 0
      },
      ws_data: {
        event_type: 0,
        room: {
          id: "",
          name: "",
          user_type: 0,
          max_member: "2"
        },
        msg: "",
        data: {}
      }
    };
  },
  computed: {
    user() {
      return this.$store.getters.user;
    }
  },
  methods: {
    createRoom() {
      if (this.ws_data.room.name) {
        this.ws_data.event_type = this.Code.EVENT_CREATE;
        this.ws_data.msg = "创建房间";
        let data = this.ws_data;
        data.room.max_member = Number(this.ws_data.room.max_member);
        console.log(JSON.stringify(data));
        this.$sendSocketMessage(data);
      } else {
        wx.showToast({
          title: "请输入房间名",
          icon: "none",
          duration: 2000
        });
      }
    },
    joinRoom(room_id) {
      this.ws_data.room.id = room_id;
      this.ws_data.msg = "请求加入"; 
      this.ws_data.event_type = this.Code.EVENT_JOIN;
      let data = this.ws_data;
      console.log(JSON.stringify(data));
      this.$sendSocketMessage(data);
    },
    sendMsg() {
      if (this.ws_data.msg) {
        console.log("this.ws_data --", this.ws_data);
        this.ws_data.room.id = this.open_room_data.id;
        this.ws_data.msg = this.ws_data.msg;
        this.ws_data.event_type = this.Code.EVENT_MESSAGE;
        let data = this.ws_data;
        this.$sendSocketMessage(data);
      } else {
        wx.showToast({
          title: "要输入东西喔~",
          icon: "none",
          duration: 2000
        });
      }
    },

    SetRoomMaxMember(n) {
      this.ws_data.room.max_member = n;
    },
    SetRoomName(e) {
      this.ws_data.room.name = e.mp.detail;
    },

    GetAllRoom() {
      let that = this;
      GetAllRoom(this.user.token).then(res => {
        console.log(res);
      });
    }
  }
};
</script>

<style scoped>
.bnt {
  margin-top: 0.5rem;
  width: 100%;
  text-align: center;
}
</style>
