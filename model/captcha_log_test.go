package model

import "testing"

func TestCaptchaLog_TableName(t *testing.T) {
	log := CaptchaLog{}
	if log.TableName() != "captcha_logs" {
		t.Errorf("TableName() = %v, want 'captcha_logs'", log.TableName())
	}
}
