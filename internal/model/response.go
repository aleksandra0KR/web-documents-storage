package model

type ErrorResponse struct {
	Status  int    `json:"code"`
	Message string `json:"text"`
}
