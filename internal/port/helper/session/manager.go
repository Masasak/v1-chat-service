package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrNotFound indicates that requested session does not exist.
	ErrNotFound = errors.New("session not found")
)

// Manager manages the real-time chat session.
// It stores user's individual session information.
// And it is also responsible for sending payloads to chats.
//
// TODO: It might be more effective to seperate these responsibility.
// You might want to change this into 2 manager, which are responsible for user's session and message exchangement.
type Manager interface {
	// Fetch fetches user's session info.
	// If not found, it will return session.ErrNotFound.
	Fetch(ctx context.Context, userID uuid.UUID) (Session, error)

	// Set sets user's session with given user id.
	Set(ctx context.Context, userID uuid.UUID, session Session) error

	// Send sends given payload to specific chat.
	// TODO: This api is unstable. It might change a lot.
	Send(ctx context.Context, chatID uuid.UUID, payload Payload) error
}
