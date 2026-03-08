package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-captcha/model"
	"gorm.io/gorm"
)

type CaptchaLogMySQLRepository struct {
	db *gorm.DB
}

func NewCaptchaLogMySQLRepository(db *gorm.DB) *CaptchaLogMySQLRepository {
	return &CaptchaLogMySQLRepository{db: db}
}

func (r *CaptchaLogMySQLRepository) Create(ctx context.Context, log *model.CaptchaLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *CaptchaLogMySQLRepository) Paginate(ctx context.Context, page, size int, filters CaptchaLogFilters) ([]*model.CaptchaLog, int64, error) {
	var total int64
	var logs []*model.CaptchaLog

	q := r.db.WithContext(ctx).Model(&model.CaptchaLog{})
	if filters.Scene != "" {
		q = q.Where("scene = ?", filters.Scene)
	}
	if filters.Success != nil {
		q = q.Where("success = ?", *filters.Success)
	}
	if filters.IP != "" {
		q = q.Where("ip = ?", filters.IP)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := q.Order("id DESC").Offset(offset).Limit(size).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
