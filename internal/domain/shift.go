package domain

import "time"

type ShiftStatus string

const (
	ShiftStatusPlanned   ShiftStatus = "planned"
	ShiftStatusCompleted ShiftStatus = "completed"
	ShiftStatusCancelled ShiftStatus = "cancelled"
)

type ShiftCreate struct {
	Date        time.Time
	EmployeeID  int
	ShiftTypeID int
	Status      ShiftStatus
}

type ShiftType struct {
	ID   int
	Name string
}

type ShiftTypePatch struct {
	Name *string
}

type ShiftTypeCreate struct {
	Name string
}

type Shift struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	ShiftType ShiftType
	Employee  EmployeeShort
	Date      time.Time
	Status    string
}

type ShiftPatch struct {
	ShiftTypeID *int
	Date        *time.Time
	Status      *ShiftStatus
	EmployeeID  *int
}

type ShiftFilter struct {
	ShiftTypeID *int
	DateFrom    *time.Time
	DateTo      *time.Time
	Status      *ShiftStatus
	EmployeeID  *int
	Limit       int
	Offset      int
}
