// Package query provides strongly typed helper functions for working with Azure Cosmos DB queries.
package query

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/query/metrics"
)

// QueryItems executes a SQL query against a Cosmos DB container and returns strongly typed results.
// Returns a slice of unmarshaled items of type T or an error if the query fails.
func QueryItems[T any](container *azcosmos.ContainerClient, query string, partitionKey azcosmos.PartitionKey, opts *azcosmos.QueryOptions) ([]T, error) {

	var items []T
	queryPager := container.NewQueryItemsPager(query, partitionKey, opts)

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		// Process each item in the page
		for _, item := range queryResponse.Items {
			var typedItem T
			if err := json.Unmarshal(item, &typedItem); err != nil {
				return nil, err
			}
			items = append(items, typedItem)
		}
	}

	return items, nil
}

// QueryItem retrieves a single item from a Cosmos DB container
// Returns the unmarshaled item of type T or an error if the item cannot be retrieved or unmarshaled.
func QueryItem[T any](container *azcosmos.ContainerClient, itemID string, partitionKey azcosmos.PartitionKey, opts *azcosmos.ItemOptions) (T, error) {

	var typedItem T

	response, err := container.ReadItem(context.Background(), partitionKey, itemID, opts)
	if err != nil {
		return typedItem, err
	}

	if err := json.Unmarshal(response.Value, &typedItem); err != nil {
		return typedItem, err
	}

	return typedItem, nil
}

type QueryResult[T any] struct {
	Items         []T
	Metrics       []metrics.QueryMetrics // per page
	RequestCharge float64                // total for all pages
}

func QueryItemsWithMetrics[T any](container *azcosmos.ContainerClient, query string, partitionKey azcosmos.PartitionKey, opts *azcosmos.QueryOptions) (QueryResult[T], error) {
	if opts == nil {
		opts = &azcosmos.QueryOptions{}
	}
	// Enable metrics collection
	opts.PopulateIndexMetrics = true

	var items []T
	var metricsList []metrics.QueryMetrics
	totalRequestCharge := 0.0

	queryPager := container.NewQueryItemsPager(query, partitionKey, opts)

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())
		if err != nil {
			return QueryResult[T]{}, err
		}

		// Process each item in the page
		for _, item := range queryResponse.Items {
			var typedItem T
			if err := json.Unmarshal(item, &typedItem); err != nil {
				return QueryResult[T]{}, err
			}
			items = append(items, typedItem)
		}

		if queryResponse.QueryMetrics != nil {
			qm, err := metrics.ParseQueryMetrics(*queryResponse.QueryMetrics)
			if err != nil {
				return QueryResult[T]{}, err
			}
			metricsList = append(metricsList, qm)
		}

		totalRequestCharge += float64(queryResponse.RequestCharge)

	}

	return QueryResult[T]{Items: items, Metrics: metricsList, RequestCharge: totalRequestCharge}, nil
}
