import WebScoketComponent from "./index.vue"

const WebScoket = {};

// 注册WebScoket
WebScoket.install = function (Vue) {
    // 生成一个Vue的子类
    // 同时这个子类也就是组件
    const WebScoketConstructor = Vue.extend(WebScoketComponent)
    // 生成一个该子类的实例
    const instance = new WebScoketConstructor();

    // 通过Vue的原型注册一个方法
    // 让所有实例共享这个方法 
    Vue.prototype.$sendSocketMessage = (data) => {
        instance.sendSocketMessage(data)
    }
}

export default WebScoket