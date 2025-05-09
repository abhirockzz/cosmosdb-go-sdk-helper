package cosmosdb_errors

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type CosmosDBError struct {
	Message string
	Status  int
}

func GetError(err error) CosmosDBError {

	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return CosmosDBError{}
	}

	return CosmosDBError{Message: err.Error(), Status: respErr.StatusCode}
}
