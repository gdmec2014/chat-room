import Fly from 'flyio/dist/npm/wx'
import {
    getToken
} from '@/utils/auto'


const fly = new Fly()
const host = "http://127.0.0.1:2332"
//添加请求拦截器
fly.interceptors.request.use((request) => {
    wx.showLoading({
        title: "加载中",
        mask: true
    });
    request.headers['content-type'] = 'application/json';
    let token = getToken()
    if (token) {
        request.headers['Authorization'] = token
    }
    request.body && Object.keys(request.body).forEach((val) => {
        if (request.body[val] === "") {
            delete request.body[val]
        };
    });
    request.body = {
        ...request.body,
    }
    return request;
});

//添加响应拦截器
fly.interceptors.response.use(
    (response) => {
        wx.hideLoading();
        return response.data;//请求成功之后将返回值返回
    },
    (err) => {
        //请求出错，根据返回状态码判断出错原因
        console.log(err);
        wx.hideLoading();
        wx.showToast({
            title: err.message,
            icon: 'none',
            duration: 2000
        })
        if (err) {
            return "请求失败";
        };
    }
);

fly.config.baseURL = host;

export default fly;