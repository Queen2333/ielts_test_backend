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

// @Summary 获取听力做题记录列表
// @Description 根据条件获取听力记录列表，并返回分页结果
// @Tags Listening
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param status query int false "状态"
// @Param type query int true "3"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.ListeningRecordsResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/list [get]
func ListeningRecords(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	
	// 执行分页查询
    results, total, err := database.PaginationQuery("listening_records", pageNo, pageLimit, conditions)
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
// @Summary 获取听力做题记录详情
// @Description 根据id获取听力做题记录详情
// @Tags Listening
// @Accept json
// @Produce json
// @Param id query int true "听力做题记录id"
// @Success 200 {object} models.ResponseData{data=models.ListeningRecordsItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/detail/{id} [get]
func ListeningRecordDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening set ID")
		return
	}

	record, err := database.GetDataById("listening_records", id)
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

// @Summary 新增听力做题记录
// @Description 新增听力做题记录
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningRecordsItem true "听力做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/add [post]
func AddListeningRecord(c *gin.Context) {
	var part models.ListeningRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	part.Status = "0"
	part.Type = "3"

	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
		return
	}
	// 将 user_id 添加到 part 中
	part.UserID = userID

	// 将数据插入数据库
	result, err := database.InsertData("listening_records", &part, "create")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert listening records")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新听力做题记录
// @Description 更新听力做题记录
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningRecordsItem true "听力做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/update [put]
func UpdateListeningRecord(c *gin.Context) {
	var part models.ListeningRecordsItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	part.Status = "0"
	part.Type = "3"
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
		return
	}
	// 将 user_id 添加到 part 中
	part.UserID = userID

	// 将数据插入数据库
	result, err := database.InsertData("listening_records", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening records")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除听力做题记录
// @Description 根据ID删除听力做题记录
// @Tags Listening
// @Accept json
// @Produce json
// @Param id path int true "听力做题记录ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/delete/{id} [delete]
func DeleteListeningRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening record ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("listening_records", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete listening record")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "Listening record not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}

// @Summary 提交听力做题记录
// @Description 提交听力做题记录
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningRecordsItem true "听力做题记录内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /record/listening/submit [post]
func SubmitListeningRecord(c *gin.Context) {
	// 获取提交的听力做题记录
	var part models.ListeningRecordsItem

	spew.Dump(part, "part")
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 根据test_id获取听力列表
	test, err := database.GetDataById("listening_list", part.TestID)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get data by id")
        return
    }
	// 获取part_list
	partListInterface, ok := test["part_list"].([]interface{})
	if !ok {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to parse part_list")
		return
	}

	// 调用 GetPartDetails 获取part详细信息
	details, err := utils.GetPartDetails(partListInterface)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "failed to get part detail")
		return
	}

	// 继续使用 parsedDetails 进行评分计算
	score := utils.CalculateScore(details, part.Answers)

	part.Status = "0"
	part.Type = "3"
	part.Score = int(score)
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
		return
	}
	// 将 user_id 添加到 part 中
	part.UserID = userID

	// 将数据插入数据库
	result, err := database.InsertData("listening_records", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening records")
		return
	}

	// response := map[string]interface{}{
	// 	// "part_list": details,
	// 	"id": part.ID, // 听力做题记录id
	// 	"test_id": part.TestID, // 听力试题id
	// 	"name": name, // 听力名称
	// 	"score": score, // 得分
	// }

	// // 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
	// utils.HandleResponse(c, http.StatusOK, response, "Success")
}