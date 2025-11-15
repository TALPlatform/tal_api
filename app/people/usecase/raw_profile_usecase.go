package usecase

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
)

func ConcatSlices[T any](a, b []T) []T {
	return append(a, b...)

}

func (u *PeopleUsecase) RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest], stream *connect.ServerStream[talv1.RawProfileListResponse]) error {
	// sqlReq := u.adapter.RawProfileListSqlFromGrpc(req.Msg)
	// records, err := u.repo.RawProfileList(ctx, sqlReq)
	// if err != nil {
	// 	return err
	// }
	//
	// rawProfilIds := make([]int32, len(*records))
	// for i, r := range *records {
	// 	grpcRow := u.adapter.RawProfileListGrpcFromSql(r)
	// 	if err := stream.Send(grpcRow); err != nil {
	// 		return err
	// 	}
	// 	rawProfilIds[i] = int32(r.PersonID)
	// }
	// err = u.repo.SourcingSessionProfileSync(ctx, &db.SourcingSessionProfileSyncParams{SourcingSessionID: req.Msg.SourcingSessionId, RawProfileIds: rawProfilIds})
	// if err != nil {
	// 	return fmt.Errorf("error streaming db records", err)
	// }
	//
	// recordsLength := len(*records)
	// if recordsLength == int(req.Msg.Limit) {
	// 	return nil
	// }
	// req.Msg.Limit = req.Msg.Limit - int32(recordsLength)
	crustDataRequest := u.adapter.RawProfileListCrustDataFromGrpc(req.Msg)
	crustDataProfilesRaw, err := u.crustDataClient.PeopleSearch(ctx, crustDataRequest)
	if err != nil {
		return err
	}

	crustDataProfiles, err := u.crustDataClient.PeopleSearchParse(crustDataProfilesRaw)
	if err != nil {
		return err
	}
	for _, r := range crustDataProfiles.Profiles {
		log.Debug().Interface("row is", r).Msg("raw profile sync")
		grpcRow := u.adapter.RawProfileListGrpcFromCrustdata(r)
		if err := stream.Send(grpcRow); err != nil {
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("error streaming db records", err)
	}
	// resp := u.adapter.RawProfileListGrpcFromSql(records, crustDataRecords)
	// err = u.SyncProfilesAsync(ctx, req.Msg.SourcingSessionId, crustDataProfiles)
	return nil

}

func (u *PeopleUsecase) SyncProfilesAsync(ctx context.Context, sessionID int32, data *crustdata.PeopleSearchResponse) error {
	payload := map[string]interface{}{
		"session_id":          sessionID,
		"crust_data_profiles": data,
	}
	info, err := u.asynqClient.Enqueue(ctx, "people:raw_profile_sync", payload)
	log.Debug().Interface("task info is", info).Msg("asynq task")
	return err
}
func (u *PeopleUsecase) RawProfileSync(ctx context.Context, sessionId int32, crustDataProfiles *crustdata.PeopleSearchResponse) error {
	log.Debug().Interface("task is fired", sessionId).Msg("asynq task")
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

func (u *PeopleUsecase) RawProfileFind(ctx context.Context, req *connect.Request[talv1.RawProfileFindRequest]) (*talv1.RawProfileFindResponse, error) {
	record, err := u.repo.RawProfileFind(ctx, req.Msg.PersonId)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.RawProfileFindGrpcFromSql(record)
	return resp, nil
}
func (u *PeopleUsecase) RawProfileListRequestBuild(
	ctx context.Context,
	req *connect.Request[talv1.RawProfileListRequestBuildRequest],
) (*talv1.RawProfileListRequestBuildResponse, error) {
	// if u.rawProfileAgent == nil {
	// 	return nil, fmt.Errorf("structured agent for raw profile is not initialized")
	// }
	// structuredResult, err := u.rawProfileAgent.GenerateContent(ctx, req.Msg.Text)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate raw profile structured content: %w", err)
	// }
	// resp := &talv1.RawProfileListRequestBuildResponse{
	// 	StructuredResponse: &talv1.RawProfileListRequest{
	// 		Query: req.Msg.Text,
	// 		// Industries: structuredResult.Industries,
	// 		YearsOfExperienceFrom: structuredResult.YearsOfExperienceFrom,
	// 		YearsOfExperienceTo:   structuredResult.YearsOfExperienceTo,
	// 		Locations:             structuredResult.Locations,
	// 		Skills:                structuredResult.Skills,
	// 		Companies:             structuredResult.Companies,
	// 		Projects:              structuredResult.Projects,
	// 	},
	// }
	// return resp, nil
	return nil, nil
}
