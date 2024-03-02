/*
* @Author: Oatmeal107
* @Date:   2023/6/19 18:03
 */

package service

import (
	"Animal_database/config"
	"Animal_database/dao"
	"Animal_database/model"
	"Animal_database/serializer"
	"Animal_database/utils"
	"strconv"
	"strings"
)

type RecordService struct {
	Page            int    `form:"page" json:"page"`
	Size            int    `form:"size" json:"size"`
	Province        string `form:"province" json:"province"`
	City            string `form:"city" json:"city"`
	ReviewedRecords string `form:"reviewedRecords" json:"reviewedRecords"`
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
	// todo
	if user.Grade == 0 {
		code = utils.UserGradeErr
		return serializer.CreateResponse(code, "无权查看", utils.GetMsg(code))
	}
	uRecords, err := dao.GetRecord(r.Page, r.Size)
	if err != nil {
		code = utils.ErrorGetUnreviewedRecord
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
	}
	return serializer.CreateResponse(code, uRecords, utils.GetMsg(code))
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
	count, err := dao.GetRecordCount()
	if err != nil {
		code = utils.ErrorDatabase
		return serializer.CreateResponse(code, nil, utils.GetMsg(code))
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
