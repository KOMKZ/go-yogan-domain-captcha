package errors

import "testing"

func TestModuleCaptchaCode(t *testing.T) {
	if ModuleCaptcha != 36 {
		t.Errorf("ModuleCaptcha = %d, want 36", ModuleCaptcha)
	}
}

func TestErrorsNotNil(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrCaptchaRequired", ErrCaptchaRequired},
		{"ErrCaptchaVerifyFailed", ErrCaptchaVerifyFailed},
		{"ErrCaptchaServiceUnavailable", ErrCaptchaServiceUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
		})
	}
}
