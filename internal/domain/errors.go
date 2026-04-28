package domain

import (
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrWrongPosition   = errors.New("position id is wrong")
	ErrWrongDepartment = errors.New("department id is wrong")
	ErrWrongEmployee   = errors.New("employee id is wrong")
	ErrWrongUser       = errors.New("user id is wrong")
	ErrWrongShiftType  = errors.New("shift_type is wrong")

	ErrWrongShiftStatus = errors.New("wrong shift status")

	ErrEmployeeHasShifts    = errors.New("cannot delete employee: it has shifts")
	ErrEmployeeHasVacations = errors.New("cannot delete employee: it has vacations")

	ErrShiftTypeInUse = errors.New("cannot delete shift type: it is using by some shifts")
)
