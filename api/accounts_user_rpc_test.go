package api

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/random"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/stretchr/testify/require"
)

var (
	userEmail       = random.RandomEmail()
	userPassword    = random.RandomString(8)
	userNewPassword = random.RandomString(10)
	userPhone       = random.RandomPhone()
	userNewPhone    = random.RandomPhone()
	loginRequest    = connect.NewRequest(&talv1.AuthLoginRequest{
		LoginCode:    userEmail,
		UserPassword: userPassword,
	})

	emptyRequest      = connect.NewRequest(&talv1.RoleListRequest{})
	userCreateRequest = connect.NewRequest(&talv1.UserCreateUpdateRequest{
		UserId:            0,
		UserEmail:         userEmail,
		UserPassword:      userPassword,
		UserPhone:         userPhone,
		UserTypeId:        1,
		UserSecurityLevel: 1,
		UserName:          random.RandomName(),
	})
)

func TestCycle(t *testing.T) {
	ctx := context.Background()

	t.Run("Login not existed user", func(t *testing.T) {
		_, err := testClient.AuthLogin(ctx, loginRequest)
		require.NotEmpty(t, err)
	})
	// // // // create new user
	t.Run("create new user with no roles", func(t *testing.T) {
		loginResp, err := testClient.AuthLogin(ctx, connect.NewRequest(&talv1.AuthLoginRequest{LoginCode: "admin@devkit.com", UserPassword: "123456"}))
		userCreateRequest.Header().Add("Authorization", fmt.Sprintf("bearer %s", loginResp.Msg.LoginInfo.AccessToken))

		resp, err := testClient.UserCreateUpdate(ctx, userCreateRequest)
		require.NoError(t, err)
		userCreateRequest.Msg.UserId = resp.Msg.User.UserId

		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // // login with wrong password
	t.Run("Login with wrong password", func(t *testing.T) {

		wrongLoginRequest := connect.NewRequest(&talv1.AuthLoginRequest{
			LoginCode:    userEmail,
			UserPassword: "wrongPassword",
		})

		_, err := testClient.AuthLogin(ctx, wrongLoginRequest)
		require.NotEmpty(t, err)
	})
	// // // // login with wrong password
	t.Run("Login with correct password", func(t *testing.T) {
		resp, err := testClient.AuthLogin(ctx, loginRequest)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Msg.LoginInfo.AccessToken)
		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // update the created user password
	t.Run("update the password of created user", func(t *testing.T) {
		userCreateRequest.Msg.UserPassword = userNewPassword
		resp, err := testClient.UserCreateUpdate(ctx, userCreateRequest)
		require.NoError(t, err)
		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // login with the old password
	t.Run("Login with old password", func(t *testing.T) {
		_, err := testClient.AuthLogin(ctx, loginRequest)
		require.NotEmpty(t, err)
		require.Contains(t, err.Error(), "invalid login")
	})
	// // login with the updated password
	t.Run("Login with updated password", func(t *testing.T) {
		loginRequest.Msg.UserPassword = userNewPassword
		resp, err := testClient.AuthLogin(ctx, loginRequest)
		require.NoError(t, err)
		emptyRequest.Header().Add("Authorization", fmt.Sprintf("bearer %s", resp.Msg.LoginInfo.AccessToken))
		require.NotEmpty(t, resp.Msg.LoginInfo.AccessToken)
		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // try to access ednpoint that needs access
	t.Run("forbidden access", func(t *testing.T) {
		req := connect.NewRequest(&talv1.RoleListRequest{})
		_, err := testClient.RoleList(context.Background(), req)
		require.NotEmpty(t, err)
	})
	// // update the created user  roles
	t.Run("update the roles of created user", func(t *testing.T) {
		userCreateRequest.Msg.Roles = []int32{1}
		resp, err := testClient.UserCreateUpdate(ctx, userCreateRequest)
		require.NoError(t, err)
		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // login with the updated roles
	t.Run("Login with updated roles", func(t *testing.T) {
		resp, err := testClient.AuthLogin(ctx, loginRequest)
		require.NoError(t, err)
		emptyRequest.Header().Add("Authorization", fmt.Sprintf("bearer %s", resp.Msg.LoginInfo.AccessToken))
		require.NotEmpty(t, resp.Msg.LoginInfo.AccessToken)
		require.Equal(t, resp.Msg.User.UserName, userCreateRequest.Msg.UserName)
		require.Equal(t, resp.Msg.User.UserEmail, userCreateRequest.Msg.UserEmail)
		require.Equal(t, resp.Msg.User.UserPhone, userCreateRequest.Msg.UserPhone)
	})
	// // try to recall the endpoint after permissions added
	// t.Run("granted access", func(t *testing.T) {
	// 	_, err := testClient.RoleList(ctx, emptyRequest)
	// 	require.NoError(t, err)
	// })
	// //
}
