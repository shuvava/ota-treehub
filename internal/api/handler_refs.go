package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shuvava/treehub/pkg/data"
	"github.com/shuvava/treehub/pkg/services"

	"github.com/labstack/echo/v4"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
)

const (
	// PathRefs is route for data.Ref operations
	PathRefs = "/refs/*"
)

// RefsUpload is endpoint uploading refs file to server from client
func RefsUpload(ctx echo.Context, svc *services.RefService) error {
	c := cmnapi.GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	refName := getRefNameFromPath(ctx)
	force := IsForcePush(ctx)
	commit, err := getCommitFromBody(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}

	err = svc.StoreRef(c, ns, refName, commit, force)
	if err != nil {
		return EchoResponse(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}

// RefDownload is endpoint download data.Ref file from server to client
func RefDownload(ctx echo.Context, svc *services.RefService) error {
	c := cmnapi.GetRequestContext(ctx)
	ns := GetNamespace(ctx)
	refName := getRefNameFromPath(ctx)
	exists, err := svc.Exists(c, ns, refName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	if !exists {
		err = fmt.Errorf("ref with namespace='%s' name='%s' does not exist", string(ns), refName)
		return ctx.JSON(http.StatusNotFound, cmnapi.NewErrorResponse(c, http.StatusNotFound, err))
	}
	ref, err := svc.GetRef(c, ns, refName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	return ctx.Blob(http.StatusOK, echo.MIMEOctetStream, []byte(ref.Value))
}

func getRefNameFromPath(ctx echo.Context) data.RefName {
	uri := ctx.Request().RequestURI
	parts := strings.Split(uri, "refs")
	ref := ""
	for inx, part := range parts {
		if inx == 0 {
			continue
		}
		ref = ref + part
	}
	return data.RefName(ref)
}

func getCommitFromBody(ctx echo.Context) (data.Commit, error) {
	if err := ValidateUploadContentType(ctx); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	if _, err := io.Copy(writer, ctx.Request().Body); err != nil {
		return "", err
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}
	return data.Commit(buf.String()), nil
}
