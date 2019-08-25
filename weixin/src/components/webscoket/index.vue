<script>
import store from "../../store";

export default {
  data() {
    return {
      url: "ws://127.0.0.1:2332/v1/room/hand?token=",
      ws: null,
      socket_open: false
    };
  },
  watch: {
    // 实时检测
    user(new_user) {
      console.log("new_user", new_user);
      if (new_user) {
        if (new_user.token) {
          this.init();
        }
      }
    }
  },
  computed: {
    user() {
      return store.getters.user;
    }
  },
  methods: {
    init() {
      wx.connectSocket({
        url: this.url + this.user.token
      });
      wx.onSocketOpen(res => this.onSocketOpen(res));
      wx.onSocketMessage(res => this.onSocketMessage(res));
      wx.onSocketError(res => this.onSocketError(res));
      wx.onSocketClose(res => this.onSocketClose(res));
    },
    onSocketOpen(res) {
      console.log("webscoket 连接成功");
      this.socket_open = true;
      store.dispatch("SetUserConn", true);
    },
    onSocketMessage(res) {
      console.log("webscoket 收到消息", res);
      let data = JSON.parse(res.data);

      switch (data.event_type) {
        case this.Code.EVENT_HAND:
          console.log("握手成功", data);
          break;
        case this.Code.EVENT_CREATE:
          this.ws_data = data;
          console.log("创建成功", data);
          store.dispatch("AddRoom", data.room);
          //如果是自己创建的，则应该自动进入该房间
          if (data.data.uid == this.user.uid) {
            console.log("自己创建的，则应该自动进入该房间")
            store.dispatch("SetOpenRoomData", data.room);
            this.SetViewRoomPage(3)
          }
          break;
        case this.Code.EVENT_JOIN:
          console.log("加入成功", data);
          //this.$message(data.msg);
          //this.getAllRoom();
          break;
        case this.Code.EVENT_MESSAGE:
          console.log("收到消息", data);
          // let msg = data.data;
          // msg["msg"] = data.msg;
          // console.log(msg);
          // this.setLocalStorageMsg(msg);
          break;
        case this.Code.EVENT_DRAW:
          //是否为绘画者，这个一定要判断，不然会无限发送事件
          console.log("绘图事件", data);
          // this.is_master = data.data.user.uid == this.user.uid;
          // this.receiveDraw(
          //   data.data.user.uid,
          //   data.room.id,
          //   data.data.position.x,
          //   data.data.position.y
          // );
          break;
        case this.Code.EVENT_BREAK_DRAW:
          //中断绘画
          console.log("中断绘图事件", data);
          // this.is_master = data.data.user.uid == this.user.uid;
          // if (!this.is_master) {
          //   this.is_draw = false;
          //   this.in_start = false;
          // }
          break;
        default:
          console.log("不知道是啥", data);
      }
    },
    sendSocketMessage(data) {
      console.log("webscoket 发送消息");
      if (this.socket_open) {
        wx.sendSocketMessage({
          data: JSON.stringify(data)
        });
      } else {
        wx.showToast({
          title: "已经与后台断开连接",
          icon: "none",
          duration: 2000
        });
      }
    },
    onSocketError(err) {
      console.error(err);
    },
    onSocketClose(res) {
      console.info("webscoket 关闭了呢");
      console.log(res);
    },
    SetViewRoomPage(in_page) {
      store.dispatch("SetRoomPage", in_page);
    }
  }
};
</script>
