package domain

import "time"

type VacationStatus string

const (
	VacationStatusPending  VacationStatus = "pending"
	VacationStatusApproved VacationStatus = "approved"
	VacationStatusRejected VacationStatus = "rejected"
)

type Vacation struct {
	ID        int
	Employee  EmployeeShort
	StartDate time.Time
	EndDate   time.Time
	Status    VacationStatus
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VacationCreate struct {
	EmployeeID int
	StartDate  time.Time
	EndDate    time.Time
	Comment    *string
}

type VacationPatch struct {
	StartDate *time.Time
	EndDate   *time.Time
	Status    *VacationStatus
	Comment   *string
}

type VacationFilter struct {
	EmployeeID *int
	FromDate   *time.Time
	ToDate     *time.Time
	Status     *VacationStatus

	Limit  int
	Offset int
}
