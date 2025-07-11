package repository

import (
	"context"
	"time"

	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/metrics"
	"github.com/google/uuid"
)

// InstrumentedUserRepository wraps UserRepository with metrics
type InstrumentedUserRepository struct {
	base    UserRepository
	metrics *metrics.Metrics
	logger  log.Logger
}

func NewInstrumentedUserRepository(base UserRepository, metrics *metrics.Metrics, logger log.Logger) UserRepository {
	return &InstrumentedUserRepository{
		base:    base,
		metrics: metrics,
		logger:  logger,
	}
}

func (r *InstrumentedUserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error) {
	start := time.Now()
	user, err := r.base.CreateUser(ctx, arg)
	r.metrics.RecordDBQuery("INSERT", "users", time.Since(start), err)
	
	if err == nil {
		r.metrics.UsersTotal.Inc()
	}
	
	return user, err
}

func (r *InstrumentedUserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error) {
	start := time.Now()
	user, err := r.base.GetUser(ctx, userID)
	r.metrics.RecordDBQuery("SELECT", "users", time.Since(start), err)
	return user, err
}

func (r *InstrumentedUserRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	start := time.Now()
	user, err := r.base.GetUserByUsername(ctx, username)
	r.metrics.RecordDBQuery("SELECT", "users", time.Since(start), err)
	return user, err
}

func (r *InstrumentedUserRepository) GetUserByHandle(ctx context.Context, handle string) (*db.User, error) {
	start := time.Now()
	user, err := r.base.GetUserByHandle(ctx, handle)
	r.metrics.RecordDBQuery("SELECT", "users", time.Since(start), err)
	return user, err
}

func (r *InstrumentedUserRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	start := time.Now()
	user, err := r.base.GetUserByEmail(ctx, email)
	r.metrics.RecordDBQuery("SELECT", "users", time.Since(start), err)
	return user, err
}

func (r *InstrumentedUserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	start := time.Now()
	err := r.base.UpdateUser(ctx, arg)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	start := time.Now()
	err := r.base.UpdateUsername(ctx, userID, username)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error {
	start := time.Now()
	err := r.base.UpdateHandle(ctx, userID, handle)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
	start := time.Now()
	err := r.base.UpdateEmail(ctx, userID, email)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error {
	start := time.Now()
	err := r.base.UpdatePremiumStatus(ctx, userID, isPremium)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error {
	start := time.Now()
	err := r.base.UpdateAdminStatus(ctx, userID, isAdmin)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) error {
	start := time.Now()
	err := r.base.UpdateOnboardedStatus(ctx, userID, onboarded)
	r.metrics.RecordDBQuery("UPDATE", "users", time.Since(start), err)
	return err
}

func (r *InstrumentedUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	start := time.Now()
	err := r.base.DeleteUser(ctx, userID)
	r.metrics.RecordDBQuery("DELETE", "users", time.Since(start), err)
	
	if err == nil {
		r.metrics.UsersTotal.Dec()
	}
	
	return err
}