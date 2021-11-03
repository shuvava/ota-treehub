package apperrors

const (
	// ErrorDataObjectIDSerialization is error type returned if data.ObjectID was not serialized successfully
	ErrorDataObjectIDSerialization = ErrorNamespaceData + ":ObjectIdSerialization"
	// ErrorDataRefValidation is error type returned if data.Ref is not pass validation
	ErrorDataRefValidation = ErrorNamespaceData + ":RefValidation"
)
