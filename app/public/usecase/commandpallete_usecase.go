package usecase

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)
func (u *PublicUsecase) CommandPalleteSync(ctx context.Context, req *connect.Request[talv1.CommandPalleteSyncRequest]) (*talv1.CommandPalleteSyncResponse, error) {
	params := u.adapter.CommandPalleteWeaviateFromGrpc(req.Msg.Record)
	if req.Msg.TriggerType == "DELETE" {
		err := u.weaviateClient.CommandPalleteDelete(ctx,req.Msg.Record.MenuKey)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve users list: %w", err)
		}
	}
	err := u.weaviateClient.CommandPalleteCreateUpdate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users list: %w", err)
	}
	return &talv1.CommandPalleteSyncResponse{Message: "created"}, nil
}

func (u *PublicUsecase) CommandPalleteSearch(ctx context.Context, req *connect.Request[talv1.CommandPalleteSearchRequest]) (*talv1.CommandPalleteSearchResponse, error) {
	tenantID , _ := contextkeys.TenantID(ctx)
	resp, err := u.weaviateClient.CommandPaletteSearch(ctx, tenantID , req.Msg.Query , 10)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users list: %w", err)
	}
	hits := make([]*talv1.CommandPallete , len(resp))
	for index , v := range resp {
		hits[index] =  u.adapter.CommandPalleteGrpcFromWeaviate(v)
	}

	return &talv1.CommandPalleteSearchResponse{Hits: hits}, nil
}
