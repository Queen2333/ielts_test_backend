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

// @Summary 获取测试套题套题列表
// @Description 根据条件获取测试套题列表，并返回分页结果
// @Tags Testing
// @Accept json
// @Produce json
// @Param name query string false "测试套题名称"
// @Param status query int false "测试套题状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.TestingListResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/testing/list [get]
func TestingList(c *gin.Context) {

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("pageLimit", "-1"))

	conditions, err := utils.ProcessRequest(c)
	if err != nil {
		// 错误处理已在 ProcessRequest 中处理
		return
	}

	// 执行分页查询
    results, total, err := database.PaginationQuery("testing_list", pageNo, pageLimit, conditions)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to execute pagination query")
        return
    }

	// 定义字段与对应表的映射关系
	fieldToTableMap := map[string]string{
		"reading_ids":   "reading_part_list",
		"listening_ids": "listening_part_list",
		"writing_ids":   "writing_part_list",
	}

	fieldToNameMap := map[string]string{
		"reading_ids":   "reading_parts",
		"listening_ids": "listening_parts",
		"writing_ids":   "writing_parts",
	}
	// 遍历 results 提取所有的 part_list ID
	for i, result := range results {
		for field, table := range fieldToTableMap {
			partListInterface, ok := result[field].([]interface{})
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

			var details []map[string]interface{}
			for _, id := range partList {
				partDetail, err := database.GetPartsByIds(table, []int{id})
				if err != nil {
					utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query testing parts")
					return
				}
				if len(partDetail) > 0 {
					details = append(details, partDetail[0])
				}
			}
			results[i][fieldToNameMap[field]] = details
		}

	}

		
	// 返回查询结果
	response := map[string]interface{}{
		"items": results,
		"total":   total,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}

// @Summary 获取测试套题详情
// @Description 根据id获取测试详情
// @Tags Testing
// @Accept json
// @Produce json
// @Param id query int true "测试id"
// @Success 200 {object} models.ResponseData{data=models.TestingItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/testing/detail/{id} [get]
func TestingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid testing set ID")
		return
	}

	record, err := database.GetDataById("testing_list", id)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to get data by id")
        return
    }

	// 定义字段与对应表的映射关系
	fieldToTableMap := map[string]string{
		"reading_ids":   "reading_part_list",
		"listening_ids": "listening_part_list",
		"writing_ids":   "writing_part_list",
	}

	fieldToNameMap := map[string]string{
		"reading_ids":   "reading_parts",
		"listening_ids": "listening_parts",
		"writing_ids":   "writing_parts",
	}

	for field, table := range fieldToTableMap {
		partListInterface, ok := record[field].([]interface{})
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

		var details []map[string]interface{}
		for _, id := range partList {
			partDetail, err := database.GetPartsByIds(table, []int{id})
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query testing parts")
				return
			}
			if len(partDetail) > 0 {
				details = append(details, partDetail[0])
			}
		}
		record[fieldToNameMap[field]] = details
	}

	// 返回查询结果
	response := map[string]interface{}{
		"data": record,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}

// @Summary 新增测试套题
// @Description 新增测试套题
// @Tags Testing
// @Accept json
// @Produce json
// @Param part body models.BasicTestingItem true "测试套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/testing/add [post]
func AddTesting(c *gin.Context) {
	var part models.BasicTestingItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	if part.Type == 3 {
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
	result, err := database.InsertData("testing_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert testing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新测试套题
// @Description 更新测试套题
// @Tags Testing
// @Accept json
// @Produce json
// @Param part body models.BasicTestingItem true "测试套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/testing/update [put]
func UpdateTesting(c *gin.Context) {
	var part models.BasicTestingItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	if part.Type == 3 {
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
	result, err := database.InsertData("testing_list", &part, "update")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update testing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除测试套题
// @Description 根据ID删除测试套题
// @Tags Testing
// @Accept json
// @Produce json
// @Param id path int true "测试套题ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/testing/delete/{id} [delete]
func DeleteTesting(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid testing set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("testing_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete testing set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "testing set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}