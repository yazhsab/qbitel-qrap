package model

import (
	"time"

	"github.com/google/uuid"
)

type Assessment struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Status         string     `json:"status"`
	OverallRisk    *string    `json:"overall_risk"`
	RiskScore      float64    `json:"risk_score"`
	TargetAssets   []string   `json:"target_assets"`
	AssetsScanned  int        `json:"assets_scanned"`
	PqcReadiness   float64    `json:"pqc_readiness"`
	StartedAt      *time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedBy      string     `json:"updated_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateAssessmentRequest struct {
	Name           string   `json:"name"`
	OrganizationID string   `json:"organization_id"`
	TargetAssets   []string `json:"target_assets"`
	CreatedBy      string   `json:"created_by"`
}

type AssessmentResponse struct {
	ID             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	OrganizationID uuid.UUID          `json:"organization_id"`
	Status         string             `json:"status"`
	OverallRisk    *string            `json:"overall_risk,omitempty"`
	RiskScore      float64            `json:"risk_score"`
	TargetAssets   []string           `json:"target_assets"`
	Summary        *AssessmentSummary `json:"summary,omitempty"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
}

type AssessmentSummary struct {
	TotalFindings    int     `json:"total_findings"`
	CriticalFindings int     `json:"critical_findings"`
	HighFindings     int     `json:"high_findings"`
	MediumFindings   int     `json:"medium_findings"`
	LowFindings      int     `json:"low_findings"`
	PqcReadiness     float64 `json:"pqc_readiness_percentage"`
	AssetsScanned    int     `json:"assets_scanned"`
}

type AssessmentListResponse struct {
	Assessments []AssessmentResponse `json:"assessments"`
	TotalCount  int                  `json:"total_count"`
	Offset      int                  `json:"offset"`
	Limit       int                  `json:"limit"`
}

func (a *Assessment) ToResponse() AssessmentResponse {
	return AssessmentResponse{
		ID:             a.ID,
		Name:           a.Name,
		OrganizationID: a.OrganizationID,
		Status:         a.Status,
		OverallRisk:    a.OverallRisk,
		RiskScore:      a.RiskScore,
		TargetAssets:   a.TargetAssets,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}
}
