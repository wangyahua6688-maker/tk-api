package httpresp

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Envelope 统一响应结构，保持与 tk-web 协议兼容。
type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// OK 返回成功响应。
func OK(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Envelope{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// Fail 返回失败响应。
func Fail(w http.ResponseWriter, statusCode, code int, msg string) {
	httpx.WriteJson(w, statusCode, Envelope{
		Code: code,
		Msg:  msg,
	})
}
