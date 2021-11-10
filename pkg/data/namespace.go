package data

import "github.com/shuvava/go-ota-svc-common/data"

// Namespace is object namespace
type Namespace string

// NewNamespace creates namespace with random uuid
func NewNamespace(id string) Namespace {
	if id == "" {
		id = data.NewNamespaceURN(data.NewCorrelationID())
	}

	return Namespace(id)
}
