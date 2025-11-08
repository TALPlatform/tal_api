package api

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Extract record_id based on user_id or <entity>_id
func (s *Server) extractRecordID(req connect.AnyRequest, permissionGroup string) int32 {
	msg, ok := req.Any().(protoreflect.ProtoMessage)
	if !ok {
		return 0
	}
	msgReflect := msg.ProtoReflect()

	if val, ok := s.getFiledFromRequest(msgReflect, "user_id"); ok {
		return int32(val.Int())
	}
	if val, ok := s.getFiledFromRequest(msgReflect, permissionGroup+"_id"); ok {
		return int32(val.Int())
	}
	return 0
}

// Extract action name from function
func (s *Server) extractAction(fn string) string {
	fn = strings.ToLower(fn)
	switch {
	case strings.Contains(fn, "create"):
		return "create"
	case strings.Contains(fn, "update"):
		return "update"
	case strings.Contains(fn, "deleterestore"):
		return "restore"
	case strings.Contains(fn, "delete"):
		return "delete"
	default:
		return "unknown"
	}
}

func (s *Server) buildLogParams(
	ctx context.Context,
	req connect.AnyRequest,
	result connect.AnyResponse,
	err error,
	duration time.Duration,
) db.LogCreateParams {
	callerID, _ := contextkeys.CallerID(ctx)
	permissionGroup, _ := contextkeys.PermissionGroup(ctx)
	permissionFunction, _ := contextkeys.PermissionFunction(ctx)

	recordID := s.extractRecordID(req, permissionGroup)
	action := s.extractAction(permissionFunction)

	statusCode := "successful"
	apiErrorMessage := ""
	if err != nil {
		statusCode = connect.CodeOf(err).String()
		apiErrorMessage = err.Error()
	}
	return db.LogCreateParams{
		LogTitle:             req.Spec().Procedure,
		ActionType:           db.StringToPgtext(action),
		StatusCode:           db.StringToPgtext(statusCode),
		UserID:               int32(callerID),
		RecordID:             db.ToPgInt(recordID),
		DuartionMilliseconds: db.ToPgInt(int32(duration.Milliseconds())),
		ApiErrorMessage:      db.StringToPgtext(apiErrorMessage),
		PermissionName:       db.StringToPgtext(permissionFunction),
	}
}

// Sync log insert for error cases
func (s *Server) syncLogInsert(params db.LogCreateParams) {
	if _, err := s.store.LogCreate(context.Background(), params); err != nil {
		log.Error().Err(err).Msg("failed to insert sync API log")
	}
}

// Async log insert for successful requests
func (s *Server) asyncLogInsert(params db.LogCreateParams) {
	go func() {
		if _, err := s.store.LogCreate(context.Background(), params); err != nil {
			log.Error().Err(err).Msg("failed to insert async API log")
		}
	}()
}

// Print dev logs for local inspection
func (s *Server) printDevLog(
	req connect.AnyRequest,
	result connect.AnyResponse,
	duration time.Duration,
) {
	log.Info().
		Str("procedure", req.Spec().Procedure).
		Interface("request", req.Any()).
		Interface("response", result).
		Dur("duration", duration).
		Msg("gRPC request handled")
}

// Interceptor for logging all gRPC calls
func (s *Server) NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startTime := time.Now()
			result, err := next(ctx, req)
			duration := time.Since(startTime)

			logParams := s.buildLogParams(ctx, req, result, err, duration)

			if err != nil {
				s.syncLogInsert(logParams)
				return result, err
			}
			if logParams.ActionType.String != "unknown" && logParams.ActionType.String != "list" {
				s.asyncLogInsert(logParams)
			}

			if s.config.State == "dev" {
				s.printDevLog(req, result, duration)
			}
			return result, nil
		})
	})
}
