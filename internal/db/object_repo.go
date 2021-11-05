package db

import (
	"context"

	"github.com/shuvava/treehub/pkg/data"
)

// ObjectRepository interface of CRUD operation with data.Object
type ObjectRepository interface {
	// Create persist new data.Object in database
	Create(ctx context.Context, obj data.Object) error
	// Find looking up data.Object in database
	Find(ctx context.Context, ns data.Namespace, id data.ObjectID) (*data.Object, error)
	// Update change data.Object properties in database
	Update(ctx context.Context, ns data.Namespace, id data.ObjectID, size int64, status data.ObjectStatus) error
	// Delete removes data.Object from database
	Delete(ctx context.Context, ns data.Namespace, id data.ObjectID) error
	// Exists checks if data.Object exists in database
	Exists(ctx context.Context, ns data.Namespace, id data.ObjectID) (bool, error)
	// SetCompleted change data.Object status to data.Uploaded
	SetCompleted(ctx context.Context, ns data.Namespace, id data.ObjectID) error
	// IsUploaded checks if data.Object was data.Uploaded
	IsUploaded(ctx context.Context, ns data.Namespace, id data.ObjectID) (bool, error)
	// FindAllByStatus returns all object with specific status
	FindAllByStatus(ctx context.Context, status data.ObjectStatus) ([]data.Object, error)
	// Usage returns space used by data.Namespace
	Usage(ctx context.Context, ns data.Namespace) (int64, error)
}
