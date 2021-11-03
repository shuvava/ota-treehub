package apperrors

const (
	// ErrorFsPath is error type related file path availability
	ErrorFsPath = ErrorNamespaceFs + ":PathNotExist"
	// ErrorFsIOOpen is error type related opening file
	ErrorFsIOOpen = ErrorNamespaceFs + ":IOOpen"
	// ErrorFsIOOperation is error type related file io operations
	ErrorFsIOOperation = ErrorNamespaceFs + ":IOOperation"
	// ErrorFsIOCreate is error type related file creating operation
	ErrorFsIOCreate = ErrorNamespaceFs + ":IOCreate"
)
