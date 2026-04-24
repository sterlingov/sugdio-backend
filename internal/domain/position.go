package domain

type Position struct {
	ID   int64
	Name string
}

type PositionCreate struct {
	Name string
}
