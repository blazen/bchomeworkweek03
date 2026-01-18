package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// 简化处理，实际应该解析 binding 错误
	errors["general"] = err.Error()
	return errors
}

func GetQueryPage(c *gin.Context) (pageNo, pageSize int) {
	pageNoStr := c.DefaultQuery("pageNo", "1") // 带默认值
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		HandleError(c, err)
		return
	}
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil {
		HandleError(c, err)
		return
	}
	return
}
