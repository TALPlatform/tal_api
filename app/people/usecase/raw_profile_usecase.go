package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func ConcatSlices[T any](a, b []T) []T {
	return append(a, b...)
}
func (u *PeopleUsecase) RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest]) (*talv1.RawProfileListResponse, error) {
	sqlReq := u.adapter.RawProfileListSqlFromGrpc(req.Msg)
	records, err := u.repo.RawProfileList(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	recordsLength := len(*records)
	if recordsLength == int(req.Msg.Limit) {
		resp := u.adapter.RawProfileListGrpcFromSql(records)
		return resp, nil
	}
	req.Msg.Limit = req.Msg.Limit - int32(recordsLength)
	crustDataRequest := u.adapter.RawProfileListCrustDataFromGrpc(req.Msg)
	crustDataProfiles, err := u.crustDataClient.PeopleSearch(ctx, crustDataRequest)
	if err != nil {
		return nil, err
	}

	crustDataRecords := u.adapter.RawProfileListDbFromCrustData(&crustDataProfiles.Profiles)
	resp := u.adapter.RawProfileListGrpcFromSql(records, crustDataRecords)
	err = u.RawProfileSync(ctx, req.Msg.SourcingSessionId, crustDataProfiles)
	return resp, nil

}

func (u *PeopleUsecase) RawProfileSync(ctx context.Context, sessionId int32, crustDataProfiles *crustdata.PeopleSearchResponse) error {
	syncParams, err := u.adapter.RawProfileListEnrichAndMarshal(ctx, sessionId, crustDataProfiles)
	if err != nil {
		return err
	}
	err = u.repo.RawProfileBulkCreateUpdate(ctx, syncParams)
	if err != nil {
		return err
	}
	return nil
}
