import Vue from 'vue'
import FastClick from 'fastclick'
import VueRouter from 'vue-router'
import App from './App'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import {Code} from './api/index'
import {
  getToken
} from './util/auto'

Vue.use(VueRouter)
Vue.use(ElementUI);

FastClick.attach(document.body)

Vue.config.productionTip = false
Vue.prototype.Code = Code

router.beforeEach((to, from, next) => {
  if(to.meta.requireAuth) {
    if(getToken()) {
      next()
    } else {
      next({
        path: '/login'
      })
    }
  } else {
    next();
  }
});

/* eslint-disable no-new */
new Vue({
  router,
  render: h => h(App)
}).$mount('#app-box')
