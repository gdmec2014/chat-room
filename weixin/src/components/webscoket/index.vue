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
        if (true) {
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
      this.socketOpen = true;
      store.dispatch("SetUserConn", true);
    },
    onSocketMessage(res) {
      console.log("webscoket 收到消息");
    },
    sendSocketMessage(data) {
      console.log("webscoket 发送消息");
      if (this.socket_open) {
        wx.sendSocketMessage(data);
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
    }
  }
};
</script>
