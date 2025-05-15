package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/auth"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/common"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/cosmosdb_errors"
	"github.com/abhirockzz/cosmosdb-go-sdk-helper/operations"
)

func defaultAzureCredentialExample(endpoint string) {
	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// use the client to perform operations
	_ = client
}

func dbAndContainerCreationExample(endpoint string) {
	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	// Create database if not exists
	db, err := common.CreateDatabaseIfNotExists(client, azcosmos.DatabaseProperties{
		ID: "tododb",
	}, nil)
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())

	// Create container if not exists
	container, err := common.CreateContainerIfNotExists(db, azcosmos.ContainerProperties{
		ID: "tasks",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/id"},
			Kind:  azcosmos.PartitionKeyKindHash,
		},
	}, nil)
	if err != nil {
		log.Fatalf("CreateContainerIfNotExists failed: %v", err)
	}
	fmt.Println("Container ready:", container.ID())
}

func getAllDBandContainersExample(endpoint string) {

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
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
	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
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
		errInfo := cosmosdb_errors.GetError(err)
		fmt.Printf("Error info: status=%d, message=%q\n", errInfo.Status, errInfo.Message)
	}
}

func emulatorADAuthExample() {
	// start the emulator first: docker run -p 8081:8081 -n linux-emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest OR

	//docker run --publish 8081:8081 --publish 1234:1234 mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:vnext-preview

	// Authenticate with Emulator using Azure AD token

	emuClient, err := auth.GetCosmosDBClient("http://localhost:8081", true, nil)
	if err != nil {
		log.Fatalf("Emulator auth failed: %v", err)
	}
	fmt.Println("Authenticated with Emulator.")

	db, err := common.CreateDatabaseIfNotExists(emuClient, azcosmos.DatabaseProperties{
		ID: "sampledb",
	}, nil)
	if err != nil {
		log.Fatalf("CreateDatabaseIfNotExists failed: %v", err)
	}
	fmt.Println("Database ready:", db.ID())
}

func queryItemsExample1(endpoint, sqlQuery, databaseName, containerName string) {

	type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	taskID := "task-" + strconv.Itoa(rand.Intn(1000))
	task := Task{
		ID:   taskID,
		Info: "Sample task",
	}

	insertedTask, err := operations.InsertItemWithResponse(container, task, azcosmos.NewPartitionKeyString(task.ID), nil)
	if err != nil {
		log.Fatalf("InsertItem failed: %v", err)
	}
	fmt.Printf("Inserted task: %s (%s)\n", insertedTask.ID, insertedTask.Info)

	tasks, err := operations.ExecuteQuery[Task](container, sqlQuery, azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("Task: %s (%s)\n", task.ID, task.Info)
	}
}

func queryItemsExample2(endpoint, databaseName, containerName string) {

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	tasks, err := operations.ExecuteQuery[map[string]any](container, "SELECT * FROM c", azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("Task: %s (%s)\n", task["id"], task["info"])
	}
}

func queryItemExample(endpoint, databaseName, containerName, itemID, partitionKey string) {

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	task, err := operations.GetItem[map[string]any](container, itemID, azcosmos.NewPartitionKeyString(partitionKey), nil)
	if err != nil {
		log.Fatalf("GetItem failed: %v", err)
	}
	fmt.Printf("Task: %s (%s)\n", task["id"], task["info"])
}

func queryItemsWithMetricsExample1(endpoint, databaseName, containerName string) {

	type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	//queryResult, err := query.QueryItemsWithMetrics[Task](container, "SELECT * FROM c", azcosmos.NewPartitionKey(), nil)
	queryResult, err := operations.ExecuteQueryWithMetrics[Task](container, "SELECT * FROM c WHERE c.id = 1", azcosmos.NewPartitionKey(), nil)
	if err != nil {
		log.Fatalf("QueryItems failed: %v", err)
	}
	for _, task := range queryResult.Items {
		fmt.Printf("Task: %s (%s)\n", task.ID, task.Info)
	}

	// Print metrics for each page
	for i, metrics := range queryResult.Metrics {
		fmt.Printf("Metrics for page %d: ", i)
		fmt.Printf("TotalExecutionTimeInMs: %f, QueryCompileTimeInMs: %f\n", metrics.TotalExecutionTimeInMs, metrics.QueryCompileTimeInMs)
	}

	// Print total request charge
	fmt.Printf("Total request charge: %f\n", queryResult.RequestCharge)
}

func insertItemExample(endpoint, databaseName, containerName string) {

	type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	task := Task{
		ID:   "44",
		Info: "Sample task",
	}

	insertedTask, err := operations.InsertItemWithResponse(container, task, azcosmos.NewPartitionKeyString(task.ID), nil)
	if err != nil {
		log.Fatalf("InsertItem failed: %v", err)
	}
	fmt.Printf("Inserted task: %s (%s)\n", insertedTask.ID, insertedTask.Info)
}

func replaceItemExample(endpoint, databaseName, containerName string) {

	type Task struct {
		ID   string `json:"id"`
		Info string `json:"info"`
	}

	client, err := auth.GetCosmosDBClient(endpoint, false, nil)
	if err != nil {
		log.Fatalf("Azure AD auth failed: %v", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatalf("NewContainer failed: %v", err)
	}

	task := Task{
		ID:   "42",
		Info: "Sample task updated_2",
	}

	replacedTask, err := operations.ReplaceItemWithResponse(container, "42", azcosmos.NewPartitionKeyString(task.ID), task, nil)
	if err != nil {
		log.Fatalf("ReplaceItem failed: %v", err)
	}
	fmt.Printf("Replaced task: %s (%s)\n", replacedTask.ID, replacedTask.Info)
}
func main() {
	endpoint := "https://ACCOUNT_NAME.documents.azure.com:443"

	//defaultAzureCredentialExample(endpoint)
	//dbAndContainerCreationExample(endpoint)
	// getAllDBandContainersExample(endpoint)
	// errorHandlingHelperExample(endpoint)
	// emulatorADAuthExample()

	//insertItemExample(endpoint, "tododb", "tasks")
	queryItemsExample1(endpoint, "select * from c", "tododb", "tasks")
	//queryItemsExample1(endpoint, "select top 3 * from c", "tododb", "tasks")

	//queryItemsExample2(endpoint, "tododb", "tasks")
	//queryItemExample(endpoint, "tododb", "tasks", "3", "3")
	//queryItemsWithMetricsExample1(endpoint, "tododb", "tasks")
	//replaceItemExample(endpoint, "tododb", "tasks")
}
