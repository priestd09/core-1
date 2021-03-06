package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestExtractPortEmpty(t *testing.T) {
	dep := Dependency{}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	dep := &Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 2)
	require.Equal(t, ports[0].Target, uint32(80))
	require.Equal(t, ports[0].Published, uint32(80))
	require.Equal(t, ports[1].Target, uint32(8080))
	require.Equal(t, ports[1].Published, uint32(3000))
}

func TestStartService(t *testing.T) {
	service := &Service{
		Name: "TestStartService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(service.GetDependencies()), len(dockerServices))
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartWith2Dependencies(t *testing.T) {
	service := &Service{
		Name: "TestStartWith2Dependencies",
		Dependencies: map[string]*Dependency{
			"testa": {
				Image: "nginx:latest",
			},
			"testb": {
				Image: "alpine:latest",
			},
		},
	}
	servicesID, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.DependenciesFromService()
	container1, _ := defaultContainer.FindContainer(deps[0].namespace())
	container2, _ := defaultContainer.FindContainer(deps[1].namespace())
	require.Equal(t, "nginx:latest", container1.Config.Image)
	require.Equal(t, "alpine:latest", container2.Config.Image)
}

func TestStartAgainService(t *testing.T) {
	service := &Service{
		Name: "TestStartAgainService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	defer service.Stop()
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestPartiallyRunningService(t *testing.T) {
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]*Dependency{
			"testa": {
				Image: "nginx",
			},
			"testb": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	defer service.Stop()
	service.DependenciesFromService()[0].Stop()
	status, _ := service.Status()
	require.Equal(t, PARTIAL, status)
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.GetDependencies()))
	status, _ = service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartDependency(t *testing.T) {
	service := &Service{
		Name: "TestStartDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	defer defaultContainer.DeleteNetwork(service.namespace())
	dep := service.DependenciesFromService()[0]
	serviceID, err := dep.Start(networkID)
	defer dep.Stop()
	require.Nil(t, err)
	require.NotEqual(t, "", serviceID)
	status, _ := dep.Status()
	require.Equal(t, container.RUNNING, status)
}

func TestNetworkCreated(t *testing.T) {
	service := &Service{
		Name: "TestNetworkCreated",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	defer service.Stop()
	network, err := defaultContainer.FindNetwork(service.namespace())
	require.Nil(t, err)
	require.NotEqual(t, "", network.ID)
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestStartStopStart(t *testing.T) {
	service := &Service{
		Name: "TestStartStopStart",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	service.Stop()
	time.Sleep(10 * time.Second)
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 1)
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}
