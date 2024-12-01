package uuid

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

func CreateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}