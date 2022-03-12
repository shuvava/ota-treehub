package data

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shuvava/go-ota-svc-common/data"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// ObjectID is OSTree object identifier
type ObjectID string

// ErrorDataSerializationObjectID is error of data.ObjectID serialization
const ErrorDataSerializationObjectID = apperrors.ErrorDataSerialization + ":ObjectID"

// Validate if ObjectId has valid format
func (objectId ObjectID) Validate() error {
	err := apperrors.NewAppError(
		ErrorDataSerializationObjectID,
		fmt.Sprintf("%s must be in format <sha256>.objectType", objectId))
	parts := strings.Split(string(objectId), ".")
	if len(parts) != 2 {
		return err
	}
	sha := parts[0]
	objectType := parts[1]
	if len(objectType) == 0 || !data.ValidHex(64, sha) {
		return err
	}

	return nil
}

// Path returns absolute path to ObjectID is storage
func (objectId ObjectID) Path(parent string) string {
	s := string(objectId)
	prefix := s[:2]
	rest := s[2:]
	return filepath.Join(parent, prefix, rest)
}

// Filename returns file name of ObjectID in storage
func (objectId ObjectID) Filename() string {
	path := objectId.Path("/")
	return filepath.Base(path)
}

// NewObjectID create new ObjectID if str is valid
func NewObjectID(str string) (ObjectID, error) {
	obj := ObjectID(str)
	if err := obj.Validate(); err != nil {
		return "", err
	}
	return obj, nil
}
