package data

// Object OSTree object definition
type Object struct {
	Namespace Namespace
	ID        ObjectID
	ByteSize  int64
	Status    ObjectStatus
}

//func NewObject(str string) *Object
