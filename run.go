package main

import (
	"context"
	"strings"

	run "cloud.google.com/go/run/apiv2"

	runpb "cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/iterator"
)

var (
	grunServices = []string{}
)

func getGrunServices() []string {
	ctx := context.Background()
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := run.NewServicesClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	// https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#ListServicesRequest
	req := &runpb.ListServicesRequest{
		Parent:      "projects/xrdm-idp/locations/us-central1",
		PageSize:    0,
		ShowDeleted: false,
	}
	it := c.ListServices(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
		// The output of resp.Name is:
		// projects/xrdm-idp/locations/us-central1/services/idp-api-graphiql-service-production-f78f84f
		name := strings.Split(resp.Name, "/")
		grunServices = append(grunServices, name[len(name)-1])
	}
	return grunServices
}
