package data

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/shuvava/go-ota-svc-common/data"
)

// DeltaID is OSTree delta
type DeltaID string

func (delta DeltaID) formattingError() error {
	return fmt.Errorf("%s is not a valid DeltaID (cc/mbase64(from.rest)-mbase64(to)", delta)
}

// Validate if DeltaID has valid format
func (delta DeltaID) Validate() error {
	parts := strings.Split(string(delta), "-")
	if len(parts) != 2 {
		return delta.formattingError()
	}
	head := parts[0]
	tail := parts[1]
	if _, e := data.ToBase64(head); e != nil {
		return delta.formattingError()
	}
	if _, e := data.ToBase64(tail); e != nil {
		return delta.formattingError()
	}

	return nil
}

// URLSafe converts DeltaID string to url safe encoding
func (delta DeltaID) URLSafe() string {
	return strings.ReplaceAll(string(delta), "+", "_")
}

// ToObjectID transforms DeltaID to ObjectID
func (delta DeltaID) ToObjectID() (ObjectID, error) {
	parts := strings.Split(string(delta), "-")
	if len(parts) != 2 {
		return "", delta.formattingError()
	}
	bytes, err := data.ToBase64(parts[1])
	if err != nil {
		return "", err
	}
	commit := hex.EncodeToString(bytes)
	return NewObjectID(fmt.Sprintf("%s.commit", commit))
}
