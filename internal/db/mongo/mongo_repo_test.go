package mongo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/shuvava/treehub/internal/apperrors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/shuvava/treehub/internal/db/mongo"
	"github.com/shuvava/treehub/internal/logger"

	"github.com/sirupsen/logrus"
)

func TestMongoDB(t *testing.T) {
	t.Run("should get ErrorConnection if connection string is Invalid", func(t *testing.T) {
		var (
			connStr = "INVALID_CONNECTION_STRING"
		)
		ctx := context.Background()
		log := logger.NewNopLogger()
		_, err := mongo.NewMongoDB(ctx, log, connStr)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != apperrors.ErrorDbConnection {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbConnection)
		}
	})
}

type SomeDoc struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `bson:"name" json:"name"`
	Surname string             `bson:"surname" json:"surname"`
	Value   int                `bson:"value" json:"value"`
}

func TestMongoDB_IntegrationTest(t *testing.T) {
	var (
		connStr = "mongodb://mongoadmin:secret@localhost:27017/db?authSource=admin"
	)
	ctx := context.Background()
	log := logger.NewLogrusLogger(logrus.DebugLevel)

	t.Run("should able to connect to server", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		err = mdb.Ping(ctx)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

	t.Run("should be able read after write", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		doc := SomeDoc{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}
		coll := mdb.GetCollection("some_docs")
		id, err := mdb.InsertOne(ctx, coll, &doc)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if id != doc.ID.Hex() {
			t.Errorf("got %s want %s", id, doc.ID.Hex())
		}

		var got SomeDoc
		err = mdb.GetOneByID(ctx, coll, id, &got)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.Name != doc.Name {
			t.Errorf("got %s, expected %s", got.Name, doc.Name)
		}
		if got.Value != doc.Value {
			t.Errorf("got %d, expected %d", got.Value, doc.Value)
		}

		err = mdb.DeleteByID(ctx, coll, id)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}

		err = mdb.GetOneByID(ctx, coll, id, &got)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != apperrors.ErrorDbConnection {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbNoDocumentFound)
		}
	})

	t.Run("should be able check if obj exist", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		doc := SomeDoc{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}
		coll := mdb.GetCollection("some_docs_exist")
		id, err := mdb.InsertOne(ctx, coll, &doc)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}

		filter := bson.D{primitive.E{Key: "surname", Value: "Shurygin"}}
		cnt, err := mdb.Count(ctx, coll, filter)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if cnt < 1 {
			t.Errorf("got %d, expected 1", cnt)
		}

		err = mdb.DeleteByID(ctx, coll, id)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

	t.Run("should be able find by object properties", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		docs := []SomeDoc{{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}, {
			ID:      primitive.NewObjectID(),
			Name:    "Alex",
			Surname: "Shurygin",
			Value:   10,
		}}
		coll := mdb.GetCollection("some_docs_finds")
		for _, doc := range docs {
			_, err := mdb.InsertOne(ctx, coll, &doc)
			if err != nil {
				t.Errorf("got %s, expected nil", err)
			}
		}
		filter := bson.D{primitive.E{Key: "surname", Value: "Shurygin"}}
		var gotDocs []*SomeDoc
		err = mdb.Find(ctx, coll, filter, &gotDocs)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if len(gotDocs) < 2 {
			t.Errorf("got %d, expected equal or great 2", len(gotDocs))
		}
		err = mdb.Delete(ctx, coll, filter)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

	t.Run("should be able update object", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		doc := SomeDoc{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}
		coll := mdb.GetCollection("some_docs_v3")
		id, err := mdb.InsertOne(ctx, coll, &doc)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if id != doc.ID.Hex() {
			t.Errorf("got %s want %s", id, doc.ID.Hex())
		}
		filter := bson.D{primitive.E{Key: "name", Value: doc.Name}}
		newName := fmt.Sprintf("%s-updated", doc.Name)
		update := bson.D{
			{"$set", bson.D{
				{"name", newName},
			}},
		}

		err = mdb.UpdateOne(ctx, coll, filter, update)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		var got SomeDoc
		err = mdb.GetOneByID(ctx, coll, id, &got)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.Name != newName {
			t.Errorf("got %s want %s", got.Name, newName)
		}
		filter = bson.D{primitive.E{Key: "surname", Value: "Shurygin"}}
		err = mdb.Delete(ctx, coll, filter)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

	t.Run("should be able replace object", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		doc := SomeDoc{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}
		coll := mdb.GetCollection("some_docs_replace")
		id, err := mdb.InsertOne(ctx, coll, &doc)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if id != doc.ID.Hex() {
			t.Errorf("got %s want %s", id, doc.ID.Hex())
		}
		filter := bson.D{primitive.E{Key: "name", Value: doc.Name}}
		docNew := SomeDoc{
			Name:    "Vlad-Replaced",
			Surname: "Shurygin",
			Value:   10,
		}

		err = mdb.ReplaceOne(ctx, coll, filter, &docNew)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		var got SomeDoc
		err = mdb.GetOneByID(ctx, coll, id, &got)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if got.Name != docNew.Name {
			t.Errorf("got %s want %s", got.Name, docNew.Name)
		}
		filter = bson.D{primitive.E{Key: "surname", Value: "Shurygin"}}
		err = mdb.Delete(ctx, coll, filter)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

	t.Run("should be able =run custom agg query", func(t *testing.T) {
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		docs := []SomeDoc{{
			ID:      primitive.NewObjectID(),
			Name:    "Vlad",
			Surname: "Shurygin",
			Value:   40,
		}, {
			ID:      primitive.NewObjectID(),
			Name:    "Alex",
			Surname: "Shurygin",
			Value:   10,
		}, {
			ID:      primitive.NewObjectID(),
			Name:    "Ivan",
			Surname: "Petrov",
			Value:   100,
		}}
		coll := mdb.GetCollection("some_docs_agg")
		for _, doc := range docs {
			_, err := mdb.InsertOne(ctx, coll, &doc)
			if err != nil {
				t.Errorf("got %s, expected nil", err)
			}
		}
		pipeline := make([]bson.M, 0)

		groupStage := bson.M{
			"$group": bson.M{
				"_id":   "$surname",
				"total": bson.M{"$sum": "$value"},
			},
		}

		matchStage := bson.M{
			"$match": bson.M{
				"surname": "Shurygin",
			},
		}
		pipeline = append(pipeline, matchStage, groupStage)

		data, err := mdb.Aggregate(ctx, coll, pipeline)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if len(data) < 1 {
			t.Errorf("got %d, expected equal or great 1", len(data))
		}
		if data[0]["total"].(int32) != 50 {
			t.Errorf("got %d, expected %d", data[0]["total"], 50)
		}
		filter := bson.D{primitive.E{
			Key: "$or",
			Value: bson.A{
				bson.D{primitive.E{Key: "surname", Value: "Shurygin"}},
				bson.D{primitive.E{Key: "surname", Value: "Petrov"}},
			},
		}}
		err = mdb.Delete(ctx, coll, filter)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})

}
