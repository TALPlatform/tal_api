package db

import (
	"context"
	"testing"

	"github.com/TALPlatform/tal_api/pkg/random"
	"github.com/stretchr/testify/require"
)

var (
	name              = random.RandomName()
	email             = random.RandomEmail()
	phone             = random.RandomPhone()
	password          = random.RandomString(8)
	validCreateParams = &UserCreateUpdateParams{
		CallerID:     1,
		UserName:     name,
		UserEmail:    email,
		UserPassword: password,
		UserPhone:    phone,
		UserTypeID:   1,
		Roles:        []int32{1},
	}
	duplicatedPhoneParams = &UserCreateUpdateParams{

		CallerID:     1,
		UserName:     random.RandomName(),
		UserEmail:    random.RandomName(),
		UserPassword: password,
		UserPhone:    phone,
		UserTypeID:   1,
		Roles:        []int32{1},
	}
	duplicatedEmailParams = &UserCreateUpdateParams{

		CallerID:     1,
		UserName:     random.RandomName(),
		UserEmail:    email,
		UserPassword: password,
		UserPhone:    random.RandomPhone(),
		UserTypeID:   1,
		Roles:        []int32{1},
	}
)

func TestUserCreateUpdate(t *testing.T) {
	// Define test input parameters.
	testcases := []struct {
		name             string
		params           *UserCreateUpdateParams
		expectedErrorMsg string
		expectErr        bool
	}{
		{
			name:      "ValidiAdmin",
			params:    validCreateParams,
			expectErr: false,
		},
		{
			name:      "DuplicatePhone",
			params:    duplicatedPhoneParams,
			expectErr: true,
		},
		{
			name:      "DuplicateEmail",
			params:    duplicatedEmailParams,
			expectErr: true,
		},
		{
			name: "InvalidRoleId",
			params: &UserCreateUpdateParams{
				UserName:     random.RandomString(20), // Exceeds max length
				CallerID:     1,
				UserEmail:    random.RandomEmail(),
				UserPassword: "123456",
				UserTypeID:   1,
				UserPhone:    random.RandomString(50),
				Roles:        []int32{1, 2, 300},
			},
			expectErr: true,
		},
		{
			name: "UserNameTooLong",
			params: &UserCreateUpdateParams{

				CallerID:     1,
				UserName:     random.RandomString(220), // Exceeds max length
				UserPassword: "123456",
				UserEmail:    random.RandomEmail(),
				UserTypeID:   1,
				UserPhone:    random.RandomString(50),
				Roles:        []int32{1},
			},
			expectErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the UserCreateUpdate function with the test parameters
			role, err := store.UserCreateUpdate(context.Background(), *tc.params)

			if tc.expectErr {
				require.Error(t, err)
				require.Empty(t, role) // Expect no role to be returned on error
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, role)
				require.Equal(t, tc.params.UserName, role.UserName)
				require.Equal(t, tc.params.UserPhone, role.UserPhone.String)
			}

		})
	}

}
