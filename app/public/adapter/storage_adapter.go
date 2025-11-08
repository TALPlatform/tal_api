package adapter

import (
	"fmt"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	storage_go "github.com/darwishdev/storage-go"
	"github.com/rs/zerolog/log"
)

func (a *PublicAdapter) StorageBucketGrpcFromSupa(resp *storage_go.Bucket) *talv1.StorageBucket {
	return &talv1.StorageBucket{
		Name:      resp.Name,
		CreatedAt: resp.CreatedAt,
		Id:        resp.Id,
		Public:    resp.Public,
	}
}

func (a *PublicAdapter) FileCreateResponseGrpcFromSupa(resp *storage_go.FileUploadResponse) *talv1.FileCreateResponse {
	return &talv1.FileCreateResponse{
		Path: resp.Key,
	}
}

func convertMetadata(meta interface{}) *talv1.FileMetadata {
	if meta == nil {
		return nil
	}

	metaMap, ok := meta.(map[string]interface{})
	if !ok {
		return nil
	}

	return &talv1.FileMetadata{
		ETag:           db.StringFindFromMap(metaMap, "eTag"),
		Mimetype:       db.StringFindFromMap(metaMap, "mimetype"),
		CacheControl:   db.StringFindFromMap(metaMap, "cacheControl"),
		LastModified:   db.TimestampFindFromMap(metaMap, "lastModified"),
		HttpStatusCode: db.Int32FindFromMap(metaMap, "httpStatusCode"),
		Size:           db.Int32FindFromMap(metaMap, "size"),
		ContentLength:  db.Int32FindFromMap(metaMap, "contentLength"),
	}
}
func (a *PublicAdapter) FileObjectGrpcFromSupa(resp *storage_go.FileObject) *talv1.FileObject {
	log.Debug().Interface("buck is", resp.BucketId).Msg("bucucucuuc")
	return &talv1.FileObject{
		Name:      fmt.Sprintf("%s/%s", resp.BucketId, resp.Name),
		UpdatedAt: resp.UpdatedAt,
		BucketId:  resp.BucketId,
		// Metadata:  convertMetadata(resp.Metadata),
		CreatedAt: resp.CreatedAt,
		Id:        resp.Id,
	}
}

func (a *PublicAdapter) FileDeleteGrpcFromSupa(resp []storage_go.FileUploadResponse) *talv1.FileDeleteResponse {
	response := make([]*talv1.FileCreateResponse, len(resp))
	for index, rec := range resp {
		response[index] = a.FileCreateResponseGrpcFromSupa(&rec)
	}
	return &talv1.FileDeleteResponse{
		Responses: response,
	}
}
func (a *PublicAdapter) FileListGrpcFromSupa(resp []storage_go.FileObject, bucketId string) *talv1.FileListResponse {
	files := make([]*talv1.FileObject, len(resp))
	for index, rec := range resp {
		rec.BucketId = bucketId
		files[index] = a.FileObjectGrpcFromSupa(&rec)
	}
	return &talv1.FileListResponse{Files: files}
}
func (a *PublicAdapter) BucketListGrpcFromSupa(resp []storage_go.Bucket) *talv1.BucketListResponse {
	buckets := make([]*talv1.StorageBucket, len(resp))
	for index, rec := range resp {
		buckets[index] = a.StorageBucketGrpcFromSupa(&rec)
	}
	return &talv1.BucketListResponse{Buckets: buckets}
}
