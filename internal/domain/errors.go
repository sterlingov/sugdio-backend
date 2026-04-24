package domain

import (
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrWrongPosition   = errors.New("position id is wrong")
	ErrWrongDepartment = errors.New("department id is wrong")
	ErrWrongUser       = errors.New("user id is wrong")

	ErrEmployeeHasShifts    = errors.New("cannot delete employee: it has shifts")
	ErrEmployeeHasVacations = errors.New("cannot delete employee: it has vacations")
)
