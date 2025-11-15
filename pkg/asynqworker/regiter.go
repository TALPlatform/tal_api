package asynqworker

import (
	"context"

	"github.com/TALPlatform/tal_api/pkg/asynqworker/handlers"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
)

// RegisterHandlers registers all asynq handlers into the worker.
func (w *Worker) RegisterHandlers(
	rawProfileSyncFn func(ctx context.Context, sessionId int32, crustDataProfiles *crustdata.PeopleSearchResponse) error,
	sourcingSessionSyncFn func(ctx context.Context, sessionID int32, crustDataProfiles *[]byte, DBProfiles *[]byte) error,

) {
	rawProfileHandler := handlers.NewRawProfileSyncHandler(rawProfileSyncFn)
	w.Register(TaskRawProfileSync, rawProfileHandler.Handle)
	sourcingSessionHandler := handlers.NewSourcingSessionSyncHandler(sourcingSessionSyncFn)
	w.Register(TaskSourcingSessionSync, sourcingSessionHandler.Handle)
}
