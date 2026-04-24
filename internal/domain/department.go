package domain

type Department struct {
	ID   int64
	Name string
}

type DepartmentCreate struct {
	Name string
}
