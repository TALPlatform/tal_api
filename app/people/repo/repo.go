package repo

import (
	// INJECT IMPORTS
	"context"

	"github.com/TALPlatform/tal_api/db"
)

type PeopleRepoInterface interface {
	// INJECT INTERFACE
	RawProfileList(ctx context.Context, req *db.RawProfileListParams) (*[]*db.RawProfileListRow, error)
	RawProfileBulkCreateUpdate(ctx context.Context, params *db.RawProfilesBulkCreateUpdateParams) error
}

type PeopleRepo struct {
	store        db.Store
	errorHandler map[string]string
}

func NewPeopleRepo(store db.Store) PeopleRepoInterface {
	errorHandler := map[string]string{
		// INJECT ERROR
	}
	return &PeopleRepo{
		store:        store,
		errorHandler: errorHandler,
	}
}
