package services

import (
	"context"
	"fmt"

	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/treehub/internal/db"
	"github.com/shuvava/treehub/pkg/data"
)

// RefService is service for interaction with data.Ref
type RefService struct {
	log logger.Logger
	db  db.RefRepository
}

// NewRefService creates new instance of ObjectService
func NewRefService(l logger.Logger, db db.RefRepository) *RefService {
	log := l.SetContext("ref-service")
	return &RefService{
		log: log,
		db:  db,
	}
}

// StoreRef persists data.Ref to database
func (svc *RefService) StoreRef(ctx context.Context, ns data.Namespace, name data.RefName, commit data.Commit, force bool) error {
	log := svc.log.WithContext(ctx)
	ref, err := data.NewRef(ns, name, commit)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDataRefValidation,
			"Ref is invalid", err)
	}
	exists, err := svc.db.Exists(ctx, ref.Namespace, ref.Name)
	if err != nil {
		return err
	}
	if exists && !force {
		err = fmt.Errorf("ref already exists")
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorSvcEntityExists,
			"Ref already exists and force push header not set", err)
	}
	if !exists {
		err = svc.db.Create(ctx, ref)
	} else {
		err = svc.db.Update(ctx, ref)
	}

	return err
}

// GetRef returns data.Ref from database
func (svc *RefService) GetRef(ctx context.Context, ns data.Namespace, name data.RefName) (*data.Ref, error) {
	return svc.db.Find(ctx, ns, name)
}

// Exists checks if data.Ref exist on storage
func (svc *RefService) Exists(ctx context.Context, ns data.Namespace, name data.RefName) (bool, error) {
	return svc.db.Exists(ctx, ns, name)
}
