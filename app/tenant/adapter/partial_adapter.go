package adapter

import (
	"encoding/json"
	"strings"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
)

func (a *TenantAdapter) PartialEntityGrpcFromSql(resp *db.TenantsSchemaPartial) *talv1.TenantsSchemaPartial {
	var partialLinks map[string]string
	if len(resp.PartialLinks) > 0 {
		if err := json.Unmarshal(resp.PartialLinks, &partialLinks); err == nil {
			log.Error().Err(err).Msg("error parsing partial links")
		}
	}
	return &talv1.TenantsSchemaPartial{
		PartialId:            int32(resp.PartialID),
		PartialName:          resp.PartialName,
		PartialNameAr:        resp.PartialNameAr.String,
		PartialTypeId:        resp.PartialTypeID,
		SectionId:            int32(resp.SectionID),
		PartialImage:         resp.PartialImage.String,
		PartialImages:        strings.Split(resp.PartialImages.String, ","),
		PartialVideo:         resp.PartialVideo.String,
		IsFeatured:           resp.IsFeatured.Bool,
		PartialBrief:         resp.PartialBrief.String,
		PartialBriefAr:       resp.PartialBriefAr.String,
		PartialContent:       resp.PartialContent.String,
		PartialContentAr:     resp.PartialContentAr.String,
		PartialButtonLabel:   resp.PartialButtonLabel.String,
		PartialButtonLabelAr: resp.PartialButtonLabelAr.String,
		PartialButtonIcon:    resp.PartialButtonIcon.String,
		PartialButtonLink:    resp.PartialButtonLink.String,

		PartialButtonPageId: int32(resp.PartialButtonPageID.Int32),
		PartialIcons:        strings.Split(resp.PartialIcons.String, ","),
		Address:             resp.Address.String,
		PartialLinks:        partialLinks,
		PartialLink:         resp.PartialLink.String,
		CreatedAt:           db.TimeToProtoTimeStamp(resp.CreatedAt.Time),
		UpdatedAt:           db.TimeToProtoTimeStamp(resp.UpdatedAt.Time),
		DeletedAt:           db.TimeToProtoTimeStamp(resp.DeletedAt.Time),
	}
}

func (a *TenantAdapter) PartialEntityListGrpcFromSql(resp *[]db.TenantsSchemaPartial) *[]*talv1.TenantsSchemaPartial {
	records := make([]*talv1.TenantsSchemaPartial, 0)
	for _, v := range *resp {
		record := a.PartialEntityGrpcFromSql(&v)
		records = append(records, record)
	}
	return &records
}

func (a *TenantAdapter) PartialListGrpcFromSql(resp *[]db.TenantsSchemaPartial) *talv1.PartialListResponse {
	records := make([]*talv1.TenantsSchemaPartial, 0)
	deletedRecords := make([]*talv1.TenantsSchemaPartial, 0)
	for _, v := range *resp {
		record := a.PartialEntityGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
		} else {
			records = append(records, record)
		}
	}
	return &talv1.PartialListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}

func (a *TenantAdapter) PartialCreateUpdateSqlFromGrpc(req *talv1.PartialCreateUpdateRequest) *db.PartialCreateUpdateParams {
	return &db.PartialCreateUpdateParams{
		PartialID:            req.GetPartialId(),
		PartialName:          req.GetPartialName(),
		PartialNameAr:        req.GetPartialNameAr(),
		PartialTypeID:        req.GetPartialTypeId(),
		SectionID:            req.GetSectionId(),
		PartialImage:         req.GetPartialImage(),
		PartialImages:        strings.Join(req.GetPartialImages(), ","),
		PartialVideo:         req.GetPartialVideo(),
		IsFeatured:           req.GetIsFeatured(),
		PartialBrief:         req.GetPartialBrief(),
		PartialBriefAr:       req.GetPartialBriefAr(),
		PartialContent:       req.GetPartialContent(),
		PartialContentAr:     req.GetPartialContentAr(),
		PartialButtonLabel:   req.GetPartialButtonLabel(),
		PartialButtonLabelAr: req.GetPartialButtonLabelAr(),
		PartialButtonIcon:    req.GetPartialButtonIcon(),
		PartialButtonLink:    req.GetPartialButtonLink(),
		PartialButtonPageID:  req.GetPartialButtonPageId(),
		PartialIcons:         req.GetPartialIcons(),
		Address:              req.GetAddress(),
		PartialLink:          req.GetPartialLink(),
	}
}

func (a *TenantAdapter) PartialFindForUpdateGrpcFromSql(resp *db.TenantsSchemaPartial) *talv1.PartialFindForUpdateResponse {
	return &talv1.PartialFindForUpdateResponse{
		Request: &talv1.PartialCreateUpdateRequest{
			PartialId:            int32(resp.PartialID),
			PartialName:          resp.PartialName,
			PartialNameAr:        resp.PartialNameAr.String,
			PartialTypeId:        resp.PartialTypeID,
			SectionId:            int32(resp.SectionID),
			PartialImage:         resp.PartialImage.String,
			PartialImages:        strings.Split(resp.PartialImages.String, ","),
			PartialVideo:         resp.PartialVideo.String,
			IsFeatured:           resp.IsFeatured.Bool,
			PartialBrief:         resp.PartialBrief.String,
			PartialBriefAr:       resp.PartialBriefAr.String,
			PartialContent:       resp.PartialContent.String,
			PartialContentAr:     resp.PartialContentAr.String,
			PartialButtonLabel:   resp.PartialButtonLabel.String,
			PartialButtonLabelAr: resp.PartialButtonLabelAr.String,
			PartialButtonIcon:    resp.PartialButtonIcon.String,
			PartialButtonLink:    resp.PartialButtonLink.String,
			PartialButtonPageId:  int32(resp.PartialButtonPageID.Int32),
			Address:              resp.Address.String,
			PartialLink:          resp.PartialLink.String,
		},
	}

}

func (a *TenantAdapter) PartialTypeListInputGrpcFromSql(resp []db.PartialTypeListInputRow) *talv1.PartialTypeListInputResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Label: v.Label,
		})
	}
	return &talv1.PartialTypeListInputResponse{
		Options: records,
	}

}
