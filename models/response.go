package models

type ResponseData struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

//听力part列表返回体
type ListeningPartListResponse struct {
    Items []ListeningPartItem `json:"items"`
    Total int                 `json:"total"`
}

//听力套题列表返回体
type ListeningListResponse struct {
    Items []ListeningItem 	`json:"items"`
    Total int               `json:"total"`
}

//阅读套题列表返回体
type ReadingListResponse struct {
    Items []ReadingItem 	`json:"items"`
    Total int               `json:"total"`
}

//阅读part列表返回体
type ReadingPartListResponse struct {
    Items []ReadingPartItem `json:"items"`
    Total int               `json:"total"`
}

//写作套题列表返回体
type WritingListResponse struct {
    Items []WritingItem     `json:"items"`
    Total int               `json:"total"`
}

//写作part列表返回体
type WritingPartListResponse struct {
    Items []WritingPartItem `json:"items"`
    Total int               `json:"total"`
}

// 上传文件返回体
type UploadResponse struct {
    Url   string            `json:"url"`
}

// 测试套题返回体
type TestingListResponse struct {
    Items []TestingItem     `json:"items"`
    Total int               `json:"total"`
}

type ListeningRecordsResponse struct {
    Items []ListeningRecordsItem 	`json:"items"`
    Total int                       `json:"total"`
}