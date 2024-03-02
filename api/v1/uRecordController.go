/*
* @Author: Oatmeal107
* @Date:   2023/6/18 20:25
 */

package v1

import (
	"Animal_database/serializer"
	"Animal_database/service"
	"Animal_database/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReviewURecord 审核记录
func ReviewURecord(c *gin.Context) {
	var recordService service.URecordService

	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.ReviewURecord(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("review Record api:", err)
	}
}

// DeleteURecord 删除记录
func DeleteURecord(c *gin.Context) {
	var recordService service.URecordService

	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.DeleteURecord(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("delete Record api:", err)
	}
}

// GetUnreviewedRecord 获取等待审核的记录
func GetUnreviewedRecord(c *gin.Context) {
	var recordService service.URecordService

	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.GetUnreviewedRecord(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get unreviewed Record api:", err)
	}

}

// GetUnreviewedRecordCount 获取等待审核的记录数量
func GetUnreviewedRecordCount(c *gin.Context) {
	var recordService service.URecordService
	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	if err = c.ShouldBind(&recordService); err == nil {
		response := recordService.GetUnreviewedRecordCount(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get unreviewed Record count api:", err)
	}
}

// UploadRecord 上传记录
func UploadRecord(c *gin.Context) {
	var uRecordService service.URecordService

	//鉴权
	claims, err := utils.VerifyToken(c.GetHeader("Authorization"))
	id := claims.ID

	//获取上传的excel文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("get upload excel file err:", err)
	}
	// 检查文件类型，确保是 Excel 文件
	if file.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		c.JSON(http.StatusBadRequest, "Invalid file type")
	}

	if err := c.ShouldBind(&uRecordService); err == nil {
		response := uRecordService.UploadURecord(id, file)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, serializer.CreateErrResponse(err))
		utils.LogrusObj.Infoln("upload Record api:", err)
	}

}
