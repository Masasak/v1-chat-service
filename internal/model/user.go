package model

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID
	Nickname        string
	ProfileImageURL string
}
