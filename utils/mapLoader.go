/*
* @Author: Oatmeal107
* @Date:   2023/11/19 16:27
 */

package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func GetMap() (provinces []string, prov2city map[string][]string, city2county map[string][]string) {
	// 打开 Excel 文件
	file, err := excelize.OpenFile("./resources/static/map.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取数据
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	prov2cityTmp := make(map[string]map[string]bool)
	city2countyTmp := make(map[string]map[string]bool)
	// 遍历行
	for _, row := range rows {
		_, found := prov2cityTmp[row[0]]
		if !found {
			prov2cityTmp[row[0]] = map[string]bool{}
		}
		prov2cityTmp[row[0]][row[1]] = true

		_, found = city2countyTmp[row[1]]
		if !found {
			city2countyTmp[row[1]] = map[string]bool{}
		}
		city2countyTmp[row[1]][row[2]] = true
	}
	prov2city = make(map[string][]string)
	city2county = make(map[string][]string)
	for key, value := range prov2cityTmp {
		provinces = append(provinces, key)
		for k := range value {
			prov2city[key] = append(prov2city[key], k)
			city2county[k] = []string{}
		}
	}
	for key, value := range city2countyTmp {
		for k := range value {
			city2county[key] = append(city2county[key], k)
		}
	}

	return provinces, prov2city, city2county
}
