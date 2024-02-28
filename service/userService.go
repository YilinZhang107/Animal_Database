/*
* @Author: Oatmeal107
* @Date:   2023/6/12 17:01
 */

package service

import (
	"Animal_database/VO"
	"Animal_database/dao"
	"Animal_database/model"
	"Animal_database/serializer"
	"Animal_database/utils"
	"context"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type UserService struct {
	ID       uint   `form:"id" json:"id"`
	UserName string `form:"username" json:"username"` // form: 表单数据, json: json数据
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
	Grade    uint8  `form:"grade" json:"grade"`
	Page     int    `form:"page" json:"page"`
	Size     int    `form:"size" json:"size"`
}

// Register 用户注册
func (u *UserService) Register(c context.Context) serializer.Response {
	code := utils.SUCCESS

	exist, err := dao.ExistOrNotByUserName(u.UserName) // 判断用户名是否存在
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	if exist {
		code = utils.ErrorUserExist
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}

	u.Password, err = utils.HashPassword(u.Password) // 密码加密
	if err != nil {
		code = utils.ErrorHashPassword
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	user := model.User{
		UserName: u.UserName,
		Password: u.Password,
		Email:    u.Email,
		Grade:    1, // 默认为1级用户
		Avatar:   "avatar.JPG",
	}
	err = dao.CreateUser(&user) // 在数据库中创建用户
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}

	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// Login 用户登录
func (u *UserService) Login(c context.Context) serializer.Response {
	code := utils.SUCCESS

	user, err := dao.GetUserByName(u.UserName) // 通过用户名获取用户信息
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}

	if utils.ComparePassword(user.Password, u.Password) != nil { // 密码校验,失败返回err
		code = utils.ErrorPasswordWrong
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}

	token, err := utils.GenerateToken(user.ID, user.UserName, 0) // 生成token
	if err != nil {
		logrus.Infoln(err)
		code = utils.ErrorGenerateToken
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	//不序列化用户的话,带有密码的信息会被返回
	return serializer.CreateResponse(
		code,
		serializer.TokenData{User: VO.BuildUserVO(user), Token: token},
		utils.GetMsg(code),
	)
}

// SelfInfo 获取用户自己的信息
func (u *UserService) SelfInfo(c context.Context, id uint) serializer.Response {
	code := utils.SUCCESS
	user, err := dao.GetUserById(id) // 通过用户名获取用户信息
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, VO.BuildUserVO(user), utils.GetMsg(code))
}

// UpdateAvatar 更新用户头像
func (u *UserService) UpdateAvatar(c context.Context, id uint, file multipart.File, fileName string) serializer.Response {
	code := utils.SUCCESS
	user, err := dao.GetUserById(id) // 通过用户名获取用户信息
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}

	path, err := utils.UploadAvatarLocal(file, fileName, id) // 上传头像
	if err != nil {
		code = utils.ErrorUploadAvatar
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	user.Avatar = path
	err = dao.UpdateUserById(id, user)
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, path, utils.GetMsg(code))
}

// ChangePassword 修改密码
func (u *UserService) ChangePassword(c context.Context, id uint) serializer.Response {
	code := utils.SUCCESS
	user, err := dao.GetUserById(id) // 通过用户名获取用户信息
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}

	user.Password, err = utils.HashPassword(u.Password) // 密码加密
	if err != nil {
		code = utils.ErrorHashPassword
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}

	err = dao.UpdateUserById(id, user)
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// UpdateEmail 更新邮箱(和修改等级)
func (u *UserService) UpdateEmail(c context.Context, id uint, function uint8) serializer.Response {
	code := utils.SUCCESS
	if function == 0 {
		user, err := dao.GetUserById(id) // 通过用户名获取用户信息
		if err != nil {
			code = utils.ErrorGetUser
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
		user.Email = u.Email
		err = dao.UpdateUserById(id, user)
		if err != nil {
			code = utils.ErrorDatabase
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
	} else if function == 1 {
		user, err := dao.GetUserById(id) // 通过用户名获取用户信息
		if err != nil {
			code = utils.ErrorGetUser
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
		if user.Grade != 4 {
			code = utils.ErrorNoPermission
			return serializer.CreateResponse(code, nil, utils.GetMsg(code))
		}
		user, err = dao.GetUserById(u.ID)
		if err != nil {
			code = utils.ErrorGetUser
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
		if u.Grade != 1 && u.Grade != 2 && u.Grade != 3 && u.Grade != 4 {
			code = utils.ErrorGrade
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
		user.Grade = u.Grade
		err = dao.UpdateUserById(u.ID, user)
		if err != nil {
			code = utils.ErrorDatabase
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// GetUserCount 获取用户数量
func (u *UserService) GetUserCount(c context.Context, id uint) serializer.Response {
	code := utils.SUCCESS
	Boss, err := dao.GetUserById(id)
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	if Boss.Grade == 0 { //没太想好,暂时让等级为1及以上的都能查到
		code = utils.ErrorNoPermission
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	count, err := dao.GetUserCount()
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, count, utils.GetMsg(code))
}

// FindUsers 按条件查找用户 todo 目前只有等级为4的才能查
func (u *UserService) FindUsers(c context.Context, id uint) serializer.Response {
	code := utils.SUCCESS
	Boss, err := dao.GetUserById(id)
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	if Boss.Grade != 4 {
		code = utils.ErrorNoPermission
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	if u.UserName != "" {
		user, err := dao.GetUserByName(u.UserName)
		if err != nil {
			code = utils.ErrorDatabase
			return serializer.CreateResponse(code, nil, utils.GetMsg(code))
		}
		return serializer.CreateResponse(code, *VO.BuildUserVO(user), utils.GetMsg(code))
	} else {
		users, err := dao.GetUsers(u.Size, u.Page)
		if err != nil {
			code = utils.ErrorDatabase
			return serializer.CreateResponse(code, nil, utils.GetMsg(code))
		}
		var data []VO.UserVO
		for _, user := range *users {
			data = append(data, *VO.BuildUserVO(&user))
		}
		return serializer.CreateResponse(code, data, utils.GetMsg(code))
	}
}

// DeleteUser 删除用户
func (u *UserService) DeleteUser(c context.Context, id uint) serializer.Response {
	code := utils.SUCCESS
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	if user.Grade != 4 {
		code = utils.ErrorNoPermission
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	user, err = dao.GetUserById(u.ID)
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, "该用户不存在", utils.GetMsg(code))
	}
	err = dao.DeleteUserById(u.ID)
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}
