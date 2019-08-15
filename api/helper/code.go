package helper

type Status int

//成功
const SUCCESS Status = 10000

//失败
const FAILED Status = 30000

//数据库错误
const SQL_ERROR Status = 50000

//参数错误
const (
	_                 Status = iota + 19999
	PARAMETER_ERROR    //20000 参数错误
	DEFAULT_FIELD      //20001 默认参数不能修改
	EXIST_FAILED       //20002 数据已存在
	REPASSWORD_FAIELD  //20003 两次密码不一致
	PASSWORD_ERROR     //20004 密码错误
	NOT_EXIST_FAILED   //20005 用户不存在
)

//系统错误
const (
	_               Status = iota + 39999
	UNKNOWE_ERROR    //40000 未知错误
	MARSHAL_ERROR    //40001 序列化错误
	UNMARSHAL_ERROR  //40002 反序列化错误
)

type RestfulReturn struct {
	Result  Status      `json:"Result"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}