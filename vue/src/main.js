import Vue from 'vue'
import FastClick from 'fastclick'
import VueRouter from 'vue-router'
import App from './App'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import {Code} from './api/index'

Vue.use(VueRouter)
Vue.use(ElementUI);

FastClick.attach(document.body)

Vue.config.productionTip = false
Vue.prototype.Code = Code

/* eslint-disable no-new */
new Vue({
  router,
  render: h => h(App)
}).$mount('#app-box')
