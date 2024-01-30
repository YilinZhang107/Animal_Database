/*
* @Author: Oatmeal107
* @Date:   2023/6/12 16:12
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

// UnreviewedRecord    未审核记录表 (后面这个xlsx暂时没用)
type UnreviewedRecord struct {
	gorm.Model
	InvestigationMethod     string    `xlsx:"0" gorm:"size:4"`                    //调查方式  todo 样线样点两种
	GridNumber              string    `xlsx:"1" gorm:"size:32"`                   //网格编号
	LineNumber              string    `xlsx:"2" gorm:"size:32"`                   //样线编号(样点不需要
	LineLength              float64   `xlsx:"3"`                                  //样线长度样点不需要
	GPSNumber               uint64    `xlsx:"4"`                                  //GPS编号
	SpeciesName             string    `xlsx:"5" gorm:"uniqueIndex:once;size:32"`  //物种名称
	Count                   uint64    `xlsx:"6" gorm:"uniqueIndex:once"`          //数量
	Distance                string    `xlsx:"7" gorm:"uniqueIndex:once;size:16"`  //距离(单位m, 或"在飞")
	PictureNumber           string    `xlsx:"8" gorm:"uniqueIndex:once;size:32"`  //照片编号(可能是"无"
	HabitatPicture          string    `xlsx:"9" gorm:"size:32"`                   //栖息地照片编号
	SpeciesThreatenedFactor string    `xlsx:"10" gorm:"size:64"`                  //物种受威胁因素
	LivingEnvironmentType   string    `xlsx:"11" gorm:"size:16"`                  //生境类型
	EnvThreatenedFactor     string    `xlsx:"12" gorm:"size:64"`                  //生境受威胁因素
	ThreatIntensity         string    `xlsx:"13" gorm:"size:16"`                  //威胁强度
	Weather                 string    `xlsx:"14" gorm:"size:16"`                  //天气
	Investigator            string    `xlsx:"15" gorm:"size:16"`                  //调查人
	InvestigationProvince   string    `xlsx:"16" gorm:"size:16"`                  //调查省
	InvestigationCity       string    `xlsx:"17" gorm:"size:16"`                  //调查市
	InvestigationCounty     string    `xlsx:"18" gorm:"size:16"`                  //调查县
	InvestigationTown       string    `xlsx:"19" gorm:"size:16"`                  //调查乡
	Latitude                string    `xlsx:"20" gorm:"uniqueIndex:once;size:64"` //纬度
	Longitude               string    `xlsx:"21" gorm:"uniqueIndex:once;size:64"` //经度
	Altitude                string    `xlsx:"22" gorm:"size:16"`                  //海拔
	DateAndTime             time.Time // 调查时间和日期
	Uploader                uint      //上传者

}

// ToRecord 转换uRecord为Record
func (u *UnreviewedRecord) ToRecord(AuditorId uint) *Record {
	record := new(Record)
	*record = Record{
		InvestigationMethod:     u.InvestigationMethod,
		GridNumber:              u.GridNumber,
		LineNumber:              u.LineNumber,
		LineLength:              u.LineLength,
		GPSNumber:               u.GPSNumber,
		SpeciesName:             u.SpeciesName,
		Count:                   u.Count,
		Distance:                u.Distance,
		PictureNumber:           u.PictureNumber,
		HabitatPicture:          u.HabitatPicture,
		SpeciesThreatenedFactor: u.SpeciesThreatenedFactor,
		LivingEnvironmentType:   u.LivingEnvironmentType,
		EnvThreatenedFactor:     u.EnvThreatenedFactor,
		ThreatIntensity:         u.ThreatIntensity,
		Weather:                 u.Weather,
		Investigator:            u.Investigator,
		InvestigationProvince:   u.InvestigationProvince,
		InvestigationCity:       u.InvestigationCity,
		InvestigationCounty:     u.InvestigationCounty,
		InvestigationTown:       u.InvestigationTown,
		Latitude:                u.Latitude,
		Longitude:               u.Longitude,
		Altitude:                u.Altitude,
		DateAndTime:             u.DateAndTime,
		Uploader:                u.Uploader,
		Auditor:                 AuditorId,
	}
	return record
}
