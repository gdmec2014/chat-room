import Vue from 'vue'
import Router from 'vue-router'
import Register from '@/views/auth/register'
import Login from '@/views/auth/login'
import Home from '@/views/home/index'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/register',
      name: 'Register',
      component: Register
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/home',
      name: 'Home',
      component: Home,
      meta: {
        requireAuth: true,
      }
    }
  ]
})
