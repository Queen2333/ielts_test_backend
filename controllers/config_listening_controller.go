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

// @Summary 获取听力套题列表
// @Description 根据条件获取听力列表，并返回分页结果
// @Tags Listening
// @Accept json
// @Produce json
// @Param name query string false "听力名称"
// @Param status query int false "听力状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.ListeningListResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/list [get]
func ListeningList(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	
	// 执行分页查询
    results, total, err := database.PaginationQuery("listening_list", pageNo, pageLimit, conditions)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to execute pagination query")
        return
    }

	result, err := utils.ProcessPartList(c, results, "listening_part_list")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	// 返回查询结果
	response := map[string]interface{}{
		"items": result,
		"total":   total,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}

// @Summary 获取听力套题详情
// @Description 根据id获取听力详情
// @Tags Listening
// @Accept json
// @Produce json
// @Param id query int true "听力id"
// @Success 200 {object} models.ResponseData{data=models.ListeningItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/detail/{id} [get]
func ListeningDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening set ID")
		return
	}

	record, err := database.GetDataById("listening_list", id)
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

// @Summary 新增听力套题
// @Description 新增听力套题
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.BasicListeningItem true "听力套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/add [post]
func AddListening(c *gin.Context) {
	var part models.BasicListeningItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	if part.Type.Int() == 3 {
		userID, err := utils.GetUserIDFromToken(c)
		if err != nil {
			// 处理获取 user_id 失败的情况
			utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
			return
		}

		// 将 user_id 添加到 part 中
		part.UserID = userID
	}

	// 将数据插入数据库
	result, err := database.InsertData("listening_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert listening part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新听力套题
// @Description 更新听力套题
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.BasicListeningItem true "听力套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/update [put]
func UpdateListening(c *gin.Context) {
	var part models.BasicListeningItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 如果 type=3（用户自定义），需要保留 user_id
	if part.Type.Int() == 3 {
		// 如果请求中没有提供 user_id，从数据库获取原有的 user_id
		if part.UserID == "" {
			existingData, err := database.GetDataById("listening_list", part.ID)
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get existing data")
				return
			}

			if existingUserID, ok := existingData["user_id"].(string); ok {
				part.UserID = existingUserID
			}
		}
	}

	// 将数据更新到数据库
	result, err := database.InsertData("listening_list", &part, "update")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening part")
		return
	}

	// 返回更新后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除听力套题
// @Description 根据ID删除听力套题
// @Tags Listening
// @Accept json
// @Produce json
// @Param id path int true "听力套题ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/delete/{id} [delete]
func DeleteListening(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("listening_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete listening set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "Listening set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}

// @Summary      获取听力篇列表
// @Description  根据条件获取听力part列表，并返回分页结果
// @Tags         Listening
// @Accept       json
// @Produce      json
// @Param        name      query  string  false  "试题名称"
// @Param        status    query  int     false  "试题状态"
// @Param        type      query  int     false  "试题类型"
// @Param        pageNo    query  int     true   "页码"
// @Param        pageLimit query  int     false   "每页条数"
// @Success      200  {object}  models.ResponseData{data=models.ListeningPartListResponse}
// @Failure      400  {object}  models.ResponseData{data=nil}
// @Failure      500  {object}  models.ResponseData{data=nil}
// @Router       /config/listening-part/list [get]
func ListeningPartList(c * gin.Context) {
	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	// 执行分页查询
    results, total, err := database.PaginationQuery("listening_part_list", pageNo, pageLimit, conditions)
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

// @Summary 获取听力part详情
// @Description 根据id获取听力part详情
// @Tags Listening
// @Accept json
// @Produce json
// @Param id query int true "听力partid"
// @Success 200 {object} models.ResponseData{data=models.ListeningPartItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/detail/{id} [get]
func ListeningPartDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening set ID")
		return
	}

	record, err := database.GetDataById("listening_part_list", id)
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

// @Summary 新增听力part
// @Description 新增听力part
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningPartItem true "听力part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/add [post]
func AddListeningPart(c *gin.Context) {
	var part models.ListeningPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 获取当前用户ID并设置到part中
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		utils.HandleResponse(c, http.StatusUnauthorized, "", err.Error())
		return
	}
	part.UserID = userID

	// 将数据插入数据库
	result, err := database.InsertData("listening_part_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert listening part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新听力part
// @Description 更新听力part
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningPartItem true "听力part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/update [put]
func UpdateListeningPart(c *gin.Context) {
	var part models.ListeningPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 获取原有数据以保留 user_id
	existingData, err := database.GetDataById("listening_part_list", part.ID)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get existing data")
		return
	}

	// 如果请求中没有 user_id，则保留原有的 user_id
	if part.UserID == "" {
		if existingUserID, ok := existingData["user_id"].(string); ok {
			part.UserID = existingUserID
		}
	}

	// 将数据更新到数据库
	result, err := database.InsertData("listening_part_list", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening part")
		return
	}

	// 返回更新后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除听力part
// @Description 根据ID删除听力part
// @Tags Listening
// @Accept json
// @Produce json
// @Param id path int true "听力partID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/delete/{id} [delete]
func DeleteListeningPart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid listening set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("listening_part_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete listening set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "Listening set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}