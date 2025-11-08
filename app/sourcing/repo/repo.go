package repo

import (
	// INJECT IMPORTS
	"context"

	"github.com/TALPlatform/tal_api/db"
)

type SourcingRepoInterface interface {
	// INJECT INTERFACE
	SourcingSessionProfileCreateUpdate(ctx context.Context, req *db.SourcingSessionProfileCreateUpdateParams) (*db.SourcingSessionProfileCreateUpdateRow, error)
	SourcingSessionFind(ctx context.Context, req int32) (*db.SourcingSessionFindRow, error)
	SourcingSessionCreateUpdate(ctx context.Context, req *db.SourcingSessionCreateUpdateParams) (*db.SourcingSessionCreateUpdateRow, error)
	ProjectInputList(ctx context.Context) (*[]db.ProjectInputListRow, error)
}

type SourcingRepo struct {
	store        db.Store
	errorHandler map[string]string
}

func NewSourcingRepo(store db.Store) SourcingRepoInterface {
	errorHandler := map[string]string{
		// INJECT ERROR
	}
	return &SourcingRepo{
		store:        store,
		errorHandler: errorHandler,
	}
}
