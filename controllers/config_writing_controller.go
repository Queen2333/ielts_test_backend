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
// @Router /config/writing/list [get]
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

// @Summary 获取写作套题详情
// @Description 根据id获取写作详情
// @Tags Writing
// @Accept json
// @Produce json
// @Param id query int true "写作id"
// @Success 200 {object} models.ResponseData{data=models.WritingItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing/detail/{id} [get]
func WritingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid writing set ID")
		return
	}

	record, err := database.GetDataById("writing_list", id)
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

// @Summary 新增写作套题
// @Description 新增写作套题
// @Tags Writing
// @Accept json
// @Produce json
// @Param part body models.BasicWritingItem true "写作套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
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

// @Summary 更新写作套题
// @Description 更新写作套题
// @Tags Writing
// @Accept json
// @Produce json
// @Param part body models.BasicWritingItem true "写作套题内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing/update [put]
func UpdateWriting(c *gin.Context) {
	var part models.BasicWritingItem
	if err := c.ShouldBindJSON(&part); err != nil {
		fmt.Println(err)
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("writing_list", &part, "update")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update writing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除写作套题
// @Description 根据ID删除写作套题
// @Tags Writing
// @Accept json
// @Produce json
// @Param id path int true "写作套题ID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing/delete/{id} [delete]
func DeleteWriting(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid writing set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("writing_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete writing set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "writing set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}

// @Summary      获取写作篇列表
// @Description  根据条件获取写作part列表，并返回分页结果
// @Tags         Writing
// @Accept       json
// @Produce      json
// @Param        name      query  string  false  "试题名称"
// @Param        status    query  int     false  "试题状态"
// @Param        type      query  int     false  "试题类型"
// @Param        pageNo    query  int     true   "页码"
// @Param        pageLimit query  int     false   "每页条数"
// @Success      200  {object}  models.ResponseData{data=models.WritingPartListResponse}
// @Failure      400  {object}  models.ResponseData{data=nil}
// @Failure      500  {object}  models.ResponseData{data=nil}
// @Router       /config/writing-part/list [get]
func WritingPartList(c * gin.Context) {
	var request struct {
		Name      string `json:"name,omitempty"`
		Status    int    `json:"status,omitempty"`
		Type      int    `json:"type,omitempty"`
		Source 	  int	 `json:"source,omitempty"`
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
    results, total, err := database.PaginationQuery("writing_part_list", request.PageNo, request.PageLimit, conditions)
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

// @Summary 获取写作part详情
// @Description 根据id获取写作part详情
// @Tags Writing
// @Accept json
// @Produce json
// @Param id query int true "写作part id"
// @Success 200 {object} models.ResponseData{data=models.WritingPartItem}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing-part/detail/{id} [get]
func WritingPartDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid writing part set ID")
		return
	}

	record, err := database.GetDataById("writing_part_list", id)
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

// @Summary 新增写作part
// @Description 新增写作part
// @Tags Writing
// @Accept json
// @Produce json
// @Param part body models.WritingPartItem true "写作part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing-part/add [post]
func AddWritingPart(c *gin.Context) {
	var part models.WritingPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("writing_part_list", &part, "create")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to insert writing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 更新写作part
// @Description 更新写作part
// @Tags Writing
// @Accept json
// @Produce json
// @Param part body models.WritingPartItem true "写作part内容"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing-part/update [put]
func UpdateWritingPart(c *gin.Context) {
	var part models.WritingPartItem
	if err := c.ShouldBindJSON(&part); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 将数据插入数据库
	result, err := database.InsertData("writing_part_list", &part, "update")
	fmt.Println(err, "err")
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to update writing part")
		return
	}

	// 返回插入后的数据
	utils.HandleResponse(c, http.StatusOK, result, "Success")
}

// @Summary 删除写作part
// @Description 根据ID删除写作part
// @Tags Writing
// @Accept json
// @Produce json
// @Param id path int true "写作partID"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 404 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /config/writing-part/delete/{id} [delete]
func DeleteWritingPart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid writing set ID")
		return
	}

	// 执行删除操作
	rowsAffected, err := database.DeleteData("writing_part_list", id)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to delete writing set")
		return
	}

	// 检查是否有记录被删除
	if rowsAffected == 0 {
		utils.HandleResponse(c, http.StatusNotFound, "", "writing set not found")
		return
	}

	// 返回成功响应
	utils.HandleResponse(c, http.StatusOK, nil, "Success")
}