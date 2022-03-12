package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	cmndata "github.com/shuvava/go-ota-svc-common/data"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/shuvava/treehub/internal/db"
	"github.com/shuvava/treehub/pkg/data"
)

const refTableName = "refs"

type refDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Namespace string             `bson:"namespace"`
	Value     string             `bson:"value"`
	ObjectID  string             `bson:"objectId"`
}

// RefMongoRepository implementations of db.RefRepository for MongoDb repo
type RefMongoRepository struct {
	db   *intMongo.Db
	coll *mongo.Collection
	log  logger.Logger
	db.RefRepository
}

// NewRefMongoRepository creates new instance of RefMongoRepository
func NewRefMongoRepository(logger logger.Logger, db *intMongo.Db) *RefMongoRepository {
	log := logger.SetOperation("RefRepo")
	return &RefMongoRepository{
		db:   db,
		coll: db.GetCollection(refTableName),
		log:  log,
	}
}

// Create persist new data.Ref in MongoDB
func (store *RefMongoRepository) Create(ctx context.Context, ref data.Ref) error {
	log := store.log.WithContext(ctx)
	log.WithField("Name", ref.Name).
		WithField("Namespace", ref.Namespace).
		Debug("Creating new Ref")
	dto := refToDTO(ref)
	exists, err := store.Exists(ctx, ref.Namespace, ref.Name)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document (Ref) with Name='%s' namespace='%s' already exist in database", ref.Name, ref.Namespace)
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	_, err = store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.WithField("Name", ref.Name).
			WithField("Namespace", ref.Namespace).
			Debug("Ref created successful")
	} else {
		log.WithField("Name", ref.Name).
			WithField("Namespace", ref.Namespace).
			Warn("Ref creation failed")
	}
	return err
}

// Find looking up data.Ref in database
func (store *RefMongoRepository) Find(ctx context.Context, ns cmndata.Namespace, name data.RefName) (*data.Ref, error) {
	log := store.log.WithContext(ctx)
	log.WithField("Name", name).
		WithField("Namespace", ns).
		Debug("Looking up ref")
	filter := getOneRefFilter(ns, name)
	var dto refDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("Name", name).
				WithField("Namespace", ns).
				Warn("Ref not found")
		}
		return nil, err
	}
	log.WithField("Name", name).
		WithField("Namespace", ns).
		Debug("Ref Found")
	model := refDtoToModel(dto)
	return &model, nil
}

// Update change data.Ref properties
func (store *RefMongoRepository) Update(ctx context.Context, ref data.Ref) error {
	log := store.log.WithContext(ctx)
	log.WithField("Name", ref.Name).
		WithField("Namespace", ref.Namespace).
		Debug("Update ref")
	filter := getOneRefFilter(ref.Namespace, ref.Name)
	upd := bson.D{primitive.E{
		Key: "$set", Value: bson.M{
			"value":    string(ref.Value),
			"objectId": string(ref.ObjectID),
		},
	}}
	err := store.db.UpdateOne(ctx, store.coll, filter, upd)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("Name", ref.Name).
				WithField("Namespace", ref.Namespace).
				Warn("Ref not found")
		} else {
			log.WithField("Name", ref.Name).
				WithField("Namespace", ref.Namespace).
				Warn("Ref updating failed")
		}
		return err
	}
	log.WithField("Name", ref.Name).
		WithField("Namespace", ref.Namespace).
		Debug("Ref updated successful")
	return nil
}

// Delete removes data.Ref from mongo database
func (store *RefMongoRepository) Delete(ctx context.Context, ns cmndata.Namespace, name data.RefName) error {
	log := store.log.WithContext(ctx)
	log.WithField("Name", name).
		WithField("Namespace", ns).
		Debug("Deleting ref")
	filter := getOneRefFilter(ns, name)
	err := store.db.Delete(ctx, store.coll, filter)
	if err != nil {
		log.WithField("Name", name).
			WithField("Namespace", ns).
			Warn("Ref delete failed")
		return err
	}
	log.WithField("Name", name).
		WithField("Namespace", ns).
		Debug("Ref deleted")
	return nil
}

// Exists checks if data.Object exists in mongo database
func (store *RefMongoRepository) Exists(ctx context.Context, ns cmndata.Namespace, name data.RefName) (bool, error) {
	log := store.log.WithContext(ctx)
	log.WithField("Name", name).
		WithField("Namespace", ns).
		Debug("Looking up ref")
	filter := getOneRefFilter(ns, name)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// refToDTO converts data.Ref to refDTO
func refToDTO(obj data.Ref) refDTO {
	dto := refDTO{
		ID:        primitive.NewObjectID(),
		Name:      string(obj.Name),
		Namespace: string(obj.Namespace),
		Value:     string(obj.Value),
		ObjectID:  string(obj.ObjectID),
	}
	return dto
}

// refDtoToModel converts refDTO to data.Ref
func refDtoToModel(dto refDTO) data.Ref {
	model := data.Ref{
		Namespace: cmndata.Namespace(dto.Namespace),
		Name:      data.RefName(dto.Name),
		Value:     data.Commit(dto.Value),
		ObjectID:  data.ObjectID(dto.ObjectID),
	}
	return model
}

func getOneRefFilter(ns cmndata.Namespace, name data.RefName) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "name", Value: name}},
			bson.D{primitive.E{Key: "namespace", Value: ns}},
		},
	}}
}
