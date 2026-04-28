package handlers

import (
	"sugdio/api"
	"sugdio/internal/domain"

	"github.com/oapi-codegen/runtime/types"
)

func toAPIEmployee(employee *domain.Employee) api.Employee {
	res := api.Employee{
		FirstName:  employee.FirstName,
		MiddleName: employee.MiddleName,
		SecondName: employee.SecondName,
		Active:     &employee.Active,
		CreatedAt:  &employee.CreatedAt,
		Id:         employee.Id,
	}

	if employee.User != nil {
		res.User = &api.UserShort{
			Id:    int(employee.User.ID),
			Email: employee.User.Email,
			Role:  employee.User.Role,
		}
	}
	if employee.Department != nil {
		res.Department = &api.DepartmentShort{
			Id:   int(employee.Department.ID),
			Name: employee.Department.Name,
		}
	}
	if employee.Position != nil {
		res.Position = &api.PositionShort{
			Id:   int(employee.Position.ID),
			Name: employee.Position.Name,
		}
	}

	return res
}

func toAPIShift(s *domain.Shift) api.Shift {
	res := api.Shift{
		Id:        s.ID,
		CreatedAt: &s.CreatedAt,
		UpdatedAt: &s.UpdatedAt,
		Date:      types.Date{Time: s.Date},
		ShiftType: &api.ShiftType{Id: &s.ShiftType.ID, Name: &s.ShiftType.Name},
		Status:    &s.Status,
		Employee:  &api.EmployeeShort{Active: &s.Employee.Active, FirstName: s.Employee.FirstName, SecondName: s.Employee.SecondName, Id: s.Employee.ID, MiddleName: s.Employee.MiddleName},
	}

	return res
}
