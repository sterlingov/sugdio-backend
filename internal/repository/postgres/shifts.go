package repository

import (
	"context"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) HasShifts(ctx context.Context, employeeID int64) (bool, error) {
	var shiftsCount int

	query := "SELECT COUNT(*) FROM shifts WHERE employee_id = $1"
	row := r.db.QueryRowContext(ctx, query, employeeID)

	err := row.Scan(&shiftsCount)
	if err != nil {
		return true, err
	}

	return shiftsCount > 0, nil
}

func (r *PostgresRepository) CreateShift(ctx context.Context, shift *domain.ShiftCreate) (int, error) {
	var insertedID int

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO shifts 
		(
		date, 
		shift_type_id, 
		employee_id, 
		status
		)
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		shift.Date,
		shift.ShiftTypeID,
		shift.EmployeeID,
		shift.Status).Scan(&insertedID)

	return insertedID, err
}
func (r *PostgresRepository) GetByIDShift(ctx context.Context, shiftID int) (*domain.Shift, error) {
	panic("unimplemented")
}
func (r *PostgresRepository) ListShift(ctx context.Context, filter domain.ShiftFilter) ([]domain.Shift, error) {
	panic("unimplemented")
}
func (r *PostgresRepository) UpdateShift(ctx context.Context, patch domain.ShiftPatch) error {
	panic("unimplemented")
}
func (r *PostgresRepository) DeleteShift(ctx context.Context, shiftID int) error {
	panic("unimplemented")
}
