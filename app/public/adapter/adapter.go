package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/weaviateclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	storage_go "github.com/darwishdev/storage-go"
	"github.com/resend/resend-go/v2"
)

type PublicAdapterInterface interface {
	// INJECT INTERFACE

	IconFindSqlFromGrpc(icon *talv1.IconFindRequest) *db.IconFindParams
	IconGrpcFromSql(icon *db.Icon) *talv1.Icon
	IconCreateUpdateBulkSqlFromGrpc(req *talv1.IconCreateUpdateBulkRequest) db.IconCreateUpdateBulkParams
	EmailSendResendFromGrpc(req *talv1.EmailSendRequest) resend.SendEmailRequest
	TranslationCreateUpdateBulkGrpcFromSql(resp []db.TranslationCreateUpdateBulkRow) talv1.TranslationListResponse
	TranslationListGrpcFromSql(resp []db.Translation) talv1.TranslationListResponse
	TranslationFindLocaleGrpcFromSql(resp []db.Translation, locale string) talv1.TranslationFindLocaleResponse
	TranslationGrpcFromSql(resp *db.Translation) *talv1.Translation
	TranslationCreateUpdateBulkSqlFromGrpc(req *talv1.TranslationCreateUpdateBulkRequest) *db.TranslationCreateUpdateBulkParams
	FileDeleteGrpcFromSupa(resp []storage_go.FileUploadResponse) *talv1.FileDeleteResponse
	FileListGrpcFromSupa(resp []storage_go.FileObject, bucketId string) *talv1.FileListResponse
	FileObjectGrpcFromSupa(resp *storage_go.FileObject) *talv1.FileObject
	FileCreateResponseGrpcFromSupa(resp *storage_go.FileUploadResponse) *talv1.FileCreateResponse
	BucketListGrpcFromSupa(resp []storage_go.Bucket) *talv1.BucketListResponse
	StorageBucketGrpcFromSupa(resp *storage_go.Bucket) *talv1.StorageBucket
	SettingUpdateSqlFromGrpc(req *talv1.SettingUpdateRequest) *db.SettingUpdateParams
	SettingEntityGrpcFromSql(resp []db.Setting) []*talv1.Setting
	SettingFindForUpdateGrpcFromSql(resp *[]db.SettingFindForUpdateRow) *talv1.SettingFindForUpdateResponse
	IconListGrpcFromSql(resp []db.Icon) *talv1.IconListResponse

// Converts from gRPC model to Weaviate model:
	CommandPalleteWeaviateFromGrpc(req *talv1.CommandPallete) *weaviateclient.CommandPallete

// Converts from Weaviate model back to gRPC model:
	CommandPalleteGrpcFromWeaviate(doc *weaviateclient.CommandPallete) *talv1.CommandPallete
}

type PublicAdapter struct {
}

func NewPublicAdapter() PublicAdapterInterface {
	return &PublicAdapter{}
}
