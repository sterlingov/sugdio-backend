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

func (r *PostgresRepository) CreateUserAuth(ctx context.Context, email, passwordHash, role string) (domain.UserShort, error) {
	query := "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, email, role"

	var us domain.UserShort

	row := r.db.QueryRowContext(ctx, query, email, passwordHash, role)
	err := row.Scan(&us.ID, &us.Email, &us.Role)
	return us, err
}
