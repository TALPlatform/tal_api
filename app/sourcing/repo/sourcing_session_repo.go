package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/rs/zerolog/log"
)

func (repo *SourcingRepo) SourcingSessionCreateUpdate(ctx context.Context, req *db.SourcingSessionCreateUpdateParams) (*db.SourcingSessionCreateUpdateRow, error) {
	callerId, _ := contextkeys.CallerID(ctx)
	req.CreatedBy = callerId
	log.Debug().Interface("sss", *req).Msg("msg us")
	resp, err := repo.store.SourcingSessionCreateUpdate(ctx, *req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionFind(ctx context.Context, req int32) (*db.SourcingSessionFindRow, error) {
	resp, err := repo.store.SourcingSessionFind(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionProfileCreateUpdate(ctx context.Context, req *db.SourcingSessionProfileCreateUpdateParams) (*db.SourcingSessionProfileCreateUpdateRow, error) {
	resp, err := repo.store.SourcingSessionProfileCreateUpdate(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionList(ctx context.Context) (*[]db.SourcingSessionListRow, error) {
	tenantID, _ := contextkeys.TenantID(ctx)
	resp, err := repo.store.SourcingSessionList(ctx, tenantID)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionProfileFindForAI(ctx context.Context, req int32) (*db.SourcingSessionProfileFindForAIRow, error) {
	resp, err := repo.store.SourcingSessionProfileFindForAI(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionCriteriaCreate(ctx context.Context, req *db.SourcingSessionCriteriaCreateParams) (*db.SourcingSchemaSourcingSessionCriterium, error) {
	resp, err := repo.store.SourcingSessionCriteriaCreate(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionCriteriaProfilesBulkInsert(
	ctx context.Context,
	req []db.SourcingSessionCriteriaProfilesBulkInsertParams,
) error {
	_, err := repo.store.SourcingSessionCriteriaProfilesBulkInsert(ctx, req)
	if err != nil {
		return repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return nil
}

func (repo *SourcingRepo) SourcingSessionProfileSync(
	ctx context.Context,
	req *db.SourcingSessionProfileSyncParams,
) error {
	err := repo.store.SourcingSessionProfileSync(ctx, *req)
	if err != nil {
		return repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return nil
}

func (repo *SourcingRepo) SourcingSessionApply(ctx context.Context, req int32) (*[]db.SourcingSessionApplyRow, error) {
	resp, err := repo.store.SourcingSessionApply(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
