package adapter

import (
	"encoding/json"
	"time"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
)

func (a *TenantAdapter) TenantDashboardGrpcFromSql(resp *[]db.TenantDashboardRow) *talv1.TenantDashboardResponse {
	partialByType := make([]*talv1.PartialByType, len(*resp))
	counts := &talv1.TenantDashboardCounts{}
	for i, r := range *resp {
		if i == 0 {
			counts.RoleCount = &talv1.CountCard{
				Link:  "/accounts/role",
				Count: int32(r.RoleCount),
			}
			counts.UserCount = &talv1.CountCard{
				Link:  "/accounts/user",
				Count: int32(r.UserCount),
			}
			counts.SectionCount = &talv1.CountCard{
				Link:  "/tenants/section",
				Count: int32(r.SectionCount),
			}
			counts.TenantCount = &talv1.CountCard{
				Link:  "/tenants/tenant",
				Count: int32(r.TenantCount),
			}
		}
		partialByType[i] = &talv1.PartialByType{PartialTypeId: r.PartialTypeID, PartialTypeName: r.PartialTypeName, PartialCount: int32(r.PartialTypeCount)}
	}
	return &talv1.TenantDashboardResponse{
		Counts:        counts,
		PartialByType: partialByType,
	}
}
func (a *TenantAdapter) TenantEntityGrpcFromSql(resp *db.TenantsSchemaTenant) *talv1.TenantsSchemaTenant {
	var tenantLinks map[string]string
	log.Debug().Interface("tenant links is", resp.TenantLinks).Msg("links is")
	if len(resp.TenantLinks) > 0 {
		if err := json.Unmarshal(resp.TenantLinks, &tenantLinks); err == nil {
			log.Error().Err(err).Msg("error parsing partial links")
		}
	}

	return &talv1.TenantsSchemaTenant{
		TenantId:               int32(resp.TenantID),
		TenantName:             resp.TenantName,
		TenantNameAr:           resp.TenantNameAr.String,
		TenantPhone:            resp.TenantPhone.String,
		TenantAddress:          resp.TenantAddress.String,
		TenantAddressAr:        resp.TenantAddressAr.String,
		TenantDescription:      resp.TenantDescription.String,
		TenantDescriptionAr:    resp.TenantDescriptionAr.String,
		TenantEmail:            resp.TenantEmail.String,
		TenantLinks:            tenantLinks,
		TenantLogo:             resp.TenantLogo.String,
		TenantValues:           resp.TenantValues.String,
		TenantVision:           resp.TenantVision.String,
		TenantMission:          resp.TenantMission.String,
		TenantLogoVertical:     resp.TenantLogoVertical.String,
		TenantLogoDark:         resp.TenantLogoDark.String,
		TenantLogoDarkVertical: resp.TenantLogoDarkVertical.String,
		CreatedAt:              db.TimeToString(resp.CreatedAt.Time),
		DeletedAt:              db.TimeToString(resp.DeletedAt.Time),
	}
}

func (a *TenantAdapter) TenantCreateUpdateSqlFromGrpc(req *talv1.TenantCreateUpdateRequest) *db.TenantCreateUpdateParams {
	tenantLinks, err := json.Marshal(req.GetTenantLinks())
	if err != nil {
		log.Error().Err(err).Msg("error parsing partial links")
	}
	log.Debug().Interface("tenat log", req.TenantLogo).Msg("show me logo")
	resp := &db.TenantCreateUpdateParams{
		TenantID:               req.TenantId,
		TenantName:             req.TenantName,
		TenantNameAr:           req.TenantNameAr,
		TenantLinks:            tenantLinks,
		TenantPhone:            req.TenantPhone,
		TenantAddress:          req.TenantAddress,
		TenantAddressAr:        req.TenantAddressAr,
		TenantValues:           req.TenantValues,
		TenantVision:           req.TenantVision,
		TenantMission:          req.TenantMission,
		TenantEmail:            req.TenantEmail,
		TenantDescription:      req.TenantDescription,
		TenantLogo:             req.TenantLogo,
		TenantDescriptionAr:    req.TenantDescriptionAr,
		TenantLogoVertical:     req.TenantLogoVertical,
		TenantLogoDark:         req.TenantLogoDark,
		TenantLogoDarkVertical: req.TenantLogoDarkVertical,
	}
	return resp
}
func (a *TenantAdapter) TenantListGrpcFromSql(resp *[]db.TenantsSchemaTenant) *talv1.TenantListResponse {
	records := make([]*talv1.TenantsSchemaTenant, 0)
	deletedRecords := make([]*talv1.TenantsSchemaTenant, 0)
	for _, v := range *resp {
		record := a.TenantEntityGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
		} else {
			records = append(records, record)
		}
	}
	return &talv1.TenantListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}
func (a *TenantAdapter) TenantFindGrpcFromSql(resp *db.TenantFindRow) *talv1.TenantFindResponse {

	var navigations []*talv1.NavigationBar
	var tenantLinks map[string]string
	if len(resp.TenantLinks) > 0 {
		if err := json.Unmarshal(resp.TenantLinks, &tenantLinks); err == nil {
			log.Error().Err(err).Msg("error parsing partial links")
		}
	}

	log.Debug().Interface("debug", tenantLinks).Msg("navigations are")

	if len(resp.Navigations) > 0 {
		err := json.Unmarshal(resp.Navigations, &navigations)
		log.Error().Err(err).Msg("error parsing the navigations")
	}
	var pages []*talv1.TenantsSchemaPageView
	if len(resp.Pages) > 0 {
		err := json.Unmarshal(resp.Pages, &pages)
		log.Error().Err(err).Msg("error parsing the pages")
	}
	return &talv1.TenantFindResponse{
		Tenant: &talv1.TenantsSchemaTenantView{
			TenantId:               resp.TenantID,
			TenantName:             resp.TenantName,
			TenantNameAr:           resp.TenantNameAr.String,
			TenantPhone:            resp.TenantPhone.String,
			TenantAddress:          resp.TenantAddress.String,
			TenantAddressAr:        resp.TenantAddressAr.String,
			TenantDescription:      resp.TenantDescription.String,
			TenantDescriptionAr:    resp.TenantDescriptionAr.String,
			TenantLinks:            tenantLinks,
			TenantEmail:            resp.TenantEmail.String,
			TenantLogo:             resp.TenantLogo.String,
			TenantLogoVertical:     resp.TenantLogoVertical.String,
			TenantLogoDark:         resp.TenantLogoDark.String,
			TenantLogoDarkVertical: resp.TenantLogoDarkVertical.String,
			TenantValues:           resp.TenantValues.String,
			TenantVision:           resp.TenantVision.String,
			TenantMission:          resp.TenantMission.String,
			CreatedAt:              resp.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:              resp.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt:              resp.DeletedAt.Time.Format(time.RFC3339),
			Pages:                  pages,
			Navigations:            navigations,
		},
	}
}
func (a *TenantAdapter) TenantDeleteRestoreGrpcFromSql(resp *[]db.TenantsSchemaTenant) *talv1.TenantDeleteRestoreResponse {
	records := make([]*talv1.TenantsSchemaTenant, 0)
	for _, v := range *resp {
		record := a.TenantEntityGrpcFromSql(&v)
		records = append(records, record)
	}
	return &talv1.TenantDeleteRestoreResponse{
		Records: records,
	}

}

func (a *TenantAdapter) TenantListInputGrpcFromSql(resp *[]db.TenantListInputRow) *talv1.TenantListInputResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Label: v.Label,
		})
	}
	return &talv1.TenantListInputResponse{
		Options: records,
	}

}
