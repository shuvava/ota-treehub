package apperrors

const (
	// ErrorSvcEntityExists is error to incorrect user operation( try accidentally replacing some entity)
	ErrorSvcEntityExists = ErrorNamespaceSvc + ":EntityAlreadyExist"
)
