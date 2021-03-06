package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverstop = new(Server)

func TestStopService(t *testing.T) {
	deployment, _ := serverstop.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service.Service{
			Name: "TestStopService",
			Dependencies: map[string]*service.Dependency{
				"test": {
					Image: "nginx",
				},
			},
		},
	})
	s, _ := services.Get(deployment.ServiceID)
	s.Start()
	reply, err := serverstop.StopService(context.Background(), &StopServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	status, _ := s.Status()
	require.Equal(t, service.STOPPED, status)
	require.Nil(t, err)
	require.NotNil(t, reply)
	services.Delete(deployment.ServiceID)
}
