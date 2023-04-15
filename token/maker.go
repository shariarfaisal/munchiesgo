package token

import "time"

type Maker interface {
	// CreateToken creates a new token for the given user ID.
	CreateToken(userID int64, duration time.Duration) (string, error)

	// VerifyToken verifies the given token and returns the user ID.
	VerifyToken(token string) (*Payload, error)
}
