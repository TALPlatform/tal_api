package handlers

import (
	"context"
	"encoding/json"

	"github.com/TALPlatform/tal_api/pkg/crustdata"
	"github.com/hibiken/asynq"
)

type RawProfileSyncPayload struct {
	SessionID         int32                           `json:"session_id"`
	CrustDataProfiles *crustdata.PeopleSearchResponse `json:"crust_data_profiles"`
}

type RawProfileSyncHandler struct {
	rawProfileSync func(ctx context.Context, sessionId int32, crustDataProfiles *crustdata.PeopleSearchResponse) error
}

func NewRawProfileSyncHandler(fn func(ctx context.Context, sessionId int32, crustDataProfiles *crustdata.PeopleSearchResponse) error) *RawProfileSyncHandler {
	return &RawProfileSyncHandler{rawProfileSync: fn}
}

func (h *RawProfileSyncHandler) Handle(ctx context.Context, t *asynq.Task) error {
	var payload RawProfileSyncPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return h.rawProfileSync(ctx, payload.SessionID, payload.CrustDataProfiles)
}
