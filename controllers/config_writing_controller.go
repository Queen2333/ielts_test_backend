package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 获取写作套题列表
// @Description 根据条件获取写作列表，并返回分页结果
// @Tags Writing
// @Accept json
// @Produce json
// @Param name query string false "写作名称"
// @Param status query int false "写作状态"
// @Param type query int false "试题类型"
// @Param pageNo query int true "页码"
// @Param pageLimit query int false "每页条数"
// @Success 200 {object} models.ResponseData{data=models.WritingListResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing/list [post]
func WritingList(c *gin.Context) {

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
    results, total, err := database.PaginationQuery("writing_list", request.PageNo, request.PageLimit, conditions)
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
			partDetail, err := database.GetPartsByIds("writing_part_list", []int{id})
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query writing parts")
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

// @Summary 新增写作套题
// @Description 新增写作套题
// @Tags Writing
// @Accept json
// @Produce json
// @Param part body models.BasicWritingItem true "写作套题内容"
// @Success 200 {object} models.ResponseData{data=models.BasicWritingItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing/add [post]
func AddWriting(c *gin.Context) {
	var part models.BasicWritingItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("writing_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert writing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}