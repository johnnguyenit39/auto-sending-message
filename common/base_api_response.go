package common

type BaseApiResponse[T any] struct {
	Success           bool   `json:"success"`
	HttpRequestStatus int    `json:"http_request_status"`
	Message           string `json:"message"`
	UserMessage       string `json:"user_message"`
	Data              T      `json:"data"`
}
