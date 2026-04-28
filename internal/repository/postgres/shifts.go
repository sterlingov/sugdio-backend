package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func ScanShift(r Scanner) (*domain.Shift, error) {
	var s domain.Shift
	var e domain.EmployeeShort
	var st domain.ShiftType

	var eMiddleName sql.NullString

	err := r.Scan(&s.ID, &s.Date, &s.Status, &s.CreatedAt, &s.UpdatedAt, &st.ID, &st.Name, &e.ID, &e.FirstName, &eMiddleName, &e.SecondName, &e.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if eMiddleName.Valid {
		e.MiddleName = &eMiddleName.String
	}
	s.Employee = e
	s.ShiftType = st

	return &s, nil
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
	query := `
	SELECT s.id, s.date, s.status, s.created_at, s.updated_at, st.id, st.name, e.id, e.first_name, e.middle_name, e.second_name, e.active
	FROM shifts s
	LEFT JOIN shift_types st ON s.shift_type_id = st.id
	LEFT JOIN employees e ON s.employee_id = e.id
	WHERE s.id = $1
	`

	row := r.db.QueryRowContext(ctx, query, shiftID)

	s, err := ScanShift(row)
	if err != nil {
		return nil, err
	}
	return s, err
}
func (r *PostgresRepository) ListShift(ctx context.Context, filter domain.ShiftFilter) ([]domain.Shift, error) {
	query := `
	SELECT s.id, s.date, s.status, s.created_at, s.updated_at, st.id, st.name, e.id, e.first_name, e.middle_name, e.second_name, e.active
	FROM shifts s
	LEFT JOIN shift_types st ON s.shift_type_id = st.id
	LEFT JOIN employees e ON s.employee_id = e.id
	WHERE 1 = 1
	`

	var args []any
	argID := 1

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND s.date >= $%d", argID)
		args = append(args, *filter.DateFrom)
		argID++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND s.date <= $%d", argID)
		args = append(args, *filter.DateTo)
		argID++
	}

	if filter.EmployeeID != nil {
		query += fmt.Sprintf(" AND e.id = $%d", argID)
		args = append(args, *filter.EmployeeID)
		argID++
	}

	if filter.ShiftTypeID != nil {
		query += fmt.Sprintf(" AND st.id = $%d", argID)
		args = append(args, *filter.ShiftTypeID)
		argID++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND s.status = $%d", argID)
		args = append(args, *filter.Status)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query shifts: %w", err)
	}
	defer rows.Close()

	shifts := make([]domain.Shift, 0, filter.Limit)

	for rows.Next() {
		s, err := ScanShift(rows)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, *s)
	}

	return shifts, nil
}

func (r *PostgresRepository) UpdateShift(ctx context.Context, id int, patch domain.ShiftPatch) error {
	pb := NewPatchBuilder()
	pb.Head("shifts")

	if patch.Date != nil {
		pb.Add("date", *patch.Date)
	}
	if patch.EmployeeID != nil {
		pb.Add("employee_id", *patch.EmployeeID)
	}
	if patch.ShiftTypeID != nil {
		pb.Add("shift_type_id", *patch.ShiftTypeID)
	}
	if patch.Status != nil {
		pb.Add("status", *patch.Status)
	}

	if pb.Len() == 0 {
		return nil
	}

	pb.Where("id", id)

	result, err := r.db.ExecContext(ctx, pb.String(), pb.Args()...)
	if err != nil {
		return fmt.Errorf("update shift: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteShift(ctx context.Context, shiftID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM shifts WHERE id = $1", shiftID)

	if err != nil {
		return fmt.Errorf("execute shift delete: %w", err)
	}
	return nil
}
