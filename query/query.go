package query

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// QueryItems executes a SQL query and returns items of type T
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
