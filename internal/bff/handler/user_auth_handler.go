package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wangyahua6688-maker/tk-common/utils/codes"
	"github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
)

// SendSMSCode 发送登录/注册短信验证码。
func (h *PublicHandler) SendSMSCode(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体（限制 1MB）。
	var reqBody struct {
		// 处理当前语句逻辑。
		Phone string `json:"phone"`
		// 处理当前语句逻辑。
		Purpose string `json:"purpose"`
	}
	// 定义并初始化当前变量。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodySendSMS, "invalid request body")
		// 返回当前处理结果。
		return
	}

	// 2) 基础参数校验。
	phone := strings.TrimSpace(reqBody.Phone)
	// 判断条件并进入对应分支逻辑。
	if phone == "" {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhoneRequired, "phone is required")
		// 返回当前处理结果。
		return
	}

	// 3) 转发到用户域认证 RPC。
	resp, rpcErr := h.svcCtx.User.SendSMSCode(r.Context(), &tkv1.AuthSendCodeRequest{
		// 处理当前语句逻辑。
		Phone: phone,
		// 调用strings.TrimSpace完成当前处理。
		Purpose: strings.TrimSpace(reqBody.Purpose),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, rpcErr)
}

// RegisterByPhone 手机号注册。
func (h *PublicHandler) RegisterByPhone(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体（限制 1MB）。
	var reqBody struct {
		// 处理当前语句逻辑。
		Phone string `json:"phone"`
		// 处理当前语句逻辑。
		Password string `json:"password"`
		// 处理当前语句逻辑。
		SMSCode string `json:"sms_code"`
		// 处理当前语句逻辑。
		Nickname string `json:"nickname"`
	}
	// 定义并初始化当前变量。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodyReg, "invalid request body")
		// 返回当前处理结果。
		return
	}

	// 2) 校验必要字段。
	phone := strings.TrimSpace(reqBody.Phone)
	// 定义并初始化当前变量。
	password := strings.TrimSpace(reqBody.Password)
	// 判断条件并进入对应分支逻辑。
	if phone == "" || password == "" {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhonePwdRequired, "phone/password required")
		// 返回当前处理结果。
		return
	}

	// 3) 转发用户域注册 RPC。
	resp, rpcErr := h.svcCtx.User.RegisterByPhone(r.Context(), &tkv1.AuthRegisterRequest{
		// 处理当前语句逻辑。
		Phone: phone,
		// 处理当前语句逻辑。
		Password: password,
		// 调用strings.TrimSpace完成当前处理。
		SmsCode: strings.TrimSpace(reqBody.SMSCode),
		// 调用strings.TrimSpace完成当前处理。
		Nickname: strings.TrimSpace(reqBody.Nickname),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, rpcErr)
}

// LoginByPassword 手机号+密码登录。
func (h *PublicHandler) LoginByPassword(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体。
	var reqBody struct {
		// 处理当前语句逻辑。
		Phone string `json:"phone"`
		// 处理当前语句逻辑。
		Password string `json:"password"`
	}
	// 定义并初始化当前变量。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodyLogin, "invalid request body")
		// 返回当前处理结果。
		return
	}

	// 2) 校验必要字段。
	phone := strings.TrimSpace(reqBody.Phone)
	// 定义并初始化当前变量。
	password := strings.TrimSpace(reqBody.Password)
	// 判断条件并进入对应分支逻辑。
	if phone == "" || password == "" {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhonePwdNeed, "phone/password required")
		// 返回当前处理结果。
		return
	}

	// 3) 调用用户域登录 RPC。
	resp, rpcErr := h.svcCtx.User.LoginByPassword(r.Context(), &tkv1.AuthPasswordLoginRequest{
		// 处理当前语句逻辑。
		Phone: phone,
		// 处理当前语句逻辑。
		Password: password,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, rpcErr)
}

// LoginBySMS 手机号+验证码登录。
func (h *PublicHandler) LoginBySMS(w http.ResponseWriter, r *http.Request) {
	// 1) 解析请求体。
	var reqBody struct {
		// 处理当前语句逻辑。
		Phone string `json:"phone"`
		// 处理当前语句逻辑。
		SMSCode string `json:"sms_code"`
	}
	// 定义并初始化当前变量。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(body) == 0 || json.Unmarshal(body, &reqBody) != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthInvalidBodySMS, "invalid request body")
		// 返回当前处理结果。
		return
	}

	// 2) 校验参数。
	phone := strings.TrimSpace(reqBody.Phone)
	// 定义并初始化当前变量。
	smsCode := strings.TrimSpace(reqBody.SMSCode)
	// 判断条件并进入对应分支逻辑。
	if phone == "" || smsCode == "" {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthPhoneCodeRequired, "phone/sms_code required")
		// 返回当前处理结果。
		return
	}

	// 3) 调用用户域短信登录 RPC。
	resp, rpcErr := h.svcCtx.User.LoginBySMS(r.Context(), &tkv1.AuthSMSLoginRequest{
		// 处理当前语句逻辑。
		Phone: phone,
		// 处理当前语句逻辑。
		SmsCode: smsCode,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, rpcErr)
}

// Profile 获取当前登录用户资料。
func (h *PublicHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// 1) 仅从 Authorization 头读取 Bearer token。
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	// 2) 参数校验。
	if token == "" {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.UserAuthAccessTokenNeed, "access token required")
		// 返回当前处理结果。
		return
	}

	// 3) 调用用户域资料 RPC。
	resp, rpcErr := h.svcCtx.User.Profile(r.Context(), &tkv1.AuthProfileRequest{
		// 处理当前语句逻辑。
		AccessToken: token,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, rpcErr)
}
