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
    Total int                 `json:"total"`
}