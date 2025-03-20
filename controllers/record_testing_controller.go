package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

// @Summary 获取测试做题记录列表
// @Description 根据条件获取测试记录列表，并返回分页结果
// @Tags Testing
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param status query int false "状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.TestingRecordsResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/list [get]
func TestingRecords(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	
	// 执行分页查询
    results, total, err := database.PaginationQuery("testing_records", pageNo, pageLimit, conditions)
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

// @Summary 获取测试做题记录详情
// @Description 根据id获取测试做题记录详情
// @Tags Testing
// @Accept json
// @Produce json
// @Param id query int true "测试做题记录id"
// @Success 200 {object} models.ResponseData{data=models.TestingRecordsItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/detail/{id} [get]
func TestingRecordDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid testing record ID")
		return
	}

	record, err := database.GetDataById("testing_records", id)
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

// @Summary 新增测试做题记录
// @Description 新增测试做题记录
// @Tags Testing
// @Accept json
// @Produce json
// @Param part body models.TestingRecordsItem true "测试做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/add [post]
func AddTestingRecord(c *gin.Context) {
	var part models.TestingRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("testing_records", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert testing records")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新测试做题记录
// @Description 更新测试做题记录
// @Tags Testing
// @Accept json
// @Produce json
// @Param part body models.TestingRecordsItem true "测试做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/update [put]
func UpdateTestingRecord(c *gin.Context) {
	var part models.TestingRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("testing_records", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update testing record")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除测试做题记录
// @Description 根据ID删除测试做题记录
// @Tags Testing
// @Accept json
// @Produce json
// @Param id path int true "测试做题记录ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/delete/{id} [delete]
func DeleteTestingRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid testing record ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("testing_records", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete testing record")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "testing record not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}

// @Summary 提交套题做题记录
// @Description 提交套题做题记录
// @Tags Testing
// @Accept json
// @Produce json
// @Param part body models.TestingRecordsItem true "套题做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/testing/submit [post]
func SubmitTestingRecord(c *gin.Context) {
	// 获取提交的套题做题记录
	var part models.TestingRecordsItem

	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 根据test_id获取套题列表
	test, err := database.GetDataById("testing_list", part.TestID)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get data by id")
        return
    }

	// 获取part_list
	listeningPartListInterface, ok := test["listening_ids"].([]interface{})
	if !ok {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to parse listening_ids")
		return
	}

	// 调用 GetPartDetails 获取part详细信息
	listeningDetails, err := utils.GetPartDetails(listeningPartListInterface, "listening_part_list")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to get listening part detail")
		return
	}
	spew.Dump(listeningPartListInterface, "listeningPartListInterface")

	// 获取part_list
	readingPartListInterface, ok := test["reading_ids"].([]interface{})
	if !ok {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to parse reading_ids")
		return
	}

	readingDetails, err := utils.GetPartDetails(readingPartListInterface, "reading_part_list")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to get reading part detail")
		return
	}
	spew.Dump(readingPartListInterface, "readingPartListInterface")

	listeningAnswers := part.Answers[:40] // 前40个答案
	readingAnswers := part.Answers[40:80] // 后40个答案（从第41个到第80个）

	listeningScore := utils.CalculateScore(listeningDetails, listeningAnswers)
	readingScore := utils.CalculateScore(readingDetails, readingAnswers)

	part.Status = "0"
	part.Type = "3"
	part.Score = []int{int(listeningScore), int(readingScore)}
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
		return
	}
	// 将 user_id 添加到 part 中
	part.UserID = userID

	// 将数据插入数据库
	result, err := database.InsertData("testing_records", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update testing records")
		return
	}

	// // 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")

}