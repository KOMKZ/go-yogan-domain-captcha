package service

import (
	"context"

	captchaerrors "github.com/KOMKZ/go-yogan-domain-captcha/errors"
	"github.com/KOMKZ/go-yogan-domain-captcha/model"
	"github.com/KOMKZ/go-yogan-domain-captcha/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

type CaptchaVerifier interface {
	Verify(ctx context.Context, param string) (success bool, requestID string, err error)
}

type CaptchaConfig struct {
	Enabled   bool
	EnableLog bool
}

type VerifyInput struct {
	CaptchaVerifyParam string
	Scene              string
	IP                 string
	UserAgent          string
}

type CaptchaService struct {
	verifier CaptchaVerifier
	repo     repository.CaptchaLogRepository
	config   CaptchaConfig
	logger   *logger.CtxZapLogger
}

func NewCaptchaService(
	verifier CaptchaVerifier,
	repo repository.CaptchaLogRepository,
	config CaptchaConfig,
	log *logger.CtxZapLogger,
) *CaptchaService {
	return &CaptchaService{
		verifier: verifier,
		repo:     repo,
		config:   config,
		logger:   log,
	}
}

func (s *CaptchaService) IsEnabled() bool {
	return s.config.Enabled
}

func (s *CaptchaService) Verify(ctx context.Context, input VerifyInput) error {
	success := true
	requestID := ""
	var verifyErr error

	if s.config.Enabled {
		if input.CaptchaVerifyParam == "" {
			return captchaerrors.ErrCaptchaRequired
		}
		success, requestID, verifyErr = s.verifier.Verify(ctx, input.CaptchaVerifyParam)
	}

	s.logger.InfoCtx(ctx, "captcha verify attempt",
		zap.String("scene", input.Scene),
		zap.String("ip", input.IP),
		zap.Bool("enabled", s.config.Enabled),
		zap.Bool("success", success),
		zap.String("request_id", requestID),
		zap.String("param", input.CaptchaVerifyParam),
	)

	if s.repo != nil && input.CaptchaVerifyParam != "" {
		logEntry := &model.CaptchaLog{
			Scene:       input.Scene,
			IP:          input.IP,
			UserAgent:   input.UserAgent,
			Success:     success,
			RequestID:   requestID,
			VerifyParam: input.CaptchaVerifyParam,
		}
		if logErr := s.repo.Create(ctx, logEntry); logErr != nil {
			s.logger.ErrorCtx(ctx, "failed to save captcha log", zap.Error(logErr))
		}
	}

	if s.config.Enabled {
		if verifyErr != nil {
			s.logger.WarnCtx(ctx, "captcha verify error",
				zap.String("scene", input.Scene),
				zap.Error(verifyErr),
			)
			return captchaerrors.ErrCaptchaServiceUnavailable
		}
		if !success {
			s.logger.WarnCtx(ctx, "captcha verify failed",
				zap.String("scene", input.Scene),
				zap.Error(verifyErr),
			)
			return captchaerrors.ErrCaptchaVerifyFailed
		}
	}

	return nil
}

type ListLogsInput struct {
	Scene   string
	Success *bool
	IP      string
}

func (s *CaptchaService) ListLogs(ctx context.Context, page, size int, input ListLogsInput) ([]*model.CaptchaLog, int64, error) {
	if s.repo == nil {
		return nil, 0, captchaerrors.ErrCaptchaServiceUnavailable
	}
	filters := repository.CaptchaLogFilters{
		Scene:   input.Scene,
		Success: input.Success,
		IP:      input.IP,
	}
	return s.repo.Paginate(ctx, page, size, filters)
}
