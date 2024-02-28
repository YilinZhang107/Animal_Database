/*
* @Author: Oatmeal107
* @Date:   2023/6/12 13:28
 */

package dao

import (
	"Animal_database/config"
	"Animal_database/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"os"
	"time"
)

var DB *gorm.DB

func InitMySQL() {
	conf := config.Conf.MysqlConfig

	//读写主从分离
	pathWrite := conf.Username + ":" + conf.Password + "@tcp(" + conf.DbHost + ":" + conf.Port + ")/" +
		conf.DbName + "?" + conf.Config
	pathRead := conf.Username + ":" + conf.Password +
		"@tcp(" + conf.DbHost + ":" + conf.Port + ")/" +
		conf.DbName + "?" + conf.Config

	//加上日志的打印
	var gormLogger logger.Interface
	if gin.Mode() == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: gormLogger, //日志打印
		NamingStrategy: schema.NamingStrategy{ //表名命名规则
			SingularTable: true, //表名不加s
		},
	})
	if err != nil {
		panic(err)
	}

	//连接池设置
	sqlDB, _ := db.DB()                        //连接池
	sqlDB.SetMaxIdleConns(20)                  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100)                 // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30) // 设置连接可以复用的最长时间
	//_db = db

	//主从配置
	_ = db.Use(dbresolver.Register(dbresolver.Config{ //读写分离配置的结构体
		// `db2` 作为 sources，`db3`、`db4` 作为 replicas
		Sources:  []gorm.Dialector{mysql.Open(pathWrite)},                      // 写操作
		Replicas: []gorm.Dialector{mysql.Open(pathRead), mysql.Open(pathRead)}, // 读操作
		Policy:   dbresolver.RandomPolicy{},                                    // sources/replicas 负载均衡策略
	}))

	err = db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.Record{},
		&model.User{},
		&model.UnreviewedRecord{},
		//todo 添加表
	)
	if err != nil {
		fmt.Println("--------------------------------------------------------------------")
		fmt.Println("register table fail")
		os.Exit(0)
	}
	DB = db
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("register table success")
}

//func NewDBClient(c context.Context) *gorm.DB {
//	return DB
//}
