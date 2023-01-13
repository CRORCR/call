package base

type Response struct {
	ErrorCode    int64       `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Succeed      bool        `json:"succeed"`
	Data         interface{} `json:"data"`
}
