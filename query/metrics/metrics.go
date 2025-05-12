package metrics

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type QueryMetrics struct {
	TotalExecutionTimeInMs         float64
	QueryCompileTimeInMs           float64
	QueryLogicalPlanBuildTimeInMs  float64
	QueryPhysicalPlanBuildTimeInMs float64
	QueryOptimizationTimeInMs      float64
	VMExecutionTimeInMs            float64
	IndexLookupTimeInMs            float64
	InstructionCount               int
	DocumentLoadTimeInMs           float64
	SystemFunctionExecuteTimeInMs  float64
	UserFunctionExecuteTimeInMs    float64
	RetrievedDocumentCount         int
	RetrievedDocumentSize          int
	OutputDocumentCount            int
	OutputDocumentSize             int
	WriteOutputTimeInMs            float64
	IndexUtilizationRatio          float64
}

func ParseQueryMetrics(metrics string) (QueryMetrics, error) {

	var qm QueryMetrics
	pairs := strings.Split(metrics, ";")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		switch key {
		case "totalExecutionTimeInMs":
			qm.TotalExecutionTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "queryCompileTimeInMs":
			qm.QueryCompileTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "queryLogicalPlanBuildTimeInMs":
			qm.QueryLogicalPlanBuildTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "queryPhysicalPlanBuildTimeInMs":
			qm.QueryPhysicalPlanBuildTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "queryOptimizationTimeInMs":
			qm.QueryOptimizationTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "VMExecutionTimeInMs":
			qm.VMExecutionTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "indexLookupTimeInMs":
			qm.IndexLookupTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "instructionCount":
			qm.InstructionCount, _ = strconv.Atoi(value)
		case "documentLoadTimeInMs":
			qm.DocumentLoadTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "systemFunctionExecuteTimeInMs":
			qm.SystemFunctionExecuteTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "userFunctionExecuteTimeInMs":
			qm.UserFunctionExecuteTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "retrievedDocumentCount":
			qm.RetrievedDocumentCount, _ = strconv.Atoi(value)
		case "retrievedDocumentSize":
			qm.RetrievedDocumentSize, _ = strconv.Atoi(value)
		case "outputDocumentCount":
			qm.OutputDocumentCount, _ = strconv.Atoi(value)
		case "outputDocumentSize":
			qm.OutputDocumentSize, _ = strconv.Atoi(value)
		case "writeOutputTimeInMs":
			qm.WriteOutputTimeInMs, _ = strconv.ParseFloat(value, 64)
		case "indexUtilizationRatio":
			qm.IndexUtilizationRatio, _ = strconv.ParseFloat(value, 64)
		}
	}
	return qm, nil
}

func GetIndexMetrics(response azcosmos.QueryItemsResponse) (string, error) {
	if response.IndexMetrics == nil {
		return "", nil
	}
	// The index metrics are base64 encoded in the response
	// and need to be decoded before they can be used
	decoded, err := base64.StdEncoding.DecodeString(*response.IndexMetrics)

	if err != nil {
		//fmt.Println("Failed to decode index metrics:", err)
		return "", err
	}

	return string(decoded), nil
}
