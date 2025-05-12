package trigger

import "encoding/json"

// Parse unmarshals the Cosmos DB trigger payload and extracts the documents.
// It performs a two-step unmarshaling process due to the nested JSON structure.
func ParseToCosmosDBDataMap(functionsTriggerPayload []byte) ([]map[string]any, error) {

	// First unmarshal step: convert the Documents field from string to []byte
	documentsRaw, err := ParseToRawString(functionsTriggerPayload)
	if err != nil {
		return nil, err
	}

	// Second unmarshal step: convert the JSON string to []map[string]any
	var documents []map[string]any
	if err := json.Unmarshal([]byte(documentsRaw), &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// ParseFunctionsTriggerPayload unmarshals the Cosmos DB trigger payload and extracts the documents in raw string format.
// It partially unmarshals the payload to get the Documents field, which is a string containing the JSON representation of the documents.
// You are free to use this to parse the raw string in your own way.
func ParseToRawString(functionsTriggerPayload []byte) (string, error) {
	var triggerPayload CosmosDBTriggerPayload
	if err := json.Unmarshal(functionsTriggerPayload, &triggerPayload); err != nil {
		return "", err
	}

	// The data.documents field is a string that contains the JSON representation of the documents
	// Convert the Documents field from string to raw string that represents Cosmos DB documents
	var documentsRaw string
	if err := json.Unmarshal([]byte(triggerPayload.Data.Documents), &documentsRaw); err != nil {
		return "", err
	}
	return documentsRaw, nil
}
