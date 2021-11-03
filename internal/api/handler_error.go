package api

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type errorProcessor func(error, echo.Context)

// ErrorHandler is a wrapper on echo.HTTPErrorHandler
type ErrorHandler struct {
	Handler    echo.HTTPErrorHandler
	processors map[string]errorProcessor
}

// NewErrorHandler sets up the mapping of error type to handler
func NewErrorHandler() *ErrorHandler {
	eh := ErrorHandler{}
	eh.Handler = eh.errorHandlerFunc
	eh.processors = make(map[string]errorProcessor)
	eh.processors[errorType(&echo.HTTPError{})] = echoHTTPErrorProcessor
	return &eh
}

func (eh *ErrorHandler) errorHandlerFunc(err error, c echo.Context) {
	p, found := eh.processors[errorType(err)]
	if !found {
		p = defaultErrorProcessor
	}
	p(err, c)
}

func echoHTTPErrorProcessor(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok && clientError(he.Code) {
		sendResponse(he.Code, he.Message.(string), c)
		return
	}
	defaultErrorProcessor(err, c)
}

func clientError(statusCode int) bool {
	return statusCode < http.StatusInternalServerError && statusCode >= http.StatusBadRequest
}

func defaultErrorProcessor(err error, c echo.Context) {
	sendResponse(http.StatusInternalServerError, err.Error(), c)
}

func sendResponse(code int, res interface{}, c echo.Context) {
	if !c.Response().Committed {
		var err error
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, res)
		}
		if err != nil {
			c.Logger().Errorf("Failed to send error response. Error: %v", err.Error())
		}
	}
}

func errorType(e error) string {
	return reflect.TypeOf(e).String()
}
