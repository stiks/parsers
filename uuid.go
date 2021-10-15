package parsers

import (
	"log"
	"testing"

	"github.com/google/uuid"

)

// UUIDFromString ...
func UUIDFromString(t *testing.T, uStr string) uuid.UUID {
	id, err := uuid.Parse(uStr)
	if err != nil {
		if t == nil {
			log.Fatalf("Unable to get UUID from string: %s", err.Error())
		} else {
			t.Fatalf("Unable to get UUID from string: %s", err.Error())
		}
	}

	return id
}

// UUIDFromStringPointer ...
func UUIDFromStringPointer(t *testing.T, uStr string) *uuid.UUID {
	id := UUIDFromString(t, uStr)

	return &id
}
