package v1

import (
	"fmt"
	"gin_chat/common/response"
	"math/rand"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

// 上传文件不需要绑定，文件存储在前端文件夹，不涉及数据库存储
// 因此后端只需要管后端接收方式、存储文件名以及存储的文件夹然后返回即可
// 暂时指定文件夹"./asset/upload/"
func UploadInfo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// fmt.Println(file.Filename)

	// 将文件名分开.隔开
	filename_slice := strings.Split(file.Filename, ".")
	suffix := filename_slice[len(filename_slice)-1]
	perfix := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31())

	newFileName := perfix + "." + suffix

	// 组装地址
	dst := "././asset/upload/" + newFileName
	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(dst, "上传成功", c)
}