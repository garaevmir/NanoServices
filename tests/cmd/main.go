package main

import (
	"fmt"
	"os"

	"github.com/nanoservices/tests/models"
	"github.com/nanoservices/tests/scenarios"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <username> <password> <email>")
		return
	}

	user := models.User{
		Username: os.Args[1],
		Password: os.Args[2],
		Email:    os.Args[3],
	}

	fmt.Println("=== Running Scenario 1: Create Post and Get Stats ===")
	scenarios.Scenario1(user)

	fmt.Println("\n=== Running Scenario 2: Browse Posts and View Dynamics ===")
	scenarios.Scenario2(user)

	fmt.Println("\n=== Running Scenario 3: Update Profile and Comment ===")
	scenarios.Scenario3(user)
}
