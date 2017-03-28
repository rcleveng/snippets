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

func CreateReportRequest(consumer_id string, op_name string) (*servicecontrol.ReportRequest) {
	t := time.Now()
	rfc3339 := t.Format(time.RFC3339)
	u := uuid.NewV4()
	op := servicecontrol.Operation{
		ConsumerId:    consumer_id,
		StartTime:     rfc3339,
		EndTime:       rfc3339,
		OperationId:   u.String(),
		OperationName: op_name}
	op.Labels = make(map[string]string, 7)
	op.Labels["cloud.googleapis.com/location"] = "global"
	op.Labels["serviceruntime.googleapis.com/api_version"] = "v1"
	// This is service (with underscores v dots) + version (ditto) + '.' + method
	op.Labels["serviceruntime.googleapis.com/api_method"] = "rob_mst_endpoints_rob_mst_201703_cloud_goog_1_0_0.Echo"
	op.Labels["cloud.googleapis.com/project"] = "rob-mst-201703"
	op.Labels["cloud.googleapis.com/service"] = "rob-mst.endpoints.rob-mst-201703.cloud.goog"
	op.Labels["cloud.googleapis.com/uid"] = "92830528305210394"

	v := int64(23)
	
	op.MetricValueSets = []*servicecontrol.MetricValueSet{{
	    MetricName: "serviceruntime.googleapis.com/api/consumer/request_count",
	    MetricValues: []*servicecontrol.MetricValue{{
	        Int64Value: &v,
	    }},
	}}

	// I'm sure there's a better way but everything else I seemed to try gives
	// compile errors like cannot use op (type servicecontrol.Operation) as type
	// *servicecontrol.Operation in array or slice literal.	
	ops := make([]*servicecontrol.Operation, 1);
	ops[0] = &op
	rr := servicecontrol.ReportRequest{Operations: ops}
	return &rr
}

func main() {
        sc, err := CreateService()
	if err != nil {
		log.Fatal(err)
	}

	req := CreateCheckRequest("project:rob-mst-201703", "ServiceControlTest")
	call := sc.Check("rob-mst.endpoints.rob-mst-201703.cloud.goog", req)

	rreq := CreateReportRequest("project:rob-mst-201703", "ServiceControlTest")
	callr := sc.Report("rob-mst.endpoints.rob-mst-201703.cloud.goog", rreq)

	cresponse, err := call.Do()
	if err != nil {
		log.Println("C:")
		log.Fatal(err)
	}

	rresponse, err := callr.Do()
	if err != nil {
		log.Println("R:")
		log.Fatal(err)
	}

	log.Println("C:")
	log.Println(cresponse)
	log.Println("R:")
	log.Println(rresponse)
	fmt.Printf("Success!\n")
}
