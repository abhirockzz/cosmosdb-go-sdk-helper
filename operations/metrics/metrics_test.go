package metrics

import (
	"encoding/base64"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stretchr/testify/assert"
)

func TestParseQueryMetrics(t *testing.T) {
	input := "totalExecutionTimeInMs=12.5;queryCompileTimeInMs=1.2;queryLogicalPlanBuildTimeInMs=0.5;queryPhysicalPlanBuildTimeInMs=0.3;queryOptimizationTimeInMs=0.7;VMExecutionTimeInMs=10.1;indexLookupTimeInMs=0.2;instructionCount=100;documentLoadTimeInMs=0.4;systemFunctionExecuteTimeInMs=0.1;userFunctionExecuteTimeInMs=0.2;retrievedDocumentCount=5;retrievedDocumentSize=500;outputDocumentCount=3;outputDocumentSize=300;writeOutputTimeInMs=0.6;indexUtilizationRatio=0.95;"
	expected := QueryMetrics{
		TotalExecutionTimeInMs:         12.5,
		QueryCompileTimeInMs:           1.2,
		QueryLogicalPlanBuildTimeInMs:  0.5,
		QueryPhysicalPlanBuildTimeInMs: 0.3,
		QueryOptimizationTimeInMs:      0.7,
		VMExecutionTimeInMs:            10.1,
		IndexLookupTimeInMs:            0.2,
		InstructionCount:               100,
		DocumentLoadTimeInMs:           0.4,
		SystemFunctionExecuteTimeInMs:  0.1,
		UserFunctionExecuteTimeInMs:    0.2,
		RetrievedDocumentCount:         5,
		RetrievedDocumentSize:          500,
		OutputDocumentCount:            3,
		OutputDocumentSize:             300,
		WriteOutputTimeInMs:            0.6,
		IndexUtilizationRatio:          0.95,
	}
	qm, err := ParseQueryMetrics(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, qm)
}

func TestParseQueryMetrics_EmptyString(t *testing.T) {
	qm, err := ParseQueryMetrics("")
	assert.NoError(t, err)
	zero := QueryMetrics{}
	assert.Equal(t, zero, qm)
}

func TestParseQueryMetrics_InvalidPairs(t *testing.T) {
	input := "foo=bar;totalExecutionTimeInMs=5.5;badpair;outputDocumentCount=2"
	qm, err := ParseQueryMetrics(input)
	assert.NoError(t, err)
	assert.Equal(t, 5.5, qm.TotalExecutionTimeInMs)
	assert.Equal(t, 2, qm.OutputDocumentCount)
}

func TestGetIndexMetrics_Base64Decoded(t *testing.T) {
	jsonStr := `{"foo":"bar","baz":123}`
	encoded := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	resp := azcosmos.QueryItemsResponse{
		IndexMetrics: &encoded,
	}
	decoded, err := GetIndexMetrics(resp)
	assert.NoError(t, err)
	assert.Equal(t, jsonStr, decoded)
}

func TestGetIndexMetrics_NilIndexMetrics(t *testing.T) {
	resp := azcosmos.QueryItemsResponse{
		IndexMetrics: nil,
	}
	decoded, err := GetIndexMetrics(resp)
	assert.NoError(t, err)
	assert.Equal(t, "", decoded)
}

func TestGetIndexMetrics_InvalidBase64(t *testing.T) {
	invalid := "not-base64!"
	resp := azcosmos.QueryItemsResponse{
		IndexMetrics: &invalid,
	}
	_, err := GetIndexMetrics(resp)
	assert.Error(t, err)
}
