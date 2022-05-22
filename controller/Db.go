package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

// UserLoginInfo 数据库表
type UserLoginInfo struct {
	//User  User
	Id       int64
	Username string `gorm:"primaryKey"`
	Pwd      string
	Token    string
}

var Db *gorm.DB

// InitDB 数据库初始化
func InitDB() {

	dsn := "root:123456@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	Db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "",
			SingularTable: true, //使用单数表名，启用该选项后，`User`表将是 `user`
			//NameReplacer:  nil,
			//NoLowerCase:   false,
		},
	})

	//if err != nil {
	//	panic(err)
	//}

	//设置连接池
	sqlDB, _ := Db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

}

func checkUserExist(username string, password string) (*gorm.DB, bool) {
	Db.AutoMigrate(&UserLoginInfo{})

	var user = UserLoginInfo{
		Username: username,
		Pwd:      password,
		Token:    username + password,
	}

	// 查找该用户是否存在
	result := Db.First(&user, "username = ? and pwd = ?", username, password)
	if result.Error != nil {
		return result, false
	} else {
		return result, true
	}
}

func addUserInDB(newUser UserLoginInfo) {

	err := Db.AutoMigrate(&UserLoginInfo{})
	if err != nil {
		return
	}

	Db.Create(&newUser)

}

func UserLogin(token string) (*gorm.DB, bool) {
	Db.AutoMigrate(&UserLoginInfo{})

	var user = UserLoginInfo{}

	// 查找该用户是否存在
	result := Db.First(&user, "token=?", token)
	if result.Error != nil {
		return result, false
	} else {
		return result, true
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
