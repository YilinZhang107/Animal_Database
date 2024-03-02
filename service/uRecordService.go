/*
* @Author: Oatmeal107
* @Date:   2023/6/18 20:27
 */

package service

import (
	"Animal_database/dao"
	"Animal_database/model"
	"Animal_database/serializer"
	"Animal_database/utils"
	"github.com/tealeg/xlsx"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

type URecordService struct {
	Page            int    `form:"page" json:"page"`
	Size            int    `form:"size" json:"size"`
	ReviewedRecords string `form:"reviewedRecords" json:"reviewedRecords"`
}

// ReviewURecord 审批记录
func (r *URecordService) ReviewURecord(id uint) serializer.Response {
	code := utils.SUCCESS
	//先判断用户是否有权限
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.UserNotExist
		return serializer.CreateResponse(code, "未查到该用户", utils.GetMsg(code))
	}
	if user.Grade < 3 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权审批", utils.GetMsg(code))
	}
	//将传入的id数组转换为uint数组
	var ids []uint
	//把ReviewedRecords字符串按逗号分开
	idsStr := strings.Split(r.ReviewedRecords, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			code = utils.StringToUintErr
			return serializer.CreateResponse(code, nil, utils.GetMsg(code))
		}
		uid := uint(id)
		ids = append(ids, uid)
	}
	// 根据id数组获取记录
	uRecords, err := dao.GetURecordByIds(ids)
	if err != nil {
		code = utils.ErrorGetURecordByIds
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	// 如果找到的记录小于传入的id数组长度,说明有些记录不存在,返回错误
	if len(*uRecords) < len(ids) {
		code = utils.ErrorGetURecordByIds
		return serializer.CreateResponse(code, "选择的记录ID不存在或已审批,请重试", utils.GetMsg(code))
	}
	//删除这些记录,(在未审批表中改为已审批)
	err = dao.DeleteURecordByIds(uRecords)
	if err != nil {
		code = utils.ErrorDelURecordByIds
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	//加上审批人的id转换为Record
	records := make([]model.Record, 0)
	for _, uRecord := range *uRecords {
		record := uRecord.ToRecord(id)
		records = append(records, *record)
	}
	//将记录存在Record表中
	err, existRecords := dao.UploadRecord(&records) //重复的信息应该提醒已经存入了, 不重复的也能存进去
	if err != nil {
		if len(existRecords) != 0 {
			code = utils.DatabaseExistError // 存入重复的会提示 Duplicate entry '红隼-37.202710015699203-102.76243597269' for key 'once'
			data := "记录:" + strings.Join(existRecords, ",") + "已存在, 其余数据已录入"
			return serializer.CreateResponse(code, data, utils.GetMsg(code))
		} else {
			code = utils.ErrorDatabase
			return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
		}
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// DeleteURecord 删除记录
func (r *URecordService) DeleteURecord(id uint) serializer.Response {
	code := utils.SUCCESS
	//先判断用户是否有权限
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.UserNotExist
		return serializer.CreateResponse(code, "未查到该用户", utils.GetMsg(code))
	}
	if user.Grade < 2 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权删除", utils.GetMsg(code))
	}
	//将传入的id数组转换为uint数组
	var ids []uint
	//把ReviewedRecords字符串按逗号分开
	idsStr := strings.Split(r.ReviewedRecords, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			code = utils.StringToUintErr
			return serializer.CreateResponse(code, nil, utils.GetMsg(code))
		}
		uid := uint(id)
		ids = append(ids, uid)
	}
	// 根据id数组获取记录
	uRecords, err := dao.GetURecordByIds(ids)
	if err != nil {
		code = utils.ErrorGetURecordByIds
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	// 如果找到的记录小于传入的id数组长度,说明有些记录不存在,返回错误
	if len(*uRecords) < len(ids) {
		code = utils.ErrorGetURecordByIds
		return serializer.CreateResponse(code, "选择的记录ID不存在或已审批,请重试", utils.GetMsg(code))
	}
	//删除这些记录,(在未审批表中改为已审批)
	err = dao.DeleteURecordByIds(uRecords)
	if err != nil {
		code = utils.ErrorDelURecordByIds
		return serializer.CreateResponse(code, "删除失败", utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// GetUnreviewedRecord 获取待审批记录
func (r *URecordService) GetUnreviewedRecord(id uint) serializer.Response {
	code := utils.SUCCESS
	//
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Size == 0 {
		r.Size = 10 //默认显示10条
	}
	//先判断用户是否有权限
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.UserNotExist
		return serializer.CreateResponse(code, "未查到该用户", utils.GetMsg(code))
	}
	if user.Grade < 3 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权审批", utils.GetMsg(code))
	}
	uRecords, err := dao.GetUnreviewedRecord(r.Page, r.Size)
	if err != nil {
		code = utils.ErrorGetUnreviewedRecord
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, uRecords, utils.GetMsg(code))
}

// GetUnreviewedRecordCount 获取待审批记录数量
func (r *URecordService) GetUnreviewedRecordCount(id uint) serializer.Response {
	code := utils.SUCCESS
	//先判断用户是否有权限
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.UserNotExist
		return serializer.CreateResponse(code, "未查到该用户", utils.GetMsg(code))
	}
	if user.Grade < 3 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权审批", utils.GetMsg(code))
	}
	count, err := dao.GetUnreviewedRecordCount()
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, count, utils.GetMsg(code))
}

// UploadURecord 上传记录
func (service *URecordService) UploadURecord(Id uint, file *multipart.FileHeader) serializer.Response {
	code := utils.SUCCESS

	//还要根据用户的等级判断是直接到记录还是待审批记录
	user, err := dao.GetUserById(Id)
	if err != nil {
		code = utils.ErrorGetUser
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	if user.Grade == 1 { //权限为1的无法上传
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	} else { //2和3和4的可以上传
		return checkAndUploadByGrade(file, user.Grade, Id)
	}

}

// 在用户有权限时,检查格式并上传到待审查或者数据库
func checkAndUploadByGrade(file *multipart.FileHeader, grade uint8, id uint) serializer.Response {
	code := utils.SUCCESS
	// 打开 Excel 文件
	xlFile, err := file.Open()
	if err != nil {
		code = utils.OpenUploadFileErr
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	defer xlFile.Close()
	// 创建一个临时文件，将上传的文件保存到本地
	tempFile, err := os.CreateTemp("./resources/static/tmp", "upload-*.xlsx")
	if err != nil {
		code = utils.CreateTempFileErr
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	defer os.Remove(tempFile.Name())
	// 将上传的文件内容拷贝到临时文件
	_, err = io.Copy(tempFile, xlFile)
	if err != nil {
		code = utils.CopyFileErr
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	// 打开 Excel 文件
	excelFile, err := xlsx.OpenFile(tempFile.Name()) //这里的tempFile.Name()是临时文件的路径
	if err != nil {
		code = utils.OpenTempFileErr
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	//检查格式
	ok, data := utils.ExcelFormatCheck(excelFile) //todo 后续这个data可能改为需要在线预览并修改的数据
	if !ok {
		return serializer.CreateResponse(code, data, utils.GetMsg(code))
	}

	//把数据打包到待审批记录或者记录model
	var uRecords []model.UnreviewedRecord
	for j, row := range excelFile.Sheets[0].Rows[1:] {
		// 可能会有多的空行
		if row.Cells[0].String() != "样点" && row.Cells[0].String() != "样线" {
			break
		}
		uRecord := model.UnreviewedRecord{}
		var dateStr, timeStr string
		// 将每一列的值赋给相应的字段
		for i, cell := range row.Cells {
			if i > 24 {
				break
			} // todo 这里不应该写成固定值
			switch i {
			case 0:
				uRecord.InvestigationMethod = cell.String()
			case 1:
				//num, err := cell.Int()
				//if err != nil {
				//	code = utils.StringToIntErr
				//	return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(i), utils.GetMsg(code))
				//}
				//uRecord.GridNumber = uint64(num)
				uRecord.GridNumber = cell.String()
			case 2:
				uRecord.LineNumber = cell.String()
			case 3:
				uRecord.LineLength, err = cell.Float()
				if err != nil {
					code = utils.StringToFloatErr
					return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(i), utils.GetMsg(code))
				}
			case 4:
				num, err := cell.Int()
				if err != nil {
					code = utils.StringToIntErr
					return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(i), utils.GetMsg(code))
				}
				uRecord.GPSNumber = uint64(num)
			case 5:
				uRecord.SpeciesName = cell.String()
			case 6:
				num, err := cell.Int()
				if err != nil {
					code = utils.StringToIntErr
					return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(i), utils.GetMsg(code))
				}
				uRecord.Count = uint64(num)
			case 7:
				uRecord.Distance = cell.String()
			case 8:
				uRecord.PictureNumber = cell.String()
			case 9:
				uRecord.HabitatPicture = cell.String()
			case 10:
				uRecord.SpeciesThreatenedFactor = cell.String()
			case 11:
				uRecord.LivingEnvironmentType = cell.String()
			case 12:
				uRecord.EnvThreatenedFactor = cell.String()
			case 13:
				uRecord.ThreatIntensity = cell.String()
			case 14:
				uRecord.Weather = cell.String()
			case 15:
				uRecord.Investigator = cell.String()
			case 16:
				uRecord.InvestigationProvince = cell.String()
			case 17:
				uRecord.InvestigationCity = cell.String()
			case 18:
				uRecord.InvestigationCounty = cell.String()
			case 19:
				uRecord.InvestigationTown = cell.String()
			case 20:
				uRecord.Latitude = cell.String()
			case 21:
				uRecord.Longitude = cell.String()
			case 22:
				uRecord.Altitude = cell.String()
			case 23:
				dateStr = cell.String()
			case 24:
				timeStr = cell.String()
			}
		}
		// 解析日期字符串
		date, err := time.Parse("01-02-06", dateStr)
		if err != nil {
			code = utils.StringToDateErr
			return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(23), utils.GetMsg(code))
		}
		// 解析时间字符串
		t, err := time.Parse("15:04:05", timeStr)
		if err != nil {
			code = utils.StringToTimeErr
			return serializer.CreateResponse(code, "row:"+string(j)+"col:"+string(24), utils.GetMsg(code))
		}
		// 将日期和时间合并为一个时间对象
		uRecord.DateAndTime = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)

		uRecords = append(uRecords, uRecord)
	}

	//上传
	switch grade {
	case 2:
		//上传到待审批记录, 上传人的id也要记录
		for i := range uRecords {
			uRecords[i].Uploader = id
		}
		err, existURecordInd := dao.UploadUnRecord(&uRecords)
		if err != nil {
			if len(existURecordInd) != 0 {
				code = utils.DatabaseExistError
				data := "记录:" + strings.Join(existURecordInd, ",") + " 已存在, 其他记录已录入"
				return serializer.CreateResponse(code, data, utils.GetMsg(code))
			} else {
				code = utils.ErrorDatabase
				return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
			}
		}
		return serializer.CreateResponse(code, "记录已上传,请等待管理员审核", utils.GetMsg(code))
	case 3, 4:
		//上传到记录, 上传人和审核人的id一样
		var records []model.Record
		for _, uRecord := range uRecords {
			record := model.Record{
				InvestigationMethod:     uRecord.InvestigationMethod,
				GridNumber:              uRecord.GridNumber,
				LineNumber:              uRecord.LineNumber,
				LineLength:              uRecord.LineLength,
				GPSNumber:               uRecord.GPSNumber,
				SpeciesName:             uRecord.SpeciesName,
				Count:                   uRecord.Count,
				Distance:                uRecord.Distance,
				PictureNumber:           uRecord.PictureNumber,
				HabitatPicture:          uRecord.HabitatPicture,
				SpeciesThreatenedFactor: uRecord.SpeciesThreatenedFactor,
				LivingEnvironmentType:   uRecord.LivingEnvironmentType,
				EnvThreatenedFactor:     uRecord.EnvThreatenedFactor,
				ThreatIntensity:         uRecord.ThreatIntensity,
				Weather:                 uRecord.Weather,
				Investigator:            uRecord.Investigator,
				InvestigationProvince:   uRecord.InvestigationProvince,
				InvestigationCity:       uRecord.InvestigationCity,
				InvestigationCounty:     uRecord.InvestigationCounty,
				InvestigationTown:       uRecord.InvestigationTown,
				Latitude:                uRecord.Latitude,
				Longitude:               uRecord.Longitude,
				Altitude:                uRecord.Altitude,
				DateAndTime:             uRecord.DateAndTime,
				Uploader:                id,
				Auditor:                 id,
			}
			records = append(records, record)
		}
		err, existRecords := dao.UploadRecord(&records)
		if err != nil {
			if len(existRecords) != 0 {
				code = utils.DatabaseExistError // 存入重复的会提示 Duplicate entry '红隼-37.202710015699203-102.76243597269' for key 'once'
				data := "记录:" + strings.Join(existRecords, ",") + "已存在, 其余数据已录入"
				return serializer.CreateResponse(code, data, utils.GetMsg(code))
			} else {
				code = utils.ErrorDatabase
				return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
			}
		}
		return serializer.CreateResponse(code, "管理员上传记录成功", utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, "未知错误,未上传数据", utils.GetMsg(code))
}
