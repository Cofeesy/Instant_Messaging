package utils

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql" // 导入 GORM V2 的 MySQL 驱动
	"gorm.io/gorm"         // 导入 GORM V2 核心包
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitMysql() {

	var (
		err                                error
		dbName, user, password, host, port string
	)

	sec, err := Cfg.GetSection("mysql")
	if err != nil {
		log.Fatalf("Fail to get section 'mysql': %v", err)
	}

	// dbType = sec.Key("TYPE").String()
	dbName = sec.Key("DBNAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()

	host = sec.Key("HOST").String()
	port = sec.Key("PORT").String()
	// // tablePrefix = sec.Key("TABLE_PREFIX").String()

	// FIXME:这是gorm-v1的连接方式
	// db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	user,
	// 	password,
	// 	host,
	// 	dbName))

	// if err !=nil{
	// 	log.Println(err)
	// }

	// // 用了这个，go代码里面的结构体映射到数据库会多一个前缀
	// // gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	// //     return defaultTableName;
	// // }

	// // 禁用复数
	// db.SingularTable(true)
	// // 详细日志模式，更详细的打印日志
	// db.LogMode(true)
	// // 设置连接池中的最大空闲连接数
	// // 连接池类似一个连接缓存
	// // 提前启动的连接缓冲区最多有x个，这里的x是10
	// db.DB().SetMaxIdleConns(10)
	// // 设置数据库中最大打开连接数
	// // 应用程序最多只能同时占用x个到数据库的连接
	// db.DB().SetMaxOpenConns(100)

	// FIXME:这是gorm-v2的连接方式
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbName)

	// --- 这是 V2 初始化方式的核心变化 ---
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 配置 GORM 日志
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info，相当于 V1 的 LogMode(true)
		// 配置命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，相当于 V1 的 SingularTable(true)
		},
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// --- 连接池设置的变化 ---
	sqlDB, err := DB.DB() // 获取底层的 *sql.DB 对象
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// 连接池设置现在是在获取到的 sqlDB 对象上进行
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

func CloseDB() {
	sqlDB, _ := DB.DB()
	defer sqlDB.Close()
}
