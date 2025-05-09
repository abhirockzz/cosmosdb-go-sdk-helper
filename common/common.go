package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/cosmosdb_errors"
)

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

func CreateContainerIfNotExists(client *azcosmos.DatabaseClient, containerName string) (*azcosmos.ContainerClient, error) {
	container, err := client.NewContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client: %v", err)
	}

	_, err = container.Read(context.Background(), nil)
	if err != nil {
		if cosmosdb_errors.GetError(err).Status == http.StatusNotFound {
			// Container doesn't exist, try to create it
			_, err = client.CreateContainer(context.Background(), azcosmos.ContainerProperties{
				ID: containerName,
			}, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create container: %v", err)
			}
			return container, nil
		}
		return nil, fmt.Errorf("failed to read container: %v", err)
	}

	return container, nil
}

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
