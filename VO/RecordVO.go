/*
* @Author: Oatmeal107
* @Date:   2024/3/21 15:39
 */

package VO

import (
	"Animal_database/model"
	"time"
)

// RecordVO 数据序列化器,用于返回给前端的数据, (因为有一些数据不用返回)
type RecordVO struct {
	ID                      uint
	InvestigationMethod     string
	GridNumber              string
	LineNumber              string
	LineLength              float64
	GPSNumber               uint64
	SpeciesName             string
	Count                   uint64
	Distance                string
	PictureNumber           string
	HabitatPicture          string
	SpeciesThreatenedFactor string
	LivingEnvironmentType   string
	EnvThreatenedFactor     string
	ThreatIntensity         string
	Weather                 string
	Investigator            string
	InvestigationProvince   string
	InvestigationCity       string
	InvestigationCounty     string
	InvestigationTown       string
	Latitude                string
	Longitude               string
	Altitude                string
	DateAndTime             time.Time
	Uploader                uint
	Auditor                 uint
}

func BuildRecordVO(u *model.Record) *RecordVO {
	return &RecordVO{
		ID:                      u.ID,
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
		Auditor:                 u.Auditor,
	}
}
