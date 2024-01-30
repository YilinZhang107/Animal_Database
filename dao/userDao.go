/*
* @Author: Oatmeal107
* @Date:   2023/6/12 15:32
 */

package dao

import (
	"Animal_database/model"
)

type UserDao struct {
}

// // NewUserDao 初始化UserDao:创建一个新的数据库连接
//
//	func NewUserDao(c context.Context) UserDao {
//		ans := new(UserDao)
//		ans.DB = Init.NewDBClient(c)
//		return UserDao{
//			DB: Init.NewDBClient(c),
//		}
//	}

func UpdateUserById(Id uint, user *model.User) error {
	return DB.Model(model.User{}).Where("id = ?", Id).Updates(&user).Error
}

func GetUserById(Id uint) (user *model.User, err error) {
	err = DB.Model(model.User{}).Where("id = ?", Id).First(&user).Error
	return user, err
}

// GetUserByName 通过用户名获取用户信息
func GetUserByName(username string) (user *model.User, err error) {
	err = DB.Model(model.User{}).Where("user_name=?", username).First(&user).Error
	return user, err
}

// ExistOrNotByUserName 判断用户名是否存在
func ExistOrNotByUserName(username string) (exist bool, err error) {
	//var user model.User
	var cnt int64
	err = DB.Model(model.User{}).Where("user_name = ?", username).Count(&cnt).Error
	if cnt > 0 {
		return true, err
	}
	return false, err
}

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return DB.Create(&user).Error
}

//// GetUserByGrade 根据用户等级查找用户
//func GetUserByGrade(grade uint) (users *[]model.User, err error) {
//	err = DB.Model(model.User{}).Where("grade = ?", grade).Find(&users).Error
//	return users, err
//}

// GetUsers 获取所有用户
func GetUsers(pageSize int, pageNum int) (users *[]model.User, err error) {
	err = DB.Model(model.User{}).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	//err = DB.Model(model.User{}).Find(&users).Error
	return users, err
}

// DeleteUserById 根据用户id删除用户
func DeleteUserById(id uint) error {
	return DB.Model(model.User{}).Unscoped().Delete(&model.User{}, id).Error //硬删除
}
