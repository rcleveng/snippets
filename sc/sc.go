package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/servicecontrol/v1"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, servicecontrol.ServicecontrolScope)
	if err != nil {
		log.Fatal(err)
	}

	service, err := servicecontrol.New(client)

	if err != nil {
		log.Fatal(err)
	}

	sc := servicecontrol.NewServicesService(service)

	t := time.Now()
	rfc3339 := t.Format(time.RFC3339)
	op := servicecontrol.Operation{
		ConsumerId:    "project:project:rob-mst-201703",
		StartTime:     rfc3339,
		OperationId:   "8356d3c5-f9b5-4274-b4f9-079a3731e611",
		OperationName: "ServiceControlTest"}
	req := servicecontrol.CheckRequest{Operation: &op}
	call := sc.Check("rob-mst.endpoints.rob-mst-201703.cloud.goog", &req)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
	fmt.Printf("Success!\n")
}
