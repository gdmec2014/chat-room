<template>
  <div>
    <div v-if="in_page == 0">
      <a class="counter" @click="SetViewRoomPage('/pages/room/create')">创建房间</a>
      <a class="counter" @click="SetViewRoomPage('/pages/room/list')">查看房间</a>
    </div>
  </div>
</template>

<script>
import roomCreate from "../room/create";
import roomList from "../room/list";

export default {
  mpType: "page",

  data() {
    return {
      in_index: true
    };
  },

  computed: {
    user() {
      return this.$store.getters.user;
    },
    in_page() {
      return this.$store.getters.room.in_page;
    }
  },

  components: {
    roomCreate,
    roomList
  },

  methods: {
    bindViewTap() {
      const url = "/packageA/logs";
      this.$router.push(url);
    },
    clickHandle(msg, ev) {
      console.log("clickHandle:", msg, ev);
    },
    ViewRoom() {
      this.SetViewRoomPage(2);
      this.$store.dispatch("GetAllRoom");
    },
    ShowCreateRoom() {
      this.SetViewRoomPage(1);
      this.$store.dispatch("GetAllRoom");
    },
    SetViewRoomPage(url) {
      this.$router.push(url);
    }
  }
};
</script>

<style scoped>
.userinfo {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.userinfo-avatar {
  width: 128rpx;
  height: 128rpx;
  margin: 20rpx;
  border-radius: 50%;
}

.userinfo-nickname {
  color: #aaa;
}

.form-control {
  display: block;
  padding: 0 12px;
  margin-bottom: 5px;
  border: 1px solid #ccc;
}

.counter {
  display: inline-block;
  margin: 10px auto;
  padding: 5px 10px;
  color: blue;
  border: 1px solid blue;
}
</style>
