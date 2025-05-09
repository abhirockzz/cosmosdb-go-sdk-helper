package cosmosdb_errors

import (
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

func TestGetError_WithResponseError(t *testing.T) {

	err := &azcore.ResponseError{
		ErrorCode:  "NotFound",
		StatusCode: 404,
	}
	result := GetError(err)
	assert.Equal(t, err.Error(), result.Message)
	assert.Equal(t, err.StatusCode, result.Status)

}

func TestGetError_WithNonResponseError(t *testing.T) {
	plainErr := errors.New("some error")
	result := GetError(plainErr)

	assert.Equal(t, 0, result.Status)
	assert.Empty(t, result.Message)
}
