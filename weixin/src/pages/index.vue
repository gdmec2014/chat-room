<template>
  <div class="container" @click="clickHandle('test click', $event)">
    <div class="userinfo" @click="bindViewTap" v-if="user">
      <img class="userinfo-avatar" v-if="user.avatar" :src="user.avatar" background-size="cover">
      <div class="userinfo-nickname" v-if="user.nick_name">
        <card :text="user.nick_name"></card>
      </div>
    </div>
    <a href="/pages/counter" class="counter">去往Vuex示例页面</a>
  </div>
</template>

<script>
import card from "@/components/card";

export default {
  mpType: "page",

  data() {
    return {
      motto: "Hello World",
      userInfo: {}
    };
  },

  computed: {
    user() {
      return this.$store.getters.user;
    }
  },

  components: {
    card
  },

  methods: {
    bindViewTap() {
      const url = "/packageA/logs";
      this.$router.push(url);
    },
    clickHandle(msg, ev) {
      console.log("clickHandle:", msg, ev);
    }
  },

  created() {}
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
