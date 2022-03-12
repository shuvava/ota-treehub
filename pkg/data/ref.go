package data

import (
	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

// RefName is name of OSTree Ref
type RefName string

// Ref is OSTree object reference
type Ref struct {
	Namespace cmndata.Namespace
	Name      RefName
	Value     Commit
	ObjectID  ObjectID
}

// Validate doing validation of Ref
func (ref Ref) Validate() error {
	if err := ref.Value.Validate(); err != nil {
		return err
	}
	if err := ref.ObjectID.Validate(); err != nil {
		return err
	}
	return nil
}

// NewRef create new instance of Ref
func NewRef(ns cmndata.Namespace, refName RefName, value Commit) (Ref, error) {
	objID, err := value.From()
	if err != nil {
		return Ref{}, err
	}
	ref := Ref{
		Namespace: ns,
		Name:      refName,
		Value:     value,
		ObjectID:  objID,
	}
	return ref, ref.Validate()
}
