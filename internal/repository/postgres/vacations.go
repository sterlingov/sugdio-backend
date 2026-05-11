package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sugdio/internal/domain"
)

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

func (r *PostgresRepository) CreateVacation(ctx context.Context, vacation *domain.VacationCreate) (int, error) {
	var insertedID int

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO vacations 
		(
		employee_id, 
		start_date, 
		end_date, 
		comment
		)
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
	).Scan(&insertedID)

	return insertedID, err
}

func ScanVacation(r Scanner) (*domain.Vacation, error) {
	var v domain.Vacation
	var e domain.EmployeeShort

	var eMiddleName sql.NullString
	var vComment sql.NullString

	err := r.Scan(&v.ID, &v.StartDate, &v.EndDate, &v.Status, &vComment, &v.CreatedAt, &v.UpdatedAt, &e.ID, &e.FirstName, &eMiddleName, &e.SecondName, &e.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if eMiddleName.Valid {
		e.MiddleName = &eMiddleName.String
	}
	v.Employee = e
	if vComment.Valid {
		v.Comment = &vComment.String
	}
	return &v, nil
}

func (r *PostgresRepository) GetByIDVacation(ctx context.Context, vacationID int) (*domain.Vacation, error) {
	query := `
	SELECT v.id, v.start_date, v.end_date, v.status, v.comment, v.created_at, v.updated_at, e.id, e.first_name, e.middle_name, e.second_name, e.active
	FROM vacations v
	LEFT JOIN employees e on v.employee_id = e.id
	WHERE v.id = $1`

	row := r.db.QueryRowContext(ctx, query, vacationID)

	return ScanVacation(row)
}

func (r *PostgresRepository) ListVacation(ctx context.Context, filter domain.VacationFilter) ([]domain.Vacation, error) {
	query := `
	SELECT v.id, v.start_date, v.end_date, v.status, v.comment, v.created_at, v.updated_at, e.id, e.first_name, e.middle_name, e.second_name, e.active
	FROM vacations v
	LEFT JOIN employees e on v.employee_id = e.id
	WHERE 1 = 1`

	var args []any
	argID := 1

	if filter.EmployeeID != nil {
		query += fmt.Sprintf(" AND v.employee_id = $%d", argID)
		args = append(args, *filter.EmployeeID)
		argID++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND v.status = $%d", argID)
		args = append(args, *filter.Status)
		argID++
	}

	query += fmt.Sprintf(" AND daterange(v.start_date, v.end_date, '[]') && daterange($%d, $%d, '[]') LIMIT $%d OFFSET $%d", argID, argID+1, argID+2, argID+3)
	args = append(args, filter.FromDate, filter.ToDate, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query vacations: %w", err)
	}
	defer rows.Close()

	vacations := make([]domain.Vacation, 0, filter.Limit)
	for rows.Next() {
		v, err := ScanVacation(rows)
		if err != nil {
			return nil, err
		}
		vacations = append(vacations, *v)
	}
	return vacations, nil
}

func (r *PostgresRepository) UpdateVacation(ctx context.Context, vacationID int, patch domain.VacationPatch) error {
	pb := NewPatchBuilder()
	pb.Head("vacations")

	if patch.StartDate != nil {
		pb.Add("start_date", *patch.StartDate)
	}
	if patch.EndDate != nil {
		pb.Add("end_date", *patch.EndDate)
	}
	if patch.Comment != nil {
		pb.Add("comment", *patch.Comment)
	}
	if patch.Status != nil {
		pb.Add("status", *patch.Status)
	}

	if pb.Len() == 0 {
		return nil
	}

	pb.Where("id", vacationID)

	result, err := r.db.ExecContext(ctx, pb.String(), pb.Args()...)
	if err != nil {
		return fmt.Errorf("update vacation: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteVacation(ctx context.Context, vacationID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM vacations WHERE id = $1", vacationID)

	if err != nil {
		return fmt.Errorf("execute vacation delete: %w", err)
	}
	return nil
}
