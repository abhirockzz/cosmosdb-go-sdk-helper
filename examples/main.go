package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/auth"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/common"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/query"
)

func defaultAzureCredentialExample(endpoint string) {
	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// use the client to perform operations
	_ = client
}

func dbAndContainerCreationExample(endpoint string) {
	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// Create database if not exists
	db, err := common.CreateDatabaseIfNotExists(client, "tododb")
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())

	// Create container if not exists
	container, err := common.CreateContainerIfNotExists(db, "tasks")
	if err != nil {
		log.Fatalf("CreateContainerIfNotExists failed: %v", err)
	}
	fmt.Println("Container ready:", container.ID())
}

func getAllDBandContainersExample(endpoint string) {

	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// List all databases
	dbs, err := common.GetAllDatabases(client)
	if err != nil {
		log.Fatalf("GetAllDatabases failed: %v", err)
	}
	fmt.Println("Databases:")
	for _, d := range dbs {
		fmt.Println("-", d.ID)
		db, err := client.NewDatabase(d.ID)
		if err != nil {
			log.Fatalf("NewDatabase failed: %v", err)
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
	}

}

func errorHandlingHelperExample(endpoint string) {
	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	db, err := client.NewDatabase("i_am_not_there")
	if err != nil {
		log.Fatalf("NewDatabase failed: %v", err)
	}
	// Error handling helper usage: simulate a non-existing database
	_, err = db.Read(context.Background(), nil)
	if err != nil {
		errInfo := common.GetError(err)
		fmt.Printf("Error info: status=%d, message=%q\n", errInfo.Status, errInfo.Message)
	}
}

func emulatorADAuthExample() {
	// start the emulator first: docker run -p 8081:8081 -n linux-emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest OR

	//docker run --publish 8081:8081 --publish 1234:1234 mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:vnext-preview

	// Authenticate with Emulator using Azure AD token

	emuClient, err := auth.GetEmulatorClientWithAzureADAuth("http://localhost:8081", nil)
	if err != nil {
		log.Fatalf("Emulator auth failed: %v", err)
	}
	fmt.Println("Authenticated with Emulator.")

	db, err := common.CreateDatabaseIfNotExists(emuClient, "sampledb")
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())
}

func queryItemsExample1(endpoint, databaseName, containerName string) {

	type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	tasks, err := query.QueryItems[Task](container, "SELECT * FROM c", azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("Task: %s (%s)\n", task.ID, task.Info)
	}
}

func queryItemsExample2(endpoint, databaseName, containerName string) {

	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	tasks, err := query.QueryItems[map[string]any](container, "SELECT * FROM c", azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("Task: %s (%s)\n", task["id"], task["info"])
	}
}

func queryItemExample(endpoint, databaseName, containerName, itemID, partitionKey string) {

	client, err := auth.GetClientWithDefaultAzureCredential(endpoint, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	task, err := query.QueryItem[map[string]any](container, itemID, azcosmos.NewPartitionKeyString(partitionKey), nil)
	if err != nil {
		log.Fatalf("QueryItem failed: %v", err)
	}
	fmt.Printf("Task: %s (%s)\n", task["id"], task["info"])
}

func main() {
	endpoint := "https://ACCOUNT_NAME.documents.azure.com:443"

	defaultAzureCredentialExample(endpoint)
	// dbAndContainerCreationExample(endpoint)
	// getAllDBandContainersExample(endpoint)
	// errorHandlingHelperExample(endpoint)
	// emulatorADAuthExample()

	//queryItemsExample1(endpoint, "tododb", "tasks")
	//queryItemsExample2(endpoint, "tododb", "tasks")
	//queryItemExample(endpoint, "tododb", "tasks", "3", "3")

}
