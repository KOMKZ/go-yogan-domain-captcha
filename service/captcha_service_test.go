package service

import (
	"context"
	"errors"
	"testing"

	captchaerrors "github.com/KOMKZ/go-yogan-domain-captcha/errors"
	"github.com/KOMKZ/go-yogan-domain-captcha/model"
	"github.com/KOMKZ/go-yogan-domain-captcha/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
)

type mockVerifier struct {
	success   bool
	requestID string
	err       error
}

func (m *mockVerifier) Verify(ctx context.Context, param string) (bool, string, error) {
	return m.success, m.requestID, m.err
}

type mockRepo struct {
	logs        []*model.CaptchaLog
	err         error
	created     bool
	paginateErr error
}

func (m *mockRepo) Create(ctx context.Context, log *model.CaptchaLog) error {
	m.created = true
	m.logs = append(m.logs, log)
	return m.err
}

func (m *mockRepo) Paginate(_ context.Context, page, size int, _ repository.CaptchaLogFilters) ([]*model.CaptchaLog, int64, error) {
	if m.paginateErr != nil {
		return nil, 0, m.paginateErr
	}
	return m.logs, int64(len(m.logs)), nil
}

func getTestLogger() *logger.CtxZapLogger {
	return logger.GetLogger("test")
}

func TestCaptchaService_Verify_Disabled(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{},
		&mockRepo{},
		CaptchaConfig{Enabled: false},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "test-param",
		Scene:              "login",
	})
	if err != nil {
		t.Errorf("Verify() should return nil when disabled, got %v", err)
	}
}

func TestCaptchaService_Verify_EmptyParam(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{},
		&mockRepo{},
		CaptchaConfig{Enabled: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "",
		Scene:              "login",
	})
	if !errors.Is(err, captchaerrors.ErrCaptchaRequired) {
		t.Errorf("Verify() should return ErrCaptchaRequired, got %v", err)
	}
}

func TestCaptchaService_Verify_Success(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{success: true, requestID: "req-123"},
		&mockRepo{},
		CaptchaConfig{Enabled: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "valid-param",
		Scene:              "login",
	})
	if err != nil {
		t.Errorf("Verify() should return nil on success, got %v", err)
	}
}

func TestCaptchaService_Verify_Failed(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{success: false, requestID: "req-456"},
		&mockRepo{},
		CaptchaConfig{Enabled: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "invalid-param",
		Scene:              "login",
	})
	if !errors.Is(err, captchaerrors.ErrCaptchaVerifyFailed) {
		t.Errorf("Verify() should return ErrCaptchaVerifyFailed, got %v", err)
	}
}

func TestCaptchaService_Verify_ServiceError(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{success: false, err: errors.New("network error")},
		&mockRepo{},
		CaptchaConfig{Enabled: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "test-param",
		Scene:              "login",
	})
	if !errors.Is(err, captchaerrors.ErrCaptchaServiceUnavailable) {
		t.Errorf("Verify() should return ErrCaptchaServiceUnavailable, got %v", err)
	}
}

func TestCaptchaService_Verify_WithLog(t *testing.T) {
	repo := &mockRepo{}
	svc := NewCaptchaService(
		&mockVerifier{success: true, requestID: "req-789"},
		repo,
		CaptchaConfig{Enabled: true, EnableLog: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "valid-param",
		Scene:              "login",
		IP:                 "127.0.0.1",
		UserAgent:          "test-agent",
	})
	if err != nil {
		t.Errorf("Verify() should return nil on success, got %v", err)
	}

	if !repo.created {
		t.Error("Verify() should create log entry when EnableLog=true")
	}
	if len(repo.logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(repo.logs))
	}
	if repo.logs[0].Scene != "login" {
		t.Errorf("log.Scene = %v, want 'login'", repo.logs[0].Scene)
	}
	if !repo.logs[0].Success {
		t.Error("log.Success should be true")
	}
	if repo.logs[0].RequestID != "req-789" {
		t.Errorf("log.RequestID = %v, want 'req-789'", repo.logs[0].RequestID)
	}
	if repo.logs[0].VerifyParam != "valid-param" {
		t.Errorf("log.VerifyParam = %v, want 'valid-param'", repo.logs[0].VerifyParam)
	}
}

func TestCaptchaService_Verify_AlwaysLogsWhenParamPresent(t *testing.T) {
	repo := &mockRepo{}
	svc := NewCaptchaService(
		&mockVerifier{success: true, requestID: "req-000"},
		repo,
		CaptchaConfig{Enabled: true, EnableLog: false},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "valid-param",
		Scene:              "login",
	})
	if err != nil {
		t.Errorf("Verify() should return nil, got %v", err)
	}
	if !repo.created {
		t.Error("Verify() should always log when param is present and repo is available")
	}
}

func TestCaptchaService_Verify_LogError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db error")}
	svc := NewCaptchaService(
		&mockVerifier{success: true, requestID: "req-err"},
		repo,
		CaptchaConfig{Enabled: true, EnableLog: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "valid-param",
		Scene:              "login",
		IP:                 "127.0.0.1",
	})
	if err != nil {
		t.Errorf("Verify() should still succeed even if log fails, got %v", err)
	}
}

func TestCaptchaService_Verify_FailedWithLog(t *testing.T) {
	repo := &mockRepo{}
	svc := NewCaptchaService(
		&mockVerifier{success: false, requestID: "req-fail"},
		repo,
		CaptchaConfig{Enabled: true, EnableLog: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "bad-param",
		Scene:              "login",
		IP:                 "10.0.0.1",
	})
	if !errors.Is(err, captchaerrors.ErrCaptchaVerifyFailed) {
		t.Errorf("Verify() should return ErrCaptchaVerifyFailed, got %v", err)
	}
	if !repo.created {
		t.Error("Verify() should log failed attempts too")
	}
	if repo.logs[0].Success {
		t.Error("log.Success should be false for failed verification")
	}
}

func TestCaptchaService_Verify_DisabledStillLogs(t *testing.T) {
	repo := &mockRepo{}
	svc := NewCaptchaService(
		&mockVerifier{},
		repo,
		CaptchaConfig{Enabled: false},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "test-param",
		Scene:              "login",
		IP:                 "127.0.0.1",
	})
	if err != nil {
		t.Errorf("Verify() should return nil when disabled, got %v", err)
	}
	if !repo.created {
		t.Error("Verify() should log even when disabled, if param is present")
	}
	if len(repo.logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(repo.logs))
	}
	if !repo.logs[0].Success {
		t.Error("log.Success should be true when disabled (no verification performed)")
	}
}

func TestCaptchaService_Verify_NoLogWhenParamEmpty(t *testing.T) {
	repo := &mockRepo{}
	svc := NewCaptchaService(
		&mockVerifier{},
		repo,
		CaptchaConfig{Enabled: false},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "",
		Scene:              "login",
	})
	if err != nil {
		t.Errorf("Verify() should return nil, got %v", err)
	}
	if repo.created {
		t.Error("Verify() should NOT log when param is empty")
	}
}

func TestCaptchaService_IsEnabled(t *testing.T) {
	svc1 := NewCaptchaService(nil, nil, CaptchaConfig{Enabled: true}, getTestLogger())
	if !svc1.IsEnabled() {
		t.Error("IsEnabled() should return true")
	}

	svc2 := NewCaptchaService(nil, nil, CaptchaConfig{Enabled: false}, getTestLogger())
	if svc2.IsEnabled() {
		t.Error("IsEnabled() should return false")
	}
}

func TestCaptchaService_Verify_NilRepo_WithLog(t *testing.T) {
	svc := NewCaptchaService(
		&mockVerifier{success: true, requestID: "req-nil"},
		nil,
		CaptchaConfig{Enabled: true, EnableLog: true},
		getTestLogger(),
	)

	err := svc.Verify(context.Background(), VerifyInput{
		CaptchaVerifyParam: "valid-param",
		Scene:              "login",
	})
	if err != nil {
		t.Errorf("Verify() should succeed even with nil repo, got %v", err)
	}
}

func TestCaptchaService_ListLogs_Success(t *testing.T) {
	repo := &mockRepo{
		logs: []*model.CaptchaLog{
			{ID: 1, Scene: "login", Success: true},
			{ID: 2, Scene: "login", Success: false},
		},
	}
	svc := NewCaptchaService(nil, repo, CaptchaConfig{Enabled: true, EnableLog: true}, getTestLogger())

	logs, total, err := svc.ListLogs(context.Background(), 1, 10, ListLogsInput{})
	if err != nil {
		t.Fatalf("ListLogs() error = %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(logs) != 2 {
		t.Errorf("len(logs) = %d, want 2", len(logs))
	}
}

func TestCaptchaService_ListLogs_NilRepo(t *testing.T) {
	svc := NewCaptchaService(nil, nil, CaptchaConfig{Enabled: true}, getTestLogger())

	_, _, err := svc.ListLogs(context.Background(), 1, 10, ListLogsInput{})
	if !errors.Is(err, captchaerrors.ErrCaptchaServiceUnavailable) {
		t.Errorf("ListLogs() should return ErrCaptchaServiceUnavailable, got %v", err)
	}
}

func TestCaptchaService_ListLogs_WithFilters(t *testing.T) {
	repo := &mockRepo{
		logs: []*model.CaptchaLog{{ID: 1, Scene: "login", Success: true}},
	}
	svc := NewCaptchaService(nil, repo, CaptchaConfig{Enabled: true, EnableLog: true}, getTestLogger())
	success := true
	logs, total, err := svc.ListLogs(context.Background(), 1, 10, ListLogsInput{
		Scene:   "login",
		Success: &success,
	})
	if err != nil {
		t.Fatalf("ListLogs() error = %v", err)
	}
	if total != 1 || len(logs) != 1 {
		t.Errorf("expected 1 result, got total=%d len=%d", total, len(logs))
	}
}

func TestCaptchaService_ListLogs_RepoError(t *testing.T) {
	repo := &mockRepo{paginateErr: errors.New("db error")}
	svc := NewCaptchaService(nil, repo, CaptchaConfig{Enabled: true, EnableLog: true}, getTestLogger())

	_, _, err := svc.ListLogs(context.Background(), 1, 10, ListLogsInput{})
	if err == nil {
		t.Error("ListLogs() should propagate repo error")
	}
}
