package utils

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cannotReachTheCore = "Cannot reach the Core"
	startCore          = "Please start the core by running: mesg-core start"
	cannotReachDocker  = "Cannot reach Docker"
	installDocker      = "Please make sure Docker is running.\nIf Docker is not installed on your machine you can install it here: https://store.docker.com/search?type=edition&offering=community"
)

// HandleError displays the error and stops the process if the error exists
func HandleError(err error) {
	if err != nil {
		fmt.Println(errorMessage(err))
		os.Exit(1)
	}
}

func errorMessage(err error) string {
	switch {
	case coreConnectionError(err):
		return fmt.Sprintf("%s\n%s", cannotReachTheCore, startCore)
	case dockerDaemonError(err):
		return fmt.Sprintf("%s\n%s", cannotReachDocker, installDocker)
	default:
		return err.Error()
	}
}

func coreConnectionError(err error) bool {
	s := status.Convert(err)
	return s.Code() == codes.Unavailable
}

func dockerDaemonError(err error) bool {
	return client.IsErrConnectionFailed(err)
}
