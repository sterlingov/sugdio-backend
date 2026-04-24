package repository

import (
	"context"
	"database/sql"
	"errors"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) GetByIDUser(ctx context.Context, id int64) (*domain.UserShort, error) {
	query := "SELECT id, email, role FROM users WHERE id=$1"

	var u domain.UserShort

	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&u.ID, &u.Email, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}
