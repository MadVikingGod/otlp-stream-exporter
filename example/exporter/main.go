package main

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	"log"
	"os"
)

func setupExport() coltracepb.TraceServiceClient {
	cc, err := grpc.Dial("localhost:55680", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to jaeger: %v", err)
	}
	return coltracepb.NewTraceServiceClient(cc)
}

func main() {
	client := setupExport()

	traceFile, err := os.Open("trace.json")
	if err != nil {
		log.Fatalf("Could not open trace file: %v\n", err)
	}
	dec := json.NewDecoder(traceFile)

	for dec.More() {
		pbRequest := &coltracepb.ExportTraceServiceRequest{}
		jsonpb.UnmarshalNext(dec, pbRequest)
		_, err := client.Export(context.Background(), pbRequest)
		if err != nil {
			log.Fatalf("jaeger returned an error: %v", err)
		}
	}
}
