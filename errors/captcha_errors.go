package errors

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

const ModuleCaptcha = 36

var (
	ErrCaptchaRequired = errcode.Register(errcode.New(
		ModuleCaptcha, 1001, "captcha",
		"error.captcha.required", "需要人机验证",
		http.StatusBadRequest,
	))
	ErrCaptchaVerifyFailed = errcode.Register(errcode.New(
		ModuleCaptcha, 1002, "captcha",
		"error.captcha.verify_failed", "人机验证未通过",
		http.StatusBadRequest,
	))
	ErrCaptchaServiceUnavailable = errcode.Register(errcode.New(
		ModuleCaptcha, 1003, "captcha",
		"error.captcha.service_unavailable", "验证服务不可用",
		http.StatusServiceUnavailable,
	))
)
