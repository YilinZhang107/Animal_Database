/*
* @Author: Oatmeal107
* @Date:   2023/6/19 15:09
 */

package dao

import (
	"Animal_database/model"
	"strconv"
	"strings"
	"time"
)

// GetRecord  分页获取记录
func GetRecord(pageNum int, pageSize int, gridNumber string, lineNumber string,
	startTime string, endTime string, province string, city string, county string,
	speciesName string, investigator string, livingEnvironmentType string) (records *[]model.Record, data interface{}, err error) {

	query := DB.Model(model.Record{})
	if gridNumber != "" {
		query = query.Where("grid_number LIKE ?", "%"+gridNumber+"%")
	}
	if lineNumber != "" {
		query = query.Where("line_number LIKE ?", "%"+lineNumber+"%")
	}
	if county != "" {
		query = query.Where("investigation_county = ?", county)
	} else if city != "" {
		query = query.Where("investigation_city = ?", city)
	} else if province != "" {
		query = query.Where("investigation_province = ?", province)
	}
	if startTime != "" || endTime != "" { //只要设置了一个时间必须设置另一个时间
		if startTime == "" {
			startTime = "631123200"
		}
		if endTime == "" {
			//更改endTime为当前时间的时间戳
			endTime = strconv.FormatInt(time.Now().Unix(), 10)
		}
		nowTime := strconv.FormatInt(time.Now().Unix(), 10)
		now, err := strconv.ParseInt(nowTime, 10, 64)

		// 字符串转为int64
		startTime, err := strconv.ParseInt(startTime, 10, 64) // 10进制, 64位
		if startTime < 631123200 || startTime > now || err != nil {
			return records, "开始时间解析错误", err
		}
		// 解析开始和结束日期与时间戳为time.Time格式
		start := time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, err := strconv.ParseInt(endTime, 10, 64)
		if endTime < 631123200 || endTime > now || err != nil {
			return records, "结束时间解析错误", err
		}
		end := time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		//如果开始时间大于结束时间,返回错误
		if startTime > endTime {
			return records, "开始时间在结束时间之后", err
		}
		query = query.Where("date_and_time BETWEEN ? AND ?", start, end)
	}
	if speciesName != "" {
		query = query.Where("species_name LIKE ?", "%"+speciesName+"%")
	}
	if investigator != "" {
		query = query.Where("investigator = ?", investigator)
	}
	if livingEnvironmentType != "" {
		query = query.Where("living_environment_type LIKE ?", "%"+livingEnvironmentType+"%")
	}

	err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&records).Error

	return records, nil, err
}

// GetRecordAll 获取所有记录(不分页,用以创建文件并供下载)
func GetRecordAll(gridNumber string, lineNumber string,
	startTime string, endTime string, province string, city string, county string,
	speciesName string, investigator string, livingEnvironmentType string) (records *[]model.Record, data interface{}, err error) {

	query := DB.Model(model.Record{})
	if gridNumber != "" {
		query = query.Where("grid_number LIKE ?", "%"+gridNumber+"%")
	}
	if lineNumber != "" {
		query = query.Where("line_number LIKE ?", "%"+lineNumber+"%")
	}
	if county != "" {
		query = query.Where("investigation_county = ?", county)
	} else if city != "" {
		query = query.Where("investigation_city = ?", city)
	} else if province != "" {
		query = query.Where("investigation_province = ?", province)
	}
	if startTime != "" || endTime != "" { //只要设置了一个时间必须设置另一个时间
		if startTime == "" {
			startTime = "631123200"
		}
		if endTime == "" {
			//更改endTime为当前时间的时间戳
			endTime = strconv.FormatInt(time.Now().Unix(), 10)
		}
		nowTime := strconv.FormatInt(time.Now().Unix(), 10)
		now, err := strconv.ParseInt(nowTime, 10, 64)

		// 字符串转为int64
		startTime, err := strconv.ParseInt(startTime, 10, 64) // 10进制, 64位
		if startTime < 631123200 || startTime > now || err != nil {
			return records, "开始时间解析错误", err
		}

		// 解析开始和结束日期与时间戳为time.Time格式
		start := time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, err := strconv.ParseInt(endTime, 10, 64)
		if endTime < 631123200 || endTime > now || err != nil {
			return records, "结束时间解析错误", err
		}
		end := time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		//如果开始时间大于结束时间,返回错误
		if startTime > endTime {
			return records, "开始时间在结束时间之后", err
		}
		query = query.Where("date_and_time BETWEEN ? AND ?", start, end)
	}
	if speciesName != "" {
		query = query.Where("species_name LIKE ?", "%"+speciesName+"%")
	}
	if investigator != "" {
		query = query.Where("investigator = ?", investigator)
	}
	if livingEnvironmentType != "" {
		query = query.Where("living_environment_type LIKE ?", "%"+livingEnvironmentType+"%")
	}

	err = query.Find(&records).Error
	return records, data, err

}

// GetRecordCount 获取记录总数
func GetRecordCount(gridNumber string, lineNumber string,
	startTime string, endTime string, province string, city string, county string,
	speciesName string, investigator string, livingEnvironmentType string) (count int64, data interface{}, err error) {

	query := DB.Model(model.Record{})
	if gridNumber != "" {
		query = query.Where("grid_number LIKE ?", "%"+gridNumber+"%")
	}
	if lineNumber != "" {
		query = query.Where("line_number LIKE ?", "%"+lineNumber+"%")
	}
	if county != "" {
		query = query.Where("investigation_county = ?", county)
	} else if city != "" {
		query = query.Where("investigation_city = ?", city)
	} else if province != "" {
		query = query.Where("investigation_province = ?", province)
	}
	if startTime != "" || endTime != "" { //只要设置了一个时间必须设置另一个时间
		if startTime == "" {
			startTime = "631123200"
		}
		if endTime == "" {
			//更改endTime为当前时间的时间戳
			endTime = strconv.FormatInt(time.Now().Unix(), 10)
		}
		nowTime := strconv.FormatInt(time.Now().Unix(), 10)
		now, err := strconv.ParseInt(nowTime, 10, 64)

		// 字符串转为int64
		startTime, err := strconv.ParseInt(startTime, 10, 64) // 10进制, 64位
		if startTime < 631123200 || startTime > now || err != nil {
			return 0, "开始时间解析错误", err
		}

		// 解析开始和结束日期与时间戳为time.Time格式
		start := time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, err := strconv.ParseInt(endTime, 10, 64)
		if endTime < 631123200 || endTime > now || err != nil {
			return 0, "结束时间解析错误", err
		}
		end := time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		//如果开始时间大于结束时间,返回错误
		if startTime > endTime {
			return 0, "开始时间在结束时间之后", err
		}
		query = query.Where("date_and_time BETWEEN ? AND ?", start, end)
	}
	if speciesName != "" {
		query = query.Where("species_name LIKE ?", "%"+speciesName+"%")
	}
	if investigator != "" {
		query = query.Where("investigator = ?", investigator)
	}
	if livingEnvironmentType != "" {
		query = query.Where("living_environment_type LIKE ?", "%"+livingEnvironmentType+"%")
	}

	err = query.Count(&count).Error
	return count, data, err
}

// GetRecordByIds 根据id数组获取记录
func GetRecordByIds(ids []uint) (records *[]model.Record, err error) {
	err = DB.Model(model.Record{}).Find(&records, ids).Error //多条查询
	//根据ids数组查询, 这三个方法都行
	//err = DB.Model(model.UnreviewedRecord{}).Where("id in (?)", ids).Find(&uRecords).Error
	//根据id进行查询
	//for _, id := range ids {
	//	err = DB.Model(model.UnreviewedRecord{}).Where("id = ?", id).First(&uRecords).Error
	//}
	return records, err
}

// DeleteRecordByIds 根据id数组删除记录
func DeleteRecordByIds(records *[]model.Record) error {
	return DB.Model(model.Record{}).Unscoped().Delete(&records).Error
}

func UploadRecord(records *[]model.Record) (returnErr error, existRecordsIndex []string) {
	for _, record := range *records {
		if err := DB.Model(model.Record{}).Create(&record).Error; err != nil {
			returnErr = err
			data := strings.Split(err.Error(), " ")[5]
			existRecordsIndex = append(existRecordsIndex, data)
		}
	}
	return returnErr, existRecordsIndex

	//return DB.Model(model.Record{}).Create(&record).Error
}

// GetRecordByArea 根据地区获取记录 todo 现在是把所有数据都查出来, 后续可能优化成只查出数量, 目前时间可以接受
func GetRecordByArea(province string, city string, county string) (records *[]model.Record, err error) {
	query := DB.Model(model.Record{}) // note 动态构建查询条件
	if province != "" {
		query = query.Where("investigation_province = ?", province)
	}
	if city != "" {
		query = query.Where("investigation_city = ?", city)
	}
	if county != "" {
		query = query.Where("investigation_county = ?", county)
	}
	err = query.Find(&records).Error
	return records, err
}
