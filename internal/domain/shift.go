package domain

type ShiftStatus string

const (
	ShiftStatusPlanned   ShiftStatus = "planned"
	ShiftStatusCompleted ShiftStatus = "completed"
	ShiftStatusCancelled ShiftStatus = "cancelled"
)

type ShiftCreate struct {
	Date        string
	EmployeeID  int
	ShiftTypeID int
	Status      *ShiftStatus
}

type ShiftType struct {
	ID   int
	Name string
}

type Shift struct {
	ID        int
	CreatedAt string
	UpdatedAt string
	ShiftType ShiftType
	Employee  EmployeeShort
	Date      string
	Status    string
}

type ShiftPatch struct {
	ShiftTypeID *int
	Date        *string
	Status      *ShiftStatus
	EmployeeID  *int
}

type ShiftFilter struct {
	ShiftTypeID *int
	DateFrom    *string
	DateTo      *string
	Status      *ShiftStatus
	EmployeeID  *int
	Limit       int
	Offset      int
}
