package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/shuvava/treehub/internal/db"
	"github.com/shuvava/treehub/pkg/data"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	cmndata "github.com/shuvava/go-ota-svc-common/data"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const objectTableName = "objects"

type objectDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ObjectID  string             `bson:"id"`
	Namespace string             `bson:"namespace"`
	ByteSize  int64              `bson:"byteSize"`
	Status    int                `bson:"status"`
}

// ObjectMongoRepository implementations of db.ObjectRepository for MongoDb repo
type ObjectMongoRepository struct {
	db   *intMongo.Db
	coll *mongo.Collection
	log  logger.Logger
	db.ObjectRepository
}

// NewObjectMongoRepository creates new instance of ObjectMongoRepository
func NewObjectMongoRepository(logger logger.Logger, db *intMongo.Db) *ObjectMongoRepository {
	log := logger.SetOperation("ObjectRepo")
	return &ObjectMongoRepository{
		db:   db,
		coll: db.GetCollection(objectTableName),
		log:  log,
	}
}

// Create persist new data.Object in MongoDB
func (store *ObjectMongoRepository) Create(ctx context.Context, obj data.Object) error {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", obj.ID).
		WithField("Namespace", obj.Namespace).
		Debug("Creating new Object")
	dto := objectToDTO(obj)
	exists, err := store.Exists(ctx, obj.Namespace, obj.ID)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document(Object) with id='%s' namespace='%s' already exist in database", obj.ID, obj.Namespace)
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	_, err = store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.
			WithField("ObjectID", obj.ID).
			WithField("Namespace", obj.Namespace).
			Debug("Object created successful")
	} else {
		log.
			WithField("ObjectID", obj.ID).
			WithField("Namespace", obj.Namespace).
			Warn("Object creation failed")
	}
	return err
}

// Find looking up data.Object in database
func (store *ObjectMongoRepository) Find(ctx context.Context, ns cmndata.Namespace, id data.ObjectID) (*data.Object, error) {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Looking up object")
	filter := getOneObjectFilter(ns, id)
	var dto objectDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("ObjectID", id).
				WithField("Namespace", ns).
				Warn("Object not found")
		}
		return nil, err
	}
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Object Found")
	model := objectDtoToModel(dto)
	return &model, nil
}

// Update change data.Object properties
func (store *ObjectMongoRepository) Update(ctx context.Context, ns cmndata.Namespace, id data.ObjectID, size int64, status data.ObjectStatus) error {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Update object")
	filter := getOneObjectFilter(ns, id)
	upd := bson.D{primitive.E{
		Key: "$set", Value: bson.M{
			"byteSize": size,
			"status":   int(status),
		},
	}}
	err := store.db.UpdateOne(ctx, store.coll, filter, upd)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("ObjectID", id).
				WithField("Namespace", ns).
				Warn("Object not found")
		} else {
			log.WithField("ObjectID", id).
				WithField("Namespace", ns).
				Warn("Object updating failed")
		}
		return err
	}
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Object updated successful")
	return nil
}

// Delete removes object in mongo database
func (store *ObjectMongoRepository) Delete(ctx context.Context, ns cmndata.Namespace, id data.ObjectID) error {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Deleting object")
	filter := getOneObjectFilter(ns, id)
	err := store.db.Delete(ctx, store.coll, filter)
	if err != nil {
		log.WithField("ObjectID", id).
			WithField("Namespace", ns).
			Warn("Object delete failed")
		return err
	}
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Object deleted")
	return nil
}

// Exists checks if data.Object exists in mongo database
func (store *ObjectMongoRepository) Exists(ctx context.Context, ns cmndata.Namespace, id data.ObjectID) (bool, error) {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Looking up object")
	filter := getOneObjectFilter(ns, id)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// SetCompleted change data.Object status to data.Uploaded
func (store *ObjectMongoRepository) SetCompleted(ctx context.Context, ns cmndata.Namespace, id data.ObjectID) error {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Update object status")
	filter := bson.D{primitive.E{
		Key: "$and", Value: bson.A{
			bson.D{primitive.E{Key: "id", Value: id}},
			bson.D{primitive.E{Key: "namespace", Value: ns}},
			bson.D{primitive.E{Key: "status", Value: int(data.ServerUploading)}},
		},
	}}
	upd := bson.D{primitive.E{
		Key: "$set", Value: bson.M{
			"status": int(data.Uploaded),
		},
	}}
	err := store.db.UpdateOne(ctx, store.coll, filter, upd)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("ObjectID", id).
				WithField("Namespace", ns).
				Warn("Object not found")
		} else {
			log.WithField("ObjectID", id).
				WithField("Namespace", ns).
				Warn("Updating object status failed")
		}
		return err
	}
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Updated successful")
	return nil
}

// IsUploaded checks if data.Object was data.Uploaded
func (store *ObjectMongoRepository) IsUploaded(ctx context.Context, ns cmndata.Namespace, id data.ObjectID) (bool, error) {
	log := store.log.WithContext(ctx)
	log.WithField("ObjectID", id).
		WithField("Namespace", ns).
		Debug("Looking up object")
	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "id", Value: id}},
			bson.D{primitive.E{Key: "namespace", Value: ns}},
			bson.D{primitive.E{Key: "status", Value: int(data.Uploaded)}},
		},
	}}
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// FindAllByStatus returns all object with specific status
func (store *ObjectMongoRepository) FindAllByStatus(ctx context.Context, status data.ObjectStatus) ([]data.Object, error) {
	log := store.log.WithContext(ctx)
	log.WithField("status", status).
		Debug("Looking up objects")

	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "status", Value: int(status)}},
		},
	}}
	var docs []objectDTO
	err := store.db.Find(ctx, store.coll, filter, &docs)
	if err != nil {
		log.WithField("status", status).
			Debug("Not Found")
		return nil, err
	}

	var res []data.Object
	for _, doc := range docs {
		obj := objectDtoToModel(doc)
		res = append(res, obj)
	}

	log.WithField("status", status).
		WithField("Count", len(res)).
		Debug("Lookup completed successful")

	return res, nil
}

// Usage returns space used by data.Namespace
func (store *ObjectMongoRepository) Usage(ctx context.Context, ns cmndata.Namespace) (int64, error) {
	log := store.log.WithContext(ctx)
	log.WithField("Namespace", ns).
		Debug("Usage stats")

	pipeline := make([]bson.M, 0)
	groupStage := bson.M{
		"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$byteSize"},
		},
	}
	matchStage := bson.M{
		"$match": bson.M{
			"namespace": ns,
		},
	}
	pipeline = append(pipeline, matchStage, groupStage)

	var res []intMongo.DBResult
	if err := store.db.Aggregate(ctx, store.coll, pipeline, nil, &res); err != nil {
		return 0, err
	}
	if len(res) < 1 {
		return 0, apperrors.NewAppError(apperrors.ErrorDbOperation, "unexpected aggregation result")
	}

	return res[0]["total"].(int64), nil
}

// objectToDTO converts data.Object to objectDTO
func objectToDTO(obj data.Object) objectDTO {
	dto := objectDTO{
		ID:        primitive.NewObjectID(),
		ObjectID:  string(obj.ID),
		Namespace: string(obj.Namespace),
		ByteSize:  obj.ByteSize,
		Status:    int(obj.Status),
	}
	return dto
}

// objectDtoToModel converts objectDTO to data.Object
func objectDtoToModel(dto objectDTO) data.Object {
	model := data.Object{
		Namespace: cmndata.Namespace(dto.Namespace),
		ID:        data.ObjectID(dto.ObjectID),
		ByteSize:  dto.ByteSize,
		Status:    data.ObjectStatus(dto.Status),
	}
	return model
}

func getOneObjectFilter(ns cmndata.Namespace, id data.ObjectID) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "id", Value: id}},
			bson.D{primitive.E{Key: "namespace", Value: ns}},
		},
	}}
}
