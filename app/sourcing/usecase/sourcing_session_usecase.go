package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/agenthub"

	// "github.com/TALPlatform/tal_api/pkg/llm"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	// "github.com/darwishdev/genaiclient/pkg/genaiconfig"
	"github.com/rs/zerolog/log"
)

func (u *SourcingUsecase) SourcingSessionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCreateUpdateRequest]) (*talv1.SourcingSessionCreateUpdateResponse, error) {
	sqlReq := u.adapter.SourcingSessionCreateUpdateSqlFromGrpc(req.Msg)
	record, err := u.repo.SourcingSessionCreateUpdate(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.SourcingSessionCreateUpdateGrpcFromSql(record)
	return resp, nil
}

func (u *SourcingUsecase) SourcingSessionFind(ctx context.Context, req *connect.Request[talv1.SourcingSessionFindRequest]) (*talv1.SourcingSessionFindResponse, error) {
	record, err := u.repo.SourcingSessionFind(ctx, req.Msg.SourcingSessionId)
	if err != nil {
		return nil, err
	}
	resp, err := u.adapter.SourcingSessionFindGrpcFromSql(record)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *SourcingUsecase) SourcingSessionProfileCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileCreateUpdateRequest]) (*talv1.SourcingSessionProfileCreateUpdateResponse, error) {
	sqlReq := u.adapter.SourcingSessionProfileCreateUpdateSqlFromGrpc(req.Msg)
	record, err := u.repo.SourcingSessionProfileCreateUpdate(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.SourcingSessionProfileCreateUpdateGrpcFromSql(record)
	return resp, nil
}

func (u *SourcingUsecase) SourcingSessionList(ctx context.Context, req *connect.Request[talv1.SourcingSessionListRequest]) (*talv1.SourcingSessionListResponse, error) {
	record, err := u.repo.SourcingSessionList(ctx)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.SourcingSessionListGrpcFromSql(record)
	return resp, nil

}
func structToMap(profile interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(profile)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported field, skip
			continue
		}
		value := v.Field(i).Interface()
		out[field.Name] = value
	}
	return out
}
func (u *SourcingUsecase) SourcingSessionProfileJustify(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileJustifyRequest], stream *connect.ServerStream[talv1.SourcingSessionProfileJustifyResponse]) error {
	sessionContext, err := u.repo.SourcingSessionProfileFindForAI(ctx, req.Msg.SourcingSessionId)
	if err != nil {
		return err
	}
	profileContext, err := u.rawProfileFinder(ctx, connect.NewRequest(&talv1.RawProfileFindRequest{PersonId: req.Msg.RawProfileId}))
	if err != nil {
		return err
	}
	if u.sessionProfileJustiferAgent == nil {
		return fmt.Errorf("sourcing session filters build agent is not registerd please register agents first")
	}
	session, err := u.sessionProfileJustiferAgent.NewRedisSession(ctx, "user", "day", u.agentsStore)
	if err != nil {
		return fmt.Errorf("failed to create agent session: %w", err)
	}

	// 4️⃣ Build sessionInfo map from SourcingSessionProfileFindForAIRow
	sessionInfo := map[string]any{
		"sourcing_session_name":  sessionContext.SourcingSessionName,
		"min_experience":         sessionContext.SMinExperience,
		"max_experience":         sessionContext.SMaxExperience,
		"required_contact_info":  sessionContext.SRequiredContactInfo,
		"timezone":               sessionContext.STimezone,
		"locations":              sessionContext.SLocations,
		"job_titles":             sessionContext.SJobTitles,
		"job_seniority":          sessionContext.SJobSeniority,
		"job_functions":          sessionContext.SJobFunctions,
		"companies":              sessionContext.SCompanies,
		"company_headcount":      sessionContext.SCompanyHeadcount,
		"company_funding":        sessionContext.SCompanyFunding,
		"industries":             sessionContext.SIndustries,
		"keywords":               sessionContext.SKeywords,
		"skills":                 sessionContext.SSkills,
		"education_levels":       sessionContext.SEducationLevels,
		"languages":              sessionContext.SLanguages,
		"filter_limit":           sessionContext.SFilterLimit,
		"sourcing_project_name":  sessionContext.SourcingProjectName,
		"sourcing_project_breif": sessionContext.SourcingProjectBreif.String,
	}

	// 5️⃣ Build rawProfile map from RawProfileFindResponse
	rawProfile := map[string]any{
		"person_id":                 profileContext.PersonId,
		"name":                      profileContext.Name,
		"first_name":                profileContext.FirstName,
		"last_name":                 profileContext.LastName,
		"region":                    profileContext.Region,
		"region_address_components": profileContext.RegionAddressComponents,
		"headline":                  profileContext.Headline,
		"summary":                   profileContext.Summary,
		"skills":                    profileContext.Skills,
		"languages":                 profileContext.Languages,
		"profile_language":          profileContext.ProfileLanguage,
		"emails":                    profileContext.Emails,
		"twitter_handle":            profileContext.TwitterHandle,
		"open_to_cards":             profileContext.OpenToCards,
		"num_of_connections":        profileContext.NumOfConnections,
		"recently_changed_jobs":     profileContext.RecentlyChangedJobs,
		"years_of_experience":       profileContext.YearsOfExperience,
		"years_of_experience_raw":   profileContext.YearsOfExperienceRaw,
		"current_employers":         profileContext.CurrentEmployers,
		"past_employers":            profileContext.PastEmployers,
		"education_background":      profileContext.EducationBackground,
		"honors":                    profileContext.Honors,
		"certifications":            profileContext.Certifications,
	}

	resp, err := session.Send(ctx, struct {
		SessionInfo map[string]any
		RawProfile  map[string]any
	}{SessionInfo: sessionInfo, RawProfile: rawProfile})
	if err != nil {
		return fmt.Errorf("error sending message : %w", err)
	}
	stream.Send(&talv1.SourcingSessionProfileJustifyResponse{Value: resp})
	return nil
}

func (u *SourcingUsecase) SourcingSessionCriteriaCreate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCriteriaCreateRequest], stream *connect.ServerStream[talv1.SourcingSessionCriteriaCreateResponse]) error {
	sqlReq := u.adapter.SourcingSessionCriteriaCreateSqlFromGrpc(req.Msg)
	_, err := u.repo.SourcingSessionCriteriaCreate(ctx, sqlReq)
	if err != nil {
		return err
	}
	// resp := u.adapter.SourcingSessionCriteriaCreateGrpcFromSql(record)
	return nil

}
func (u *SourcingUsecase) SourcingSessionApply(ctx context.Context, req *connect.Request[talv1.SourcingSessionApplyRequest], stream *connect.ServerStream[talv1.SourcingSessionApplyResponse]) error {
	// sqlReq := u.adapter.SourcingSessionApplySqlFromGrpc(req.Msg)
	records, err := u.repo.SourcingSessionApply(ctx, req.Msg.SourcingSessionId)
	if err != nil {
		return fmt.Errorf("error getting the data from db: %w", err)
	}
	rawProfilIds := make([]int32, len(*records))
	for i, r := range *records {
		grpcRow := u.adapter.SourcingSessionApplyGrpcFromSql(&r)
		if err := stream.Send(grpcRow); err != nil {
			return fmt.Errorf("error sending the stream from the ai: %w", err)
		}
		rawProfilIds[i] = int32(r.PersonID)
	}
	// recordsLength := len(*records)
	// if recordsLength == int(req.Msg.Limit) {
	// 	return nil
	// }
	if len(*records) == 50 {
		return nil
	}
	session, err := u.repo.SourcingSessionFind(ctx, req.Msg.SourcingSessionId)
	if err != nil {
		return fmt.Errorf("error getting profile info: %w", err)
	}
	crustDataRequest := u.adapter.SourcingSessionApplyCrustDataFromSql(session)

	crustDataProfilesRaw, err := u.crustDataClient.PeopleSearch(ctx, crustDataRequest)
	if err != nil {
		return fmt.Errorf("error fetching profiles from crustdata : %w", err)
	}
	crustDataProfiles, err := u.crustDataClient.PeopleSearchParse(crustDataProfilesRaw)
	if err != nil {
		return fmt.Errorf("error fetching profiles from crustdata : %w", err)
	}
	for _, r := range crustDataProfiles.Profiles {
		grpcRow := u.adapter.SourcingSessionApplyGrpcFromCrustdata(r)
		if err := stream.Send(grpcRow); err != nil {
			return err
		}

		grpcRow.SourcingSessionId = req.Msg.SourcingSessionId
	}
	dbProfiles, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("error marshalling db records into json : %w", err)
	}
	payload := map[string]interface{}{
		"session_id":          req.Msg.SourcingSessionId,
		"crust_data_profiles": crustDataProfilesRaw,
		"db_profiles":         dbProfiles,
	}
	_, err = u.asynqClient.Enqueue(ctx, "sourcing:sourcig_session_profile_sync", payload)
	return err
}

func (s *SourcingUsecase) SyncSourcingSessionProfiles(
	ctx context.Context,
	sessionID int32,
	crustDataProfiles *[]byte,
	DBProfiles *[]byte,
) error {
	err := s.repo.SourcingSessionProfileSync(ctx, &db.SourcingSessionProfileSyncParams{
		SessionID:     sessionID,
		CrustProfiles: *crustDataProfiles,
		DbProfiles:    *DBProfiles,
	})
	log.Info().
		Int32("session_id", sessionID).
		Int("db_profiles", len(*DBProfiles)).
		Int("crust_data_profiles", len(*crustDataProfiles)).
		Msg("Finished syncing sourcing session profiles")
	return err
}

func (u *SourcingUsecase) SourcingSessionFiltersBuilder(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersBuilderRequest]) (*talv1.SourcingSessionFiltersBuilderResponse, error) {
	if u.sourcingSessionFiltersBuildAgent == nil {
		return nil, fmt.Errorf("sourcing session filters build agent is not registerd please register agents first ")
	}
	session, err := u.sourcingSessionFiltersBuildAgent.NewRedisSession(ctx, "user", "day", u.agentsStore)
	if err != nil {
		return nil, err
	}
	resp, err := session.Send(ctx, struct {
		Prompt  string
		Infered agenthub.SourcingSessionFiltersInfered
	}{Prompt: req.Msg.Prompt})

	return &talv1.SourcingSessionFiltersBuilderResponse{
		Locations:     resp.Locations,
		Skills:        resp.Skills,
		Companies:     resp.Companies,
		MinExperience: resp.YearsOfExpMin,
		MaxExperience: resp.YearsOfExpMax,
	}, nil
	// extractedFilters, err := u.filtersAgent.GenerateContent(ctx, req.Msg.Prompt)
	// if err != nil {
	// 	return nil, err
	// }
	// r, err := u.repo.SourcingSessionCreateUpdate(ctx, &db.SourcingSessionCreateUpdateParams{
	// 	Locations:         extractedFilters.Locations,
	// 	Skills:            extractedFilters.Skills,
	// 	Companies:         extractedFilters.Companies,
	// 	SourcingSessionID: req.Msg.SourcingSessionId,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return &talv1.SourcingSessionFiltersBuilderResponse{
	// 	SourcingSessionId:   req.Msg.SourcingSessionId,
	// 	SourcingSessionName: r.SourcingSessionName,
	// 	MinExperience:       r.MinExperience,
	// 	MaxExperience:       r.MaxExperience,
	// 	RequiredContactInfo: r.RequiredContactInfo,
	// 	Locations:           r.Locations,
	// 	Companies:           r.Companies,
	// 	Keywords:            r.Keywords,
	// 	CompanyHeadcount:    r.CompanyHeadcount,
	// 	EducationLevels:     r.EducationLevels,
	// 	SessionCreatedAt:    db.PgtimeStampToString(r.CreatedAt), // You'll need to define u.formatTimestamp
	// 	SessionUpdatedAt:    db.PgtimeStampToString(r.UpdatedAt),
	// 	SessionDeletedAt:    db.PgtimeStampToString(r.DeletedAt),
	// 	SourcingProjectId:   r.SourcingProjectID,
	// 	CreatedBy:           r.CreatedBy,
	// 	Languages:           r.Languages,
	// 	FilterLimit:         r.FilterLimit,
	// 	CompanyFunding:      r.CompanyFunding,
	// 	Industries:          r.Industries,
	// 	JobTitles:           r.JobTitles,
	// 	JobSeniority:        r.JobSeniority,
	// 	JobFunctions:        r.JobFunctions,
	// 	Timezone:            r.Timezone,
	// 	Skills:              r.Skills,
	// }, nil
}

func (u *SourcingUsecase) SourcingSessionFiltersInfered(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersInferedRequest]) (*talv1.SourcingSessionFiltersInferedResponse, error) {
	if u.sourcingSessionFiltersGuessAgent == nil {
		return nil, fmt.Errorf("sourcing session filters guess agent not registerd please register agents first")
	}
	session, err := u.sourcingSessionFiltersGuessAgent.NewRedisSession(ctx, "user", "day", u.agentsStore)
	if err != nil {
		return nil, err
	}
	agentResp, err := session.Send(ctx, req.Msg.Prompt)
	resp := u.adapter.SourcingSessionFiltersInferedGrpcFromAgent(&agentResp)
	return resp, nil
}
