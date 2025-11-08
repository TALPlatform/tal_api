package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/auth"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/TALPlatform/tal_api/pkg/headerkeys"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
)

func (api *Api) authenticateRequest(ctx context.Context, req connect.AnyRequest) (*auth.Payload, error) {

	return nil, nil
}

func (api *Api) authorizeRequestHeader(header http.Header) (*auth.Payload, error) {
	return nil, nil
}
func (api *Api) authorizedUserPermissions(ctx context.Context, payload *auth.Payload) (redisclient.PermissionsMap, error) {
	return nil, nil

}

func (api *Api) CurrentUserPermissionsMap(ctx context.Context) (redisclient.PermissionsMap, error) {
	callerId, ok := contextkeys.CallerID(ctx)
	if !ok {
		return nil, nil
	}
	permissionsMap, err := api.redisClient.UserPermissionFind(ctx, callerId)
	if err != nil {
		permissions, err := api.store.UserPermissionsMap(ctx, callerId)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		for _, rec := range permissions {
			groupPermissions := make(map[string]bool)
			err = json.Unmarshal(rec.Permissions, &groupPermissions)
			if err != nil {
				return nil, err
			}
			permissionsMap[rec.PermissionGroup] = groupPermissions
		}

		err = api.redisClient.UserPermissionCreate(ctx, callerId, &permissionsMap)
		if err != nil {
			return nil, err
		}

	}
	return permissionsMap, nil

}
func (api *Api) CheckForAccess(ctx context.Context, function string, group string) (*map[string]bool, error) {
	permissionsMap, err := api.CurrentUserPermissionsMap(ctx)
	if err != nil {
		return nil, err
	}

	permissionGroup, ok := permissionsMap[group]
	if !ok {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not have the required permission for this group %s", group))
	}
	isPermissionGranted, ok := permissionGroup[function]
	if !ok || !isPermissionGranted {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not have the required permission for this permission %s on this group %s", function, group))
	}

	return &permissionGroup, nil

}

func (api *Api) getAvailableOptions(header http.Header, variants ...string) *talv1.AvailableOptions {
	variant := "list"
	if len(variants) > 0 {
		variant = variants[0]
	}
	var (
		group                        = headerkeys.PermissionGroup(&header)
		findEndpoint                 = strcase.ToLowerCamel(fmt.Sprintf("%s_find_for_update", group))
		findRequestProperty          = "recordId"
		findResponseProperty         = "request"
		deleteRestoreRequestProperty = "records"
		redirectRoute                = strcase.ToCamel(fmt.Sprintf("%s_list", group))
		update                       = strcase.ToCamel(fmt.Sprintf("%s_update", group))
		create                       = strcase.ToCamel(fmt.Sprintf("%s_create", group))
		createUpdate                 = strcase.ToCamel(fmt.Sprintf("%s_create_update", group))
		deleteKey                    = strcase.ToCamel(fmt.Sprintf("%s_delete", group))
		deleteRestore                = strcase.ToCamel(fmt.Sprintf("%s_delete_restore", group))
		result                       = &talv1.AvailableOptions{
			Title:       fmt.Sprintf("%s_%s", variant, group),
			Description: fmt.Sprintf("%s_%s_description", variant, group),
		}
	)

	permittedActions := headerkeys.PermittedActions(&header)

	log.Debug().Interface("gr", permittedActions).Msg("hola")
	if variant == "list" {
		isCreatePermitted, ok := permittedActions[create]
		if isCreatePermitted && ok {
			result.CreateHandler = &talv1.CreateHandler{
				RedirectRoute: redirectRoute,
				Title:         create,
				Endpoint:      strcase.ToLowerCamel(createUpdate),
				RouteName:     strcase.ToSnake(create),
			}
		}
	}
	isUpdatePermitted, ok := permittedActions[update]
	if isUpdatePermitted && ok {
		result.UpdateHandler = &talv1.UpdateHandler{
			RedirectRoute:        redirectRoute,
			Title:                update,
			Endpoint:             strcase.ToLowerCamel(createUpdate),
			RouteName:            strcase.ToSnake(update),
			FindEndpoint:         findEndpoint,
			FindRequestProperty:  findRequestProperty,
			FindResponseProperty: findResponseProperty,
		}
	}
	isDeletePermitted, ok := permittedActions[deleteKey]
	if isDeletePermitted && ok {
		result.DeleteHandler = &talv1.DeleteHandler{
			Endpoint:        strcase.ToLowerCamel(deleteKey),
			RequestProperty: deleteRestoreRequestProperty,
		}
	}
	isDeleteRestorePermitted, ok := permittedActions[deleteRestore]
	if isDeleteRestorePermitted && ok {
		result.DeleteRestoreHandler = &talv1.DeleteRestoreHandler{
			Endpoint:        strcase.ToLowerCamel(deleteRestore),
			RequestProperty: deleteRestoreRequestProperty,
		}
	}

	return result
}
