import fly from '../utils/request'
import {
    Service
} from './index'

import {getToken} from "../utils/auto"

/**
 * 获取所有的房间   
 * @param {*} body 
 */
export function GetAllRoom() {
    return fly.request({
        method: "get", //post/get 请求方式
        url: Service.GetAllRoom + getToken()
    })
}