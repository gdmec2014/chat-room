import fly from '../utils/request'
import {
  Service
} from './index'

/**
 * 用户注册
 * @param {*} data 
 */
export function Register(body) {
  return fly.request({
    url: Service.Register,
    method: 'POST',
    body
  })
}

/**
 * 用户登录
 * @param {*} data 
 */
export function Login(body) {
  return fly.request({
    url: Service.Login,
    method: 'POST',
    body
  })
}

/**
 * 获取用户信息
 */
export function GetUserByToken() {
  return fly.request({
    url: Service.GetUserByToken,
    method: 'GET'
  })
}