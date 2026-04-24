package repository

import "context"

func (r *PostgresRepository) HasVacations(ctx context.Context, employeeID int64) (bool, error) {
	var vacationsCount int

	query := "SELECT COUNT(*) FROM vacations WHERE employee_id = $1"
	row := r.db.QueryRowContext(ctx, query, employeeID)

	err := row.Scan(&vacationsCount)
	if err != nil {
		return true, err
	}

	return vacationsCount > 0, nil
}
