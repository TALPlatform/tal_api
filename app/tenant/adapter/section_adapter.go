package adapter

import (
	"strings"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *TenantAdapter) SectionFindForUpdateSqlFromGrpc(req *talv1.SectionFindForUpdateRequest) *db.SectionFindParams {
	return &db.SectionFindParams{
		SectionID: req.RecordId,
	}
}

func (a *TenantAdapter) SectionFindForUpdateGrpcFromSql(resp *db.TenantsSchemaSection) *talv1.SectionFindForUpdateResponse {
	return &talv1.SectionFindForUpdateResponse{
		Request: &talv1.SectionCreateUpdateRequest{
			SectionId:            int32(resp.SectionID),
			SectionName:          resp.SectionName,
			SectionNameAr:        resp.SectionNameAr.String,
			SectionHeader:        resp.SectionHeader.String,
			SectionHeaderAr:      resp.SectionHeaderAr.String,
			SectionButtonLabel:   resp.SectionButtonLabel.String,
			SectionButtonLabelAr: resp.SectionButtonLabelAr.String,
			SectionButtonPageId:  int32(resp.SectionButtonPageID.Int32),
			SectionDescription:   resp.SectionDescription.String,
			SectionDescriptionAr: resp.SectionDescriptionAr.String,
			TenantId:             int32(resp.TenantID.Int32),
			SectionBackground:    resp.SectionBackground.String,
			SectionImages:        strings.Split(resp.SectionImages.String, ","),
			SectionIcon:          resp.SectionIcon.String,
		},
	}
}

func (a *TenantAdapter) SectionEntityGrpcFromSql(resp *db.TenantsSchemaSection) *talv1.TenantsSchemaSection {
	return &talv1.TenantsSchemaSection{
		SectionId:            int32(resp.SectionID),
		SectionName:          resp.SectionName,
		SectionNameAr:        resp.SectionNameAr.String,
		SectionHeader:        resp.SectionHeader.String,
		SectionHeaderAr:      resp.SectionHeaderAr.String,
		SectionButtonLabel:   resp.SectionButtonLabel.String,
		SectionButtonLabelAr: resp.SectionButtonLabelAr.String,
		SectionButtonPageId:  int32(resp.SectionButtonPageID.Int32),
		SectionDescription:   resp.SectionDescription.String,
		SectionDescriptionAr: resp.SectionDescriptionAr.String,
		TenantId:             int32(resp.TenantID.Int32),
		SectionBackground:    resp.SectionBackground.String,
		SectionImages:        resp.SectionImages.String,
		SectionImagesArray:   strings.Split(resp.SectionImages.String, ","),
		SectionIcon:          resp.SectionIcon.String,
		CreatedAt:            db.TimeToProtoTimeStamp(resp.CreatedAt.Time),
		UpdatedAt:            db.TimeToProtoTimeStamp(resp.UpdatedAt.Time),
		DeletedAt:            db.TimeToProtoTimeStamp(resp.DeletedAt.Time),
	}
}

func (a *TenantAdapter) SectionEntityListGrpcFromSql(resp *[]db.TenantsSchemaSection) *[]*talv1.TenantsSchemaSection {
	records := make([]*talv1.TenantsSchemaSection, 0)
	for _, v := range *resp {
		record := a.SectionEntityGrpcFromSql(&v)
		records = append(records, record)
	}
	return &records
}
func (a *TenantAdapter) SectionListGrpcFromSql(resp *[]db.TenantsSchemaSection) *talv1.SectionListResponse {
	records := make([]*talv1.TenantsSchemaSection, 0)
	deletedRecords := make([]*talv1.TenantsSchemaSection, 0)
	for _, v := range *resp {
		record := a.SectionEntityGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
		} else {
			records = append(records, record)
		}
	}
	return &talv1.SectionListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}

func (a *TenantAdapter) SectionCreateUpdateSqlFromGrpc(req *talv1.SectionCreateUpdateRequest) *db.SectionCreateUpdateParams {
	return &db.SectionCreateUpdateParams{
		SectionID:            req.GetSectionId(),
		SectionName:          req.GetSectionName(),
		SectionNameAr:        req.GetSectionNameAr(),
		SectionHeader:        req.GetSectionHeader(),
		SectionHeaderAr:      req.GetSectionHeaderAr(),
		SectionButtonLabel:   req.GetSectionButtonLabel(),
		SectionButtonLabelAr: req.GetSectionButtonLabelAr(),
		SectionButtonPageID:  req.GetSectionButtonPageId(),
		SectionDescription:   req.GetSectionDescription(),
		SectionDescriptionAr: req.GetSectionDescriptionAr(),
		TenantID:             req.GetTenantId(),
		SectionBackground:    req.GetSectionBackground(),
		SectionImages:        strings.Join(req.GetSectionImages(), ","),
		SectionIcon:          req.GetSectionIcon(),
	}
}

func (a *TenantAdapter) SectionListInptGrpcFromSql(resp *[]db.SectionListInptRow) *talv1.SectionListInptResponse {

	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Label: v.Label,
		})
	}
	return &talv1.SectionListInptResponse{
		Options: records,
	}

}
