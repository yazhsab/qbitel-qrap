package model

import (
	"time"

	"github.com/google/uuid"
)

type Finding struct {
	ID                   uuid.UUID `json:"id"`
	AssessmentID         uuid.UUID `json:"assessment_id"`
	Category             string    `json:"category"`
	RiskLevel            string    `json:"risk_level"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	AffectedAsset        string    `json:"affected_asset"`
	CurrentAlgorithm     *string   `json:"current_algorithm"`
	RecommendedAlgorithm *string   `json:"recommended_algorithm"`
	Remediation          *string   `json:"remediation"`
	DiscoveredAt         time.Time `json:"discovered_at"`
	CreatedAt            time.Time `json:"created_at"`
}

type FindingResponse struct {
	ID                   uuid.UUID `json:"id"`
	AssessmentID         uuid.UUID `json:"assessment_id"`
	Category             string    `json:"category"`
	RiskLevel            string    `json:"risk_level"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	AffectedAsset        string    `json:"affected_asset"`
	CurrentAlgorithm     *string   `json:"current_algorithm,omitempty"`
	RecommendedAlgorithm *string   `json:"recommended_algorithm,omitempty"`
	Remediation          *string   `json:"remediation,omitempty"`
	DiscoveredAt         string    `json:"discovered_at"`
}

type FindingListResponse struct {
	Findings   []FindingResponse `json:"findings"`
	TotalCount int               `json:"total_count"`
	Offset     int               `json:"offset"`
	Limit      int               `json:"limit"`
}

func (f *Finding) ToResponse() FindingResponse {
	return FindingResponse{
		ID:                   f.ID,
		AssessmentID:         f.AssessmentID,
		Category:             f.Category,
		RiskLevel:            f.RiskLevel,
		Title:                f.Title,
		Description:          f.Description,
		AffectedAsset:        f.AffectedAsset,
		CurrentAlgorithm:     f.CurrentAlgorithm,
		RecommendedAlgorithm: f.RecommendedAlgorithm,
		Remediation:          f.Remediation,
		DiscoveredAt:         f.DiscoveredAt.Format(time.RFC3339),
	}
}
