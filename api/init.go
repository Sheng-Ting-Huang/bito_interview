package api

import "github.com/bito_interview/storage"

var IdGenerator storage.IDGenerator

func init() {
	IdGenerator = storage.UUIDGenerator{}
}
