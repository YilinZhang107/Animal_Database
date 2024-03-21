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
	"strconv"
	"strings"
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
