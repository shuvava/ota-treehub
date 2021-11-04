package app

import (
	"context"
	"strings"

	"github.com/shuvava/treehub/internal/blobs"
	"github.com/shuvava/treehub/internal/blobs/localfs"
	"github.com/shuvava/treehub/pkg/services"

	"github.com/shuvava/treehub/internal/db"

	intDb "github.com/shuvava/treehub/internal/db/mongo"
)

func (s *Server) initDbService() {
	log := s.log.SetContext("server-init-db")
	if s.svc.Db != nil {
		if err := s.svc.Db.Disconnect(context.Background()); err != nil {
			log.WithError(err).
				Fatal("Error on Db service distracting")
		}
	}
	switch db.Type(strings.ToLower(s.config.Db.Type)) {
	case db.MongoDb:
		mongoDB, err := intDb.NewMongoDB(context.Background(), s.log, s.config.Db.ConnectionString)
		if err != nil {
			log.WithError(err).
				Fatal("Error on Db service creating")
		}
		s.svc.Db = mongoDB
		s.svc.ObjectRepo = intDb.NewObjectMongoRepository(s.log, mongoDB)
		s.svc.RefRepo = intDb.NewRefMongoRepository(s.log, mongoDB)
	default:
		log.WithField("type", s.config.Db.Type).
			Fatal("Unsupported mongoDB type")
	}
}

func (s *Server) initStorage() {
	log := s.log.SetContext("server-init-storage")
	if s.svc.ObjectStore != nil {
		log.Warn("Blob storage subsystem reloading")
	}
	switch blobs.Type(strings.ToLower(s.config.Storage.Type)) {
	case blobs.LocalFs:
		if s.config.Storage.Root == "" {
			s.config.Storage.Root = "/tmp"
			log.Warn("Blob storage root directory will be ", s.config.Storage.Root)
		}
		store, err := localfs.NewLocalFsBlobStore(s.config.Storage.Root, s.log)
		if err != nil {
			log.WithError(err).
				Fatal("Error on Storage service creating")
		}
		s.svc.ObjectStore = store
	default:
		log.WithField("type", s.config.Storage.Type).
			Fatal("Unsupported blob storage type")
	}
}

// create all application services
func (s *Server) initServices() {
	s.initDbService()
	s.initStorage()
	s.svc.Objects = services.NewObjectService(s.log, s.svc.ObjectRepo, s.svc.ObjectStore)
	s.svc.Refs = services.NewRefService(s.log, s.svc.RefRepo)
}
