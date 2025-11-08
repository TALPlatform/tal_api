package usecase

import (
	"context"

	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (s *PublicUsecase) TranslationDelete(ctx context.Context, req *talv1.TranslationDeleteRequest) (*talv1.TranslationDeleteResponse, error) {
	resp, err := s.repo.TranslationDelete(ctx, req.Keys)
	if err != nil {
		return nil, err
	}
	response := s.adapter.TranslationListGrpcFromSql(resp)
	return &talv1.TranslationDeleteResponse{
		Translations: response.Translations,
	}, nil

}

func (s *PublicUsecase) TranslationCreateUpdateBulk(ctx context.Context, req *talv1.TranslationCreateUpdateBulkRequest) (*talv1.TranslationCreateUpdateBulkResponse, error) {
	params := s.adapter.TranslationCreateUpdateBulkSqlFromGrpc(req)
	resp, err := s.repo.TranslationCreateUpdateBulk(ctx, *params)
	if err != nil {
		return nil, err
	}
	response := s.adapter.TranslationCreateUpdateBulkGrpcFromSql(resp)
	return &talv1.TranslationCreateUpdateBulkResponse{
		Translations: response.Translations,
	}, nil
}

func (u *PublicUsecase) TranslationFindLocale(ctx context.Context, req *talv1.TranslationFindLocaleRequest) (*talv1.TranslationFindLocaleResponse, error) {
	settings, err := u.repo.TranslationList(ctx)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.TranslationFindLocaleGrpcFromSql(settings, req.Locale)
	return &resp, nil
}
func (u *PublicUsecase) TranslationList(ctx context.Context) (*talv1.TranslationListResponse, error) {
	settings, err := u.repo.TranslationList(ctx)

	if err != nil {
		return nil, err
	}
	resp := u.adapter.TranslationListGrpcFromSql(settings)
	return &resp, nil
}
