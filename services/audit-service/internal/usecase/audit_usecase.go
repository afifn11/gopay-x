package usecase

import (
	"context"
	"time"

	"github.com/afifn11/gopay-x/services/audit-service/internal/domain"
)

type AuditUsecase interface {
	RecordLog(ctx context.Context, log *domain.AuditLog) error
	GetByActorID(ctx context.Context, actorID string, page, limit int) ([]domain.AuditLog, int64, error)
	GetByResourceID(ctx context.Context, resourceID string, page, limit int) ([]domain.AuditLog, int64, error)
	QueryLogs(ctx context.Context, req *QueryRequest) ([]domain.AuditLog, int64, error)
}

type QueryRequest struct {
	Service string
	Event   string
	From    string
	To      string
	Page    int
	Limit   int
}

type auditUsecase struct {
	auditRepo domain.AuditLogRepository
}

func NewAuditUsecase(auditRepo domain.AuditLogRepository) AuditUsecase {
	return &auditUsecase{auditRepo: auditRepo}
}

func (uc *auditUsecase) RecordLog(ctx context.Context, log *domain.AuditLog) error {
	return uc.auditRepo.Create(ctx, log)
}

func (uc *auditUsecase) GetByActorID(ctx context.Context, actorID string, page, limit int) ([]domain.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.auditRepo.FindByActorID(ctx, actorID, limit, (page-1)*limit)
}

func (uc *auditUsecase) GetByResourceID(ctx context.Context, resourceID string, page, limit int) ([]domain.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.auditRepo.FindByResourceID(ctx, resourceID, limit, (page-1)*limit)
}

func (uc *auditUsecase) QueryLogs(ctx context.Context, req *QueryRequest) ([]domain.AuditLog, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	var from, to time.Time
	if req.From != "" {
		from, _ = time.Parse("2006-01-02", req.From)
	}
	if req.To != "" {
		to, _ = time.Parse("2006-01-02", req.To)
	}

	return uc.auditRepo.FindByServiceAndEvent(ctx, req.Service, req.Event, from, to, req.Limit, (req.Page-1)*req.Limit)
}