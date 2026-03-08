package do

import (
	"github.com/KOMKZ/go-yogan-domain-captcha/repository"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

// ---- Repository Providers ----

func ProvideCaptchaLogRepository(i do.Injector) (repository.CaptchaLogRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewCaptchaLogMySQLRepository(db), nil
}
