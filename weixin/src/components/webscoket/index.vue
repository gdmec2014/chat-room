<script>
import store from "../../store";
export default {
  data() {
    return {
      url: "ws://192.168.0.108:2332/v1/room/hand?token=",
      ws: null,
      socket_open: false,
      re_data: {}
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
      const url = this.url + this.user.token;
      wx.connectSocket({
        url
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
      store.dispatch("GetWsData", data);
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
    }
  }
};
</script>
