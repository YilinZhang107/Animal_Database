/*
* @Author: Oatmeal107
* @Date:   2023/6/16 15:51
 */

package VO

import "Animal_database/model"

// UserVO 用户序列化器,用于返回给前端的用户信息
type UserVO struct {
	ID       uint   `json:"id"` //加了json的tag, 会以这个名字为准, 不加的话会以变量名为准
	UserName string `json:"user_name"`
	Grade    uint8  `json:"grade"` //用户等级, 1仅能读取,2能上传, 3审核, 4管理员
	Email    string `json:"email"`
	Avatar   string `json:"avatar"` //头像
	CreateAt int64  `json:"create_at"`
}

func BuildUserVO(u *model.User) *UserVO {
	return &UserVO{
		ID:       u.ID,
		UserName: u.UserName,
		Grade:    u.Grade,
		Email:    u.Email,
		Avatar:   u.Avatar,
		CreateAt: u.CreatedAt.Unix(),
	}
}
