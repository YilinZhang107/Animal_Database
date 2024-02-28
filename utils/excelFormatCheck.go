/*
* @Author: Oatmeal107
* @Date:   2023/6/19 14:15
 */

package utils

import (
	"Animal_database/config"
	"github.com/tealeg/xlsx"
)

// ExcelFormatCheck 检查上传的 Excel 文件格式是否正确
// 检查的有物种名必须是字符串，海拔，坐标等必须是浮点，并且不能超出限值等等。
func ExcelFormatCheck(excelFile *xlsx.File) (bool, string) {
	//先查表头是否正确,必须使用标准格式的表来上传数据
	//读取Excel文件的第一行，获取表头
	var header []string
	for _, row := range excelFile.Sheets[0].Rows {
		for _, cell := range row.Cells {
			if cell.String() != "" {
				header = append(header, cell.String())
			}
		}
		break
	}
	// 打开 模板Excel 文件,与其表头进行对比
	templateFile, err := xlsx.OpenFile(config.UploadTemplateExcelPath)
	if err != nil {
		return false, "打开模板文件失败"
	}
	var headerT []string
	for _, row := range templateFile.Sheets[0].Rows {
		for _, cell := range row.Cells {
			headerT = append(headerT, cell.String())
		}
		break
	}
	// 检查表头是否一致
	if len(header) != len(headerT) {
		return false, "表头数量不一致"
	}
	for i := 0; i < len(header); i++ {
		if header[i] != headerT[i] {
			return false, "表头不一致"
		}
	}

	// todo 检查表格内容是否符合要求
	return true, ""
}

// checkAnimalName 检查物种名是否符合要求
func checkAnimalName() {

}
