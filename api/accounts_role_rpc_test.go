package api

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/random"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type roleCreateUpdateTest struct {
	name      string
	params    *talv1.RoleCreateUpdateRequest
	expectErr bool
}

func getValidRole() *talv1.RoleCreateUpdateRequest {
	return &talv1.RoleCreateUpdateRequest{
		RoleName:          random.RandomName(),
		RoleDescription:   random.RandomString(50),
		RoleSecurityLevel: 1,
		Permissions:       []int32{1, 2, 3},
	}
}
func TestRoleCreateUpdate(t *testing.T) {
	validRole := getValidRole()
	// Define a slice of test cases
	testcases := []roleCreateUpdateTest{
		// Test for a valid role creation.
		{
			name: "ValidRole",
			params: &talv1.RoleCreateUpdateRequest{
				RoleName:          validRole.RoleName,
				RoleDescription:   validRole.RoleDescription,
				RoleSecurityLevel: 1,
				Permissions:       validRole.Permissions,
			},
			expectErr: false,
		},
		{
			name: "ValidRoleUpdate",
			params: &talv1.RoleCreateUpdateRequest{
				RoleId:            3,
				RoleName:          "updated role name",
				RoleSecurityLevel: 1,
				RoleDescription:   validRole.RoleDescription,
				Permissions:       validRole.Permissions,
			},
			expectErr: false,
		},

		{
			name: "InValidNameTooShort",
			params: &talv1.RoleCreateUpdateRequest{
				RoleName:          random.RandomString(1),
				RoleSecurityLevel: 1,
				RoleDescription:   validRole.RoleDescription,
				Permissions:       validRole.Permissions,
			},
			expectErr: true,
		},
		{
			name: "InValidNameToLong",
			params: &talv1.RoleCreateUpdateRequest{
				RoleName:          random.RandomString(220),
				RoleSecurityLevel: 1,
				RoleDescription:   validRole.RoleDescription,
				Permissions:       validRole.Permissions,
			},
			expectErr: true,
		},
		{
			name: "InValidDescriptionToLong",
			params: &talv1.RoleCreateUpdateRequest{
				RoleName:          random.RandomString(120),
				RoleSecurityLevel: 1,
				RoleDescription:   random.RandomString(220),
				Permissions:       validRole.Permissions,
			},
			expectErr: true,
		},

		{
			name: "InvalideDuplicatedPermissions",
			params: &talv1.RoleCreateUpdateRequest{
				RoleName:          random.RandomString(120),
				RoleDescription:   random.RandomString(22),
				RoleSecurityLevel: 1,
				Permissions:       []int32{1, 1},
			},
			expectErr: true,
		},
	}

	// Loop through the test cases and test each one
	ctx := context.Background()
	// store := mockdb.NewMockStore(storeCtrl)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			roleReq := connect.NewRequest(tc.params)
			loginResp, err := testClient.AuthLogin(ctx, connect.NewRequest(&talv1.AuthLoginRequest{LoginCode: "admin@devkit.com", UserPassword: "123456"}))
			roleReq.Header().Add("Authorization", fmt.Sprintf("bearer %s", loginResp.Msg.LoginInfo.AccessToken))
			_, err = testClient.RoleCreateUpdate(context.Background(), roleReq)
			// If the current test case expects an error and no error occurred, fail the test
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but got none %s", tc.name)
			}
			// If the current test case does not expect an error and an error occurred, fail the test
			if !tc.expectErr && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
			// 	//
		})
	}
}
