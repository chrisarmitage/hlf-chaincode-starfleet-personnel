package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/api/internal/fabricgateway"
	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/api/internal/personnelclient"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	gateway := fabricgateway.NewGateway()
	defer gateway.Close()

	contract, err := gateway.GetContract()
	if err != nil {
		log.Fatalf("failed to get contract: %v", err)
	}

	client := personnelclient.NewPersonnelClient(contract)

	command := os.Args[1]

	switch command {
	case "get-personnel":
		handleGetPersonnel(client, os.Args[2:])
	case "enroll-cadet":
		handleEnrollCadet(client, os.Args[2:])
	case "complete-training":
		handleCompleteTraining(client, os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run . get-personnel <personnel-id>")
	fmt.Println("  go run . enroll-cadet <personnel-id> <name> <campus>")
	fmt.Println("  go run . complete-training <record-id> <personnel-id> <campus> <training-code> <completed-at> <issued-by>")
	fmt.Println("\nExamples:")
	fmt.Println("  go run . get-personnel SF-001")
	fmt.Println(`  go run . enroll-cadet SF-001 "Malcom Reynolds" Engineering`)
	fmt.Println(`  go run . complete-training TR-001 SF-001 Engineering ENG-WARP-201 2024-06-01T12:00:00Z "Captain Janeway"`)
}

func handleGetPersonnel(client *personnelclient.PersonnelClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: personnel-id is required")
		fmt.Println("Usage: go run . get-personnel <personnel-id>")
		os.Exit(1)
	}

	personnelID := args[0]

	personnel, err := client.GetPersonnel(personnelID)
	if err != nil {
		log.Fatalf("failed to get personnel: %v", err)
	}

	fmt.Printf("Personnel Info:\n")
	fmt.Printf("  ID:     %s\n", personnel.PersonnelID)
	fmt.Printf("  Name:   %s\n", personnel.Name)
	fmt.Printf("  Rank:   %s\n", personnel.Rank)
	fmt.Printf("  Campus: %s\n", personnel.Campus)
	fmt.Printf("  Status: %s\n", personnel.Status)
}

func handleEnrollCadet(client *personnelclient.PersonnelClient, args []string) {
	if len(args) < 3 {
		fmt.Println("Error: personnel-id, name, and campus are required")
		fmt.Println(`Usage: go run . enroll-cadet <personnel-id> <name> <campus>`)
		os.Exit(1)
	}

	personnelID := args[0]
	name := args[1]
	campus := args[2]

	personnel, err := client.EnrollCadet(personnelID, name, campus)
	if err != nil {
		log.Fatalf("failed to enroll cadet: %v", err)
	}

	fmt.Printf("Cadet enrolled successfully:\n")
	fmt.Printf("  ID:     %s\n", personnel.PersonnelID)
	fmt.Printf("  Name:   %s\n", personnel.Name)
	fmt.Printf("  Rank:   %s\n", personnel.Rank)
	fmt.Printf("  Campus: %s\n", personnel.Campus)
	fmt.Printf("  Status: %s\n", personnel.Status)
}

func handleCompleteTraining(client *personnelclient.PersonnelClient, args []string) {
	if len(args) < 6 {
		fmt.Println("Error: record-id, personnel-id, campus, training-code, completed-at, and issued-by are required")
		fmt.Println(`Usage: go run . complete-training <record-id> <personnel-id> <campus> <training-code> <completed-at> <issued-by>`)
		os.Exit(1)
	}

	recordID := args[0]
	personnelID := args[1]
	campus := args[2]
	trainingCode := args[3]
	completedAt := args[4]
	issuedBy := args[5]

	training, err := client.CompleteTraining(recordID, personnelID, campus, trainingCode, completedAt, issuedBy)
	if err != nil {
		log.Fatalf("failed to complete training: %v", err)
	}

	fmt.Printf("Training completed successfully:\n")
	fmt.Printf("  Record ID:     %s\n", training.RecordID)
	fmt.Printf("  Personnel ID:  %s\n", training.PersonnelID)
	fmt.Printf("  Campus:        %s\n", training.Campus)
	fmt.Printf("  Training Code: %s\n", training.TrainingCode)
	fmt.Printf("  Completed At:  %s\n", training.CompletedAt)
	fmt.Printf("  Issued By:     %s\n", training.IssuedBy)
	fmt.Printf("  Status:        %s\n", training.Status)
}
