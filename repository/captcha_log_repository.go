package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-captcha/model"
)

type CaptchaLogFilters struct {
	Scene   string
	Success *bool
	IP      string
}

type CaptchaLogRepository interface {
	Create(ctx context.Context, log *model.CaptchaLog) error
	Paginate(ctx context.Context, page, size int, filters CaptchaLogFilters) ([]*model.CaptchaLog, int64, error)
}
