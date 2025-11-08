package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *SourcingAdapter) SourcingSessionCreateUpdateSqlFromGrpc(req *talv1.SourcingSessionCreateUpdateRequest) *db.SourcingSessionCreateUpdateParams {
	params := &db.SourcingSessionCreateUpdateParams{
		SourcingSessionID:   req.SourcingSessionId,
		SourcingSessionName: req.SourcingSessionName,
		SourcingProjectID:   req.SourcingProjectId,
	}

	// Convert JSON string to []byte for JSONB
	if req.InitialFilters != "" {
		params.InitialFilters = []byte(req.InitialFilters)
	}

	return params
}

func (a *SourcingAdapter) SourcingSessionCreateUpdateGrpcFromSql(req *db.SourcingSessionCreateUpdateRow) *talv1.SourcingSessionCreateUpdateResponse {
	resp := &talv1.SourcingSessionCreateUpdateResponse{
		SourcingSessionId:   req.SourcingSessionID,
		SourcingSessionName: req.SourcingSessionName,
		SourcingProjectId:   req.SourcingProjectID,
		CreatedBy:           req.CreatedBy,
	}
	return resp
}

func (a *SourcingAdapter) SourcingSessionFindGrpcFromSql(req *db.SourcingSessionFindRow) (*talv1.SourcingSessionFindResponse, error) {
	resp := &talv1.SourcingSessionFindResponse{
		// Session fields
		SourcingSessionId:   req.SourcingSessionID,
		SourcingSessionName: req.SourcingSessionName,
		InitialFilters:      string(req.InitialFilters),
		SessionCreatedAt:    db.PgtimeStampToString(req.SessionCreatedAt),
		SessionUpdatedAt:    db.PgtimeStampToString(req.SessionUpdatedAt),
		SessionDeletedAt:    db.PgtimeStampToString(req.SessionDeletedAt),

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
