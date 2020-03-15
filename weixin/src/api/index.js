export const Code = {
  //成功
  SUCCESS: 10000, //请求成功
  //参数错误
  PARAMETER_ERROR: 20000, //20000 参数错误
  DEFAULT_FAIELD: 20001, //20001 默认参数不能修改
  EXIST_FAILED: 20002, //20002 数据已存在
  REPASSWORD_FAIELD: 20003, //20003 两次密码不一致
  PASSWORD_ERROR: 20004, //20004 密码错误
  NOT_EXIST_FAILED: 20005, //20005 用户不存在
  TOKEN_ERROR: 20006, //20006 token错误
  //操作失败
  FAILED: 30000, //30000 操作失败
  LOGIN_EXPIRACTION: 30001, //30001 登录失效
  //系统错误
  UNKNOWE_ERROR: 40000, //40000 未知错误
  MARSHAL_ERROR: 40001, //40001 序列化错误
  UNMARSHAL_ERROR: 40002, //40002 反序列化错误
  //数据库错误
  SQL_ERROR: 50000, //数据库错误

  //webscoket
  EVENT_HAND: 10, //握手事件
  EVENT_CREATE: 11, //创房事件
  EVENT_JOIN: 12, //加房事件
  EVENT_LEAVE: 13, //离线事件
  EVENT_MESSAGE: 14, //消息事件
  EVENT_INVAILD: 15, //无效事件
  EVENT_DRAW: 16, //绘图事件
  EVENT_BREAK_DRAW: 17, //中断绘画事件
  EVENT_GIVE_IDENTITY: 18, //赋予游戏身份事件
  EVENT_NO_PLACE: 19, //房间满人了额
  EVENT_GAME_NO_START: 20, //还不能开始游戏事件
  EVENT_GAME_CAN_START: 21, //开始游戏事件
  EVENT_GAME_ANSWER: 22, //回答问题事件
  EVENT_GAME_BONUS: 23, //答对问题加分事件
  EVENT_GAME_TIME: 24, //游戏计时事件
  EVENT_GAME_OVER: 25, //游戏结束事件
  EVENT_GAME_IS_START: 26, //遊戲正在進行，不能重複開始
  EVENT_GAME_MEMBER_NOT_ENOUGH: 27, //人數不夠，不能開始遊戲
  EVENT_GAME_RE_START: 28, //重新開始遊戲
  EVENT_NEW_DRAW: 29, //新的绘图事件
  EVENT_SYSTEM_MESSAGE: 30, //系统消息

  VIEWER: 100, //观众
  PLAYER: 101, //猜一方
  MASTER: 102, //出题一方
  NO_MASTER: 103 //不可以再赋予出题一方，当全部玩家身份为103时候，游戏结束，结算分数
}

export const Service = {
  //TODO auth
  Register: '/v1/auth/register',
  Login: '/v1/auth/login',
  GetUserByToken: '/api/getuserbytoken',

  //TODO wexin
  WxLogin: '/v1/wx/login',
  WxRegist: '/v1/wx/regist',

  //TODO game
  GameRestart: '/v1/room/restart',
  GameConn: '/v1/room/hand',
  GetAllRoom: '/v1/room/get_all_room'
}
