package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	// Generate UUIDv7 if not provided
	if role.ID == uuid.Nil {
		role.ID = uuid.Must(uuid.NewV7())
	}

	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	return err
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL
	`

	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL
	`

	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles
		WHERE deleted_at IS NULL
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role := &entity.Role{}
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, role.Name, role.Description, role.UpdatedAt, role.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE roles
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *roleRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *roleRepository) AssignToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

func (r *roleRepository) RemoveFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

func (r *roleRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at, r.deleted_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.deleted_at IS NULL
		ORDER BY r.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role := &entity.Role{}
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *roleRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *roleRepository) Restore(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE roles
		SET deleted_at = NULL, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NOT NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
