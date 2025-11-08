package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (u *AccountsUsecase) RoleDelete(ctx context.Context, req *connect.Request[talv1.RoleDeleteRequest]) (*talv1.RoleDeleteResponse, error) {
	response := make([]*talv1.AccountsSchemaRole, 0)
	for _, role := range req.Msg.Records {
		params := db.RoleDeleteParams{
			RoleID: role,
		}
		deletedRole, err := u.repo.RoleDelete(ctx, params)
		if err != nil {
			return nil, err
		}

		response = append(response, u.adapter.RoleEntityGrpcFromSql(deletedRole))
	}
	return &talv1.RoleDeleteResponse{
		Records: response,
	}, nil
}
func (u *AccountsUsecase) RoleDeleteRestore(ctx context.Context, req *connect.Request[talv1.RoleDeleteRestoreRequest]) (*talv1.RoleDeleteRestoreResponse, error) {
	response := make([]*talv1.AccountsSchemaRole, 0)
	for _, rec := range req.Msg.Records {
		params := db.RoleDeleteRestoreParams{
			RoleID: rec,
		}
		resp, err := u.repo.RoleDeleteRestore(ctx, params)
		if err != nil {
			return nil, err
		}
		response = append(response, u.adapter.RoleEntityGrpcFromSql(resp))
	}
	return &talv1.RoleDeleteRestoreResponse{
		Records: response,
	}, nil
}
func (u *AccountsUsecase) RoleFindForUpdate(ctx context.Context, req *connect.Request[talv1.RoleFindForUpdateRequest]) (*talv1.RoleFindForUpdateResponse, error) {
	role, err := u.repo.RoleFindForUpdate(ctx, req.Msg.RecordId)
	if err != nil {
		return nil, err
	}
	request := u.adapter.RoleFindForUpdateUpdateGrpcFromSql(role)
	return &talv1.RoleFindForUpdateResponse{
		Request: request,
	}, nil
}
func (u *AccountsUsecase) RoleListInput(ctx context.Context) (*talv1.RoleListInputResponse, error) {
	roles, err := u.repo.RoleListInput(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.RoleListInputGrpcFromSql(roles)
	return response, nil
}

func (u *AccountsUsecase) RoleList(ctx context.Context, req *connect.Request[talv1.RoleListRequest]) (*talv1.RoleListResponse, error) {
	roles, err := u.repo.RoleList(ctx)
	if err != nil {
		return nil, err
	}
	response := u.adapter.RoleListGrpcFromSql(roles)
	return response, nil
}

func (u *AccountsUsecase) RoleCreateUpdate(ctx context.Context, req *connect.Request[talv1.RoleCreateUpdateRequest]) (*talv1.RoleCreateUpdateResponse, error) {
	roleCreateParams := u.adapter.RoleCreateUpdateSqlFromGrpc(req.Msg)
	role, err := u.repo.RoleCreateUpdate(ctx, *roleCreateParams)
	if err != nil {
		return nil, err
	}
	response := u.adapter.RoleCreateUpdateGrpcFromSql(role)
	return response, nil
}
