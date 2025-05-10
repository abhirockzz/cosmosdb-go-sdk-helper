package common

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/cosmosdb_errors"
)

// CreateDatabaseIfNotExists returns a DatabaseClient for the given database, creating the database if it does not exist.
// This ensures idempotent database creation and simplifies setup for Cosmos DB resources.
func CreateDatabaseIfNotExists(client *azcosmos.Client, dbName string) (*azcosmos.DatabaseClient, error) {
	db, err := client.NewDatabase(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %v", err)
	}

	_, err = db.Read(context.Background(), nil)
	if err != nil {
		if cosmosdb_errors.GetError(err).Status == http.StatusNotFound {
			// Database doesn't exist, try to create it
			_, err = client.CreateDatabase(context.Background(), azcosmos.DatabaseProperties{
				ID: dbName,
			}, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create database: %v", err)
			}
			return db, nil
		}
		return nil, fmt.Errorf("failed to read database: %v", err)
	}

	return db, nil
}

// CreateContainerIfNotExists returns a ContainerClient for the given container, creating the container if it does not exist.
// This is useful for idempotent container setup in Cosmos DB databases.
func CreateContainerIfNotExists(db *azcosmos.DatabaseClient, containerName string) (*azcosmos.ContainerClient, error) {
	container, err := db.NewContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client: %v", err)
	}

	_, err = container.Read(context.Background(), nil)
	if err != nil {
		if cosmosdb_errors.GetError(err).Status == http.StatusNotFound {
			// Container doesn't exist, try to create it
			_, err = db.CreateContainer(context.Background(), azcosmos.ContainerProperties{
				ID: containerName,
				PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
					Paths: []string{"/partitionKey"},
					Kind:  azcosmos.PartitionKeyKindHash,
				},
			}, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create container: %v", err)
			}
			return db.NewContainer(containerName)
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
