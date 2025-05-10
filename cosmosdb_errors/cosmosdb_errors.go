package cosmosdb_errors

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// CosmosDBError represents a structured Cosmos DB error with message and status code.
type CosmosDBError struct {
	Message string
	Status  int
}

// GetError extracts a CosmosDBError from a generic error, if possible.
// Returns an empty CosmosDBError if the error is not a Cosmos DB response error.
func GetError(err error) CosmosDBError {

	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return CosmosDBError{}
	}

	return CosmosDBError{Message: err.Error(), Status: respErr.StatusCode}
}
