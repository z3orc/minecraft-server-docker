package mojang

import "github.com/google/uuid"

type Profile struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
