package models

type ResponseData struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Message string `json:"message"`
}