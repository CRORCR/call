package base

type Response struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Succeed      bool        `json:"succeed"`
	Data         interface{} `json:"data"`
}
