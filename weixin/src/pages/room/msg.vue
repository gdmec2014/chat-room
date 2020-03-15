<template>
  <div></div>
</template> 

<script>
export default {
  name: "game",
  data() {
    return {
      img: "",
      members: [],
      all_over: false,
      answer: "",
      user_type: VIEWER,
      user_type_tip: "正在打醬油",
      g_time: 0,
      room_msg: [],
      in_room: false,
      open_room_data: {
        id: "",
        name: "",
        user_type: 0
      },
      msg: "",
      websock: null,
      ws_data: {
        event_type: 0,
        room: {
          id: "",
          name: "",
          user_type: 0
        },
        msg: "",
        data: {}
      },
      rule_form: {
        name: "111",
        password: "111"
      },
      user: null,
      all_room: [],
      //绘图
      context: null, //2d对象
      is_draw: false, //是否正在绘制，松开鼠标时候应该为false
      the_canvas: null, //画布
      is_master: true, //是否当前绘画者
      in_start: false, //是否绘制了起点
      msg_data: {
        uid: -1,
        name: "系统消息",
        nick_name: "小灰",
        wx_id: "gdmec_2331333",
        last_login: "2020-03-12T18:39:48+08:00",
        delete_time: "0001-01-01T00:00:00Z",
        conn: {},
        mutex: {},
        msg: "这是系统消息"
      }
    };
  },
  computed: {
    user() {
      return this.$store.getters.user;
    },
    ws_data() {
      return this.$store.getters.room.ws_data;
    }
  },
  methods: {
    updatePoint(data) {
      var marks = data.room.mark;
      if (marks) {
        for (var i = 0; i < this.members.length; i++) {
          this.members[i].point = 0;
          for (var j in marks) {
            var ps = marks[j];
            for (var k = 0; k < ps.length; k++) {
              if (this.members[i].user_id == ps[k].id) {
                this.members[i].point += ps[k].point;
              }
            }
          }
        }
      }
    },
    updateIdentity(data) {
      //更新用戶身份
      if (data.room.member) {
        var an = 0;
        var members = data.room.member;
        for (var i = 0; i < members.length; i++) {
          //找到自己
          if (members[i].user_id == this.user.uid) {
            var t = members[i].user_type;
            this.user_type = t;
            switch (t) {
              case VIEWER:
                //观众
                this.user_type_tip = "正在打醬油";
                break;
              case PLAYER:
                //猜一方
                this.user_type_tip = "正在觀摩對手的傑作";
                break;
              case MASTER:
                //出题一方
                this.user_type_tip = "現在是你的表演時間";
                break;
              case NO_MASTER:
                //不可以再赋予出题一方，当全部玩家身份为103时候，游戏结束，结算分数
                this.user_type_tip = "正在觀摩對手的傑作";
                break;
            }
          }
          //判斷是不是全部人都答了一輪
          if (members[i].user_type == NO_MASTER) {
            an++;
          }
          if (an == members.length) {
            console.info("全部人答了一輪，應該結算分數");
            this.user_type_tip = "看看你多少分~";
            this.all_over = true;
          }
        }
      }
    },
    sendMsg() {
      if (this.ws_data.msg) {
        console.log("this.ws_data --", this.ws_data);
        this.ws_data.room.id = this.open_room_data.id;
        this.ws_data.msg = this.ws_data.msg;
        this.ws_data.event_type = EVENT_GAME_ANSWER;
        this.websocketSend(JSON.stringify(this.ws_data));
      } else {
        this.$message("要输入东西喔~");
      }
    },
    sendSystemMsg() {
      if (this.ws_data.msg) {
        console.log("this.ws_data --", this.ws_data);
        this.ws_data.room.id = this.open_room_data.id;
        this.ws_data.msg = this.ws_data.msg;
        this.ws_data.event_type = EVENT_GAME_ANSWER;
        this.websocketSend(JSON.stringify(this.ws_data));
      } else {
        this.$message("要输入东西喔~");
      }
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
