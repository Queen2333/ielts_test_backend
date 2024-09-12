package controllers

import (
	"net/http"
	"strconv"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 获取阅读做题记录列表
// @Description 根据条件获取阅读记录列表，并返回分页结果
// @Tags Reading
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param status query int false "状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.ReadingRecordsResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/reading/list [get]
func ReadingRecords(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	
	// 执行分页查询
    results, total, err := database.PaginationQuery("reading_records", pageNo, pageLimit, conditions)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to execute pagination query")
        return
    }

	// 返回查询结果
	response := map[string]interface{}{
		"items": results,
		"total":   total,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}