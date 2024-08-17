package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"google.golang.org/api/option"
)

func main() {
	projectID := flag.String("project_id", os.Getenv("PROJECT_ID"), "Cloud Project ID")
	location := flag.String("location", os.Getenv("LOCASTION"), "The Processor location")
	// Create a Processor before running sample
	processorID := flag.String("processor_id", os.Getenv("PROCESSOR_ID"), "The Processor ID")
	filePath := flag.String("file_path", "invoice.pdf", "The path to the file to parse")
	mimeType := flag.String("mime_type", "application/pdf", "The mimeType of the file")
	flag.Parse()

	ctx := context.Background()

	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", *location)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
		return
	}
	defer client.Close()

	// Open local file.
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("os.ReadFile: %w", err))
		return
	}

	req := &documentaipb.ProcessRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", *projectID, *location, *processorID),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: *mimeType,
			},
		},
	}
	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
		return
	}

	// Handle the results.
	document := resp.GetDocument()
	fmt.Print(document.GetText())
}
