<template>
  <div>
    <ws></ws>123
    <canvas
      type="2d"
      @touchstart="startEvent"
      @touchmove="moveEvent"
      @touchend="endEvent"
      id="canvas"
      canvas-id="canvas"
    ></canvas>
    <button @click="draw">画笔</button>
    <button @click="reset">重置</button>
  </div>
</template> 

<script>
import ws from "./ws";

export default {
  name: "game",
  mpType: "page",

  components: { ws },
  data() {
    return {
      showCanvas: true,
      shareImage: "",
      painting: {},
      width: 0,
      height: 0,
      img: "",
      members: [],
      all_over: false,
      answer: "",
      user_type: this.Code.VIEWER,
      user_type_tip: "正在打醬油",
      g_time: 0,
      room_msg: [],
      in_room: false,
      msg: "",
      websock: null,
      rule_form: {
        name: "111",
        password: "111"
      },
      all_room: [],
      //绘图
      context: null, //2d对象
      is_draw: false, //是否正在绘制，松开鼠标时候应该为false
      canvas: null, //画布
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
      },
      startPoint: {},
      endPoint: {}
    };
  },
  computed: {
    user() {
      return this.$store.getters.user;
    },
    ws_data() {
      return this.$store.getters.room.ws_data;
    },
    open_room_data() {
      return this.$store.getters.room.open_room_data;
    }
  },
  mounted() {
    let that = this;
    const query = wx.createSelectorQuery();
    query
      .select("#canvas")
      .node()
      .exec(res => {
        console.log(res);
        const canvas = res[0].node;
        that.canvas = canvas;
        that.context = canvas.getContext("2d");
        that.width = canvas._width;
        that.height = canvas._height;

        that.initCanvas();
      });
  },
  methods: {
    draw() {
      console.log(this.canvas);
    },
    drawLine(ctx, startPoint, endPoint) {
      ctx.moveTo(startPoint.x, startPoint.y);
      ctx.lineTo(endPoint.x, endPoint.y);
      ctx.stroke();
      console.log(this.canvas);
      console.log(this.ctx);
    },
    reset() {
      // Round One
      // this.context.clearRect(0, 0, this.canvas.width, this.canvas.height)
      // Round Two
      // this.canvas.height = this.canvas.height
      this.initCanvas();
      // Round Three
      this.context.globalCompositeOperation = "destination-out";
      this.context.beginPath();
      this.context.fillRect(0, 0, this.canvas.width, this.canvas.height);
      this.context.fill();
      this.context.globalCompositeOperation = "source-over";
    },
    startEvent(e) {
      console.log("startEvent", e);
      let touchkPoint = e.touches[0];
      this.startPoint = touchkPoint;
    },
    moveEvent(e) {
      this.endPoint = e.touches[0];
      this.drawLine(this.context, this.startPoint, this.endPoint);
      this.startPoint = this.endPoint;
    },
    endEvent(e) {
      this.startPoint = {};
      this.endPoint = {};
    },
    initCanvas() {
      console.log("init");

      const dpr = wx.getSystemInfoSync().pixelRatio;
      console.log("dpr", dpr);
      console.log("this.width", this.width);
      console.log("this.height", this.height);
      this.canvas.width = this.width * dpr;
      this.canvas.height = this.height * dpr;
      this.context.fillStyle = "#4194d1";
      this.context.strokeStyle = "#4194d1";
      this.context.lineCap = "round";
      this.context.lineJoin = "round";
      this.context.scale(dpr, dpr);
    },
    restart() {
      let that = this;
      //更新用戶信息
      // $.ajax({
      //   type: "GET",
      //   url: "/v1/room/restart?token=" + token,
      //   contentType: "application/json; charset=utf-8",
      //   success: function(res) {
      //     console.log(res);
      //     if (res.Result == 10000) {
      //     } else {
      //       that.$message(res.Message);
      //     }
      //   },
      //   error: function(res) {
      //     console.error(res);
      //   }
      // });
    },

    //绘图初始化
    drawInit() {
      this.the_canvas = document.querySelector("#the_canvas");
      if (!this.the_canvas || !this.the_canvas.getContext) {
        alert("你的浏览器不支持画画喔~");
        return;
      }
      //获取2D对象
      this.context = this.the_canvas.getContext("2d");

      this.the_canvas.onmousedown = e => this.onMouseDown(e);
      this.the_canvas.onmouseup = e => this.onMouseup(e);
    },
    //求出缩放的倍数
    //每个人的屏幕物理像素大小都可能不一样，所以需要重新设计
    windowToCanvas(canvas, x, y) {
      //canvas 是移动中不断更新的，所以要传过来
      //获取canvas元素距离窗口的一些属性，MDN上有解释
      let rect = canvas.getBoundingClientRect();
      //x和y参数分别传入的是鼠标距离窗口的坐标，然后减去canvas距离窗口左边和顶部的距离。
      return {
        x: x - rect.left * (canvas.width / rect.width),
        y: y - rect.top * (canvas.height / rect.height)
      };
    },
    //按下鼠标
    onMouseDown(e) {
      if (this.user_type != MASTER) {
        console.log("你不可以画1");
        return;
      }
      this.is_draw = true;
      //获得鼠标按下的点相对canvas的坐标。
      let ele = this.windowToCanvas(this.the_canvas, e.clientX, e.clientY);
      let { x, y } = ele;

      //绘制起点。
      this.context.beginPath();
      this.context.moveTo(x, y);
      //鼠标移动事件
      this.the_canvas.onmousemove = e => {
        this.the_canvas.onmouseup = e => this.onMouseup(e);
        this.draw(e.clientX, e.clientY);
      };
    },
    //松开鼠标
    onMouseup() {
      if (this.user_type != MASTER) {
        return;
      }
      //告诉其他画家，哥哥我中断了绘画
      this.breakDraw();
      this.is_draw = false;
    },
    //绘制到屏幕
    draw(old_x, old_y) {
      if (this.is_draw) {
        //移动时获取新的坐标位置，用lineTo记录当前的坐标，然后stroke绘制上一个点到当前点的路径
        let ele = this.windowToCanvas(this.the_canvas, old_x, old_y);
        //console.log(ele)
        let { x, y } = ele;

        this.context.lineTo(x, y);
        this.context.stroke();
        //是绘制者才应该发送数据去
        if (this.is_master) {
          //this.newSendDraw()
          //this.sendDraw(x, y)
        }
      }
    },
    //接收绘图事件
    receiveDraw(draw_user_id, draw_room_id, x, y) {
      //自己绘制的图，不需要再画在屏幕
      //不是自己的房间，不要画图
      if (
        draw_user_id == this.user.uid ||
        draw_room_id != this.open_room_data.id
      ) {
        return;
      }
      //绘制起点
      if (!this.in_start) {
        this.in_start = true;
        this.context.moveTo(x, y);
        return;
      }
      this.is_draw = true;
      this.draw(x, y);
    },
    newReceiveDraw(draw_user_id, draw_room_id, pic_url) {
      //自己绘制的图，不需要再画在屏幕
      //不是自己的房间，不要画图
      if (
        draw_user_id == this.user.uid ||
        draw_room_id != this.open_room_data.id
      ) {
        return;
      }
      //console.log(pic_url)
      this.img = pic_url;
    },
    //发送绘图事件
    sendDraw(x, y) {
      this.ws_data.room.id = this.open_room_data.id;
      this.ws_data.msg = "绘图";
      this.ws_data.event_type = EVENT_DRAW;
      this.ws_data.data = { x, y };
      this.websocketSend(JSON.stringify(this.ws_data));
    },
    newSendDraw() {
      this.ws_data.room.id = this.open_room_data.id;
      this.ws_data.msg = "绘图";
      //this.ws_data.event_type = EVENT_DRAW
      this.ws_data.event_type = EVENT_NEW_DRAW;

      var d = this.the_canvas.toDataURL("image/png", 0.2);
      // console.log("d --", d)
      this.ws_data.data = d;
      this.websocketSend(JSON.stringify(this.ws_data));
    },
    breakDraw() {
      this.newSendDraw();
      // this.ws_data.room.id = this.open_room_data.id
      // this.ws_data.msg = "中断绘图"
      // this.ws_data.event_type = EVENT_BREAK_DRAW
      // this.websocketSend(JSON.stringify(this.ws_data));
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
