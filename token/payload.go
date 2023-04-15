package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token is expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"userId"`
	Role 	string    `json:"role"` 
	VendorID int64     `json:"vendorId"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewPayload(userID int64, duration time.Duration) *Payload {
	now := time.Now()

	payload := Payload{
		ID:        uuid.New(),
		UserID:    userID,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}

	return &payload
}

func (payload *Payload) Valid() error {
	if payload.ExpiredAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}
