package data

import (
	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

// Object OSTree object definition
type Object struct {
	Namespace cmndata.Namespace
	ID        ObjectID
	ByteSize  int64
	Status    ObjectStatus
}

//func NewObject(str string) *Object
