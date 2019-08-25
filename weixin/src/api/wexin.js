import fly from '../utils/request'
import {
    Service
} from './index'

/**
 * 微信登录
 * @param {*} body 
 */
export function WxLogin(body) {
    return fly.request({
        method: "post", //post/get 请求方式
        url: Service.WxLogin,
        body
    })
}

/**
 * 微信注册
 * @param {*} body 
 */
export function WxRegist(body) {
    return fly.request({
        method: "post", //post/get 请求方式
        url: Service.WxRegist,
        body
    })
}