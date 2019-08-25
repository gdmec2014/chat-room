import Vue from 'vue'
import MpvueRouterPatch from 'mpvue-router-patch'
import App from '@/App'
import store from '@/store'
import WebScoket from '@/components/webscoket'
import {Code} from './api/index'

Vue.use(WebScoket)
Vue.use(MpvueRouterPatch)
Vue.config.productionTip = false
Vue.prototype.Code = Code

const app = new Vue({
  mpType: 'app',
  store,
  ...App
})
app.$mount()
