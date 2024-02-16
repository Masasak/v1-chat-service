package session

import (
	"time"

	"github.com/google/uuid"
)

// Session is user's session.
// It includes the user's current chat, and the time the user has connected.
type Session struct {
	ChatID uuid.UUID
	Since  time.Time
}
