/*
* @Author: Oatmeal107
* @Date:   2023/6/12 11:05
 */

package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"` //用户名
	Grade    uint8  //用户等级, 1仅能读取,2能上传, 3审核, 4管理员
	Password string //密码,存前端加密后的
	Email    string //邮箱
	Avatar   string `gorm:"size:1000"` //头像
	//todo 是否还需要别的信息
}
