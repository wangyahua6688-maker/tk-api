package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"tk-common/utils/codes"
	"tk-common/utils/httpresp"
	tkv1 "tk-proto/tk/v1"
)

// SendSMSCode 发送登录/注册短信验证码。
func (h *PublicHandler) SendSMSCode(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体（限制 1MB）。
	var reqBody struct {
		Phone   string `json:"phone"`
		Purpose string `json:"purpose"`
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodySendSMS, "invalid request body")
		return
	}

	// 2) 基础参数校验。
	phone := strings.TrimSpace(reqBody.Phone)
	if phone == "" {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhoneRequired, "phone is required")
		return
	}

	// 3) 转发到用户域认证 RPC。
	resp, rpcErr := h.svcCtx.User.SendSMSCode(r.Context(), &tkv1.AuthSendCodeRequest{
		Phone:   phone,
		Purpose: strings.TrimSpace(reqBody.Purpose),
	})
	writeRPCReply(w, resp, rpcErr)
}

// RegisterByPhone 手机号注册。
func (h *PublicHandler) RegisterByPhone(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体（限制 1MB）。
	var reqBody struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		SMSCode  string `json:"sms_code"`
		Nickname string `json:"nickname"`
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodyReg, "invalid request body")
		return
	}

	// 2) 校验必要字段。
	phone := strings.TrimSpace(reqBody.Phone)
	password := strings.TrimSpace(reqBody.Password)
	if phone == "" || password == "" {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhonePwdRequired, "phone/password required")
		return
	}

	// 3) 转发用户域注册 RPC。
	resp, rpcErr := h.svcCtx.User.RegisterByPhone(r.Context(), &tkv1.AuthRegisterRequest{
		Phone:    phone,
		Password: password,
		SmsCode:  strings.TrimSpace(reqBody.SMSCode),
		Nickname: strings.TrimSpace(reqBody.Nickname),
	})
	writeRPCReply(w, resp, rpcErr)
}

// LoginByPassword 手机号+密码登录。
func (h *PublicHandler) LoginByPassword(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体。
	var reqBody struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodyLogin, "invalid request body")
		return
	}

	// 2) 校验必要字段。
	phone := strings.TrimSpace(reqBody.Phone)
	password := strings.TrimSpace(reqBody.Password)
	if phone == "" || password == "" {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhonePwdNeed, "phone/password required")
		return
	}

	// 3) 调用用户域登录 RPC。
	resp, rpcErr := h.svcCtx.User.LoginByPassword(r.Context(), &tkv1.AuthPasswordLoginRequest{
		Phone:    phone,
		Password: password,
	})
	writeRPCReply(w, resp, rpcErr)
}

// LoginBySMS 手机号+验证码登录。
func (h *PublicHandler) LoginBySMS(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体。
	var reqBody struct {
		Phone   string `json:"phone"`
		SMSCode string `json:"sms_code"`
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodySMS, "invalid request body")
		return
	}

	// 2) 校验参数。
	phone := strings.TrimSpace(reqBody.Phone)
	smsCode := strings.TrimSpace(reqBody.SMSCode)
	if phone == "" || smsCode == "" {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhoneCodeRequired, "phone/sms_code required")
		return
	}

	// 3) 调用用户域短信登录 RPC。
	resp, rpcErr := h.svcCtx.User.LoginBySMS(r.Context(), &tkv1.AuthSMSLoginRequest{
		Phone:   phone,
		SmsCode: smsCode,
	})
	writeRPCReply(w, resp, rpcErr)
}

// Profile 获取当前登录用户资料。
func (h *PublicHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// 1) 优先从 Authorization 头读取 Bearer token。
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	// 2) 兼容 query 传入 access_token（便于联调）。
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("access_token"))
	}
	// 3) 参数校验。
	if token == "" {
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthAccessTokenNeed, "access token required")
		return
	}

	// 4) 调用用户域资料 RPC。
	resp, rpcErr := h.svcCtx.User.Profile(r.Context(), &tkv1.AuthProfileRequest{
		AccessToken: token,
	})
	writeRPCReply(w, resp, rpcErr)
}
