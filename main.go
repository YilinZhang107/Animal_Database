/*
* @Author: Oatmeal107
* @Date:   2023/6/4 8:53
 */

package main

import (
	"Animal_database/Init"
	"Animal_database/config"
	"Animal_database/dao"
	"Animal_database/utils"
)

func main() {
	provinces, prov2city, city2county := utils.GetMap()
	config.InitViper(provinces, prov2city, city2county) //(在这里初始化了全局的配置文件
	dao.InitMySQL()
	utils.InitLog()
	Init.NewRouter()
}
