package localfs

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	cmndata "github.com/shuvava/go-ota-svc-common/data"
	intdata "github.com/shuvava/go-ota-svc-common/data"
	"github.com/sirupsen/logrus"

	"github.com/shuvava/treehub/pkg/data"
)

const useNopLogger = false

func TestObjectLocalFsStore(t *testing.T) {
	checkOnNil := func(err error) {
		t.Helper()
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	}
	checkBool := func(got, want bool) {
		t.Helper()
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	}
	checkInt64 := func(got, want int64) {
		t.Helper()
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}
	checkStr := func(got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dir, e := ioutil.TempDir("", "treehub_object")
	checkOnNil(e)
	var log logger.Logger
	if useNopLogger {
		log = logger.NewNopLogger()
	} else {
		log = logger.NewLogrusLogger(logrus.DebugLevel)
	}
	store, e := NewLocalFsBlobStore(dir, log)
	checkOnNil(e)
	ns := cmndata.Namespace("test")

	t.Run("Exists func should return false if file not exists", func(t *testing.T) {
		id := data.ObjectID("1")
		got, err := store.Exists(ctx, ns, id)
		checkOnNil(err)
		checkBool(got, false)
	})
	t.Run("Exists func should return true if file exists", func(t *testing.T) {
		id := data.ObjectID("112345")
		path, err := store.objectPath(ctx, ns, id)
		checkOnNil(err)
		file, err := os.Create(path)
		d1 := []byte("hello\ngo\n")
		_, _ = file.Write(d1)
		checkOnNil(err)
		err = file.Close()
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		got, err := store.Exists(ctx, ns, id)
		checkOnNil(err)
		checkBool(got, true)
	})
	t.Run("StoreStream should be able create new file and put content", func(t *testing.T) {
		text := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."
		want := int64(len(text))
		id := data.ObjectID(intdata.NewCorrelationID().String())
		reader := strings.NewReader(text)
		got, err := store.StoreStream(ctx, ns, id, reader)
		checkOnNil(err)
		checkInt64(got, want)
	})
	t.Run("StoreStream should be able rewrite file content", func(t *testing.T) {
		textOrigin := "Lorem non."
		text := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."
		id := data.ObjectID(intdata.NewCorrelationID().String())
		reader := strings.NewReader(textOrigin)
		got, err := store.StoreStream(ctx, ns, id, reader)
		checkOnNil(err)
		checkInt64(got, int64(len(textOrigin)))
		reader = strings.NewReader(text)
		got, err = store.StoreStream(ctx, ns, id, reader)
		checkOnNil(err)
		checkInt64(got, int64(len(text)))
	})
	t.Run("ReadFull should be able read written content", func(t *testing.T) {
		text := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."
		id := data.ObjectID(intdata.NewCorrelationID().String())
		reader := strings.NewReader(text)
		_, err := store.StoreStream(ctx, ns, id, reader)
		checkOnNil(err)
		var buf bytes.Buffer
		stream := bufio.NewWriter(&buf)
		checkOnNil(store.ReadFull(ctx, ns, id, stream))
		_ = stream.Flush()
		got := buf.String()
		checkStr(got, text)
	})
	t.Run("ReadFull should return error if file not exists", func(t *testing.T) {
		id := data.ObjectID(intdata.NewCorrelationID().String())
		var buf bytes.Buffer
		stream := bufio.NewWriter(&buf)
		err := store.ReadFull(ctx, ns, id, stream)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != apperrors.ErrorFsIOOpen {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorFsIOOpen)
		}
	})
}
