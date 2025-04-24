package model

import (
	"github.com/google/uuid"
)

// ID represents a unique identifier
type ID uuid.UUID

// NewID creates a new ID
func NewID() ID {
	return ID(uuid.New())
}

// FromString creates an ID from a string
func FromString(id string) (ID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}
	return ID(parsedUUID), nil
}

// String returns the string representation of the ID
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalBSON implements the bson.Marshaler interface
func (id ID) MarshalBSON() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalBSON implements the bson.Unmarshaler interface
func (id *ID) UnmarshalBSON(data []byte) error {
	parsedUUID, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = ID(parsedUUID)
	return nil
}
