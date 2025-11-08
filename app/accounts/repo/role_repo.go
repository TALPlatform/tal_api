package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
)

func (repo *AccountsRepo) RoleFindForUpdate(ctx context.Context, req int32) (*db.RoleFindForUpdateRow, error) {
	resp, err := repo.store.RoleFindForUpdate(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *AccountsRepo) RoleListInput(ctx context.Context) (*[]db.RoleListInputRow, error) {
	callerID, _ := contextkeys.CallerID(ctx)
	resp, err := repo.store.RoleListInput(ctx, callerID)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *AccountsRepo) RoleList(ctx context.Context) (*[]db.RoleListRow, error) {
	resp, err := repo.store.RoleList(ctx)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *AccountsRepo) RoleCreateUpdate(ctx context.Context, req db.RoleCreateUpdateParams) (*db.AccountsSchemaRole, error) {
	callerId, ok := contextkeys.CallerID(ctx)
	if ok {
		req.CallerID = callerId
	}
	resp, err := repo.store.RoleCreateUpdate(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}

func (repo *AccountsRepo) RoleDelete(ctx context.Context, req db.RoleDeleteParams) (*db.AccountsSchemaRole, error) {
	callerId, ok := contextkeys.CallerID(ctx)
	if ok {
		req.CallerID = callerId
	}
	role, err := repo.store.RoleDelete(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)

	}
	return &role, nil
}
func (repo *AccountsRepo) RoleDeleteRestore(ctx context.Context, req db.RoleDeleteRestoreParams) (*db.AccountsSchemaRole, error) {
	callerId, ok := contextkeys.CallerID(ctx)
	if ok {
		req.CallerID = callerId
	}
	role, err := repo.store.RoleDeleteRestore(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)

	}
	return &role, nil
}
