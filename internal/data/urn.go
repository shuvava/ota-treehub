package data

import (
	"fmt"
)

const (
	urnPrefix = "here-ota"
)

// NewURN create new URN
func NewURN(urntype string, id CorrelationID) string {
	return fmt.Sprintf("urn:%s:%s:%s", urnPrefix, urntype, id.String())
}

// NewNamespaceURN creates Namespace URN
func NewNamespaceURN(id CorrelationID) string {
	return NewURN("namespace", id)
}

// NewMultiTargetUpdateURN creates CorrelationID URN
func NewMultiTargetUpdateURN(id CorrelationID) string {
	return NewURN("mtu", id)
}
