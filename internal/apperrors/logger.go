package apperrors

import "github.com/shuvava/go-logging/logger"

// CreateErrorAndLogIt create log record and throw an error
func CreateErrorAndLogIt(log logger.Logger, code AppErrorCode, descr string, err error) error {
	log.WithError(err).
		Error(descr)
	return CreateError(code, descr, err)
}
