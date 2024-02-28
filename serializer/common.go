/*
* @Author: Oatmeal107
* @Date:   2023/6/15 16:20
 */

package serializer

// Response 基础响应结构体
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"` // 用于前端提示信息,信息或错误都在这里
}

// CreateResponse 创建响应, 用于service层返回响应
func CreateResponse(status int, data interface{}, msg string) Response {
	return Response{
		Status: status,
		Data:   data,
		Msg:    msg,
	}
}

func CreateErrResponse(err error) Response {
	return Response{
		Status: 400,
		Data:   nil,
		Msg:    err.Error(),
	}
}

// TokenData 带有token的Data结构
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
