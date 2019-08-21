<template>
  <div>
    <el-row>
      <el-col :span="4">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span round>
              <i class="el-icon-warning"></i>
            </span>
            <span v-text="userInfo.name">用户昵称</span>
            <el-button style="padding: 3px 0" type="text">
              <i class="el-icon-edit"></i>
            </el-button>
          </div>
          <div class="text item">
            <el-row>
              <el-col :span="8">等级</el-col>
              <el-col v-text="userInfo.level" :span="16">4</el-col>
            </el-row>
            <el-row>
              <el-col :span="8">积分</el-col>
              <el-col v-text="userInfo.score" :span="16">360</el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <el-col v-for="o in 24" :key="o" :xs="12" :sm="12" :md="8" :lg="6" :xl="4" class="room-box">
        <div class="room">
          <el-row>
            <el-col :span="12" align="right" class="user-box left">
              <span class="head"></span>
            </el-col>
            <el-col :span="12" align="left" class="user-box right">
              <span class="head"></span>
            </el-col>
          </el-row>
          <el-row class="room-center">
            <el-col :span="8">
              <div class="user-box left">
                <span class="head"></span>
              </div>
            </el-col>
            <el-col :span="8" align="middle">
              <span type="flex" justify="center" class="room-num">{{o}}</span>
            </el-col>
            <el-col :span="8">
              <div class="user-box right">
                <span class="head"></span>
              </div>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12" align="right" class="user-box left">
              <span class="head"></span>
            </el-col>
            <el-col :span="12" align="left" class="user-box right">
              <span class="head"></span>
            </el-col>
          </el-row>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "home",
  data() {
    return {
      userInfo: {
        name: "sinkr",
        level: 4,
        score: 3600
      }
    };
  },
  mounted(){

  },

  methods: {

  }
};
</script>
<style lang="less">
.room-box {
  margin-top: 2rem;
}
.room .head {
  display: inline-block;
  height: 3rem;
  width: 3rem;
  border-radius: 50%;
  background-color: aquamarine;
}
.room .room-num {
  display: inline-block;
  height: 6rem;
  width: 6rem;
  border-radius: 50%;
  background-color: rgb(255, 200, 127);
  line-height: 6rem;
  font-size: 1.75rem;
  font-weight: bold;
}
.room .room-center {
  height: 6rem;
}
.room .room-center .user-box {
  height: 6rem;
  display: flex;
  align-items: center;
}
.user-box {
  font-size: 0;
}
.user-box.left {
  padding-right: 1rem;
}
.user-box.right {
  padding-left: 1rem;
}
.room-center .user-box.left {
  justify-content: flex-end;
  padding-right: .5rem;
}
.room-center .user-box.right {
  justify-content: flex-start;
  padding-left: .5rem;
}
</style>