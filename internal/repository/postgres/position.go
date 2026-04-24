package repository

import (
	"context"
	"database/sql"
	"errors"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) GetByIDPosition(ctx context.Context, id int) (*domain.Position, error) {
	var p domain.Position

	query := "SELECT id, name FROM positions WHERE id=$1"

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

func (r *PostgresRepository) CreatePosition(ctx context.Context, position *domain.PositionCreate) (int, error) {
	var insertedID int

	query := "INSERT INTO positions (name) VALUES ($1) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, position.Name).Scan(insertedID)

	return insertedID, err
}
