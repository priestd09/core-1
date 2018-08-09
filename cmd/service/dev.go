package service

import (
	"context"
	"fmt"
	"log"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Dev command will run the service from a path in dev mode
// It will also listen for all events and outputs from the tasks
var Dev = &cobra.Command{
	Use:               "dev",
	Short:             "Run your service in development mode",
	Example:           "mesg-core service dev PATH",
	Run:               devHandler,
	DisableAutoGenTag: true,
}

func init() {
	Dev.Flags().StringP("event-filter", "e", "*", "Only log the data of the given event")
	Dev.Flags().StringP("task-filter", "t", "", "Only log the result of the given task")
	Dev.Flags().StringP("output-filter", "o", "", "Only log the data of the given output of a task result. If set, you also need to set the task in --task-filter")
}

func devHandler(cmd *cobra.Command, args []string) {
	serviceID, err := createService(defaultPath(args))
	utils.HandleError(err)
	fmt.Printf("✔ Service started with success\n")
	fmt.Printf("Service ID: %s\n", serviceID)

	go listenEvents(serviceID, cmd.Flag("event-filter").Value.String())
	go listenResults(serviceID, cmd.Flag("task-filter").Value.String(), cmd.Flag("output-filter").Value.String())

	closeReaders := showLogs(serviceID, "*")
	defer closeReaders()

	<-utils.WaitForCancel()

	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Deleting test service..."}, func() {
		cli().DeleteService(context.Background(), &core.DeleteServiceRequest{ // Delete service. This will automatically stop the service too
			ServiceID: serviceID,
		})
	})
}

func createService(path string) (string, error) {
	service := prepareService(path)
	deployment, err := cli().DeployService(context.Background(), &core.DeployServiceRequest{
		Service: service,
	})
	if err != nil {
		return "", err
	}
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Starting service..."}, func() {
		_, err = cli().StartService(context.Background(), &core.StartServiceRequest{
			ServiceID: deployment.ServiceID,
		})
	})
	return deployment.ServiceID, err
}

func listenEvents(serviceID string, filter string) {
	stream, err := cli().ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID:   serviceID,
		EventFilter: filter,
	})
	utils.HandleError(err)
	fmt.Println("Listening for events from the service...")
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Receive event", event.EventKey, ":", event.EventData)
	}
}

func listenResults(serviceID string, result string, output string) {
	stream, err := cli().ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:    serviceID,
		TaskFilter:   result,
		OutputFilter: output,
	})
	utils.HandleError(err)
	fmt.Println("Listening for results from the service...")
	for {
		result, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Receive result", result.TaskKey, result.OutputKey, "with data", result.OutputData)
	}
}
