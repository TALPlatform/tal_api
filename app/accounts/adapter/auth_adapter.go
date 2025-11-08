package adapter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/iancoleman/strcase"
	"github.com/supabase-community/auth-go/types"
)

type NavigationsMap map[int32]*talv1.NavigationBarItem

// The UserNavigationBarFindGrpcFromSql function builds a hierarchical structure of navigation bar items from a sorted database response (dbResponse).
// This hierarchy supports any number of levels, as long as the input list is sorted in ascending order by level.
// The function uses navigationBarItemId as a primary key and parentId as a foreign key to link items with their parents.
// The function logic breakdown is listed as comments inside the function body.
func (a *AccountsAdapter) UserNavigationBarFindGrpcFromSql(dbResponse []db.UserNavigationBarFindRow) ([]*talv1.NavigationBarItem, error) {
	// 1. Get the maximum level in the tree by accessing the last element in the array, since the response is sorted by level.
	maxLevel := dbResponse[len(dbResponse)-1].Level

	// 2. Declare rootItems, which will store the top-level items and serve as the function's return value.
	rootItems := make([]*talv1.NavigationBarItem, 0)

	// 3. Initialize levelItemsMap with empty maps. This will store each level's items separately in its own hashmap.
	//    We populate levelItemsMap in reverse order so that the highest level appears first.
	//    For example, if maxLevel is 3, levelItemsMap will look like:
	//    [
	//      {map with all items of level 3},
	//      {map with all items of level 2},
	//      {map with all items of level 1},
	//    ]
	levelItemsMap := make([]NavigationsMap, maxLevel)
	for i := range levelItemsMap {
		levelItemsMap[i] = NavigationsMap{}
	}

	// 4. Populate levelItemsMap from dbResponse.
	for _, dbItem := range dbResponse {
		primaryKeyValue := dbItem.NavigationBarItemID
		grpcItem := a.NavigationBarItemGrpcFromSql(&dbItem)

		// Calculate levelIndex to store items in reverse order in levelItemsMap.
		// If dbItem.Level equals maxLevel, it will go to the first index in levelItemsMap,
		// and if dbItem.Level equals 1, it will go to the last index.
		levelIndex := maxLevel - dbItem.Level
		levelItemsMap[levelIndex][primaryKeyValue] = grpcItem
	}

	// 5. Construct the response by linking items from levelItemsMap.
	for levelIndex, grpcItemsMap := range levelItemsMap {
		for _, grpcItem := range grpcItemsMap {
			// If we're at the last level in levelItemsMap (root level),
			// we add the item directly to rootItems, as it has no parent to attach to.
			if levelIndex == int(maxLevel)-1 {
				rootItems = append(rootItems, grpcItem)
				continue
			}

			// Retrieve the parent level, which is always the next item in levelItemsMap.
			// This works because levelItemsMap is ordered in reverse.
			// The previous check prevents out-of-bounds errors by ensuring we're not at the root level.
			parentLevel := levelItemsMap[levelIndex+1]

			// Find the parentNode of the current item by looking up its parentId in the parent level map.
			parentNode, ok := parentLevel[grpcItem.ParentId]
			if !ok {
				return nil, fmt.Errorf("item %s with parent id %d couldn't be found the parent level - check that data is sorted by level", grpcItem.Label, grpcItem.ParentId)
			}
			// Attach the current item as a child of its parentNode.
			parentNode.Items = append(parentNode.Items, grpcItem)

		}
	}
	return rootItems, nil
}
func (a *AccountsAdapter) AuthSessionListGrpcFromRedis(resp []*redisclient.AuthSession) *talv1.AuthSessionListResponse {
	records := make([]*talv1.AuthSession, 0)
	for _, session := range resp {
		records = append(records, &talv1.AuthSession{
			SessionKey:                    session.SessionKey,
			TokenId:                       session.TokenID,
			UserId:                        session.UserID,
			IpAddress:                     session.IPAddress,
			IsBlocked:                     session.IsBlocked,
			UserAgent:                     session.UserAgent,
			AccessToken:                   session.AccessToken,
			RefreshToken:                  session.RefreshToken,
			SupabaseAccessToken:           session.SupabaseAccessToken,
			SupabaseRefreshToken:          session.SupabaseRefreshToken,
			AccessTokenExpiresAt:          session.AccessTokenExpiresAt.Format(time.RFC3339),
			RefreshTokenExpiresAt:         session.RefreshTokenExpiresAt.Format(time.RFC3339),
			SupabaseAccessTokenExpiresAt:  session.SupabaseAccessTokenExpiresAt.Format(time.RFC3339),
			SupabaseRefreshTokenExpiresAt: session.SupabaseRefreshTokenExpiresAt.Format(time.RFC3339),
		})
	}
	return &talv1.AuthSessionListResponse{
		Records: records,
	}
}
func (a *AccountsAdapter) UserCreateUpdateRequestFromAuthRegister(req *talv1.AuthRegisterRequest) *talv1.UserCreateUpdateRequest {
	resp := &talv1.UserCreateUpdateRequest{
		UserName:     req.UserName,
		UserTypeId:   req.UserTypeId,
		UserPhone:    req.UserPhone,
		UserEmail:    req.UserEmail,
		UserPassword: req.UserPassword,
	}
	return resp
}

func (a *AccountsAdapter) UserPermissionsMapRedisFromSql(resp *[]db.UserPermissionsMapRow) (*redisclient.PermissionsMap, error) {
	respMap := make(redisclient.PermissionsMap)
	for _, rec := range *resp {
		perms := make(map[string]bool)
		err := json.Unmarshal(rec.Permissions, &perms)
		if err != nil {
			return nil, err
		}
		respMap[rec.PermissionGroup] = perms
	}
	return &respMap, nil
}
func (a *AccountsAdapter) AuttSessionRedisFromGrpc(response *talv1.AuthLoginResponse, ipAddress string, userAgent string) (*redisclient.AuthSession, error) {
	supabaseAccessExpiredAt, err := db.StringToTime(response.LoginInfo.SupabaseTokenExpiresAt)
	supabaseRefreshExpiredAt, err := db.StringToTime(response.LoginInfo.SupabaseRefreshTokenExpiresAt)
	accessTokenExpiredAt, err := db.StringToTime(response.LoginInfo.AccessTokenExpiresAt)
	refreshTokenExpiredAt, err := db.StringToTime(response.LoginInfo.RefreshTokenExpiresAt)
	if err != nil {
		return nil, err
	}

	return &redisclient.AuthSession{
		UserID:                        response.User.UserId,
		AccessToken:                   response.LoginInfo.AccessToken,
		SupabaseAccessTokenExpiresAt:  supabaseAccessExpiredAt,
		SupabaseRefreshTokenExpiresAt: supabaseRefreshExpiredAt,
		RefreshToken:                  response.LoginInfo.RefreshToken,
		SupabaseAccessToken:           response.LoginInfo.SupabaseToken,
		AccessTokenExpiresAt:          accessTokenExpiredAt,
		IPAddress:                     ipAddress,
		UserAgent:                     userAgent,
		RefreshTokenExpiresAt:         refreshTokenExpiredAt,
		SupabaseRefreshToken:          response.LoginInfo.SupabaseRefreshToken,
	}, nil
}
func (a *AccountsAdapter) AuthLoginSqlFromGrpc(req *talv1.AuthLoginRequest) (*db.UserFindForAuthParams, *types.TokenRequest) {
	isEmail := strings.Contains(req.LoginCode, "@") && strings.Contains(req.LoginCode, ".")
	supabseRequest := &types.TokenRequest{Password: req.UserPassword}
	normalizedCode := strings.ToLower(strings.TrimSpace(req.LoginCode))
	if isEmail {
		supabseRequest.Email = normalizedCode
	} else {
		supabseRequest.Phone = normalizedCode
	}
	supabseRequest.GrantType = "password"
	return &db.UserFindForAuthParams{
		SearchKey: normalizedCode,
	}, supabseRequest
}

func (a *AccountsAdapter) AuthLoginGrpcFromSql(resp *db.AccountsSchemaUserView) *talv1.AuthLoginResponse {
	return &talv1.AuthLoginResponse{
		User: a.UserViewEntityGrpcFromSql(resp),
	}
}
func (a *AccountsAdapter) NavigationBarItemGrpcFromSql(resp *db.UserNavigationBarFindRow) *talv1.NavigationBarItem {
	if !resp.LabelAr.Valid {
		resp.LabelAr.String = resp.Label
	}
	return &talv1.NavigationBarItem{
		Key:                 resp.MenuKey,
		Level:               resp.Level,
		NavigationBarItemId: resp.NavigationBarItemID,
		ParentId:            resp.ParentID.Int32,
		Label:               resp.Label,
		LabelAr:             resp.LabelAr.String,
		Icon:                resp.Icon.String,
		Route:               resp.Route.String,
	}
}

func (a *AccountsAdapter) AuthResetPasswordSupaFromGrpc(req *talv1.AuthResetPasswordRequest) *types.VerifyForUserRequest {
	return &types.VerifyForUserRequest{
		Type:       types.VerificationTypeRecovery,
		Token:      req.ResetToken,
		Email:      req.Email,
		RedirectTo: req.RedirectUrl,
	}
}

func (a *AccountsAdapter) PermissionGroupFromFunctionName(funName string) string {
	functionNameSnake := strcase.ToSnake(funName)
	functionNameParts := strings.Split(functionNameSnake, "_")
	group := functionNameParts[0]

	return group
}
