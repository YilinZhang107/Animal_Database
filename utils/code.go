/*
* @Author: Oatmeal107
* @Date:   2023/6/15 17:38
 */

package utils

// 定义状态码
const (
	SUCCESS = 200
	ERROR   = 500

	//user
	ErrorUserExist     = 10001
	ErrorPasswordWrong = 10002
	ErrorUploadAvatar  = 10003
	ErrorGetUser       = 10004
	UserNotExist       = 10005
	ErrorNoPermission  = 10006
	ErrorGrade         = 10007
	//数据库错误
	ErrorDatabase            = 40001
	ErrorGetUnreviewedRecord = 40006
	DatabaseExistError       = 40007
	ErrorGetRecordByArea     = 40008
	//记录
	OpenUploadFileErr    = 20001
	CreateTempFileErr    = 20002
	CopyFileErr          = 20003
	OpenTempFileErr      = 20004
	UserGradeErr         = 20005
	ExcelFormatCheckErr  = 20006
	StringToIntErr       = 20007
	StringToFloatErr     = 20008
	StringToDateErr      = 20009
	StringToTimeErr      = 20010
	StringToUintErr      = 20011
	ErrorGetURecordByIds = 20012
	ErrorDelURecordByIds = 20013
	//ErrorHashPassword
	ErrorHashPassword = 40002
	//jwt
	ErrorGenerateToken         = 40003
	ErrorAuthCheckTokenFail    = 40004
	ErrorAuthCheckTokenTimeout = 40005
)

var msgs = map[int]string{
	SUCCESS: "ok",
	ERROR:   "fail",

	//user
	ErrorUserExist:     "用户名已存在",
	ErrorPasswordWrong: "密码错误",
	ErrorUploadAvatar:  "头像上传错误",
	ErrorGetUser:       "获取用户错误",
	UserNotExist:       "用户不存在",
	ErrorNoPermission:  "用户无权限",
	ErrorGrade:         "需要修改的权限错误",
	//数据库错误
	ErrorDatabase:            "数据库错误",
	ErrorGetUnreviewedRecord: "获取待审批记录错误",
	DatabaseExistError:       "数据库中数据已存在",
	ErrorGetRecordByArea:     "按地区获取记录错误",
	//记录
	OpenUploadFileErr:    "打开上传文件错误",
	CreateTempFileErr:    "创建临时文件错误",
	CopyFileErr:          "拷贝文件错误",
	OpenTempFileErr:      "打开临时文件错误",
	UserGradeErr:         "用户权限不足",
	ExcelFormatCheckErr:  "Excel格式检查错误",
	StringToIntErr:       "字符串转换为整数错误",
	StringToFloatErr:     "字符串转换为浮点数错误",
	StringToDateErr:      "字符串转换为日期错误",
	StringToTimeErr:      "字符串转换为时间错误",
	StringToUintErr:      "字符串转换为无符号整数错误",
	ErrorGetURecordByIds: "获取记录错误",
	ErrorDelURecordByIds: "删除记录错误",

	//
	ErrorHashPassword: "密码加密错误",
	//jwt
	ErrorGenerateToken:         "token生成错误",
	ErrorAuthCheckTokenFail:    "token验证失败",
	ErrorAuthCheckTokenTimeout: "token已过期",
}

func GetMsg(code int) string {
	msg, ok := msgs[code]
	if ok {
		return msg
	}
	return msgs[ERROR]
}
