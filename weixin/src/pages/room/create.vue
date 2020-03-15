<template>
  <div>
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

    <ws></ws>
  </div>
</template> 

<script>
import ws from "./ws";

export default {
  name: "room_create",
  components: { ws },

  data() {
    return {};
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
    SetRoomMaxMember(n) {
      this.ws_data.room.max_member = n;
    },
    SetRoomName(e) {
      this.ws_data.room.name = e.mp.detail;
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
