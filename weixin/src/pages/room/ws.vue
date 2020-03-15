<script>
export default {
  name: "ws",
  data() {
    return {};
  },
  watch: {
    re_ws_data(val) {
      this.controllers(val);
    }
  },
  computed: {
    re_ws_data() {
      return this.$store.getters.room.re_ws_data;
    }
  },
  methods: {
    controllers(data) {
      console.log(data);
      switch (data.event_type) {
        case this.Code.EVENT_HAND:
          console.log("握手成功", data);
          break;
        case this.Code.EVENT_CREATE:
          this.ws_data = data;
          console.log("创建成功", data);
          this.$store.dispatch("SetOpenRoomData", data);
          this.$router.push("/pages/room/game");
          break;
        case this.Code.EVENT_JOIN:
          console.log("加入成功", data);
          break;
        case this.Code.EVENT_SYSTEM_MESSAGE:
          this.msg_data.msg = data.msg;
          //this.setLocalStorageMsg(this.msg_data);
          break;
        case this.Code.EVENT_MESSAGE:
        case this.Code.EVENT_GAME_ANSWER:
        case this.Code.EVENT_GAME_BONUS:
          //console.log("收到消息", data)
          let msg = data.data;
          msg["msg"] = data.msg;
          //this.setLocalStorageMsg(msg);
          break;
        case this.Code.EVENT_NEW_DRAW:
          console.log("新的绘图事件");
          this.is_master = data.data.user.uid == this.user.uid;
          //   this.newReceiveDraw(
          //     data.data.user.uid,
          //     data.room.id,
          //     data.data.position
          //   );
          break;
        case this.Code.EVENT_DRAW:
          //是否为绘画者，这个一定要判断，不然会无限发送事件
          //console.log("绘图事件", data)
          this.is_master = data.data.user.uid == this.user.uid;
          //   this.receiveDraw(
          //     data.data.user.uid,
          //     data.room.id,
          //     data.data.position.x,
          //     data.data.position.y
          //   );
          break;
        case this.Code.EVENT_BREAK_DRAW:
          //中断绘画
          //console.log("中断绘图事件", data)
          this.is_master = data.data.user.uid == this.user.uid;
          if (!this.is_master) {
            this.is_draw = false;
            this.in_start = false;
          }
          break;
        case this.Code.EVENT_GAME_NO_START:
          //玩家加入房間，但人數沒全
          console.log("玩家加入房間，但人數沒全");
          console.log(data);
          //   this.openRoom(data.room.id);
          //   this.initMember(data);
          //   this.joinRoomMsg(data);
          break;
        case this.Code.EVENT_GAME_CAN_START:
          //玩家加入房間，人數ok,開打
          console.log("開打");
          this.the_canvas.height = this.the_canvas.height;
          //   this.openRoom(data.room.id);
          //   this.initMember(data);
          break;
        case this.Code.EVENT_GAME_ANSWER:
          //回答问题事件
          console.log("回答问题事件");
          break;
        case this.Code.EVENT_GAME_BONUS:
          //答对问题加分事件
          console.log("答对问题加分事件");
          break;
        case this.Code.EVENT_GAME_TIME:
          //游戏计时事件
          console.log("游戏计时事件");
          //   this.g_time = data.data;
          //   if (data.room.key_word) {
          //     this.answer = data.room.key_word;
          //   }
          //   if (this.g_time < 2) {
          //     this.user_type_tip = "你的時間要到了喔~";
          //     this.user_type = VIEWER;
          //   }
          //   if (this.g_time < 2) {
          //     console.log("時間結束");
          //     this.answer = "";
          //   }
          //   this.updateIdentity(data);

          //   if (this.members.length < 1) {
          //     this.initMember(data);
          //   }
          //   this.updatePoint(data);

          break;
        case this.Code.EVENT_GAME_OVER:
          //游戏结束事件
          console.log("結束");
          //this.user_type_tip = "看看你多少分~";
          //this.updateIdentity(data);
          break;
        case this.Code.EVENT_GIVE_IDENTITY:
          //赋予游戏身份
          console.log("赋予游戏身份");
          //   this.the_canvas.height = this.the_canvas.height;
          //   this.img = "";
          break;
        case this.Code.EVENT_GAME_BONUS:
          console.log("加分事件");
          break;
        default:
          console.log("不知道是啥", data);
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

.van-card {
  margin-bottom: 0.5rem;
}
</style>
