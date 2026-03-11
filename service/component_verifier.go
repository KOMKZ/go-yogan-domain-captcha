package service

import (
	"context"

	captchacomp "github.com/KOMKZ/go-yogan-component-captcha"
)

type ComponentVerifier struct {
	manager *captchacomp.Manager
}

func NewComponentVerifier(manager *captchacomp.Manager) *ComponentVerifier {
	return &ComponentVerifier{manager: manager}
}

func (v *ComponentVerifier) Verify(ctx context.Context, param string) (bool, string, error) {
	result, err := v.manager.Verify(ctx, param)
	if err != nil {
		if result != nil {
			return result.Success, result.RequestID, err
		}
		return false, "", err
	}
	return result.Success, result.RequestID, nil
}
