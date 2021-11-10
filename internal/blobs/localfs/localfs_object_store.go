package localfs

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/shuvava/treehub/internal/blobs"
	"github.com/shuvava/treehub/internal/utils/fshelper"

	"github.com/shuvava/treehub/pkg/data"
)

// ObjectLocalFsStore implementation of ObjectRepo interface for local store
type ObjectLocalFsStore struct {
	blobs.ObjectStore
	log  logger.Logger
	root string
}

// NewLocalFsBlobStore creates ObjectLocalFsStore object
func NewLocalFsBlobStore(root string, log logger.Logger) (*ObjectLocalFsStore, error) {
	if err := fshelper.EnsureDir(root); err != nil {
		return nil, apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorFsPath,
			"Failed to get root directory", err)
	}
	return &ObjectLocalFsStore{
		root: root,
		log:  log,
	}, nil
}

// StoreStream persist Object in local path
func (store *ObjectLocalFsStore) StoreStream(ctx context.Context, ns data.Namespace, id data.ObjectID, reader io.Reader) (int64, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Persisting object to file system")
	path, err := store.objectPath(ctx, ns, id)
	if err != nil {
		return 0, err
	}
	log.
		WithField("filename", path).
		Debug("Persisting stream into blob")
	written, err := safeStoreStream(path, reader)
	if err != nil {
		return 0, apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorFsIOOperation,
			"Failed to persist file stream", err)
	}
	log.
		WithField("filename", path).
		WithField("size", written).
		Debug("Blob created")
	return written, nil
}

// ReadFull read the whole file into memory and return content
func (store *ObjectLocalFsStore) ReadFull(ctx context.Context, ns data.Namespace, id data.ObjectID, writer io.Writer) error {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Reading object from file system")
	path, err := store.objectPath(ctx, ns, id)
	if err != nil {
		return err
	}
	log.
		WithField("filename", path).
		Debug("Reading blob")

	file, err := os.Open(path)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorFsIOOpen,
			"Failed to open file", err)
	}
	defer func() { _ = file.Close() }()

	written, err := io.Copy(writer, file)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorFsIOOperation,
			"Failed to read file", err)
	}
	log.
		WithField("filename", path).
		WithField("size", written).
		Debug("Object reading completed")

	return nil
}

// Exists checks if object exist on local storage
func (store *ObjectLocalFsStore) Exists(ctx context.Context, ns data.Namespace, id data.ObjectID) (bool, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Looking up object in file system")
	path, err := store.objectPath(ctx, ns, id)
	if err != nil {
		return false, err
	}
	exists := fshelper.IsPathExist(path)
	return exists, nil
}

func (store *ObjectLocalFsStore) namespacePath(ns data.Namespace) string {
	return filepath.Join(store.root, string(ns))
}

func (store *ObjectLocalFsStore) objectPath(ctx context.Context, ns data.Namespace, id data.ObjectID) (string, error) {
	log := store.log.WithContext(ctx)
	path := filepath.Join(store.namespacePath(ns), string(id))
	parent := filepath.Dir(path)

	err := fshelper.EnsureDir(parent)
	if err != nil {
		return "", apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorFsPath,
			"Failed to create object directory", err)
	}

	return path, nil
}
