package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quantun-opensource/qrap/api/internal/model"
)

type FindingRepository struct {
	pool *pgxpool.Pool
}

func NewFindingRepository(pool *pgxpool.Pool) *FindingRepository {
	return &FindingRepository{pool: pool}
}

func (r *FindingRepository) Create(ctx context.Context, f *model.Finding) error {
	query := `
		INSERT INTO findings (id, assessment_id, category, risk_level, title, description,
		                      affected_asset, current_algorithm, recommended_algorithm, remediation, discovered_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.pool.Exec(ctx, query,
		f.ID, f.AssessmentID, f.Category, f.RiskLevel, f.Title, f.Description,
		f.AffectedAsset, f.CurrentAlgorithm, f.RecommendedAlgorithm, f.Remediation, f.DiscoveredAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert finding: %w", err)
	}
	return nil
}

func (r *FindingRepository) CreateBatch(ctx context.Context, findings []model.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i := range findings {
		f := &findings[i]
		query := `
			INSERT INTO findings (id, assessment_id, category, risk_level, title, description,
			                      affected_asset, current_algorithm, recommended_algorithm, remediation, discovered_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		_, err := tx.Exec(ctx, query,
			f.ID, f.AssessmentID, f.Category, f.RiskLevel, f.Title, f.Description,
			f.AffectedAsset, f.CurrentAlgorithm, f.RecommendedAlgorithm, f.Remediation, f.DiscoveredAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert finding %s: %w", f.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *FindingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Finding, error) {
	query := `
		SELECT id, assessment_id, category, risk_level, title, description,
		       affected_asset, current_algorithm, recommended_algorithm, remediation, discovered_at, created_at
		FROM findings WHERE id = $1
	`
	f := &model.Finding{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.AssessmentID, &f.Category, &f.RiskLevel, &f.Title, &f.Description,
		&f.AffectedAsset, &f.CurrentAlgorithm, &f.RecommendedAlgorithm, &f.Remediation, &f.DiscoveredAt, &f.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("finding not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get finding: %w", err)
	}
	return f, nil
}

func (r *FindingRepository) ListByAssessment(ctx context.Context, assessmentID uuid.UUID, riskLevel, category string, offset, limit int) ([]model.Finding, int, error) {
	countQuery := `SELECT COUNT(*) FROM findings WHERE assessment_id = $1`
	listQuery := `
		SELECT id, assessment_id, category, risk_level, title, description,
		       affected_asset, current_algorithm, recommended_algorithm, remediation, discovered_at, created_at
		FROM findings WHERE assessment_id = $1
	`
	args := []interface{}{assessmentID}
	argIdx := 2

	if riskLevel != "" {
		filter := fmt.Sprintf(" AND risk_level = $%d", argIdx)
		countQuery += filter
		listQuery += filter
		args = append(args, riskLevel)
		argIdx++
	}
	if category != "" {
		filter := fmt.Sprintf(" AND category = $%d", argIdx)
		countQuery += filter
		listQuery += filter
		args = append(args, category)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count findings: %w", err)
	}

	listQuery += fmt.Sprintf(" ORDER BY risk_level ASC, discovered_at DESC OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, offset, limit)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list findings: %w", err)
	}
	defer rows.Close()

	var findings []model.Finding
	for rows.Next() {
		var f model.Finding
		if err := rows.Scan(
			&f.ID, &f.AssessmentID, &f.Category, &f.RiskLevel, &f.Title, &f.Description,
			&f.AffectedAsset, &f.CurrentAlgorithm, &f.RecommendedAlgorithm, &f.Remediation, &f.DiscoveredAt, &f.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan finding: %w", err)
		}
		findings = append(findings, f)
	}
	return findings, total, rows.Err()
}

func (r *FindingRepository) CountByAssessment(ctx context.Context, assessmentID uuid.UUID) (*model.AssessmentSummary, error) {
	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE risk_level = 'CRITICAL') AS critical,
			COUNT(*) FILTER (WHERE risk_level = 'HIGH') AS high,
			COUNT(*) FILTER (WHERE risk_level = 'MEDIUM') AS medium,
			COUNT(*) FILTER (WHERE risk_level = 'LOW') AS low
		FROM findings WHERE assessment_id = $1
	`
	s := &model.AssessmentSummary{}
	err := r.pool.QueryRow(ctx, query, assessmentID).Scan(
		&s.TotalFindings, &s.CriticalFindings, &s.HighFindings, &s.MediumFindings, &s.LowFindings,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count findings: %w", err)
	}
	return s, nil
}
