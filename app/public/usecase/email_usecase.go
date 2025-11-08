package usecase

import (
	"context"

	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (s *PublicUsecase) EmailSend(ctx context.Context, req *talv1.EmailSendRequest) (*talv1.EmailSendResponse, error) {
	params := s.adapter.EmailSendResendFromGrpc(req)
	resp, err := s.resendClient.SendEmail(&params)
	if err != nil {
		return nil, err
	}
	return &talv1.EmailSendResponse{
		Id: resp.Id,
	}, nil

}
