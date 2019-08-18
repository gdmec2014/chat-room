import request from '../util/request'
import {
  Service
} from './index'

/**
 * 用户注册
 * @param {*} data 
 */
export function Register(data) {
  return request({
    url: Service.Register,
    method: 'POST',
    data
  })
}

/**
 * 用户登录
 * @param {*} data 
 */
export function Login(data) {
  return request({
    url: Service.Login,
    method: 'POST',
    data
  })
}
