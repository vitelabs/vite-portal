package idutil

import (
	"github.com/google/uuid"
)

func NewGuid() string {
	return uuid.New().String()
}