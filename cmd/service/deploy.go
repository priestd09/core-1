package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/docker/docker/pkg/archive"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/spinner"
	"github.com/spf13/cobra"
)

// Deploy a service to the marketplace.
var Deploy = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"publish"},
	Short:   "Deploy a service",
	Long: `Deploy a service.

To get more information, see the [deploy page from the documentation](https://docs.mesg.com/guide/service/deploy-a-service.html)`,
	Example:           `mesg-core service deploy PATH_TO_SERVICE`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	path := defaultPath(args)
	serviceID, isValid, err := deployService(path)
	if !isValid {
		os.Exit(1)
	}
	utils.HandleError(err)

	fmt.Println("Service deployed with ID:", aurora.Green(serviceID))
	fmt.Printf("To start it, run the command: mesg-core service start %s\n", serviceID)

}

func deployService(path string) (serviceID string, isValid bool, err error) {
	stream, err := cli().DeployService(context.Background())
	if err != nil {
		return "", false, err
	}

	deployment := make(chan deploymentResult)
	go readDeployReply(stream, deployment)

	if govalidator.IsURL(path) {
		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Url{Url: path},
		}); err != nil {
			return "", false, err
		}
	} else {
		if err := deployServiceSendServiceContext(path, stream); err != nil {
			return "", false, err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return "", false, err
	}

	result := <-deployment
	return result.serviceID, result.isValid, result.err
}

func deployServiceSendServiceContext(path string, stream core.Core_DeployServiceClient) error {
	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := archive.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Chunk{Chunk: buf[:n]},
		}); err != nil {
			return err
		}
	}

	return nil
}

type deploymentResult struct {
	serviceID string
	err       error
	isValid   bool
}

func readDeployReply(stream core.Core_DeployServiceClient, deployment chan deploymentResult) {
	sp := spinner.New(utils.SpinnerCharset, utils.SpinnerDuration)
	result := deploymentResult{isValid: true}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			deployment <- deploymentResult{err: err}
			return
		}

		var (
			status          = message.GetStatus()
			serviceID       = message.GetServiceID()
			validationError = message.GetValidationError()
		)

		switch {
		case status != "":
			sp.Start()
			sp.Suffix = " " + status
		case serviceID != "":
			sp.Stop()
			result.serviceID = serviceID
			deployment <- result
			return
		case validationError != "":
			sp.Stop()
			fmt.Println(aurora.Red(validationError))
			fmt.Println("Run the command 'service validate' for more details")
			result.isValid = false
			deployment <- result
			return
		}
	}
}
