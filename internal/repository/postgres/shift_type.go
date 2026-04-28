package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sugdio/internal/domain"
)

func ScanShiftType(s Scanner) (*domain.ShiftType, error) {
	var st domain.ShiftType

	err := s.Scan(&st.ID, &st.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &st, nil
}

func (r *PostgresRepository) CreateShiftType(ctx context.Context, shiftType *domain.ShiftTypeCreate) (int, error) {
	var insertedID int
	query := `INSERT INTO shift_types (name) VALUES ($1)`
	err := r.db.QueryRowContext(ctx, query, shiftType.Name).Scan(&insertedID)
	return insertedID, err
}

func (r *PostgresRepository) GetByIDShiftType(ctx context.Context, id int) (*domain.ShiftType, error) {
	query := `SELECT id, name FROM shift_types WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	return ScanShiftType(row)
}
func (r *PostgresRepository) ListShiftType(ctx context.Context) ([]domain.ShiftType, error) {
	query := `SELECT id, name FROM shift_types`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query shift_types: %w", err)
	}
	defer rows.Close()
	arr := make([]domain.ShiftType, 0)
	for rows.Next() {
		st, err := ScanShiftType(rows)
		if err != nil {
			return nil, err
		}
		arr = append(arr, *st)
	}
	return arr, nil
}
func (r *PostgresRepository) UpdateShiftType(ctx context.Context, id int, shift domain.ShiftTypePatch) error {
	pb := NewPatchBuilder()
	pb.Head("shift_types")
	if shift.Name != nil {
		pb.Add("name", shift.Name)
	}
	pb.Where("id", id)

	result, err := r.db.ExecContext(ctx, pb.String(), pb.Args()...)
	if err != nil {
		return fmt.Errorf("update shift_type: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
func (r *PostgresRepository) DeleteShiftType(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM shift_types WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("execute shift_type delete: %w", err)
	}
	return nil
}
