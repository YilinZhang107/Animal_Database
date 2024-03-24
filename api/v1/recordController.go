/*
* @Author: Oatmeal107
* @Date:   2023/6/12 15:29
 */

package v1

import (
	"Animal_database/serializer"
	"Animal_database/service"
	"Animal_database/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetRecord 获取现存记录
func GetRecord(c *gin.Context) {
	var recordService service.RecordService

	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.GetRecord(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get Record api:", err)
	}
}

// GetByArea 根据给定地区获取地图显示数据
func GetByArea(c *gin.Context) {
	var recordService service.RecordService
	// 无需鉴权
	if err := c.ShouldBind(&recordService); err == nil {
		response := recordService.GetByArea()
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get Record api:", err)
	}
}

// GetRecordCount 获取现存记录数量
func GetRecordCount(c *gin.Context) {
	var recordService service.RecordService
	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.GetRecordCount(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get Record api:", err)
	}
}

// DeleteRecord 删除记录
func DeleteRecord(c *gin.Context) {
	var recordService service.RecordService
	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.DeleteRecord(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get Record api:", err)
	}
}

// Download 下载记录
func Download(c *gin.Context) {
	var recordService service.RecordService
	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID
	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.Download(id, c)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get Record api:", err)
	}
}
