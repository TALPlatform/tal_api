package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Payload for sourcing session sync tasks
type SourcingSessionSyncPayload struct {
	SessionID         int32   `json:"session_id"`
	CrustDataProfiles *[]byte `json:"crust_data_profiles"`
	DBProfiles        *[]byte `json:"db_profiles"`
}

// Handler for processing sourcing session sync tasks
type SourcingSessionSyncHandler struct {
	syncFn func(ctx context.Context, sessionID int32, crustDataProfiles *[]byte, DBProfiles *[]byte) error
}

// Constructor
func NewSourcingSessionSyncHandler(
	fn func(ctx context.Context, sessionID int32, crustDataProfiles *[]byte, DBProfiles *[]byte) error,
) *SourcingSessionSyncHandler {
	return &SourcingSessionSyncHandler{syncFn: fn}
}

// Handle processes the task payload
func (h *SourcingSessionSyncHandler) Handle(ctx context.Context, t *asynq.Task) error {
	var payload SourcingSessionSyncPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return h.syncFn(ctx, payload.SessionID, payload.CrustDataProfiles, payload.DBProfiles)
}
