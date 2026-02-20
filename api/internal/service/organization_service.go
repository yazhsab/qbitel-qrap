package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/model"
	"github.com/quantun-opensource/qrap/api/internal/repository"
)

type OrganizationService struct {
	repo   *repository.OrganizationRepository
	logger *zap.Logger
}

func NewOrganizationService(repo *repository.OrganizationRepository, logger *zap.Logger) *OrganizationService {
	return &OrganizationService{repo: repo, logger: logger}
}

func (s *OrganizationService) Create(ctx context.Context, req *model.CreateOrganizationRequest) (*model.Organization, error) {
	org := &model.Organization{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, org); err != nil {
		s.logger.Error("failed to create organization", zap.Error(err))
		return nil, err
	}

	s.logger.Info("organization created", zap.String("id", org.ID.String()), zap.String("name", org.Name))
	return org, nil
}

func (s *OrganizationService) Get(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationService) List(ctx context.Context, offset, limit int) ([]model.Organization, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.List(ctx, offset, limit)
}
