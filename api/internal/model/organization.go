package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
}

type OrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	TotalCount    int                    `json:"total_count"`
	Offset        int                    `json:"offset"`
	Limit         int                    `json:"limit"`
}

func (o *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   o.UpdatedAt.Format(time.RFC3339),
	}
}
