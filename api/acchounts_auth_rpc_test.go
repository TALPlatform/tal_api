package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/random"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

var req *talv1.AuthRegisterRequest = &talv1.AuthRegisterRequest{
	UserName:     random.RandomName(),
	UserEmail:    random.RandomEmail(),
	UserTypeId:   1,
	UserPassword: random.RandomString(8),
}

func TestAuthRegister(t *testing.T) {
	resp, err := testClient.AuthRegister(
		context.Background(),
		connect.NewRequest(req),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Msg.LoginInfo == nil || resp.Msg.LoginInfo.AccessToken == "" {
		t.Fatalf("expected tokens, got none")
	}
}

func TestAuthLogin(t *testing.T) {
	req := &talv1.AuthLoginRequest{
		LoginCode:    req.UserEmail,
		UserPassword: req.UserPassword,
	}
	resp, err := testClient.AuthLogin(context.Background(), connect.NewRequest(req))
	if err != nil {
		t.Fatalf("expected login to succeed, got error: %v", err)
	}
	if resp.Msg.LoginInfo == nil || resp.Msg.LoginInfo.AccessToken == "" {
		t.Fatalf("expected tokens, got none")
	}
}
