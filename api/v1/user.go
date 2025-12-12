// 接口处理
package v1

import (
	"ZustChat/model/response"
	"ZustChat/model"
	"ZustChat/model/request"
	"ZustChat/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Success <code> {<type>} <model or object> "<description>"

// GetUserList
// @Summary 获取用户列表
// @Tag 获取用户列表
// @Success 200 {string} json{"code","data"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	userList, err := service.GetUserList()
	if err != nil {
		response.FailWithMessage("查询失败", c)
	}
	response.OkWithData(userList, c)
}


// UpdateUser
// @Summary 更新用户
// @Tag 更新用户
// @Success 200 {string} json{"code","data"}
// @Router /user/updateUser [put]
func UpdateUserInfo(c *gin.Context) {
	var updateuserinfo request.UpdateUserInfo
	if err := c.ShouldBindJSON(&updateuserinfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	validate := validator.New()
	if err := validate.Var(updateuserinfo.Username, "omitempty,min=2,max=100"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else if err := validate.Var(updateuserinfo.Phone, "omitempty"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else if err := validate.Var(updateuserinfo.Email, "omitempty,email"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.UpdateUserInfo(&updateuserinfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)

}


// DeleteUser
// @Summary 注销
// @Tag 删除用户
// @Success 200 {string} json{"code","data"}
// @Router //user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var user model.User_Basic
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	validate := validator.New()
	if err := validate.Var(user.Username, "omitempty,min=2,max=100"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.DeleteUser(user.Username); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)

}

// 查找某个用户
func Finduser(c *gin.Context) {
	var user request.FindUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	founduser, err := service.FindUserByID(user.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(founduser, "查找成功", c)
}

// 查找用户所有好友
func LoadFriends(c *gin.Context) {
	var friendpayload request.LoadFriendsPayload
	// userid, err := strconv.Atoi(c.Query("userid"))
	err := c.ShouldBindJSON(&friendpayload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	users, err := service.FindFriendsByUserID(friendpayload.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 封装rows
	data := gin.H{
		"Rows": users,
	}
	response.OkWithDetailed(data, "返回所有好友成功", c)
}

// 用户添加好友
func AddFriend(c *gin.Context) {
	var addfriend request.AddFriend

	err := c.ShouldBindJSON(&addfriend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.AddFrend(&addfriend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("添加成功", c)
}

func FindFriend(c *gin.Context) {
	var findFriend request.FindFriend
	err := c.ShouldBindJSON(&findFriend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	contact, err := service.FindFrend(findFriend.UserId, findFriend.FriendId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(contact, "查找成功", c)
}
