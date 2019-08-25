import Vue from 'vue'
import Vuex from 'vuex'

import user from './modules/user'
import room from './modules/room'

import getters from './getters'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    user,
    room
  },
  getters
})

export default store
