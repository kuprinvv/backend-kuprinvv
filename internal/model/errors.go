package model

import "errors"

var ErrForbidden = errors.New("forbidden")

var (
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrPastSlot          = errors.New("cannot book slot")
)

var ErrRoomNotFound = errors.New("room not found")

var ErrBookingNotFound = errors.New("booking not found")

var ErrScheduleAlreadyExists = errors.New("schedule already exists")

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
