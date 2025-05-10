# Azure Cosmos DB Go SDK Helper

A simple helper package providing convenience functions for working with the [Go SDK for Azure Cosmos DB NoSQL API](https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/sdk-go). This includes utilities for authentication, querying, database and container operations, error handling, etc.

Packages:

- [auth](auth): Authentication utilities for Azure Cosmos DB
- [common](common): Common utilities for database and container operations
- [query](query): Query utilities to retrieve data using generic types
- [cosmosdb_errors](cosmosdb_errors): Error handling utilities for Cosmos DB operations

## Installation

```bash
go get github.com/abhirockzz/cosmosdb-go-sdk-helper
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/abhirockzz/cosmosdb-go-sdk-helper/auth"
    "github.com/abhirockzz/cosmosdb-go-sdk-helper/common"
    "github.com/abhirockzz/cosmosdb-go-sdk-helper/query"
)

func main() {
    // Connect using Azure AD authentication
    client, err := auth.GetClientWithDefaultAzureCredential("your-cosmos-endpoint", nil)
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    // Create database if not exists
    db, err := common.CreateDatabaseIfNotExists(client, "mydb")
    if err != nil {
        log.Fatalf("Database creation failed: %v", err)
    }
    
    // Create container if not exists
    container, err := common.CreateContainerIfNotExists(db, "mycontainer")
    if err != nil {
        log.Fatalf("Container creation failed: %v", err)
    }
    
    fmt.Println("Database and container ready!")

    type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

    tasks, err := query.QueryItems[Task](container, "SELECT * FROM c", azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("Task: %s (%s)\n", task.ID, task.Info)
	}

}
```

> Also, refer to [examples](examples).

## Authentication Options

### Azure AD Authentication

```go
// Using DefaultAzureCredential (supports multiple authentication methods)
client, err := auth.GetClientWithDefaultAzureCredential("https://your-account.documents.azure.com:443", nil)
```

### Local Emulator

```go
// For local development with the Cosmos DB Emulator
client, err := auth.GetEmulatorClientWithAzureADAuth("https://localhost:8081", nil)
```

## Available Helper Functions

### Database Operations

- `CreateDatabaseIfNotExists`: Creates a database only if it doesn't already exist
- `GetAllDatabases`: Retrieves a list of all databases in your Cosmos account

### Container Operations

- `CreateContainerIfNotExists`: Creates a container only if it doesn't already exist
- `GetAllContainers`: Retrieves a list of all containers in a database

## Query operations (using Generics)

- `QueryItems`: Executes a SQL query against a container and returns the results
- `QueryItem`: Retrieves a single item from a container using its ID and partition key

## Error Handling

Error handling for Cosmos DB operations:

```go
if err != nil {
    cosmosError := cosmosdb_errors.GetError(err)
    if cosmosError.Status == http.StatusNotFound {
        // Handle resource not found
    }
    // Handle other errors
}
```