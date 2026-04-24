package repository

import "context"

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
