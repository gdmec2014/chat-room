const TokenName = 'Authorization'

export function getToken() {
    return wx.getStorageSync(TokenName)
}

export function setToken(token) {
    try {
        wx.setStorageSync(TokenName, token)
    } catch (e) {
        console.error(e)
    }  
}