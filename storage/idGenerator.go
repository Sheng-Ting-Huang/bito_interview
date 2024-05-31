package storage

import "github.com/google/uuid"

type IDGenerator interface {
	GenerateKey() string
}

type UUIDGenerator struct{}

func (UUIDGenerator) GenerateKey() string {
	return uuid.New().String()
}

type FakeIDGenerator struct {
	FakeID string
}

func (f FakeIDGenerator) GenerateKey() string {
	return f.FakeID
}
