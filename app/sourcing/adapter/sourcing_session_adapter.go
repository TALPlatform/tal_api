package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/agenthub"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *SourcingAdapter) SourcingSessionCreateUpdateSqlFromGrpc(req *talv1.SourcingSessionCreateUpdateRequest) *db.SourcingSessionCreateUpdateParams {
	params := &db.SourcingSessionCreateUpdateParams{
		SourcingSessionID:   req.SourcingSessionId,
		SourcingSessionName: req.SourcingSessionName,
		SourcingProjectID:   req.SourcingProjectId,
		MinExperience:       req.MinExperience,
		MaxExperience:       req.MaxExperience,
		RequiredContactInfo: req.RequiredContactInfo,
		Timezone:            req.Timezone,
		Locations:           req.Locations,
		JobTitles:           req.JobTitles,
		JobSeniority:        req.JobSeniority,
		JobFunctions:        req.JobFunctions,
		Companies:           req.Companies,
		CompanyHeadcount:    req.CompanyHeadcount,
		CompanyFunding:      req.CompanyFunding,
		Industries:          req.Industries,
		Keywords:            req.Keywords,
		Skills:              req.Skills,
		EducationLevels:     req.EducationLevels,
		Languages:           req.Languages,
		FilterLimit:         req.FilterLimit,
	}
	return params
}

func (a *SourcingAdapter) SourcingSessionCreateUpdateGrpcFromSql(req *db.SourcingSessionCreateUpdateRow) *talv1.SourcingSessionCreateUpdateResponse {
	resp := &talv1.SourcingSessionCreateUpdateResponse{
		SourcingSessionId:   req.SourcingSessionID,
		SourcingSessionName: req.SourcingSessionName,
		MinExperience:       req.MinExperience,
		MaxExperience:       req.MaxExperience,
		RequiredContactInfo: req.RequiredContactInfo,
		Timezone:            req.Timezone,
		Locations:           req.Locations,
		JobTitles:           req.JobTitles,
		JobSeniority:        req.JobSeniority,
		JobFunctions:        req.JobFunctions,
		Companies:           req.Companies,
		CompanyHeadcount:    req.CompanyHeadcount,
		CompanyFunding:      req.CompanyFunding,
		Industries:          req.Industries,
		Keywords:            req.Keywords,
		Skills:              req.Skills,
		EducationLevels:     req.EducationLevels,
		Languages:           req.Languages,
		SourcingProjectId:   req.SourcingProjectID,
		CreatedAt:           db.PgtimeStampToString(req.CreatedAt),
		UpdatedAt:           db.PgtimeStampToString(req.UpdatedAt),
		DeletedAt:           db.PgtimeStampToString(req.DeletedAt),
	}
	return resp
}

func (a *SourcingAdapter) SourcingSessionFindGrpcFromSql(req *db.SourcingSessionFindRow) (*talv1.SourcingSessionFindResponse, error) {
	resp := &talv1.SourcingSessionFindResponse{
		// Session fields
		SourcingSessionId:   req.SourcingSessionID,
		SourcingSessionName: req.SourcingSessionName,
		// InitialFilters is removed/replaced
		// InitialFilters:      string(req.InitialFilters), // REMOVED

		// --- New Filter Fields Mapping ---
		MinExperience:       req.MinExperience,
		MaxExperience:       req.MaxExperience,
		RequiredContactInfo: req.RequiredContactInfo,
		Timezone:            req.Timezone,
		Locations:           req.Locations,
		JobTitles:           req.JobTitles,
		JobSeniority:        req.JobSeniority,
		JobFunctions:        req.JobFunctions,
		Companies:           req.Companies,
		CompanyHeadcount:    req.CompanyHeadcount,
		CompanyFunding:      req.CompanyFunding,
		Industries:          req.Industries,
		Keywords:            req.Keywords,
		Skills:              req.Skills,
		EducationLevels:     req.EducationLevels,
		Languages:           req.Languages,
		FilterLimit:         req.FilterLimit,

		SessionCreatedAt: db.PgtimeStampToString(req.SessionCreatedAt),
		SessionUpdatedAt: db.PgtimeStampToString(req.SessionUpdatedAt),
		SessionDeletedAt: db.PgtimeStampToString(req.SessionDeletedAt),

		// Project fields
		SourcingProjectId:    req.SourcingProjectID,
		SourcingProjectName:  req.SourcingProjectName,
		SourcingProjectBreif: req.SourcingProjectBreif,
		TenantId:             req.TenantID,
		ProjectCreatedAt:     db.PgtimeStampToString(req.ProjectCreatedAt),
		ProjectUpdatedAt:     db.PgtimeStampToString(req.ProjectUpdatedAt),
		ProjectDeletedAt:     db.PgtimeStampToString(req.ProjectDeletedAt),

		// Creator fields
		CreatedBy:    req.CreatedBy,
		CreatorName:  req.CreatorName,
		CreatorEmail: req.CreatorEmail,
	}

	// Parse profiles JSONB array directly into the struct
	if len(req.Profiles) > 0 {
		if err := json.Unmarshal(req.Profiles, &resp.Profiles); err != nil {
			return nil, fmt.Errorf("failed to unmarshal profiles: %w", err)
		}
	}

	return resp, nil
}
func (a *SourcingAdapter) SourcingSessionProfileCreateUpdateSqlFromGrpc(req *talv1.SourcingSessionProfileCreateUpdateRequest) *db.SourcingSessionProfileCreateUpdateParams {
	return &db.SourcingSessionProfileCreateUpdateParams{
		SourcingSessionID: req.SourcingSessionId,
		RawProfileID:      req.RawProfileId,
		Score:             req.Score,
		OrderIndex:        req.OrderIndex,
		Note:              req.Note,
		ReportSummary:     req.ReportSummary,
		IsShortListed:     req.IsShortListed,
		SummaryBullets:    req.SummaryBullets,
		Justification:     req.Justification,
	}
}

func (a *SourcingAdapter) SourcingSessionProfileCreateUpdateGrpcFromSql(req *db.SourcingSessionProfileCreateUpdateRow) *talv1.SourcingSessionProfileCreateUpdateResponse {
	return &talv1.SourcingSessionProfileCreateUpdateResponse{
		SourcingSessionId: req.SourcingSessionID,
		RawProfileId:      req.RawProfileID,
		Score:             req.Score,
		OrderIndex:        req.OrderIndex,
		Note:              req.Note,
		ReportSummary:     req.ReportSummary,
		IsShortListed:     req.IsShortListed,
		SummaryBullets:    req.SummaryBullets,
		Justification:     req.Justification,
	}
}

func (a *SourcingAdapter) SourcingSessionListGrpcFromSql(req *[]db.SourcingSessionListRow) *talv1.SourcingSessionListResponse {
	rows := make([]*talv1.SessionListRow, 0, len(*req))
	for _, r := range *req {
		rows = append(rows, &talv1.SessionListRow{
			SourcingSessionId:   int32(r.SourcingSessionID),
			SourcingProjectName: r.SourcingProjectName,
			TenantName:          r.TenantName,
			UserName:            r.UserName,
			SourcingSessionName: r.SourcingSessionName,
			MinExperience:       r.MinExperience.Int32,
			MaxExperience:       r.MaxExperience.Int32,
			RequiredContactInfo: r.RequiredContactInfo,
			Timezone:            r.Timezone.String,
			Locations:           r.Locations,
			JobTitles:           r.JobTitles,
			JobSeniority:        r.JobSeniority.String,
			JobFunctions:        r.JobFunctions,
			Companies:           r.Companies,
			CompanyHeadcount:    r.CompanyHeadcount.String,
			CompanyFunding:      r.CompanyFunding.String,
			Industries:          r.Industries,
			Keywords:            r.Keywords,
			Skills:              r.Skills,
			EducationLevels:     r.EducationLevels,
			Languages:           r.Languages,
			FilterLimit:         r.FilterLimit.Int32,
			SourcingProjectId:   int32(r.SourcingProjectID),
			CreatedBy:           int32(r.CreatedBy),
			CreatedAt:           db.PgtimeStampToString(r.CreatedAt),
			UpdatedAt:           db.PgtimeStampToString(r.UpdatedAt),
			DeletedAt:           db.PgtimeStampToString(r.DeletedAt),
		})
	}

	return &talv1.SourcingSessionListResponse{
		Records: rows,
	}
}

// func (a *SourcingAdapter) SourcingSessionProfileFindForAISqlFromGrpc(req *talv1.SourcingSessionProfileJustifyRequest) *db.SourcingSessionProfileFindForAIParams {
// 	return &db.SourcingSessionProfileFindForAIParams{
// 		SourcingSessionID: req.SourcingSessionId,
// 		RawProfileID:      req.RawProfileId,
// 	}
// }

func (a *SourcingAdapter) SourcingSessionCriteriaCreateSqlFromGrpc(req *talv1.SourcingSessionCriteriaCreateRequest) *db.SourcingSessionCriteriaCreateParams {
	return &db.SourcingSessionCriteriaCreateParams{
		SourcingSessionID: req.SourcingSessionId,
		Steps:             req.Steps,
	}
}

func (a *SourcingAdapter) SourcingSessionCriteriaCreateGrpcFromSql(req *db.SourcingSessionCriteriaProfilesBulkInsertParams) *talv1.SourcingSessionCriteriaCreateResponse {
	return &talv1.SourcingSessionCriteriaCreateResponse{
		SourcingSessionCriteriaId: req.SourcingSessionCriteriaID,
		RawProfileId:              req.RawProfileID,
		Score:                     req.Score.Int32, // assuming sql.NullInt32
		OrderIndex:                req.OrderIndex,
	}

}

func (a *SourcingAdapter) SourcingSessionApplyGrpcFromSql(req *db.SourcingSessionApplyRow) *talv1.SourcingSessionApplyResponse {
	if req == nil {
		return nil
	}

	return &talv1.SourcingSessionApplyResponse{
		Profile: &talv1.RawProfileDetails{
			PersonId:       int32(req.PersonID),
			Name:           req.Name,
			Headline:       req.Headline,
			Location:       req.Location,
			CurrentTitle:   req.CurrentTitle,
			CurrentCompany: req.CurrentCompany,
			Industry:       req.Industry,
			Summary:        req.Summary,
			// LinkedinUrl:       req.LinkedinProfileUrl,
			ProfilePictureUrl: req.ProfilePictureUrl,
			// YearsOfExperience: req.YearsOfExperience,
			NumOfConnections: req.NumOfConnections,
			Skills:           req.Skills,
			// Languages:         req.Languages,
		},
		SourcingSessionId: 0,
		RawProfileId:      int32(req.PersonID),
		Score:             0,
		OrderIndex:        0,
		Note:              "",
		ReportSummary:     "",
		IsShortListed:     false,
		SummaryBullets:    []string{},
		Justification:     "",
	}
}

func (a *SourcingAdapter) SourcingSessionApplyCrustDataFromSql(req *db.SourcingSessionFindRow) *crustdata.PeopleSearchRequest {
	if req == nil {
		return &crustdata.PeopleSearchRequest{}
	}

	var conditions []interface{}

	// Helper to add a condition if non-empty
	addCond := func(column, condType string, value interface{}) {
		switch v := value.(type) {
		case string:
			if v != "" {
				conditions = append(conditions, crustdata.FilterCondition{
					Column: column,
					Type:   condType,
					Value:  v,
				})
			}
		case []string:
			if len(v) > 0 {
				conditions = append(conditions, crustdata.FilterCondition{
					Column: column,
					Type:   condType,
					Value:  v,
				})
			}
		case int32, int64, float64, int:
			conditions = append(conditions, crustdata.FilterCondition{
				Column: column,
				Type:   condType,
				Value:  v,
			})
		}
	}

	// 🎯 Map each sourcing session filter
	addCond("years_of_experience_raw", "=>", req.MinExperience)
	addCond("years_of_experience_raw", "=<", req.MaxExperience)

	if req.Timezone != "" && req.Timezone != "any" {
		addCond("region", "=", req.Timezone)
	}

	if len(req.Locations) > 0 {
		addCond("region_address_components", "in", req.Locations)
	}

	if len(req.JobTitles) > 0 {
		addCond("current_employers.title", "in", req.JobTitles)
	}

	if req.JobSeniority != "" {
		addCond("current_employers.seniority_level", "=", req.JobSeniority)
	}

	if len(req.JobFunctions) > 0 {
		addCond("current_employers.function_category", "in", req.JobFunctions)
	}

	// if len(req.Companies) > 0 {
	// 	addCond("current_employers.name", "in", req.Companies)
	// }

	if req.CompanyHeadcount != "" {
		addCond("current_employers.company_headcount_range", "=", req.CompanyHeadcount)
	}

	if req.CompanyFunding != "" {
		addCond("current_employers.company_type", "=", req.CompanyFunding)
	}

	if len(req.Industries) > 0 {
		addCond("current_employers.company_industries", "in", req.Industries)
	}

	if len(req.Skills) > 0 {
		addCond("skills", "in", req.Skills)
	}

	if len(req.Languages) > 0 {
		addCond("languages", "in", req.Languages)
	}

	if len(req.EducationLevels) > 0 {
		addCond("education_background.degree_name", "in", req.EducationLevels)
	}

	// if len(req.Keywords) > 0 {
	// 	addCond("summary", "in", req.Keywords)
	// }

	// 🧠 Build final filter group
	filters := crustdata.FilterGroup{
		Op:         "and",
		Conditions: conditions,
	}

	// 🔧 Construct the request
	searchReq := &crustdata.PeopleSearchRequest{
		Filters: &filters,
		Limit:   int(req.FilterLimit),
		PostProcessing: &crustdata.PostProcessing{
			ExcludeProfiles: []string{},
			ExcludeNames:    []string{},
		},
		Preview: false,
	}

	return searchReq
}

func (a *SourcingAdapter) SourcingSessionApplyGrpcFromCrustdata(row *crustdata.PersonDBProfile) *talv1.SourcingSessionApplyResponse {
	if row == nil {
		return &talv1.SourcingSessionApplyResponse{}
	}
	// Determine current title/company
	var currentTitle, currentCompany string
	if len(row.CurrentEmployers) > 0 {
		currentTitle = row.CurrentEmployers[0].Title
		currentCompany = row.CurrentEmployers[0].Name
	} else if len(row.PastEmployers) > 0 {
		currentTitle = row.PastEmployers[0].Title
		currentCompany = row.PastEmployers[0].Name
	}

	// Map RawProfileDetails
	profile := &talv1.RawProfileDetails{
		PersonId: int32(row.PersonID),
		// FirstName:         row.FirstName,
		// LastName:          row.LastName,
		Name:           row.Name,
		Headline:       row.Headline,
		Location:       row.Region,
		CurrentTitle:   currentTitle,
		CurrentCompany: currentCompany,
		Industry:       "", // optional: could extract from row.CurrentEmployers[0].CompanyIndustries
		Summary:        row.Summary,
		// LinkedinProfile:   "", // map row.Emails or row.TwitterHandle if needed
		// ProfilePicture:    "", // populate if you have a field
		// YearsOfExperience: int32(row.YearsOfExperienceRaw),
		NumOfConnections: int32(row.NumOfConnections),
		Skills:           row.Skills,
		// Languages:         row.Languages,
	}
	resp := &talv1.SourcingSessionApplyResponse{
		RawProfileId: int32(row.PersonID),
		Profile:      profile,
	}
	return resp
}

func (a *SourcingAdapter) SourcingSessionFiltersInferedGrpcFromAgent(req *agenthub.SourcingSessionFiltersInfered) *talv1.SourcingSessionFiltersInferedResponse {
	return &talv1.SourcingSessionFiltersInferedResponse{
		Locations:  req.Locations,
		Skills:     req.Skills,
		Companies:  req.Companies,
		Projects:   req.Projects,
		YearsOfExp: req.YearsOfExp,
	}

}
