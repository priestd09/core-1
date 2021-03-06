package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	hash, _ := Save(&service.Service{Name: "Service1"})
	defer Delete(hash)
	services, err := All()
	founded := false
	for _, s := range services {
		if s.Name == "Service1" {
			founded = true
			break
		}
	}
	require.Nil(t, err)
	require.True(t, founded)
}
