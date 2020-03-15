import Vue from 'vue'
import MpvueRouterPatch from 'mpvue-router-patch'
import App from '@/App'
import store from '@/store'
import WebScoket from '@/components/webscoket'
import { Code } from './api/index'
import RenderCanvas from 'vnode2canvas'

Vue.use(RenderCanvas)
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

export default {
  config: {
    // 这儿添加要用的小程序组件
    usingComponents: {
      'canvasdrawer': '../static/canvasdrawer/canvasdrawer'
    }
  }
}