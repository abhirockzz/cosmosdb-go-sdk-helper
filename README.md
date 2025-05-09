# Cosmos DB Go SDK Helper

A simple helper library providing convenience functions for working with the Azure Cosmos DB Go SDK. This includes utilities for authentication, database and container management, error handling, etc.

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
}
```

> Also, refer to [examples](examples).


## Authentication Options

### Azure AD Authentication

```go
// Using DefaultAzureCredential (supports multiple authentication methods)
client, err := auth.GetClientWithDefaultAzureCredential("https://your-account.documents.azure.com", nil)
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