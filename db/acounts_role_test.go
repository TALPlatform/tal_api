package db

import (
	"context"
	"testing"

	"github.com/TALPlatform/tal_api/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestRoleCreateUpdate(t *testing.T) {
	validName := random.RandomName()
	testcases := []struct {
		name             string
		params           *RoleCreateUpdateParams
		expectedErrorMsg string
		expectErr        bool
	}{
		{
			name: "ValidRole",
			params: &RoleCreateUpdateParams{
				RoleName:          validName,
				RoleSecurityLevel: 1,
				CallerID:          2,
				RoleDescription:   random.RandomString(50),
				Permissions:       []int32{1, 2, 3},
			},
			expectErr: false,
		},
		{
			name: "ValidRoleUpdate",
			params: &RoleCreateUpdateParams{
				RoleID:            1,
				RoleName:          random.RandomName(),
				RoleSecurityLevel: 1,
				CallerID:          2,
				RoleDescription:   random.RandomString(50),
				Permissions:       []int32{1, 2, 3},
			},
			expectErr: false,
		},

		{
			name: "RoleNameTooLong",
			params: &RoleCreateUpdateParams{
				RoleName:          random.RandomString(220), // Exceeds max length
				RoleSecurityLevel: 1,
				RoleDescription:   random.RandomString(50),
				CallerID:          2,
				Permissions:       []int32{1, 2, 3},
			},
			expectErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the RoleCreateUpdate function with the test parameters
			role, err := store.RoleCreateUpdate(context.Background(), *tc.params)

			if tc.expectErr {
				require.Error(t, err)
				require.Empty(t, role) // Expect no role to be returned on error
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, role)
				require.Equal(t, tc.params.RoleName, role.RoleName)
				require.Equal(t, tc.params.RoleDescription, role.RoleDescription.String)
			}
		})
	}

}
