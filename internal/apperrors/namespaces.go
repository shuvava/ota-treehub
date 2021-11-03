package apperrors

const (
	// ErrorNamespaceDB is error namespace for error related to database operations
	ErrorNamespaceDB = "db"
	// ErrorNamespaceData is error namespace for error related for data transformation/validation
	ErrorNamespaceData = "data"
	// ErrorNamespaceFs is error namespace for error related for file read/write operations
	ErrorNamespaceFs = "fs"
	// ErrorNamespaceSvc is error namespace for error related igh level business logic
	ErrorNamespaceSvc = "svc"

	// ErrorGeneric is untyped error
	ErrorGeneric = "generic-error"
)
