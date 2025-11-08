package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *PublicAdapter) TranslationCreateUpdateBulkSqlFromGrpc(req *talv1.TranslationCreateUpdateBulkRequest) *db.TranslationCreateUpdateBulkParams {
	keys := make([]string, len(req.Records))
	enValues := make([]string, len(req.Records))
	arValues := make([]string, len(req.Records))
	for index, v := range req.Records {
		keys[index] = v.TranslationKey
		enValues[index] = v.EnglishValue
		arValues[index] = v.ArabicValue
	}
	return &db.TranslationCreateUpdateBulkParams{
		Keys:          keys,
		ArabicValues:  arValues,
		EnglishValues: enValues,
	}
}

func (a *PublicAdapter) TranslationCreateUpdateBulkRowGrpcFromSql(resp *db.TranslationCreateUpdateBulkRow) *talv1.Translation {
	return &talv1.Translation{
		TranslationKey: resp.TranslationKey,
		EnglishValue:   resp.EnglishValue,
		ArabicValue:    resp.ArabicValue,
	}
}

func (a *PublicAdapter) TranslationGrpcFromSql(resp *db.Translation) *talv1.Translation {
	return &talv1.Translation{
		TranslationKey: resp.TranslationKey,
		EnglishValue:   resp.EnglishValue,
		ArabicValue:    resp.ArabicValue,
	}
}

func (a *PublicAdapter) TranslationCreateUpdateBulkGrpcFromSql(resp []db.TranslationCreateUpdateBulkRow) talv1.TranslationListResponse {
	translations := make([]*talv1.Translation, len(resp))
	for index, t := range resp {
		translations[index] = a.TranslationCreateUpdateBulkRowGrpcFromSql(&t)
	}
	return talv1.TranslationListResponse{
		Translations: translations,
	}
}
func (a *PublicAdapter) TranslationListGrpcFromSql(resp []db.Translation) talv1.TranslationListResponse {
	translations := make([]*talv1.Translation, len(resp))
	for index, t := range resp {
		translations[index] = a.TranslationGrpcFromSql(&t)
	}
	return talv1.TranslationListResponse{
		Translations: translations,
	}
}
func (a *PublicAdapter) TranslationFindLocaleGrpcFromSql(resp []db.Translation, locale string) talv1.TranslationFindLocaleResponse {
	translations := make(map[string]string, len(resp))
	for _, t := range resp {
		value := t.EnglishValue
		if locale == "ar" {
			value = t.ArabicValue
		}
		translations[t.TranslationKey] = value
	}
	return talv1.TranslationFindLocaleResponse{
		Translations: translations,
	}
}
