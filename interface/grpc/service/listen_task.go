package service

import (
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenTask creates a stream that will send data for every task to execute.
func (s *Server) ListenTask(request *ListenTaskRequest, stream Service_ListenTaskServer) error {
	service, err := services.Get(request.Token)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	channel := service.TaskSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-subscription:
			execution := data.(*execution.Execution)
			inputs, _ := json.Marshal(execution.Inputs)
			if err := stream.Send(&TaskData{
				ExecutionID: execution.ID,
				TaskKey:     execution.Task,
				InputData:   string(inputs),
			}); err != nil {
				return err
			}
		}
	}
}
