package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
//var usersLoginInfo = map[string]User{
//	"zhangleidouyin": {
//		Id:            1,
//		Name:          "zhanglei",
//		FollowCount:   10,
//		FollowerCount: 5,
//		IsFollow:      true,
//	},
//}

var usersLoginInfo = map[string]User{}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	// 通过 Query 方法接收表单传输过来的 username 和 password 字段数据
	username := c.Query("username")
	password := c.Query("password")

	// 设置 token 为username + password
	//token := username + password

	//如果数据库存在该 token，则 checkUserExist 返回ture
	if _, exist := checkUserExist(username, password); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 不存在该 token ，可以注册，
		atomic.AddInt64(&userIdSequence, 1)
		newUser := UserLoginInfo{
			Id:       userIdSequence,
			Username: username,
			Pwd:      password,
			Token:    username + password,
		}
		addUserInDB(newUser)

		//usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	//用户存在
	if _, exist := checkUserExist(username, password); exist {
		// POST 传值，转换为JSON数据
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "User Login Success"},
			//UserId:   username,
			Token: token,
		})
	} else { // 用户不存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	//username := c.Query("username")
	//password := c.Query("password")

	if _, exist := UserLogin(token); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "User Login Success"},
			//User:     username,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
