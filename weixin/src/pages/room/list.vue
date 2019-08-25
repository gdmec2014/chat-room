<template>
  <div>
    <div class="van-card" v-for="(item,index) in all_room" :key="index">
      <van-card currency=" " :price="item.member.length+'/'+item.max_member" desc="一起来玩" :title="item.name">
        <view slot="footer">
          <van-button size="mini" @click="joinRoom(item.id)">加入</van-button>
        </view>
      </van-card>
    </div>
  </div>
</template> 

<script>
export default {
  name: "room_list",
  data() {
    return {};
  },
  computed: {
    user() {
      return this.$store.getters.user;
    },
    ws_data() {
      return this.$store.getters.room.ws_data;
    },
    all_room() {
      return this.$store.getters.room.all_room;
    }
  },
  methods: {
    joinRoom(room_id) {
      this.ws_data.room.id = room_id;
      this.ws_data.msg = "请求加入";
      this.ws_data.event_type = this.Code.EVENT_JOIN;
      let data = this.ws_data;
      data.room.max_member = 0
      console.log(JSON.stringify(data));
      this.$sendSocketMessage(data);
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
