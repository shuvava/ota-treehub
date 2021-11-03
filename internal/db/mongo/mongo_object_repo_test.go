package mongo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shuvava/treehub/internal/apperrors"

	"github.com/shuvava/treehub/internal/db/mongo"

	intmongo "github.com/shuvava/treehub/internal/db/mongo"
	"github.com/shuvava/treehub/internal/logger"
	"github.com/shuvava/treehub/pkg/data"

	"github.com/sirupsen/logrus"
)

func TestMongoObjectStore_IntegrationTest(t *testing.T) {
	const connStr = "mongodb://mongoadmin:secret@localhost:27017/db?authSource=admin"
	ctx := context.Background()
	log := logger.NewLogrusLogger(logrus.DebugLevel)
	mdb, err := intmongo.NewMongoDB(ctx, log, connStr)
	if err != nil {
		t.Errorf("got %s, expected nil", err)
	}
	err = mdb.Ping(ctx)
	if err != nil {
		t.Errorf("got %s, expected nil", err)
	}
	store := mongo.NewObjectMongoRepository(log, mdb)

	t.Run("should fail if object already exists", func(t *testing.T) {
		ns := data.Namespace("obj_create")
		id := data.ObjectID("1")
		status := data.CLIENT_UPLOADING
		doc := data.Object{
			ID:        id,
			Namespace: ns,
			ByteSize:  10,
			Status:    status,
		}
		err := store.Create(ctx, doc)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		doc.ByteSize = doc.ByteSize * 2
		var typedErr apperrors.AppError
		err = store.Create(ctx, doc)
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != apperrors.ErrorDbAlreadyExist {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbAlreadyExist)
		}
		// cleanup
		err = store.Delete(ctx, doc.Namespace, doc.ID)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})
	t.Run("should be able read after write", func(t *testing.T) {
		id1 := data.ObjectID("1")
		id2 := data.ObjectID("2")
		ns := data.Namespace("test1")
		status := data.CLIENT_UPLOADING
		docs := []data.Object{
			{
				ID:        id1,
				Namespace: ns,
				ByteSize:  10,
				Status:    status,
			},
			{
				ID:        id2,
				Namespace: ns,
				ByteSize:  15,
				Status:    status,
			},
			{
				ID:        data.ObjectID("3"),
				Namespace: data.Namespace("test2"),
				ByteSize:  15,
				Status:    data.SERVER_UPLOADING,
			},
		}
		for _, doc := range docs {
			err := store.Create(ctx, doc)
			if err != nil {
				t.Errorf("got %s, expected nil", err)
			}
		}
		got, err := store.Find(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.Status != status {
			t.Errorf("got %d want %d", got.Status, status)
		}
		err = store.Update(ctx, ns, id1, 20, status)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		got, err = store.Find(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.ByteSize != 20 {
			t.Errorf("got %d want %d", got.ByteSize, 20)
		}
		exs, err := store.Exists(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if exs != true {
			t.Errorf("got %t want %t", exs, true)
		}
		exs, err = store.Exists(ctx, ns, "100")
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if exs != false {
			t.Errorf("got %t want %t", exs, false)
		}
		err = store.SetCompleted(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		got, err = store.Find(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.Status != data.UPLOADED {
			t.Errorf("got %d want %d", got.Status, data.UPLOADED)
		}
		exs, err = store.IsUploaded(ctx, ns, id1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if exs != true {
			t.Errorf("got %t want %t", exs, true)
		}
		exs, err = store.IsUploaded(ctx, ns, id2)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if exs != false {
			t.Errorf("got %t want %t", exs, false)
		}
		gotDocs, err := store.FindAllByStatus(ctx, status)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if len(gotDocs) != 1 {
			t.Errorf("got %d want %d", len(gotDocs), 1)
		} else {
			if gotDocs[0].ID != id2 {
				t.Errorf("got %s want %s", gotDocs[0].ID, id2)
			}
		}
		gotUsage, err := store.Usage(ctx, ns)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if gotUsage != 35 {
			t.Errorf("got %d want %d", gotUsage, 35)
		}
		for _, doc := range docs {
			err = store.Delete(ctx, doc.Namespace, doc.ID)
			if err != nil {
				t.Errorf("got %s, expected nil", err)
			}
		}
	})
}
