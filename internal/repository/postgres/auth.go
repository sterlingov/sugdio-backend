package repository

import (
	"context"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) GetByEmailAuth(ctx context.Context, email string) (domain.UserCredentials, error) {
	query := "SELECT id, email, password_hash, role FROM users WHERE email=$1"

	var uc domain.UserCredentials

	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&uc.ID, &uc.Email, &uc.PasswordHash, &uc.Role)
	return uc, err
}
