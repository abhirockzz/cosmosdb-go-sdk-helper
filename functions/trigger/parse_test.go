package trigger

import "testing"

func TestParseToCosmosDBDataMap(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"dfa26d32-f876-44a3-b107-369f1f48c689\\\",\\\"description\\\":\\\"Setup monitoring\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCUAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCUAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f007efc-0000-0800-0000-67f5fb920000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744173970,\\\"_lsn\\\":160}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:46:10.723203Z","RandGuid":"0d00378b-6426-4af1-9fc0-0793f4ce3745"}}}`

	result, err := ParseToCosmosDBDataMap([]byte(payload))
	if err != nil {
		t.Fatalf("parse function returned an error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 document, got %d", len(result))
	}

	doc := result[0]
	if doc["id"] != "dfa26d32-f876-44a3-b107-369f1f48c689" {
		t.Errorf("expected id to be 'dfa26d32-f876-44a3-b107-369f1f48c689', got %v", doc["id"])
	}

	if doc["description"] != "Setup monitoring" {
		t.Errorf("expected description to be 'Setup monitoring', got %v", doc["description"])
	}
}

func TestParseToCosmosDBDataMapForMultipleDocuments(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"51e0c1b0-87d3-4611-ac41-7ac3e77d9920\\\",\\\"description\\\":\\\"Schedule team meeting\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCVAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCVAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a3fd-0000-0800-0000-67f5fc640000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174180,\\\"_lsn\\\":161},{\\\"id\\\":\\\"cfbf42b9-48e8-449b-9cff-17c6fbd00f83\\\",\\\"description\\\":\\\"Update dependencies\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCWAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCWAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a9fd-0000-0800-0000-67f5fc670000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174183,\\\"_lsn\\\":162}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:49:45.157601Z","RandGuid":"304980d9-584d-4323-98e3-b46bb1eebded"}}}`

	result, err := ParseToCosmosDBDataMap([]byte(payload))
	if err != nil {
		t.Fatalf("parse function returned an error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(result))
	}

	doc1 := result[0]
	if doc1["id"] != "51e0c1b0-87d3-4611-ac41-7ac3e77d9920" {
		t.Errorf("expected id of first document to be '51e0c1b0-87d3-4611-ac41-7ac3e77d9920', got %v", doc1["id"])
	}

	if doc1["description"] != "Schedule team meeting" {
		t.Errorf("expected description of first document to be 'Schedule team meeting', got %v", doc1["description"])
	}

	doc2 := result[1]
	if doc2["id"] != "cfbf42b9-48e8-449b-9cff-17c6fbd00f83" {
		t.Errorf("expected id of second document to be 'cfbf42b9-48e8-449b-9cff-17c6fbd00f83', got %v", doc2["id"])
	}

	if doc2["description"] != "Update dependencies" {
		t.Errorf("expected description of second document to be 'Update dependencies', got %v", doc2["description"])
	}
}

func TestParseToRawString(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"test-id\\\",\\\"description\\\":\\\"test-description\\\"}]\""},"Metadata":{"sys":{"MethodName":"test-method","UtcNow":"2025-05-12T00:00:00Z","RandGuid":"test-guid"}}}`
	expectedDocumentsRaw := `[{"id":"test-id","description":"test-description"}]`

	result, err := ParseToRawString([]byte(payload))
	if err != nil {
		t.Fatalf("ParseToRawString returned an error: %v", err)
	}

	if result != expectedDocumentsRaw {
		t.Errorf("expected documentsRaw to be '%s', got '%s'", expectedDocumentsRaw, result)
	}
}

func TestParseToRawStringForMultipleDocuments(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"test-id-1\\\",\\\"description\\\":\\\"test-desc-1\\\"},{\\\"id\\\":\\\"test-id-2\\\",\\\"description\\\":\\\"test-desc-2\\\"}]\""},"Metadata":{"sys":{"MethodName":"test-method","UtcNow":"2025-05-12T00:00:00Z","RandGuid":"test-guid"}}}`
	expectedDocumentsRaw := `[{"id":"test-id-1","description":"test-desc-1"},{"id":"test-id-2","description":"test-desc-2"}]`

	result, err := ParseToRawString([]byte(payload))
	if err != nil {
		t.Fatalf("ParseToRawStringForMultipleDocuments returned an error: %v", err)
	}

	if result != expectedDocumentsRaw {
		t.Errorf("expected documentsRaw to be '%s', got '%s'", expectedDocumentsRaw, result)
	}
}
