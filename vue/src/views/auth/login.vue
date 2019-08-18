<template>
  <el-row type="flex" class="row-bg register-box" justify="center" align="middle">
    <el-card class="box-card register-form">
      <el-row slot="header" type="flex" class="row-bg clearfix" justify="center">
        <span>用户登录</span>
      </el-row>
      <el-form
        :label-position="labelPosition"
        label-width="80px"
        :model="loginData"
        :rules="rules"
        ref="loginData"
        status-icon
      >
        <el-form-item label="用户名" prop="name">
          <el-input v-model="loginData.name" placeholder="用户名" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            type="password"
            v-model="loginData.password"
            placeholder="密码"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item class="submit-item">
          <el-button type="primary" @click="login()">登录</el-button>
          <router-link to="register" class="el-button el-button--success">去注册</router-link>
        </el-form-item>
      </el-form>
    </el-card>
  </el-row>
</template>

<script>
import { Login } from "../../api/auth";
import { setToken } from "../../util/auto";

export default {
  name: "login",
  data() {
    return {
      labelPosition: "left",
      loginData: {
        name: "",
        password: "",
      },
      rules: {
        name: [
          { required: true, message: "用户名", trigger: "blur" },
          { min: 2, max: 12, message: "请输入用户名", trigger: "blur" }
        ],
        password: [
          { required: true, message: "密码", trigger: "blur" },
          { message: "请输入密码", trigger: "blur" }
        ]
      }
    };
  },
  methods: {
    login() {
      let that = this;
      this.$refs["loginData"].validate(valid => {
        if (valid) {
          Login(this.loginData).then(data => {
            if (data.Result == 10000) {
              setToken(data.Data.token);
              that.$message({
                type: "success",
                message: data.Message
              });
            } else {
              that.$message({
                type: "error",
                message: data.Message
              });
            }
          });
        }
      });
    }
  }
};
</script>
<style lang="less">
html,
body,
#app,
.register-box {
  height: 100%;
  background-color: #e1f3d8;
}
.register-form {
  width: 25rem;
}
.submit-item {
  margin-bottom: 0;
}
.el-input__prefix,
.el-input__suffix {
  color: #67c23a;
}
</style>