import Vue from 'vue'
import Vuex from 'vuex'

import user from './modules/user'
import room from './modules/room'

import getters from './getters'

Vue.use(Vuex)

import createPersistedState from "vuex-persistedstate"

const store = new Vuex.Store({
  modules: {
    user,
    room
  },
  plugins: [
    createPersistedState({
      storage: {
        getItem: key => wx.getStorageSync(key),
        setItem: (key, value) => wx.setStorageSync(key, value),
        removeItem: key => null// wx.clearStorage()
      }
    })
  ],
  getters
})

export default store
