import { WxRegist } from "../../api/wexin"
import { setToken } from "../../utils/auto"
import { GetUserByToken } from "../../api/auth"
import store from "../../store"

let user = {
  state: {
    data: {
      conn: null
    }
  },

  mutations: {
    SET_USER_DATA: (state, data) => {
      console.log("设置user data ", data)
      state.data = data
    },
    SET_WX_USER_DATA: (state, data) => {
      console.log("设置wx user data ", data)
      state.data = data
    },
    SET_USER_CONN: (state, conn) => {
      console.log("设置user conn ", conn)
      state.data.conn = conn
    },
  },

  actions: {
    SetUserData({
      commit
    }, data) {
      commit('SET_USER_DATA', data)
    },
    //将微信注册到我们服务器
    SetWXUserData({
      commit
    }) {
      wx.login({
        success: () => {
          //登陆微信
          wx.getUserInfo({
            success: (res1) => {
              wx.login({
                success: (res2) => {
                  res1.code = res2.code
                  WxRegist(res1).then(res3 => {
                    console.log(res3)
                    if (res3.Result == 10000) {
                      commit('SET_USER_DATA', res3.Data)
                      setToken(res3.Data.token)
                    } else {
                      //登陆失败处理
                    }
                  })
                },
                error: (err) => {
                  console.error(err)
                  //错误处理
                }
              })
            }
          })
        },
        error: (e) => {
          console.error(e)
        }
      })
    },
    //刷新页面时候重新拉取用户数据
    CheckUserLogin({
      commit
    }) {
      GetUserByToken().then(res => {
        if (res.Result == 10000) {
          commit('SET_USER_DATA', res.Data)
          setToken(res.Data.token)
        } else {
          //登陆失败处理
          store.dispatch("SetWXUserData")
        }
      })
    },
    //更新用户连接
    SetUserConn({
      commit
    }, conn) {
      commit('SET_USER_CONN', conn)
    },
  }
}

export default user