package domain

import "time"

type Employee struct {
	Id         int64
	FirstName  string
	MiddleName *string
	SecondName string
	Active     bool
	CreatedAt  time.Time
	Department *Department
	Position   *Position
	User       *UserShort
}

type EmployeeCreate struct {
	FirstName  string
	MiddleName *string
	SecondName string
	Active     *bool

	DepartmentId *int
	PositionId   *int
	UserId       *int64
}

type EmployeePatch struct {
	Active       *bool   `json:"active,omitempty"`
	DepartmentId *int    `json:"department_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	MiddleName   *string `json:"middle_name,omitempty"`
	PositionId   *int    `json:"position_id,omitempty"`
	SecondName   *string `json:"second_name,omitempty"`
	UserId       *int64  `json:"user_id,omitempty"`
}

type EmployeeFilter struct {
	FirstName    *string
	SecondName   *string
	Active       *bool
	DepartmentId *int
	Limit        int
	Offset       int
}
