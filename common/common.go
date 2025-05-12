package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/cosmosdb_errors"
)

// CreateDatabaseIfNotExists returns a DatabaseClient for the given database, creating the database if it does not exist.
// This ensures idempotent database creation and simplifies setup for Cosmos DB resources.
func CreateDatabaseIfNotExists(client *azcosmos.Client, props azcosmos.DatabaseProperties, opts *azcosmos.CreateDatabaseOptions) (*azcosmos.DatabaseClient, error) {
	db, err := client.NewDatabase(props.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %v", err)
	}

	_, err = db.Read(context.Background(), nil)
	if err != nil {
		if cosmosdb_errors.GetError(err).Status == http.StatusNotFound {
			// Database doesn't exist, try to create it
			_, err = client.CreateDatabase(context.Background(), props, opts)
			if err != nil {
				cosmosErr := cosmosdb_errors.GetError(err)
				if cosmosErr.Status == http.StatusConflict {
					// Database was created by another process, treat as success
					return client.NewDatabase(props.ID)
				}
				return nil, fmt.Errorf("failed to create database: %v", err)
			}
			return client.NewDatabase(props.ID)
		}
		return nil, fmt.Errorf("failed to read database: %v", err)
	}

	return db, nil
}

// CreateContainerIfNotExists returns a ContainerClient for the given container, creating the container if it does not exist.
// This is useful for idempotent container setup in Cosmos DB databases.
func CreateContainerIfNotExists(db *azcosmos.DatabaseClient, props azcosmos.ContainerProperties, opts *azcosmos.CreateContainerOptions) (*azcosmos.ContainerClient, error) {
	//func CreateContainerIfNotExists(db *azcosmos.DatabaseClient, containerName string) (*azcosmos.ContainerClient, error) {
	container, err := db.NewContainer(props.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client: %v", err)
	}

	_, err = container.Read(context.Background(), nil)
	if err != nil {
		if cosmosdb_errors.GetError(err).Status == http.StatusNotFound {
			// Container doesn't exist, try to create it
			_, err = db.CreateContainer(context.Background(), props, opts)
			if err != nil {
				cosmosErr := cosmosdb_errors.GetError(err)
				if cosmosErr.Status == http.StatusConflict {
					// Container was created by another process, treat as success
					return db.NewContainer(props.ID)
				}
				return nil, fmt.Errorf("failed to create container: %v", err)
			}
			return db.NewContainer(props.ID)
		}
		return nil, fmt.Errorf("failed to read container: %v", err)
	}

	return container, nil
}

// GetAllDatabases retrieves all database properties in the Cosmos DB account.
// Use this to enumerate or inspect all databases in the account.
func GetAllDatabases(client *azcosmos.Client) ([]azcosmos.DatabaseProperties, error) {
	pager := client.NewQueryDatabasesPager("select * from c", nil)
	var databases []azcosmos.DatabaseProperties
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get databases: %v", err)
		}
		databases = append(databases, page.Databases...)
	}
	return databases, nil
}

// GetAllContainers retrieves all container properties in the specified database.
// Use this to enumerate or inspect all containers in a database.
func GetAllContainers(client *azcosmos.DatabaseClient) ([]azcosmos.ContainerProperties, error) {
	pager := client.NewQueryContainersPager("select * from c", nil)
	var containers []azcosmos.ContainerProperties
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get containers: %v", err)
		}
		containers = append(containers, page.Containers...)
	}
	return containers, nil
}

// InsertItemWithResponse inserts an item into the specified container and returns the inserted item.
func InsertItemWithResponse[T any](container *azcosmos.ContainerClient, item T, partitionKey azcosmos.PartitionKey, opts *azcosmos.ItemOptions) (T, error) {
	if opts == nil {
		opts = &azcosmos.ItemOptions{}
	}
	opts.EnableContentResponseOnWrite = true

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return item, err
	}

	response, err := container.CreateItem(context.Background(), partitionKey, itemBytes, opts)
	if err != nil {
		return item, err
	}

	if err := json.Unmarshal(response.Value, &item); err != nil {
		return item, err
	}

	return item, nil
}
