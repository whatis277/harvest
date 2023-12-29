package entity

import (
	"time"
)

type Subscription struct {
	ID     string
	UserID int

	Amount        int
	Frequency     int
	PaymentMethod string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID int

	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Token struct {
	ID     string
	UserID int

	AccessToken string

	CreatedAt time.Time
	ExpiresAt time.Time
}