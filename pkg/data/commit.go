package data

import (
	"fmt"

	"github.com/shuvava/go-ota-svc-common/data"
)

// Commit is OTSTree commit object
type Commit string

// CommitManifest is OSTree commit manifest
type CommitManifest struct {
	Namespace Namespace
	Commit    ObjectID
	Contents  string
}

//Validate if Commit has valid format
func (obj Commit) Validate() error {
	err := fmt.Errorf("%s is not a sha-256 commit hash", obj)
	sha := string(obj)
	if !data.ValidHex(64, sha) {
		return err
	}

	return nil
}

// From converts Commit to ObjectID
func (obj Commit) From() (ObjectID, error) {
	s := fmt.Sprintf("%s.commit", string(obj))
	return NewObjectID(s)
}

// NewCommit validate str and create Commit object
func NewCommit(str string) (Commit, error) {
	obj := Commit(str)
	if err := obj.Validate(); err != nil {
		return "", err
	}
	return obj, nil
}

// NewCommitFromBytes validate content and create Commit object
func NewCommitFromBytes(content []byte) (Commit, error) {
	str := data.ByteDigest(content)
	return NewCommit(str)
}
