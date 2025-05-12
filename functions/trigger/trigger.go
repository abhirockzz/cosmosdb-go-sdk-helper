// Package trigger provides shared functionality for processing Cosmos DB documents.
package trigger

// CosmosDBTriggerPayload represents the structure of the Cosmos DB trigger payload.
type CosmosDBTriggerPayload struct {
	Data     Data     `json:"Data"`
	Metadata Metadata `json:"Metadata"`
}

// Data represents the data field in the Cosmos DB trigger payload.
type Data struct {
	Documents string `json:"documents"`
}

type SysMetadata struct {
	MethodName string `json:"MethodName"`
	UtcNow     string `json:"UtcNow"`
	RandGuid   string `json:"RandGuid"`
}

type Metadata struct {
	Sys SysMetadata `json:"sys"`
}

// InvokeResponse represents the structure of the response returned by the handler.
type InvokeResponse struct {
	Outputs     map[string]any `json:"outputs"`
	Logs        []string       `json:"logs"`
	ReturnValue any            `json:"returnValue,omitempty"`
}
