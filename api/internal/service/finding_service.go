package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/model"
	"github.com/quantun-opensource/qrap/api/internal/repository"
)

type FindingService struct {
	repo   *repository.FindingRepository
	logger *zap.Logger
}

func NewFindingService(repo *repository.FindingRepository, logger *zap.Logger) *FindingService {
	return &FindingService{repo: repo, logger: logger}
}

func (s *FindingService) Get(ctx context.Context, id uuid.UUID) (*model.Finding, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FindingService) ListByAssessment(ctx context.Context, assessmentID uuid.UUID, riskLevel, category string, offset, limit int) ([]model.Finding, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListByAssessment(ctx, assessmentID, riskLevel, category, offset, limit)
}
