package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/auth"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/TALPlatform/tal_api/pkg/headerkeys"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (s *Server) InjectRefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(s.config.RefreshTokenCookieName); err == nil {
			ctx := contextkeys.WithRefreshToken(r.Context(), cookie.Value)
			r = r.WithContext(ctx)
		}

		if cookie, err := r.Cookie(s.config.SupabaseTokenCookieName); err == nil {
			ctx := contextkeys.WithSupabaseToken(r.Context(), cookie.Value)
			r = r.WithContext(ctx)
		}

		if cookie, err := r.Cookie(s.config.SupabaseRefreshTokenCookieName); err == nil {
			ctx := contextkeys.WithSupabaseRefreshToken(r.Context(), cookie.Value)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
func (s *Server) proccessProcedureName(procName string) (string, string) {
	procParts := strings.Split(procName, "/")
	procedureName := strings.TrimLeft(procParts[len(procParts)-1], "Input")
	functionNameSnake := strcase.ToSnake(procedureName)
	functionNameParts := strings.Split(functionNameSnake, "_")
	group := functionNameParts[0]
	log.Debug().Interface("err", group).Msg("group")
	return procedureName, group
}

func (s *Server) getFiledFromRequest(msgReflect protoreflect.Message, fieledName string) (*protoreflect.Value, bool) {
	field := msgReflect.Descriptor().Fields().ByName(protoreflect.Name(fieledName))
	if field == nil || !msgReflect.Has(field) {
		return nil, false
	}
	value := msgReflect.Get(field)
	return &value, true
}

// this function should handle the create update request and chech for the record id on the request to pass either [create , update]
// as permission variation
func (s *Server) createUpdateMethodPermissionName(msgReflect protoreflect.Message, group string) string {
	permissionName := strcase.ToCamel(fmt.Sprintf("%s_create", group))
	recordIdValue, ok := s.getFiledFromRequest(msgReflect, "record_id")
	if !ok {
		recordIdValue, ok = s.getFiledFromRequest(msgReflect, fmt.Sprintf("%s_id", group))
		if !ok {
			return permissionName
		}
	}
	recordID := recordIdValue.Int() // assuming role_id is an int32 or int64 field
	if recordID > 0 {
		permissionName = strcase.ToCamel(fmt.Sprintf("%s_update", group))
	}
	return strcase.ToCamel(permissionName)

}
func (s *Server) AuthorizeRequest(req connect.AnyRequest) (*auth.Payload, error) {
	authHeader := req.Header().Get("Authorization")
	if authHeader == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing metadata"))
	}
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid authorization header format"))
	}

	authType := strings.ToLower(fields[0])
	if authType != "bearer" {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unsupported authorization type: %s", authType))
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	return payload, nil
}

// checkRecordID checks if a connect.AnyRequest contains a field named "recordId"
// (or a dynamically generated group-based field) with a value greater than zero.
func (s *Server) NewAuthenticationInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			methodDesc, ok := req.Spec().Schema.(protoreflect.MethodDescriptor)
			if !ok {
				return nil, fmt.Errorf("Schema is not a MethodDescriptor")
			}
			// Access the method options
			options := methodDesc.Options()
			var payload *auth.Payload
			// var err error

			procedureName, group := s.proccessProcedureName(req.Spec().Procedure)
			if options != nil {
				skipAuth, ok := proto.GetExtension(options, talv1.E_SkipAuthentication).(bool)
				if skipAuth && ok {
					return next(ctx, req)
				}

				permissionGroup, ok := proto.GetExtension(options, talv1.E_PermissionGroup).(string)
				if permissionGroup != "" && ok {
					group = permissionGroup
				}

				permissionName, ok := proto.GetExtension(options, talv1.E_PermissionName).(string)
				if permissionName != "" && ok {
					procedureName = permissionName
				}
				var err error
				payload, err = s.AuthorizeRequest(req)
				if err != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token : %w", err))
				}

				session, err := s.redisClient.AuthSessionFindByAccessTokenID(ctx, payload.ID.String())
				if err != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("session not stored : %w", err))
				}
				if session.IsBlocked {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("session is blocked"))
				}
				// Inject the callerID into the context
				ctx = contextkeys.WithCallerID(ctx, payload.UserId)

				skipAuthorization, ok := proto.GetExtension(options, talv1.E_SkipAuthorization).(bool)
				if skipAuthorization && ok {
					return next(ctx, req)
				}

			}
			isCreateUpdate := strings.Contains(procedureName, "CreateUpdate")

			ctx = contextkeys.WithiTenantID(ctx, payload.TenantId)
			if isCreateUpdate {
				message, ok := req.Any().(protoreflect.ProtoMessage)
				if !ok {
					return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("request is not a ProtoMessage"))
				}
				procedureName = s.createUpdateMethodPermissionName(message.ProtoReflect(), group)
			}
			ctx = contextkeys.WithPermissionGroup(ctx, group)

			ctx = contextkeys.WithPermissionFunction(ctx, procedureName)
			log.Debug().Interface("err", group).Msg("group")

			headerkeys.WithPermissionGroup(req.Header(), group)
			return next(ctx, req)

		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
