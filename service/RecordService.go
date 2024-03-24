/*
* @Author: Oatmeal107
* @Date:   2023/6/19 18:03
 */

package service

import (
	"Animal_database/VO"
	"Animal_database/config"
	"Animal_database/dao"
	"Animal_database/model"
	"Animal_database/serializer"
	"Animal_database/utils"
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type RecordService struct {
	Page                  int    `form:"page" json:"page"`
	Size                  int    `form:"size" json:"size"`
	Province              string `form:"province" json:"province"`
	City                  string `form:"city" json:"city"`
	County                string `form:"county" json:"county"`
	ReviewedRecords       string `form:"reviewedRecords" json:"reviewedRecords"`
	GridNumber            string `form:"gridNumber" json:"gridNumber"` //样线样点模糊匹配
	LineNumber            string `form:"lineNumber" json:"lineNumber"`
	StartTime             string `form:"startTime" json:"startTime"` // todo 时间格式
	EndTime               string `form:"endTime" json:"endTime"`
	SpeciesName           string `form:"speciesName" json:"speciesName"`                     //物种名称精确匹配
	Investigator          string `form:"investigator" json:"investigator"`                   //调查人精确匹配
	LivingEnvironmentType string `form:"livingEnvironmentType" json:"livingEnvironmentType"` //生境类型模糊匹配
}

// GetRecord 获取现有记录 todo 可能还需要根据权限修改能看到的数据详情
func (r *RecordService) GetRecord(id uint) serializer.Response {
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
	// todo 只有高等级用户能下载
	if user.Grade == 0 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权查看", utils.GetMsg(code))
	}

	records, data, err := dao.GetRecord(r.Page, r.Size, r.GridNumber, r.LineNumber,
		r.StartTime, r.EndTime, r.Province, r.City, r.County, r.SpeciesName, r.Investigator, r.LivingEnvironmentType)
	if err != nil || data != nil {
		code = utils.ErrorGetRecordByIds
		return serializer.CreateResponse(code, data, utils.GetMsg(code))
	}
	var VOdata []VO.RecordVO
	for _, record := range *records {
		VOdata = append(VOdata, *VO.BuildRecordVO(&record))
	}
	return serializer.CreateResponse(code, VOdata, utils.GetMsg(code))
}

func (r *RecordService) GetByArea() serializer.Response {
	code := utils.SUCCESS
	//创建一个map，key为省份，value为该省份的记录
	eachRecord := make(map[string]interface{})
	var err error
	// 用于统计数据
	StatisticalData := func(oneArea map[string]*[]model.Record, areaName string) {
		recordNum := len(*oneArea[areaName])
		animalMap := make(map[string]bool)
		for _, record := range *oneArea[areaName] {
			animalMap[record.SpeciesName] = true
		}
		animalNum := len(animalMap)
		eachRecord[areaName] = map[string]int{"recordNum": recordNum, "speciesNum": animalNum}
	}

	if r.Province == "" {
		//返回全国各省
		for _, province := range config.Provinces {
			oneProv := make(map[string]*[]model.Record)
			oneProv[province], err = dao.GetRecordByArea(province, "", "")
			if err != nil {
				code = utils.ErrorGetRecordByArea
				return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
			}
			// 统计数据
			StatisticalData(oneProv, province)
		}
	} else if r.City == "" {
		//返回该省各市
		citys := config.Prov2city[r.Province]
		if citys == nil {
			code = utils.ErrorGetRecordByArea
			return serializer.CreateResponse(code, "省份名称错误", utils.GetMsg(code))
		}
		for _, city := range citys {
			oneCity := make(map[string]*[]model.Record)
			oneCity[city], err = dao.GetRecordByArea(r.Province, city, "")
			if err != nil {
				code = utils.ErrorGetRecordByArea
				return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
			}
			// 统计数据
			StatisticalData(oneCity, city)
		}
	} else {
		//返回该市各县
		countys := config.City2county[r.City]
		if countys == nil {
			code = utils.ErrorGetRecordByArea
			return serializer.CreateResponse(code, "省份或城市名称错误", utils.GetMsg(code))
		}
		for _, county := range countys {
			oneCounty := make(map[string]*[]model.Record)
			oneCounty[county], err = dao.GetRecordByArea(r.Province, r.City, county)
			if err != nil {
				code = utils.ErrorGetRecordByArea
				return serializer.CreateResponse(code, err.Error(), utils.GetMsg(code))
			}
			// 统计数据
			StatisticalData(oneCounty, county)
		}
	}
	return serializer.CreateResponse(code, eachRecord, utils.GetMsg(code))
}

// GetRecordCount 获取现存记录数量
func (r *RecordService) GetRecordCount(id uint) serializer.Response {
	code := utils.SUCCESS

	count, data, err := dao.GetRecordCount(r.GridNumber, r.LineNumber,
		r.StartTime, r.EndTime, r.Province, r.City, r.County, r.SpeciesName, r.Investigator, r.LivingEnvironmentType)
	if err != nil || data != nil {
		code = utils.ErrorGetRecordByIds
		return serializer.CreateResponse(code, data, utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, count, utils.GetMsg(code))
}

// DeleteRecord 删除记录
func (r *RecordService) DeleteRecord(id uint) serializer.Response {
	code := utils.SUCCESS
	//先判断用户是否有权限
	user, err := dao.GetUserById(id)
	if err != nil {
		code = utils.UserNotExist
		return serializer.CreateResponse(code, "未查到该用户", utils.GetMsg(code))
	}
	if user.Grade < 3 {
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
	records, err := dao.GetRecordByIds(ids)
	if err != nil {
		code = utils.ErrorGetRecordByIds
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	// 如果找到的记录小于传入的id数组长度,说明有些记录不存在,返回错误
	if len(*records) < len(ids) {
		code = utils.ErrorGetURecordByIds
		return serializer.CreateResponse(code, "选择的记录ID不存在或已删除,请重试", utils.GetMsg(code))
	}
	//删除这些记录
	err = dao.DeleteRecordByIds(records)
	if err != nil {
		code = utils.ErrorDelURecordByIds
		return serializer.CreateResponse(code, "删除失败", utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

// Download 下载记录
func (r *RecordService) Download(id uint, c *gin.Context) serializer.Response {
	code := utils.SUCCESS
	// todo 是否需要根据权限修改下载范围
	// 根据传入的参数获取记录,并生成csv文件
	records, data, err := dao.GetRecordAll(r.GridNumber, r.LineNumber,
		r.StartTime, r.EndTime, r.Province, r.City, r.County, r.SpeciesName, r.Investigator, r.LivingEnvironmentType)
	if err != nil || data != nil {
		code = utils.ErrorGetRecordByIds
		return serializer.CreateResponse(code, data, utils.GetMsg(code))
	}
	// 根据当前时间生成文件名
	fileName := fmt.Sprintf("Records_%s.csv", time.Now().Format("20060102150405"))
	path := "./resources/static/tmp/download/" + fileName
	err = writeRecordsToCSV(records, path)
	if err != nil {
		code = utils.ErrorWriteFile
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}

	file, err := os.Open(path)
	if err != nil {
		code = utils.ErrorOpenFile
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	defer file.Close()
	//fileName := filepath.Base(path)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.File(file.Name())
	//// 删除文件  todo 暂时没做,需要前端下载后删除
	//err = os.Remove(path)
	//if err != nil {
	//	code = utils.ErrorDelFile
	//	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	//}
	return serializer.CreateResponse(code, nil, utils.GetMsg(code))
}

func writeRecordsToCSV(records *[]model.Record, path string) error {
	// 创建一个带UTF-8 BOM的writer
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入UTF-8 BOM
	if _, err := file.Write(utf8bom); err != nil {
		return err
	}
	//file, err := os.Create(path)
	//if err != nil {
	//	return fmt.Errorf("failed to create file: %w", err)
	//}
	//defer file.Close()
	//writer := csv.NewWriter(file)
	//defer writer.Flush()

	utf8Writer := csv.NewWriter(bufio.NewWriter(io.MultiWriter(file)))
	// 定义CSV列标题（与Record结构体字段对应）
	cols := []string{
		"ID",
		"调查方式",
		"网格编号",
		"样线编号",
		"样线长度",
		"GPS编号",
		"物种名称",
		"数量",
		"距离(m)",
		"照片编号",
		"栖息地照片",
		"物种受威胁因素",
		"生境类型",
		"生境受威胁因素",
		"受威胁强度",
		"天气",
		"调查人",
		"调查省",
		"调查市",
		"调查县",
		"调查乡",
		"维度(Lat)",
		"经度(Lon)",
		"海拔",
		"日期与时间",
		"上传者编号",
		"审核者编号",
	}
	// 写入列标题
	if err := utf8Writer.Write(cols); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	// 遍历记录并写入行数据
	for _, record := range *records {
		row := []string{
			fmt.Sprintf("%d", record.ID), // 注意转换为字符串
			record.InvestigationMethod,
			record.GridNumber,
			record.LineNumber,
			fmt.Sprintf("%.2f", record.LineLength), // 注意转换为字符串并限制精度
			fmt.Sprintf("%d", record.GPSNumber),
			record.SpeciesName,
			fmt.Sprintf("%d", record.Count),
			record.Distance,
			record.PictureNumber,
			record.HabitatPicture,
			record.SpeciesThreatenedFactor,
			record.LivingEnvironmentType,
			record.EnvThreatenedFactor,
			record.ThreatIntensity,
			record.Weather,
			record.Investigator,
			record.InvestigationProvince,
			record.InvestigationCity,
			record.InvestigationCounty,
			record.InvestigationTown,
			record.Latitude,
			record.Longitude,
			record.Altitude,
			record.DateAndTime.Format("2006-01-02 15:04:05"), // 格式化日期时间
			fmt.Sprintf("%d", record.Uploader),
			fmt.Sprintf("%d", record.Auditor),
		}

		if err := utf8Writer.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}
	utf8Writer.Flush()
	return nil
}
