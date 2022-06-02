package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var usersLoginInfo = map[string]User{}
var userIdSequence = int64(1)

func Register(c *gin.Context) {
	DB := GetDB()
	// 通过 Query 方法接收表单传输过来的 username 和 password 字段数据
	name := c.Query("username")
	password := c.Query("password")

	// 数据验证
	if len(name) < 6 || len(name) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "用户名必须介于6到32位之间",
		})
		return
	}
	if len(password) < 6 || len(password) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码必须介于6到32位之间",
		})
		return
	}

	// 判断用户是否存在
	if _, exist := isUserExist(DB, name); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	} else {
		// 不存在该用户，可以注册，

		// 密码加密
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500, "msg": "加密错误",
			})
			return
		}

		// 发放 token
		//atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Name:     name,
			Password: string(hasedPassword),
		}
		token, err := ReleaseToken(newUser)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
			log.Printf("token generate error : %v", err)
			return
		}
		newUser.Token = token

		Db.Create(&newUser)
		usersLoginInfo[token] = newUser
		fmt.Printf("创建的新用户为：%v", newUser)
		// 返回信息
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.ID,
			Token:    token,
		})
		//保存当前请求信息到上下文c中
		c.Set("token", token)
		//继续执行后续的请求
		c.Next()
	}

}

func Login(c *gin.Context) {
	DB := GetDB()
	// 获取参数
	name := c.Query("username")
	password := c.Query("password")

	//数据验证
	if len(name) < 6 || len(name) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "用户名必须介于6到32位之间",
		})
		return
	}
	if len(password) < 6 || len(password) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码必须介于6到32位之间",
		})
		return
	}

	// 判断用户是否存在
	var user User
	DB.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	} else {
		//返回结果
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    user.Token,
		})
		usersLoginInfo[user.Token] = user
		//fmt.Printf("当前登录用户是：%v", user)
	}

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}

func isUserExist(db *gorm.DB, name string) (User, bool) {
	var user User
	db.Where("name = ?", name).First(&user)
	if user.ID != 0 {
		return user, true
	}
	return user, false

}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

//func UserInfo(c *gin.Context) {
//	db := GetDB()
//	token := c.Query("token")
//	var user User
//	db.Where("token = ?", token).First(&user)
//	usersLoginInfo[token] = user
//
//	fmt.Printf("当前登录的用户是:%v", user.Name)
//	if _, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 0},
//			User:     user,
//		})
//	} else {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
//		})
//	}
//}
