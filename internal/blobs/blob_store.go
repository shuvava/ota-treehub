package blobs

import (
	"context"
	"io"

	"github.com/shuvava/treehub/pkg/data"
)

// ObjectStore is common interface different implementation of Object stores
type ObjectStore interface {
	// StoreStream save file in store
	StoreStream(ctx context.Context, namespace data.Namespace, id data.ObjectID, reader io.Reader) (int64, error)
	// ReadFull read file content into memory
	ReadFull(ctx context.Context, namespace data.Namespace, id data.ObjectID, writer io.Writer) error
	// Exists checks if object exist on storage
	Exists(ctx context.Context, namespace data.Namespace, id data.ObjectID) (bool, error)
}
