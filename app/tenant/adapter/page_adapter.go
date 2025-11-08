package adapter

import (
	"strings"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *TenantAdapter) PageEntityGrpcFromSql(resp *db.TenantsSchemaPage) *talv1.TenantsSchemaPage {
	return &talv1.TenantsSchemaPage{
		PageId:              int32(resp.PageID),
		PageName:            resp.PageName,
		PageNameAr:          resp.PageNameAr.String,
		PageDescription:     resp.PageDescription.String,
		PageDescriptionAr:   resp.PageDescriptionAr.String,
		PageBreadcrumb:      resp.PageBreadcrumb.String,
		TenantId:            int32(resp.TenantID.Int32), // Handle nullable int
		PageRoute:           resp.PageRoute,
		PageCoverImage:      resp.PageCoverImage.String,
		PageCoverVideo:      resp.PageCoverVideo.String,
		PageKeyWords:        resp.PageKeyWords.String,
		PageMetaDescription: resp.PageMetaDescription.String,
		PageIcon:            resp.PageIcon.String,
		CreatedAt:           db.TimeToProtoTimeStamp(resp.CreatedAt.Time),
		UpdatedAt:           db.TimeToProtoTimeStamp(resp.UpdatedAt.Time),
		DeletedAt:           db.TimeToProtoTimeStamp(resp.DeletedAt.Time),
	}
}

func (a *TenantAdapter) PageEntityListGrpcFromSql(resp *[]db.TenantsSchemaPage) *[]*talv1.TenantsSchemaPage {
	records := make([]*talv1.TenantsSchemaPage, 0)
	for _, v := range *resp {
		record := a.PageEntityGrpcFromSql(&v)
		records = append(records, record)
	}
	return &records
}
func (a *TenantAdapter) PageListGrpcFromSql(resp *[]db.TenantsSchemaPage) *talv1.PageListResponse {
	records := make([]*talv1.TenantsSchemaPage, 0)
	deletedRecords := make([]*talv1.TenantsSchemaPage, 0)
	for _, v := range *resp {
		record := a.PageEntityGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
		} else {
			records = append(records, record)
		}
	}
	return &talv1.PageListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}

func (a *TenantAdapter) PageCreateUpdateSqlFromGrpc(req *talv1.PageCreateUpdateRequest) *db.PageCreateUpdateParams {
	return &db.PageCreateUpdateParams{
		PageID:              req.GetPageId(),
		PageName:            req.GetPageName(),
		PageNameAr:          req.GetPageNameAr(),
		PageDescription:     req.GetPageDescription(),
		PageDescriptionAr:   req.GetPageDescriptionAr(),
		PageBreadcrumb:      req.GetPageBreadcrumb(),
		TenantID:            req.GetTenantId(),
		PageRoute:           req.GetPageRoute(),
		PageCoverImage:      req.GetPageCoverImage(),
		PageCoverVideo:      req.GetPageCoverVideo(),
		PageKeyWords:        strings.Join(req.GetPageKeyWords(), ","),
		PageMetaDescription: req.GetPageMetaDescription(),
		PageIcon:            req.GetPageIcon(),
	}
}

func (a *TenantAdapter) PageFindForUpdateGrpcFromSql(resp *db.TenantsSchemaPage) *talv1.PageFindForUpdateResponse {
	return &talv1.PageFindForUpdateResponse{
		Request: &talv1.PageCreateUpdateRequest{
			PageId:              int32(resp.PageID),
			PageName:            resp.PageName,
			PageNameAr:          resp.PageNameAr.String,
			PageDescription:     resp.PageDescription.String,
			PageDescriptionAr:   resp.PageDescriptionAr.String,
			PageBreadcrumb:      resp.PageBreadcrumb.String,
			TenantId:            int32(resp.TenantID.Int32), // Handle nullable int
			PageRoute:           resp.PageRoute,
			PageCoverImage:      resp.PageCoverImage.String,
			PageCoverVideo:      resp.PageCoverVideo.String,
			PageMetaDescription: resp.PageMetaDescription.String,
			PageIcon:            resp.PageIcon.String,
		},
	}

}
