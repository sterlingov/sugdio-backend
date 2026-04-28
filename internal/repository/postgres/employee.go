package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sugdio/internal/domain"
)

func (r *PostgresRepository) CreateEmployee(ctx context.Context, emp *domain.EmployeeCreate) (int64, error) {
	var insertedID int64

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO employees 
		(
		first_name, 
		middle_name, 
		second_name, 
		position_id,
		department_id,
		active,
		user_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`,
		emp.FirstName,
		emp.MiddleName,
		emp.SecondName,
		emp.PositionId,
		emp.DepartmentId,
		emp.Active,
		emp.UserId).Scan(&insertedID)

	return insertedID, err
}

func (r *PostgresRepository) DeleteEmployee(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM employees WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("execute delete: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByIDEmployee(ctx context.Context, id int64) (*domain.Employee, error) {
	query := `
	SELECT ee.id, ee.first_name, ee.middle_name, ee.second_name, ee.active, ee.created_at, p.id, p.name, d.id, d.name, u.id, u.email, u.role
	FROM employees ee
	LEFT JOIN positions p ON ee.position_id = p.id
	LEFT JOIN departments d ON ee.department_id = d.id
	LEFT JOIN users u ON ee.user_id = u.id
	WHERE ee.id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	e, err := ScanEmployee(row)
	if err != nil {
		return nil, err
	}
	return e, err
}

func ScanEmployee(r Scanner) (*domain.Employee, error) {
	var e domain.Employee

	var userID sql.NullInt64
	var userEmail sql.NullString
	var userRole sql.NullString

	var positionID sql.NullInt64
	var positionName sql.NullString

	var departmentID sql.NullInt64
	var departmentName sql.NullString

	err := r.Scan(
		&e.Id, &e.FirstName, &e.MiddleName, &e.SecondName, &e.Active, &e.CreatedAt,
		&positionID, &positionName,
		&departmentID, &departmentName,
		&userID, &userEmail, &userRole,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if userID.Valid {
		e.User = &domain.UserShort{
			ID:    userID.Int64,
			Email: userEmail.String,
			Role:  userRole.String,
		}
	}

	if departmentID.Valid {
		e.Department = &domain.Department{
			ID:   departmentID.Int64,
			Name: departmentName.String,
		}
	}

	if positionID.Valid {
		e.Position = &domain.Position{
			ID:   positionID.Int64,
			Name: positionName.String,
		}
	}

	return &e, nil
}

func (r *PostgresRepository) ListEmployee(ctx context.Context, f domain.EmployeeFilter) ([]domain.Employee, error) {
	query := `
	SELECT ee.id, ee.first_name, ee.middle_name, ee.second_name, ee.active, ee.created_at, p.id, p.name, d.id, d.name, u.id, u.email, u.role
	FROM employees ee
	LEFT JOIN positions p ON ee.position_id = p.id
	LEFT JOIN departments d ON ee.department_id = d.id
	LEFT JOIN users u ON ee.user_id = u.id
	WHERE 1=1`

	var args []any
	argID := 1

	if f.FirstName != nil {
		query += fmt.Sprintf(" AND first_name ILIKE $%d", argID)
		args = append(args, "%"+*f.FirstName+"%")
		argID++
	}

	if f.SecondName != nil {
		query += fmt.Sprintf(" AND second_name ILIKE $%d", argID)
		args = append(args, "%"+*f.SecondName+"%")
		argID++
	}

	if f.Active != nil {
		query += fmt.Sprintf(" AND active = $%d", argID)
		args = append(args, *f.Active)
		argID++
	}

	if f.DepartmentId != nil {
		query += fmt.Sprintf(" AND department_id = $%d", argID)
		args = append(args, *f.DepartmentId)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query employees: %w", err)
	}
	defer rows.Close()

	employees := make([]domain.Employee, 0, f.Limit)

	for rows.Next() {
		e, err := ScanEmployee(rows)
		if err != nil {
			return nil, err
		}
		employees = append(employees, *e)
	}
	return employees, nil
}

func (r *PostgresRepository) UpdateEmployee(ctx context.Context, id int64, p domain.EmployeePatch) error {
	pb := NewPatchBuilder()
	pb.Head("employees")

	if p.FirstName != nil {
		pb.Add("first_name", *p.FirstName)
	}
	if p.MiddleName != nil {
		pb.Add("middle_name", p.MiddleName)
	}
	if p.SecondName != nil {
		pb.Add("second_name", *p.SecondName)
	}
	if p.Active != nil {
		pb.Add("active", *p.Active)
	}
	if p.DepartmentId != nil {
		pb.Add("department_id", *p.DepartmentId)
	}
	if p.PositionId != nil {
		pb.Add("position_id", *p.PositionId)
	}
	if p.UserId != nil {
		pb.Add("user_id", p.UserId)
	}

	if pb.Len() == 0 {
		return nil
	}

	pb.Where("id", id)

	result, err := r.db.ExecContext(ctx, pb.String(), pb.Args()...)
	if err != nil {
		return fmt.Errorf("update employee: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
