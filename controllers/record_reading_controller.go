package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
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

// @Summary 获取阅读做题记录详情
// @Description 根据id获取阅读做题记录详情
// @Tags Reading
// @Accept json
// @Produce json
// @Param id query int true "阅读做题记录id"
// @Success 200 {object} models.ResponseData{data=models.ReadingRecordsItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/reading/detail/{id} [get]
func ReadingRecordDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading record ID")
		return
	}

	record, err := database.GetDataById("reading_records", id)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get data by id")
        return
    }

	// 返回查询结果
	response := map[string]interface{}{
		"data": record,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}

// @Summary 新增阅读做题记录
// @Description 新增阅读做题记录
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.ReadingRecordsItem true "阅读做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/reading/add [post]
func AddReadingRecord(c *gin.Context) {
	var part models.ReadingRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("reading_records", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert reading records")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新阅读做题记录
// @Description 更新阅读做题记录
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.ReadingRecordsItem true "阅读做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/reading/update [put]
func UpdateReadingRecord(c *gin.Context) {
	var part models.ReadingRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("reading_records", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update reading record")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除阅读做题记录
// @Description 根据ID删除阅读做题记录
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "阅读做题记录ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/reading/delete/{id} [delete]
func DeleteReadingRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading record ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("reading_records", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete reading record")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "reading record not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}