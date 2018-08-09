package service

import (
	"sort"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestSort(t *testing.T) {
	status := []serviceStatus{
		{status: service.PARTIAL, service: &service.Service{Name: "Partial"}},
		{status: service.RUNNING, service: &service.Service{Name: "Running"}},
		{status: service.STOPPED, service: &service.Service{Name: "Stopped"}},
	}
	sort.Sort(byStatus(status))
	assert.Equal(t, status[0].status, service.RUNNING)
	assert.Equal(t, status[1].status, service.PARTIAL)
	assert.Equal(t, status[2].status, service.STOPPED)
}

func TestServicesWithStatus(t *testing.T) {
	services := append([]*service.Service{}, &service.Service{Name: "TestServicesWithStatus"})
	status, err := servicesWithStatus(services)
	assert.Nil(t, err)
	assert.Equal(t, status[0].status, service.STOPPED)
}
