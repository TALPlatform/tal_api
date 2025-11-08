package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/TALPlatform/tal_api/pkg/headerkeys"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (s *Server) NewAuthorizationInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// this ok will be false if endpoint has the skip authorization option

			group, ok := contextkeys.PermissionGroup(ctx)
			if !ok {
				return next(ctx, req)
			}

			permissionFunction, ok := contextkeys.PermissionFunction(ctx)
			if !ok {
				return next(ctx, req)
			}
			// this ok will be false if the user type is not tenant and the logged in user don't have the attribute tenant_id set on the db
			tenantId, ok := contextkeys.TenantID(ctx)
			if ok {

				// here we will check if the logged in user has certain tenant id to return error if he passed diffrent tenant id on the request
				message, ok := req.Any().(protoreflect.ProtoMessage)
				if !ok {
					return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("request is not a ProtoMessage"))
				}
				msgReflect := message.ProtoReflect()
				// get the value of the tenant_id key on the incoming request
				requestTenantIdValue, ok := s.getFiledFromRequest(msgReflect, "tenant_id")
				if ok {
					requestTenantId := requestTenantIdValue.Int()
					if requestTenantId != int64(tenantId) && tenantId > 0 {
						return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("this user attached to TenantId number %d , not allowed to take any actions under another tenant %d", tenantId, requestTenantId))
					}
				}
			}
			permissionGroup, err := s.api.(*Api).CheckForAccess(ctx, permissionFunction, group)
			if err != nil {
				return nil, err
			}
			// permissionsMap, err := s.redisClient.UserPermissionFind(ctx, callerId)
			// if err != nil {
			// 	permissions, err := s.store.UserPermissionsMap(ctx, callerId)
			// 	if err != nil {
			// 		return nil, connect.NewError(connect.CodeUnauthenticated, err)
			// 	}
			// 	for _, rec := range permissions {
			// 		groupPermissions := make(map[string]bool)
			// 		err := json.Unmarshal(rec.Permissions, &groupPermissions)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		permissionsMap[rec.PermissionGroup] = groupPermissions
			// 	}
			//
			// 	err = s.redisClient.UserPermissionCreate(ctx, callerId, &permissionsMap)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			//
			// }
			// permissionGroup, ok := permissionsMap[group]
			// if !ok {
			// 	return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not have the required permission for this group %s", group))
			// }
			// isPermissionGranted, ok := permissionGroup[permissionFunction]
			// if !ok || !isPermissionGranted {
			// 	return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not have the required permission for this permission %s on this group %s", permissionFunction, group))
			// }
			//
			// log.Debug().Interface("local", permissionGroup).Msg("local")
			headerkeys.WithPermittedActions(req.Header(), *permissionGroup)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
