// service/instrumented_auth_service.go
package service

import (
	"context"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/metrics"
	"github.com/0xsj/mios.io/pkg/token"
)

// InstrumentedAuthService wraps AuthService with metrics
type InstrumentedAuthService struct {
	base    AuthService
	metrics *metrics.Metrics
	logger  log.Logger
}

func NewInstrumentedAuthService(base AuthService, metrics *metrics.Metrics, logger log.Logger) AuthService {
	return &InstrumentedAuthService{
		base:    base,
		metrics: metrics,
		logger:  logger,
	}
}

func (s *InstrumentedAuthService) Register(ctx context.Context, input RegisterInput) (*UserDTO, error) {
	user, err := s.base.Register(ctx, input)
	
	if err != nil {
		s.metrics.RecordError("registration_failure", "auth_service", "warning")
	} else {
		s.metrics.UsersTotal.Inc()
	}
	
	return user, err
}

func (s *InstrumentedAuthService) Login(ctx context.Context, input LoginInput) (*TokenResponse, error) {
	response, err := s.base.Login(ctx, input)
	
	if err != nil {
		s.metrics.RecordError("login_failure", "auth_service", "info")
	}
	
	return response, err
}

func (s *InstrumentedAuthService) RefreshToken(ctx context.Context, input RefreshTokenRequest) (*TokenResponse, error) {
	response, err := s.base.RefreshToken(ctx, input)
	
	if err != nil {
		s.metrics.RecordError("token_refresh_failure", "auth_service", "warning")
	}
	
	return response, err
}

func (s *InstrumentedAuthService) GenerateResetToken(ctx context.Context, email string) error {
	err := s.base.GenerateResetToken(ctx, email)
	
	if err != nil {
		s.metrics.RecordError("reset_token_generation_failure", "auth_service", "warning")
	}
	
	return err
}

func (s *InstrumentedAuthService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	err := s.base.ResetPassword(ctx, input)
	
	if err != nil {
		s.metrics.RecordError("password_reset_failure", "auth_service", "warning")
	}
	
	return err
}

func (s *InstrumentedAuthService) Logout(ctx context.Context, userID string) error {
	err := s.base.Logout(ctx, userID)
	
	if err != nil {
		s.metrics.RecordError("logout_failure", "auth_service", "warning")
	}
	
	return err
}

func (s *InstrumentedAuthService) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	claims, err := s.base.ValidateToken(ctx, tokenStr)
	
	if err != nil {
		s.metrics.RecordError("token_validation_failure", "auth_service", "info")
	}
	
	return claims, err
}

func (s *InstrumentedAuthService) IsEmailVerified(ctx context.Context, userID string) (bool, error) {
	verified, err := s.base.IsEmailVerified(ctx, userID)
	
	if err != nil {
		s.metrics.RecordError("email_verification_check_failure", "auth_service", "warning")
	}
	
	return verified, err
}

func (s *InstrumentedAuthService) SendVerificationEmail(ctx context.Context, email, username, token string) error {
	err := s.base.SendVerificationEmail(ctx, email, username, token)
	
	status := "success"
	if err != nil {
		status = "failure"
		s.metrics.RecordError("email_send_failure", "auth_service", "warning")
	}
	
	s.metrics.RecordEmailSent("verification", status)
	return err
}

func (s *InstrumentedAuthService) SendPasswordResetEmail(ctx context.Context, email, username, token string) error {
	err := s.base.SendPasswordResetEmail(ctx, email, username, token)
	
	status := "success"
	if err != nil {
		status = "failure"
		s.metrics.RecordError("email_send_failure", "auth_service", "warning")
	}
	
	s.metrics.RecordEmailSent("password_reset", status)
	return err
}

func (s *InstrumentedAuthService) SendPasswordChangedEmail(ctx context.Context, email, username string) error {
	err := s.base.SendPasswordChangedEmail(ctx, email, username)
	
	status := "success"
	if err != nil {
		status = "failure"
		s.metrics.RecordError("email_send_failure", "auth_service", "warning")
	}
	
	s.metrics.RecordEmailSent("password_changed", status)
	return err
}

func (s *InstrumentedAuthService) SendAccountLockedEmail(ctx context.Context, email, username, unlockTime string) error {
	err := s.base.SendAccountLockedEmail(ctx, email, username, unlockTime)
	
	status := "success"
	if err != nil {
		status = "failure"
		s.metrics.RecordError("email_send_failure", "auth_service", "warning")
	}
	
	s.metrics.RecordEmailSent("account_locked", status)
	return err
}