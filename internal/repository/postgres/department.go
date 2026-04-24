package repository

import (
	"context"
	"database/sql"
	"errors"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) GetByIDDepartment(ctx context.Context, id int) (*domain.Department, error) {
	var p domain.Department

	query := "SELECT id, name FROM departments WHERE id=$1"

	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&p.ID, &p.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PostgresRepository) CreateDepartment(ctx context.Context, department *domain.DepartmentCreate) (int, error) {
	var insertedID int

	query := "INSERT INTO departments (name) VALUES ($1) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, department.Name).Scan(insertedID)

	return insertedID, err
}
