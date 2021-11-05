package services

import (
	"context"
	"io"

	"github.com/shuvava/go-logging/logger"

	objstore "github.com/shuvava/treehub/internal/blobs"
	"github.com/shuvava/treehub/internal/db"
	"github.com/shuvava/treehub/pkg/data"
)

// ObjectService is service for interaction with data.Object
type ObjectService struct {
	log logger.Logger
	db  db.ObjectRepository
	fs  objstore.ObjectStore
}

// NewObjectService creates new instance of ObjectService
func NewObjectService(l logger.Logger, db db.ObjectRepository, fs objstore.ObjectStore) *ObjectService {
	log := l.SetContext("object-service")
	return &ObjectService{
		log: log,
		db:  db,
		fs:  fs,
	}
}

// SetCompleted change data.Object status to data.UPLOADED
func (svc *ObjectService) SetCompleted(ctx context.Context, ns data.Namespace, id data.ObjectID) error {
	return svc.db.SetCompleted(ctx, ns, id)
}

// Exists checks if data.Object exist on storage
func (svc *ObjectService) Exists(ctx context.Context, ns data.Namespace, id data.ObjectID) (bool, error) {
	fsExists, err := svc.fs.Exists(ctx, ns, id)
	if err != nil {
		return false, err
	}
	dbExists, err := svc.db.Exists(ctx, ns, id)
	if err != nil {
		return false, err
	}
	return dbExists && fsExists, nil
}

// StoreStream save data.Object
func (svc *ObjectService) StoreStream(ctx context.Context, ns data.Namespace, id data.ObjectID, size int64, reader io.Reader) error {
	log := svc.log.WithContext(ctx)
	obj := data.Object{
		Namespace: ns,
		ID:        id,
		ByteSize:  size,
		Status:    data.SERVER_UPLOADING,
	}
	exists, err := svc.db.Exists(ctx, ns, id)
	if err != nil {
		return err
	}
	if !exists {
		if err = svc.db.Create(ctx, obj); err != nil {
			return err
		}
	}
	written, err := svc.fs.StoreStream(ctx, ns, id, reader)
	if err != nil {
		_ = svc.db.Delete(ctx, ns, id) // skip error, because of was logged on repo level
		return err
	}
	if written != size {
		log.WithField("ObjectID", id).
			WithField("Namespace", ns).
			Warn("Uploaded size(", written, ") does not match expected(", size, ")")
	}
	if err = svc.db.Update(ctx, ns, id, written, data.UPLOADED); err != nil {
		_ = svc.db.Delete(ctx, ns, id) // skip error, because of was logged on repo level
		return err
	}
	return nil
}

// ReadFull read data.Object
func (svc *ObjectService) ReadFull(ctx context.Context, ns data.Namespace, id data.ObjectID, writer io.Writer) error {
	return svc.fs.ReadFull(ctx, ns, id, writer)
}
