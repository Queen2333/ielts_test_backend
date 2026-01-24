package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 获取阅读套题列表
// @Description 根据条件获取阅读列表，并返回分页结果
// @Tags Reading
// @Accept json
// @Produce json
// @Param name query string false "阅读名称"
// @Param status query int false "阅读状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.ReadingListResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading/list [get]
func ReadingList(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	// 执行分页查询
    results, total, err := database.PaginationQuery("reading_list", pageNo, pageLimit, conditions)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to execute pagination query")
        return
    }

	// 遍历 results 提取所有的 part_list ID
	for i, result := range results {

		// 处理 []interface{} 类型的 part_list
		partListInterface, ok := result["part_list"].([]interface{})
		if !ok {
			utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to parse part_list")
			return
		}

		// 转换为字符串数组
		var partListStrArray []string
		for _, part := range partListInterface {
			partListStrArray = append(partListStrArray, fmt.Sprint(part))
		}

		var partListStr = strings.Join(partListStrArray, ",")

		partList := utils.StringToList(partListStr)
		// 查询 part_list 中的详细信息
		var details []map[string]interface{}
		for _, id := range partList {
			partDetail, err := database.GetPartsByIds("reading_part_list", []int{id})
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query reading parts")
				return
			}
			if len(partDetail) > 0 {
				details = append(details, partDetail[0])
			}
		}
		// 将查询结果放回到对应的 part_list 中
		results[i]["part_list"] = details
		// }
	}

	// 返回查询结果
	response := map[string]interface{}{
		"items": results,
		"total":   total,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}

// @Summary 获取阅读套题详情
// @Description 根据id获取阅读详情
// @Tags Reading
// @Accept json
// @Produce json
// @Param id query int true "阅读id"
// @Success 200 {object} models.ResponseData{data=models.ReadingItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading/detail/{id} [get]
func ReadingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading set ID")
		return
	}

	record, err := database.GetDataById("reading_list", id)
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

// @Summary 新增阅读套题
// @Description 新增阅读套题
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.BasicReadingItem true "阅读套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading/add [post]
func AddReading(c *gin.Context) {
	var part models.BasicReadingItem
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
	result, err := database.InsertData("reading_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert reading part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新阅读套题
// @Description 更新阅读套题
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.BasicReadingItem true "阅读套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading/update [put]
func UpdateReading(c *gin.Context) {
	var part models.BasicReadingItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 如果 type=3（用户自定义），需要保留 user_id
	if part.Type.Int() == 3 {
		// 如果请求中没有提供 user_id，从数据库获取原有的 user_id
		if part.UserID == "" {
			existingData, err := database.GetDataById("reading_list", part.ID)
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
	result, err := database.InsertData("reading_list", &part, "update")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update reading part")
		return
	}

	// 返回更新后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除阅读套题
// @Description 根据ID删除阅读套题
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "阅读套题ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading/delete/{id} [delete]
func DeleteReading(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("reading_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete reading set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "reading set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}

// @Summary      获取阅读篇列表
// @Description  根据条件获取阅读part列表，并返回分页结果
// @Tags         Reading
// @Accept       json
// @Produce      json
// @Param        name      query  string  false  "试题名称"
// @Param        status    query  int     false  "试题状态"
// @Param        type      query  int     false  "试题类型"
// @Param        pageNo    query  int     true   "页码"
// @Param        pageLimit query  int     false   "每页条数"
// @Success      200  {object}  models.ResponseData{data=models.ReadingPartListResponse}
// @Failure      400  {object}  models.ResponseData{data=nil}
// @Failure      500  {object}  models.ResponseData{data=nil}
// @Router       /config/reading-part/list [get]
func ReadingPartList(c * gin.Context) {
	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	// 执行分页查询
    results, total, err := database.PaginationQuery("reading_part_list", pageNo, pageLimit, conditions)
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

// @Summary 获取阅读part详情
// @Description 根据id获取阅读part
// @Tags Reading
// @Accept json
// @Produce json
// @Param id query int true "阅读part id"
// @Success 200 {object} models.ResponseData{data=models.ReadingPartItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading-part/detail/{id} [get]
func ReadingPartDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading part set ID")
		return
	}

	record, err := database.GetDataById("reading_part_list", id)
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

// @Summary 新增阅读part
// @Description 新增阅读part
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.ReadingPartItem true "阅读part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading-part/add [post]
func AddReadingPart(c *gin.Context) {
	var part models.ReadingPartItem
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
	result, err := database.InsertData("reading_part_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert reading part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新阅读part
// @Description 更新阅读part
// @Tags Reading
// @Accept json
// @Produce json
// @Param part body models.ReadingPartItem true "阅读part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading-part/update [put]
func UpdateReadingPart(c *gin.Context) {
	var part models.ReadingPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 获取原有数据以保留 user_id
	existingData, err := database.GetDataById("reading_part_list", part.ID)
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
	result, err := database.InsertData("reading_part_list", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update reading part")
		return
	}

	// 返回更新后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除阅读part
// @Description 根据ID删除阅读part
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "阅读partID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/reading-part/delete/{id} [delete]
func DeleteReadingPart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid reading set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("reading_part_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete reading set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "reading set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}