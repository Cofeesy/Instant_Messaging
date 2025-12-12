package v1

import (
	"ZustChat/model/response"
	"ZustChat/model/request"
	"ZustChat/service"
	"github.com/gin-gonic/gin"
)

// 用户创建群组
func CreateGroup(c *gin.Context) {
	// owner, err := strconv.Atoi(c.Query("ownerid"))
	// groupName := c.Query("groupname")
	var creategroup request.CreatGroup
	err := c.ShouldBindJSON(&creategroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	group, err := service.CreateGroup(creategroup)
	if err != nil {
		response.FailWithDetailed(group, err.Error(), c)
		return
	}

	response.OkWithDetailed(group, "创建群组成功", c)
}


// 返回群列表
func LoadGroups(c *gin.Context) {
	var loadgroups request.LoadGroups
	err := c.ShouldBindJSON(&loadgroups)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	groups, err := service.FindGroupsByUserID(loadgroups.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 封装Rows
	data := gin.H{
		"Rows": groups,
	}
	response.OkWithDetailed(data, "查找群组成功", c)
}

// 加入群组
func AddGroup(c *gin.Context) {
	var addGroup request.AddGroup
	err := c.ShouldBindJSON(&addGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = service.AddGroup(&addGroup); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("加入成功", c)
}

