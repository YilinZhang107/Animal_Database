/*
* @Author: Oatmeal107
* @Date:   2023/6/19 15:09
 */

package dao

import (
	"Animal_database/model"
	"strconv"
)

// DeleteURecordByIds 根据id数组删除记录
func DeleteURecordByIds(uRecords *[]model.UnreviewedRecord) error {
	return DB.Model(model.UnreviewedRecord{}).Unscoped().Delete(&uRecords).Error
}

// GetURecordByIds 根据id数组获取记录
func GetURecordByIds(ids []uint) (uRecords *[]model.UnreviewedRecord, err error) {
	err = DB.Model(model.UnreviewedRecord{}).Find(&uRecords, ids).Error //多条查询
	//根据ids数组查询, 这三个方法都行
	//err = DB.Model(model.UnreviewedRecord{}).Where("id in (?)", ids).Find(&uRecords).Error
	//根据id进行查询
	//for _, id := range ids {
	//	err = DB.Model(model.UnreviewedRecord{}).Where("id = ?", id).First(&uRecords).Error
	//}
	return uRecords, err
}

// GetUnreviewedRecord 分页获取待审批记录
func GetUnreviewedRecord(pageNum int, pageSize int) (uRecords *[]model.UnreviewedRecord, err error) {
	err = DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&uRecords).Error

	return uRecords, err
}

// UploadUnRecord 上传待审批记录
func UploadUnRecord(records *[]model.UnreviewedRecord) (returnErr error, existRecordsIndex []string) {
	for i, record := range *records {
		if err := DB.Model(model.UnreviewedRecord{}).Create(&record).Error; err != nil {
			returnErr = err
			existRecordsIndex = append(existRecordsIndex, strconv.Itoa(i+1))
		}
	}
	//return DB.Model(model.UnreviewedRecord{}).Create(&records).Error
	return returnErr, existRecordsIndex
}
