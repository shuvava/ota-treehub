package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/shuvava/treehub/pkg/services"

	"github.com/labstack/echo/v4"
)

const (
	// PathObject is route for data.Object operations
	PathObject = "/objects/:" + pathOPrefix + "/:" + pathOSuffix
)

// ObjectExists handler check if object exists
func ObjectExists(ctx echo.Context, svc *services.ObjectService) error {
	c := GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	id, err := GetObjectId(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}
	var exists bool
	exists, err = svc.Exists(c, ns, id)
	if exists {
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.NoContent(http.StatusNotFound)
}

// ObjectUploadCompleted handler updating object status
func ObjectUploadCompleted(ctx echo.Context, svc *services.ObjectService) error {
	c := GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	id, err := GetObjectId(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.SetCompleted(c, ns, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

// ObjectUpload is endpoint uploading data.Object file to server from client
func ObjectUpload(ctx echo.Context, svc *services.ObjectService) error {
	c := GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	id, err := GetObjectId(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}
	if err = ValidateUploadContentType(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}

	size := GetContentSize(ctx)
	if size == 0 {
		err := errors.New("Content-Length header is required to upload a file")
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.StoreStream(c, ns, id, size, ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

// ObjectDownload is endpoint download data.Object file from server to client
func ObjectDownload(ctx echo.Context, svc *services.ObjectService) error {
	c := GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	id, err := GetObjectId(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	}
	exists, err := svc.Exists(c, ns, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	if !exists {
		err = fmt.Errorf("object with namespace='%s' id='%s' does not exist", string(ns), id)
		return ctx.JSON(http.StatusNotFound, NewErrorResponse(c, http.StatusNotFound, err))
	}
	err = svc.ReadFull(c, ns, id, ctx.Response())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	ctx.Response().WriteHeader(http.StatusOK)
	ctx.Response().Flush()
	return nil
}
