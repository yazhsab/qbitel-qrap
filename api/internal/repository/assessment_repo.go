package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quantun-opensource/qrap/api/internal/model"
)

type AssessmentRepository struct {
	pool *pgxpool.Pool
}

func NewAssessmentRepository(pool *pgxpool.Pool) *AssessmentRepository {
	return &AssessmentRepository{pool: pool}
}

func (r *AssessmentRepository) Create(ctx context.Context, a *model.Assessment) error {
	query := `
		INSERT INTO assessments (id, name, organization_id, status, risk_score, target_assets, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		a.ID, a.Name, a.OrganizationID, a.Status, a.RiskScore, a.TargetAssets, a.CreatedBy, a.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert assessment: %w", err)
	}
	return nil
}

func (r *AssessmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Assessment, error) {
	query := `
		SELECT id, name, organization_id, status, overall_risk, risk_score,
		       target_assets, assets_scanned, pqc_readiness, started_at, completed_at,
		       created_by, created_at, updated_by, updated_at
		FROM assessments WHERE id = $1
	`
	a := &model.Assessment{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.Name, &a.OrganizationID, &a.Status, &a.OverallRisk, &a.RiskScore,
		&a.TargetAssets, &a.AssetsScanned, &a.PqcReadiness, &a.StartedAt, &a.CompletedAt,
		&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("assessment not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get assessment: %w", err)
	}
	return a, nil
}

func (r *AssessmentRepository) List(ctx context.Context, orgID *uuid.UUID, status string, offset, limit int) ([]model.Assessment, int, error) {
	countQuery := `SELECT COUNT(*) FROM assessments WHERE 1=1`
	listQuery := `
		SELECT id, name, organization_id, status, overall_risk, risk_score,
		       target_assets, assets_scanned, pqc_readiness, started_at, completed_at,
		       created_by, created_at, updated_by, updated_at
		FROM assessments WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if orgID != nil {
		filter := fmt.Sprintf(" AND organization_id = $%d", argIdx)
		countQuery += filter
		listQuery += filter
		args = append(args, *orgID)
		argIdx++
	}
	if status != "" {
		filter := fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += filter
		listQuery += filter
		args = append(args, status)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count assessments: %w", err)
	}

	listQuery += fmt.Sprintf(" ORDER BY created_at DESC OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, offset, limit)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list assessments: %w", err)
	}
	defer rows.Close()

	var assessments []model.Assessment
	for rows.Next() {
		var a model.Assessment
		if err := rows.Scan(
			&a.ID, &a.Name, &a.OrganizationID, &a.Status, &a.OverallRisk, &a.RiskScore,
			&a.TargetAssets, &a.AssetsScanned, &a.PqcReadiness, &a.StartedAt, &a.CompletedAt,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan assessment: %w", err)
		}
		assessments = append(assessments, a)
	}
	return assessments, total, rows.Err()
}

func (r *AssessmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, updatedBy string) error {
	query := `UPDATE assessments SET status = $1, updated_by = $2, updated_at = $3 WHERE id = $4`
	result, err := r.pool.Exec(ctx, query, status, updatedBy, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update assessment status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("assessment not found: %s", id)
	}
	return nil
}

func (r *AssessmentRepository) UpdateResults(ctx context.Context, id uuid.UUID, overallRisk string, riskScore, pqcReadiness float64, assetsScanned int) error {
	now := time.Now().UTC()
	query := `
		UPDATE assessments
		SET overall_risk = $1, risk_score = $2, pqc_readiness = $3, assets_scanned = $4,
		    status = 'COMPLETED', completed_at = $5, updated_at = $5
		WHERE id = $6
	`
	result, err := r.pool.Exec(ctx, query, overallRisk, riskScore, pqcReadiness, assetsScanned, now, id)
	if err != nil {
		return fmt.Errorf("failed to update assessment results: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("assessment not found: %s", id)
	}
	return nil
}
