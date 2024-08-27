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
			partDetail, err := database.GetPartsByIds("listening_part_list", []int{id})
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to query listening parts")
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

// @Summary 新增听力套题
// @Description 新增听力套题
// @Tags Listening
// @Accept json
// @Produce json
// @Param part body models.BasicListeningItem true "听力套题内容"
// @Success 200 {object} models.ResponseData{data=models.BasicListeningItem}
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
// @Success 200 {object} models.ResponseData{data=models.BasicListeningItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/listening/update [post]
func UpdateListening(c *gin.Context) {
	var part models.BasicListeningItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("listening_list", &part, "update")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update listening part")
		return
	}

	// 返回插入后的数据
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