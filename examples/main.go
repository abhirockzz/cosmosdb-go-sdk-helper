package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abhirockzz/cosmosdb-go-sdk-helper/auth"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/common"
)

func main() {
	// Authenticate with default credential
	client, err := auth.GetClientWithDefaultAzureCredential("https://ACCOUNT_NAME.documents.azure.com", nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// Create database if not exists
	db, err := common.CreateDatabaseIfNotExists(client, "sampledb")
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())

	// Create container if not exists
	container, err := common.CreateContainerIfNotExists(db, "samplecontainer")
	if err != nil {
		log.Fatalf("CreateContainerIfNotExists failed: %v", err)
	}
	fmt.Println("Container ready:", container.ID())

	// List all databases
	dbs, err := common.GetAllDatabases(client)
	if err != nil {
		log.Fatalf("GetAllDatabases failed: %v", err)
	}
	fmt.Println("Databases:")
	for _, d := range dbs {
		fmt.Println("-", d.ID)
	}

	// List all containers in the database
	containers, err := common.GetAllContainers(db)
	if err != nil {
		log.Fatalf("GetAllContainers failed: %v", err)
	}
	fmt.Println("Containers in database:")
	for _, c := range containers {
		fmt.Println("-", c.ID)
	}

	db, err = client.NewDatabase("i_am_not_there")
	if err != nil {
		log.Fatalf("NewDatabase failed: %v", err)
	}
	// Error handling helper usage: simulate a non-existing database
	_, err = db.Read(context.Background(), nil)
	if err != nil {
		errInfo := common.GetError(err)
		fmt.Printf("Error info: status=%d, message=%q\n", errInfo.Status, errInfo.Message)
	}

	// 6. Error handling helper usage: simulate a non-existing database
	// _, err = common.CreateDatabaseIfNotExists(client, "definitely_not_exist_db")
	// if err != nil {
	// 	errInfo := cosmosdb_errors.GetError(err)
	// 	fmt.Printf("Simulated error info: status=%d, message=%q\n", errInfo.Status, errInfo.Message)
	// }

	// Authenticate with Emulator using Azure AD token
	// start the emulator first: docker run -p 8081:8081 -n linux-emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest

	//docker run --publish 8081:8081 --publish 1234:1234 mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:vnext-preview

	emuClient, err := auth.GetEmulatorClientWithAzureADAuth("http://localhost:8081", nil)
	if err != nil {
		log.Fatalf("Emulator auth failed: %v", err)
	}
	fmt.Println("Authenticated with Emulator.")

	db, err = common.CreateDatabaseIfNotExists(emuClient, "sampledb")
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())
}
