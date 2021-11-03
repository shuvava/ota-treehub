package db

import (
	"context"

	"github.com/shuvava/treehub/pkg/data"
)

// RefRepository interface of CRUD operation with data.Ref
type RefRepository interface {
	// Create persist new data.Ref in database
	Create(ctx context.Context, ref data.Ref) error
	// Find looking up data.Ref in database
	Find(ctx context.Context, ns data.Namespace, name data.RefName) (*data.Ref, error)
	// Update change data.Ref properties in database
	Update(ctx context.Context, ref data.Ref) error
	// Delete removes data.Ref from database
	Delete(ctx context.Context, ns data.Namespace, name data.RefName) error
	// Exists checks if data.Ref exists in database
	Exists(ctx context.Context, ns data.Namespace, name data.RefName) (bool, error)
}
