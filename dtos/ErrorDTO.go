package dtos

type ErrorDTO struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
}
