package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
)

type PublicRepoInterface interface {
	IconCreateUpdateBulk(ctx context.Context, req db.IconCreateUpdateBulkParams) (*[]db.Icon, error)
	IconList(ctx context.Context) (*[]db.Icon, error)
	IconFind(ctx context.Context, req db.IconFindParams) (*db.Icon, error)
	TranslationDelete(ctx context.Context, req []string) ([]db.Translation, error)
	TranslationList(ctx context.Context) ([]db.Translation, error)
	TranslationCreateUpdateBulk(ctx context.Context, req db.TranslationCreateUpdateBulkParams) ([]db.TranslationCreateUpdateBulkRow, error)
	SettingUpdate(ctx context.Context, req *db.SettingUpdateParams) error
	FileDelete(ctx context.Context, req []string) (string, error)
	FileDeleteByBucket(ctx context.Context, req []string, bucket string) (string, error)
	SettingFindForUpdate(ctx context.Context) (*[]db.SettingFindForUpdateRow, error)
}

type PublicRepo struct {
	store        db.Store
	errorHandler map[string]string
}

func NewPublicRepo(store db.Store) PublicRepoInterface {
	errorHandler := map[string]string{
		"settings_setting_key_key": "roleName",
		"icons_icon_name_key":      "userName",
	}
	return &PublicRepo{
		store:        store,
		errorHandler: errorHandler,
	}
}
