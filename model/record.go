/*
* @Author: Oatmeal107
* @Date:   2023/6/12 11:06
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

// Record    记录表
type Record struct {
	gorm.Model
	InvestigationMethod     string    `gorm:"size:4"`  //调查方式  todo 样线样点两种
	GridNumber              string    `gorm:"size:32"` //网格编号
	LineNumber              string    `gorm:"size:32"` //样线编号(样点不需要
	LineLength              float64   //样线长度样点不需要
	GPSNumber               uint64    //GPS编号
	SpeciesName             string    `gorm:"uniqueIndex:once;size:32"` //物种名称
	Count                   uint64    `gorm:"uniqueIndex:once"`         //数量
	Distance                string    `gorm:"uniqueIndex:once;size:16"` //距离(单位m, 或"在飞")
	PictureNumber           string    `gorm:"uniqueIndex:once;size:32"` //照片编号
	HabitatPicture          string    `gorm:"size:32"`                  //栖息地照片编号
	SpeciesThreatenedFactor string    `gorm:"size:64"`                  //物种受威胁因素
	LivingEnvironmentType   string    `gorm:"size:16"`                  //生境类型
	EnvThreatenedFactor     string    `gorm:"size:64"`                  //生境受威胁因素
	ThreatIntensity         string    `gorm:"size:16"`                  //威胁强度
	Weather                 string    `gorm:"size:16"`                  //天气
	Investigator            string    `gorm:"size:16"`                  //调查人
	InvestigationProvince   string    `gorm:"size:16; index:province"`  //调查省
	InvestigationCity       string    `gorm:"size:16; index:city"`      //调查市
	InvestigationCounty     string    `gorm:"size:16; index:county"`    //调查县
	InvestigationTown       string    `gorm:"size:16; index:town"`      //调查乡
	Latitude                string    `gorm:"uniqueIndex:once;size:64"` //纬度
	Longitude               string    `gorm:"uniqueIndex:once;size:64"` //经度
	Altitude                string    `gorm:"size:16"`                  //海拔
	DateAndTime             time.Time // 调查时间和日期
	Uploader                uint      //上传者
	Auditor                 uint      //审核者

}
