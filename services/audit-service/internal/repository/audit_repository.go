package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/audit-service/internal/domain"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	var log domain.AuditLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &log, err
}

func (r *auditLogRepository) FindByActorID(ctx context.Context, actorID string, limit, offset int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64
	r.db.WithContext(ctx).Model(&domain.AuditLog{}).Where("actor_id = ?", actorID).Count(&total)
	err := r.db.WithContext(ctx).Where("actor_id = ?", actorID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *auditLogRepository) FindByResourceID(ctx context.Context, resourceID string, limit, offset int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64
	r.db.WithContext(ctx).Model(&domain.AuditLog{}).Where("resource_id = ?", resourceID).Count(&total)
	err := r.db.WithContext(ctx).Where("resource_id = ?", resourceID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *auditLogRepository) FindByServiceAndEvent(ctx context.Context, service, event string, from, to time.Time, limit, offset int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.AuditLog{})
	if service != "" {
		query = query.Where("service_name = ?", service)
	}
	if event != "" {
		query = query.Where("event_type = ?", event)
	}
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}