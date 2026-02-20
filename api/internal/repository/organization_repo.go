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

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	query := `
		INSERT INTO organizations (id, name, description, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, org.ID, org.Name, org.Description, org.CreatedBy, org.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	query := `
		SELECT id, name, description, created_by, created_at, updated_by, updated_at
		FROM organizations WHERE id = $1
	`
	var org model.Organization
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Description,
		&org.CreatedBy, &org.CreatedAt, &org.UpdatedBy, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *OrganizationRepository) List(ctx context.Context, offset, limit int) ([]model.Organization, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM organizations`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	query := `
		SELECT id, name, description, created_by, created_at, updated_by, updated_at
		FROM organizations ORDER BY created_at DESC OFFSET $1 LIMIT $2
	`
	rows, err := r.pool.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []model.Organization
	for rows.Next() {
		var org model.Organization
		if err := rows.Scan(
			&org.ID, &org.Name, &org.Description,
			&org.CreatedBy, &org.CreatedAt, &org.UpdatedBy, &org.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, total, rows.Err()
}

func (r *OrganizationRepository) Update(ctx context.Context, id uuid.UUID, name, description, updatedBy string) error {
	query := `
		UPDATE organizations SET name = $1, description = $2, updated_by = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.pool.Exec(ctx, query, name, description, updatedBy, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("organization not found: %s", id)
	}
	return nil
}
