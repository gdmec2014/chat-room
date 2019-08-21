<template>
  <el-row type="flex" class="row-bg register-box" justify="center" align="middle">
    <el-card class="box-card register-form">
      <el-row slot="header" type="flex" class="row-bg clearfix" justify="center">
        <span>用户注册</span>
      </el-row>
      <el-form
        :label-position="labelPosition"
        label-width="80px"
        :model="registerData"
        :rules="rules"
        ref="registerData"
        status-icon
      >
        <el-form-item label="用户名" prop="name">
          <el-input v-model="registerData.name" placeholder="用户名" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            type="password"
            v-model="registerData.password"
            placeholder="密码"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="re_password" class="is-required">
          <el-input
            type="password"
            v-model="registerData.re_password"
            placeholder="确认密码"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item class="submit-item">
          <el-button type="primary" @click="register()">提交</el-button>
          <router-link to="login" class="el-button el-button--success">去登录</router-link>
        </el-form-item>
      </el-form>
    </el-card>
  </el-row>
</template>

<script>
import { Register } from "../../api/auth";
import { setToken } from "../../util/auto";

export default {
  name: "register",
  data() {
    var checkPassword = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请再次输入密码"));
      } else if (value !== this.registerData.password) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    };
    return {
      labelPosition: "left",
      registerData: {
        name: "",
        password: "",
        re_password: ""
      },
      rules: {
        name: [
          { required: true, message: "用户名", trigger: "blur" },
          { min: 2, max: 12, message: "长度在 2 到 12 个字符", trigger: "blur" }
        ],
        password: [
          { required: true, message: "密码", trigger: "blur" },
          { min: 6, max: 18, message: "长度在 6 到 18 个字符", trigger: "blur" }
        ],
        re_password: [{ validator: checkPassword, trigger: "blur" }]
      }
    };
  },
  methods: {
    register() {
      let that = this;
      this.$refs["registerData"].validate(valid => {
        if (valid) {
          Register(this.registerData).then(data => {
            if (data.Result == 10000) {
              setToken(data.Data.token);
              that.$message({
                type: "success",
                message: data.Message
              });
              this.$router.push({
                path: '/home'
              })
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