//格式化go的时间
export function FormatGoTime(t) {
  var dateee = new Date(t).toJSON();
  var date = new Date(+new Date(dateee) + 3600000 * 8).toISOString().replace(/T/g, ' ').replace(/.[\d]{3}Z/, '')
  return date
}

//剪切字符
export function GetStringSub(t, l) {
  if (t.length > l) {
    return t.substring(0, l) + "..."
  } else {
    return t
  }
}

//获取当前时间
export function GetNowDate() {
  return GetDate((new Date()).getTime())
}

//获取当前配置
export function GetConfig() {
  let config = {
    base_url: "http://127.0.0.1:2332"
  }
  let is_dev = true
  if (process.env.NODE_ENV) {
    if (process.env.NODE_ENV != "development") {
      is_dev = false
    }
  } else {
    is_dev = false
  }
  if (!is_dev) {
    config = {
      base_url: "https://www.haibarai.com"
    }
  }
  return config
}