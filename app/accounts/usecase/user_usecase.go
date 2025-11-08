package usecase

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/dateutils"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/google/uuid"
	"github.com/supabase-community/auth-go/types"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (u *AccountsUsecase) UserDelete(ctx context.Context, req *connect.Request[talv1.UserDeleteRequest]) (*talv1.UserDeleteResponse, error) {
	var resp []*talv1.AccountsSchemaUser
	for _, recordId := range req.Msg.Records {
		params := db.UserDeleteParams{
			UserID: recordId,
		}
		user, err := u.repo.UserDelete(ctx, params)
		if err != nil {
			return nil, err
		}
		err = u.redisClient.AuthSessionClearAll(ctx, recordId)
		if err != nil {
			return nil, err
		}
		are := u.adapter.UserEntityGrpcFromSql(user)
		resp = append(resp, are)
	}
	return &talv1.UserDeleteResponse{Records: resp}, nil
}

func (u *AccountsUsecase) UserDeleteRestore(ctx context.Context, req *connect.Request[talv1.UserDeleteRestoreRequest]) (*talv1.UserDeleteRestoreResponse, error) {
	response := make([]*talv1.AccountsSchemaUser, 0)
	for _, rec := range req.Msg.Records {
		params := db.UserDeleteRestoreParams{
			UserID: rec,
		}
		resp, err := u.repo.UserDeleteRestore(ctx, params)
		if err != nil {
			return nil, err
		}
		response = append(response, u.adapter.UserEntityGrpcFromSql(resp))
	}
	return &talv1.UserDeleteRestoreResponse{
		Records: response,
	}, nil
}
func (u *AccountsUsecase) UserFindForUpdate(ctx context.Context, req *connect.Request[talv1.UserFindForUpdateRequest]) (*talv1.UserFindForUpdateResponse, error) {
	user, err := u.repo.UseriFindForUpdate(ctx, req.Msg.RecordId)
	if err != nil {
		return nil, err
	}
	request := u.adapter.UserFindForUpdateUpdateGrpcFromSql(user)
	return &talv1.UserFindForUpdateResponse{
		Request: request,
	}, nil
}

func (u *AccountsUsecase) UserFind(ctx context.Context, req *connect.Request[talv1.UserFindRequest]) (*talv1.UserFindResponse, error) {
	user, err := u.repo.UserFind(ctx, req.Msg.RecordId)
	if err != nil {
		return nil, err
	}

	record := u.adapter.UserFindRowGrpcFromSql(user)
	sessions, err := u.redisClient.AuthSessionListByUser(ctx, req.Msg.RecordId)
	if err != nil {
		return nil, err
	}
	if len(sessions) > 0 {
		latestSession := sessions[0]
		record.LastLogIn = dateutils.DateTimeToStringDigit(latestSession.CreatedAt)
	}
	return &talv1.UserFindResponse{
		Record:   record,
		Sessions: u.adapter.UserSessionsGrpcFropmSql(sessions),
	}, nil
}
func (u *AccountsUsecase) UserTypeListInput(ctx context.Context) (*talv1.UserTypeListInputResponse, error) {
	users, err := u.repo.UserTypeListInput(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.UserTypeListInputGrpcFromSql(users)
	return response, nil
}

func (u *AccountsUsecase) UserPermissionListInput(ctx context.Context, req *connect.Request[talv1.UserPermissionListInputRequest]) (*talv1.UserPermissionListInputResponse, error) {
	users, err := u.repo.UserPermissionListInput(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.UserPermissionListInputGrpcFromSql(users)
	return response, nil
}
func (u *AccountsUsecase) UserListInput(ctx context.Context) (*talv1.UserListInputResponse, error) {
	users, err := u.repo.UserListInput(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.UserListInputGrpcFromSql(users)
	return response, nil
}
func (u *AccountsUsecase) UserList(ctx context.Context) (*talv1.UserListResponse, error) {
	users, err := u.repo.UserList(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.UserListGrpcFromSql(users)
	return response, nil
}
func addBadRequestDetail(cErr *connect.Error, br *errdetails.BadRequest) {
	if br == nil {
		return
	}
	if detail, derr := connect.NewErrorDetail(br); derr == nil {
		cErr.AddDetail(detail)
	}
}

func brSingle(field, msg string) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{{
			Field:       field,
			Description: msg,
		}},
	}
}
func (u *AccountsUsecase) UserCreateUpdate(ctx context.Context, req *connect.Request[talv1.UserCreateUpdateRequest]) (*talv1.UserCreateUpdateResponse, error) {
	supabasRequest := types.AdminUpdateUserRequest{
		Email:    req.Msg.UserEmail,
		Password: req.Msg.UserPassword,
	}
	if req.Msg.UserId != 0 {
		userID, err := u.repo.AuthUserIDFindByEmail(ctx, req.Msg.UserEmail)
		if err != nil {
			return nil, err
		}
		uuid, err := uuid.Parse(*userID)
		if err != nil {
			return nil, err
		}
		supabasRequest.UserID = uuid
	}
	_, err := u.supaapi.UserCreateUpdate(supabasRequest)
	if err != nil {
		if err.Error() != "" {
			if strings.Contains(err.Error(), "user_already_exists") {
				cErr := connect.NewError(connect.CodeAlreadyExists,
					fmt.Errorf("duplicate value for %s", "userEmail"))
				addBadRequestDetail(cErr, brSingle("userEmail", "value already exists"))
				return nil, cErr
			}
		}
		return nil, err
	}
	userCreateParams := u.adapter.UserCreateUpdateSqlFromGrpc(req.Msg)
	user, err := u.repo.UserCreateUpdate(ctx, *userCreateParams)
	if err != nil {
		return nil, err
	}

	resp := u.adapter.UserCreateUpdateGrpcFromSql(user)
	return resp, nil
}
