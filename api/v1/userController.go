/*
* @Author: Oatmeal107
* @Date:   2023/6/12 16:50
 */

package v1

import (
	"Animal_database/serializer"
	"Animal_database/service"
	"Animal_database/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var userRegisterService service.UserService //相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。

	if err := c.ShouldBind(&userRegisterService); err == nil { // 将请求中的数据绑定到userRegisterService对象中
		res := userRegisterService.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user register api", err)
	}
}

// UserLogin 用户登陆
func UserLogin(c *gin.Context) {
	var userLoginService service.UserService
	if err := c.ShouldBind(&userLoginService); err == nil {
		res := userLoginService.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user Login err", err)
	}
}

// SelfInfo 用户获取自己的信息
func SelfInfo(c *gin.Context) {
	var userSelfInfoService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	res := userSelfInfoService.SelfInfo(c.Request.Context(), id)
	c.JSON(http.StatusOK, res)
}

//用户更新个人信息(由等级4用户来修改用户等级)

// UserUpdateAvatar 用户更新头像
func UserUpdateAvatar(c *gin.Context) {
	var userUpdateService service.UserService

	file, fileHeader, _ := c.Request.FormFile("file")
	fileName := fileHeader.Filename
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	//todo 需要做什么吗
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	//id := uint(idFloat.(float64))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		response := userUpdateService.UpdateAvatar(c.Request.Context(), id, file, fileName)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user update avatar err", err)
	}
}

// UserUpdatePassword 用户更新密码
func UserChangePassword(c *gin.Context) {
	var userUpdateService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&userUpdateService); err == nil {
		response := userUpdateService.ChangePassword(c.Request.Context(), id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user update password err", err)
	}
}

// UserUpdateEmail 用户更新邮箱
func UserUpdateEmail(c *gin.Context) {
	var userUpdateService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&userUpdateService); err == nil {
		response := userUpdateService.UpdateEmail(c.Request.Context(), id, 0)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user update email err", err)
	}
}

// FindUsers 按条件查找用户, 用于等级4用户调整等级( 目前写了按名字查找, 不写名字时会返回所有用户)
func FindUsers(c *gin.Context) {
	var userFindService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&userFindService); err == nil {
		response := userFindService.FindUsers(c.Request.Context(), id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("find users err", err)
	}
}

// GetUserCount 查看全部用户的数量
func GetUserCount(c *gin.Context) {
	var getUserCount service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get total user num err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&getUserCount); err == nil {
		response := getUserCount.GetUserCount(c.Request.Context(), id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("find users err", err)
	}

}

// UpdateUserGrade 更改用户等级
func UpdateUserGrade(c *gin.Context) {
	var userUpdateService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&userUpdateService); err == nil {
		response := userUpdateService.UpdateEmail(c.Request.Context(), id, 1) // 借用修改邮箱的函数
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user update email err", err)
	}
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var userDeleteService service.UserService
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("user verify err", err)
		return
	}
	id := claims.ID
	if err := c.ShouldBind(&userDeleteService); err == nil {
		response := userDeleteService.DeleteUser(c.Request.Context(), id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("delete user err", err)
	}
}
