package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/servicecontrol/v1"
	"log"
	"time"
	"github.com/satori/go.uuid"
)

func CreateService() (*servicecontrol.ServicesService, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, servicecontrol.ServicecontrolScope)
	if err != nil {
		log.Fatal(err)
	}

	service, err := servicecontrol.New(client)

	if err != nil {
		return nil, err
	}

	s := servicecontrol.NewServicesService(service)
	return s, nil

}

func CreateCheckRequest(consumer_id string, op_name string) (*servicecontrol.CheckRequest) {
	t := time.Now()
	rfc3339 := t.Format(time.RFC3339)
	u := uuid.NewV4()
	op := servicecontrol.Operation{
		ConsumerId:    consumer_id,
		StartTime:     rfc3339,
		OperationId:   u.String(),
		OperationName: op_name}
	return &servicecontrol.CheckRequest{Operation: &op}
}

func main() {
        sc, err := CreateService()
	if err != nil {
		log.Fatal(err)
	}

	req := CreateCheckRequest("project:project:rob-mst-201703", "ServiceControlTest")
	call := sc.Check("rob-mst.endpoints.rob-mst-201703.cloud.goog", req)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
	fmt.Printf("Success!\n")
}
