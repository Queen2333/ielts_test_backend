package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Queen2333/ielts_test_backend/database"
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

	var request struct {
		Name      string `json:"name,omitempty"`
		Status    int    `json:"status,omitempty"`
		Type      int    `json:"type,omitempty"`
		PageNo    int    `json:"pageNo"`
		PageLimit int    `json:"pageLimit,omitempty"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 构建查询条件
    conditions := make(map[string]interface{})
	if request.Name != "" {
        conditions["name"] = request.Name
    }

	// 执行分页查询
    results, total, err := database.PaginationQuery("testing_list", request.PageNo, request.PageLimit, conditions)
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

	// 返回查询结果
	response := map[string]interface{}{
		"data": record,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
}