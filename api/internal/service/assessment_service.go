package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/model"
	"github.com/quantun-opensource/qrap/api/internal/repository"
)

type AssessmentService struct {
	assessmentRepo *repository.AssessmentRepository
	findingRepo    *repository.FindingRepository
	logger         *zap.Logger
}

func NewAssessmentService(
	assessmentRepo *repository.AssessmentRepository,
	findingRepo *repository.FindingRepository,
	logger *zap.Logger,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: assessmentRepo,
		findingRepo:    findingRepo,
		logger:         logger,
	}
}

func (s *AssessmentService) Create(ctx context.Context, req *model.CreateAssessmentRequest) (*model.Assessment, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}

	assessment := &model.Assessment{
		ID:             uuid.New(),
		Name:           req.Name,
		OrganizationID: orgID,
		Status:         "DRAFT",
		RiskScore:      0,
		TargetAssets:   req.TargetAssets,
		CreatedBy:      req.CreatedBy,
		CreatedAt:      time.Now().UTC(),
	}

	if assessment.TargetAssets == nil {
		assessment.TargetAssets = []string{}
	}

	if err := s.assessmentRepo.Create(ctx, assessment); err != nil {
		s.logger.Error("failed to create assessment", zap.Error(err))
		return nil, err
	}

	s.logger.Info("assessment created", zap.String("id", assessment.ID.String()))
	return assessment, nil
}

func (s *AssessmentService) Get(ctx context.Context, id uuid.UUID) (*model.Assessment, error) {
	return s.assessmentRepo.GetByID(ctx, id)
}

func (s *AssessmentService) GetWithSummary(ctx context.Context, id uuid.UUID) (*model.AssessmentResponse, error) {
	a, err := s.assessmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := a.ToResponse()

	summary, err := s.findingRepo.CountByAssessment(ctx, id)
	if err != nil {
		s.logger.Warn("failed to get finding summary", zap.Error(err))
	} else {
		summary.PqcReadiness = a.PqcReadiness
		summary.AssetsScanned = a.AssetsScanned
		resp.Summary = summary
	}

	return &resp, nil
}

func (s *AssessmentService) List(ctx context.Context, orgID *uuid.UUID, status string, offset, limit int) ([]model.Assessment, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.assessmentRepo.List(ctx, orgID, status, offset, limit)
}

func (s *AssessmentService) Run(ctx context.Context, id uuid.UUID) (*model.Assessment, error) {
	a, err := s.assessmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if a.Status != "DRAFT" && a.Status != "COMPLETED" {
		return nil, fmt.Errorf("assessment cannot be run in status %s", a.Status)
	}

	if err := s.assessmentRepo.UpdateStatus(ctx, id, "IN_PROGRESS", "system"); err != nil {
		return nil, err
	}

	// Simulate scanning target assets and producing findings
	findings := s.analyzeAssets(id, a.TargetAssets)

	if err := s.findingRepo.CreateBatch(ctx, findings); err != nil {
		s.logger.Error("failed to persist findings", zap.Error(err))
		return nil, err
	}

	overallRisk, riskScore, pqcReadiness := s.calculateRisk(findings)

	if err := s.assessmentRepo.UpdateResults(ctx, id, overallRisk, riskScore, pqcReadiness, len(a.TargetAssets)); err != nil {
		return nil, err
	}

	s.logger.Info("assessment completed",
		zap.String("id", id.String()),
		zap.String("risk", overallRisk),
		zap.Float64("score", riskScore),
		zap.Int("findings", len(findings)),
	)

	return s.assessmentRepo.GetByID(ctx, id)
}

func (s *AssessmentService) analyzeAssets(assessmentID uuid.UUID, assets []string) []model.Finding {
	var findings []model.Finding
	now := time.Now().UTC()

	for _, asset := range assets {
		// Default finding: flag every asset as potentially missing PQC
		rsa := "RSA-2048"
		rec := "ML-KEM-768"
		rem := "Migrate to post-quantum key encapsulation mechanism ML-KEM-768 (FIPS 203)"

		findings = append(findings, model.Finding{
			ID:                   uuid.New(),
			AssessmentID:         assessmentID,
			Category:             "MISSING_PQC",
			RiskLevel:            "HIGH",
			Title:                fmt.Sprintf("No PQC protection on %s", asset),
			Description:          fmt.Sprintf("Asset %s uses classical cryptography without post-quantum protection", asset),
			AffectedAsset:        asset,
			CurrentAlgorithm:     &rsa,
			RecommendedAlgorithm: &rec,
			Remediation:          &rem,
			DiscoveredAt:         now,
		})

		// Check for HNDL risk
		hndlRem := "Prioritise migration of long-lived secrets; data encrypted today can be captured and decrypted later by quantum computers"
		findings = append(findings, model.Finding{
			ID:                   uuid.New(),
			AssessmentID:         assessmentID,
			Category:             "HARVEST_NOW_DECRYPT_LATER",
			RiskLevel:            "CRITICAL",
			Title:                fmt.Sprintf("HNDL risk on %s", asset),
			Description:          fmt.Sprintf("Asset %s is vulnerable to harvest-now-decrypt-later attacks", asset),
			AffectedAsset:        asset,
			CurrentAlgorithm:     &rsa,
			RecommendedAlgorithm: &rec,
			Remediation:          &hndlRem,
			DiscoveredAt:         now,
		})
	}

	return findings
}

func (s *AssessmentService) calculateRisk(findings []model.Finding) (overallRisk string, riskScore float64, pqcReadiness float64) {
	if len(findings) == 0 {
		return "LOW", 0, 100.0
	}

	var critCount, highCount, medCount int
	for _, f := range findings {
		switch f.RiskLevel {
		case "CRITICAL":
			critCount++
		case "HIGH":
			highCount++
		case "MEDIUM":
			medCount++
		}
	}

	// Score: critical=10, high=5, medium=2, low=1
	riskScore = float64(critCount*10+highCount*5+medCount*2) / float64(len(findings)) * 10

	switch {
	case critCount > 0:
		overallRisk = "CRITICAL"
	case highCount > 0:
		overallRisk = "HIGH"
	case medCount > 0:
		overallRisk = "MEDIUM"
	default:
		overallRisk = "LOW"
	}

	// PQC readiness is inverse of missing-PQC findings
	missingPqc := 0
	for _, f := range findings {
		if f.Category == "MISSING_PQC" {
			missingPqc++
		}
	}
	totalAssets := len(findings) / 2 // each asset generates 2 findings
	if totalAssets > 0 {
		pqcReadiness = float64(totalAssets-missingPqc) / float64(totalAssets) * 100
	}

	return overallRisk, riskScore, pqcReadiness
}
