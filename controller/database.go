package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var Db *gorm.DB

// InitDB 数据库初始化
func InitDB() *gorm.DB {

	dsn := "root:123456@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "",
			SingularTable: true, //使用单数表名，启用该选项后，`User`表将是 `user`
			//NameReplacer:  nil,
			//NoLowerCase:   false,
		},
	})
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}

	//将局部变量赋值给全局变量
	Db = db
	// 自动迁移到数据库
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to create database,err:" + err.Error())
	}

	//设置连接池
	sqlDB, _ := Db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func GetDB() *gorm.DB {
	return Db
}
