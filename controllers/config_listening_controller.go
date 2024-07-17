package controllers

import (
	"fmt"
	"net/http"

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
// @Router /config/listening/list [post]
func ListeningList(c *gin.Context) {

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
    results, total, err := database.PaginationQuery("listening_list", request.PageNo, request.PageLimit, conditions)
    if err != nil {
        utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to execute pagination query")
        return
    }

	// 遍历 results 提取所有的 part_list ID
	var partIDs []int
	for _, result := range results {
		// 从结果中提取 part_list 字段并解析为整数数组
		if partListStr, ok := result["part_list"].(string); ok {
			partList := utils.StringToList(partListStr)
			// 将 partList 中的 ID 添加到 partIDs 中
			partIDs = append(partIDs, partList...)
		}
	}


	// 查询 part_list 中的详细信息
	if len(partIDs) > 0 {
		partDetails, err := database.GetPartsByIds("listening_part_list", partIDs)
		if err != nil {
			utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query listening parts")
			return
		}

		/// 将查询结果放回到对应的 part_list 中
		for i := range results {
			var details []map[string]interface{}
			for _, partID := range partIDs {
				for _, partDetail := range partDetails {
					id, ok := partDetail["id"].(int)
					if ok && id == partID {
						details = append(details, partDetail)
						break
					}
				}
			}
			results[i]["part_list"] = details
		}
	}


	// 返回查询结果
	response := map[string]interface{}{
		"items": results,
		"total":   total,
	}
	utils.HandleResponse(c, http.StatusOK, response, "Success")
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
// @Router       /config/listening-part/list [post]
func ListeningPartList(c * gin.Context) {
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
    results, total, err := database.PaginationQuery("listening_part_list", request.PageNo, request.PageLimit, conditions)
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

// @Summary 新增听力part
// @Description 新增听力part
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.ListeningPartItem true "听力part内容"
// @Success 200 {object} models.ResponseData{data=models.ListeningPartItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/add [post]
func AddListeningPart(c *gin.Context) {
	var part models.ListeningPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

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
// @Success 200 {object} models.ResponseData{data=models.ListeningPartItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening-part/update [put]
func UpdateListeningPart(c *gin.Context) {
	var part models.ListeningPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("listening_part_list", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}