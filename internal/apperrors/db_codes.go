package apperrors

const (
	// ErrorDbConnection is error type related to establishing connection to db server
	ErrorDbConnection = ErrorNamespaceDB + ":ConnectionError"
	// ErrorDbOperation is error type related to execution db CRUD operations
	ErrorDbOperation = ErrorNamespaceDB + ":OperationError"
	// ErrorDbNoDocumentFound is error type returned if no document found
	ErrorDbNoDocumentFound = ErrorNamespaceDB + ":DocumentNotFound"
	// ErrorDbAlreadyExist is error type returned on creation of document if document with doc id already exists
	ErrorDbAlreadyExist = ErrorNamespaceDB + ":DocumentAlreadyExist"
)
