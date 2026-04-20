package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}

const fileDataName = "expenses.json"

func runAddCommand(cmd *cobra.Command, args []string) {
	description, _ := cmd.Flags().GetString("description")
	amount, _ := cmd.Flags().GetInt("amount")

	var expenses []Expense
	fileData, err := os.ReadFile(fileDataName)
	if err == nil && len(fileData) > 0 {
		json.Unmarshal(fileData, &expenses)
	}

	newID := 1
	if len(expenses) > 0 {
		newID = expenses[len(expenses)-1].ID + 1
	}

	newExpanse := Expense{
		ID:          newID,
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}

	expenses = append(expenses, newExpanse)
	updatedData, _ := json.MarshalIndent(expenses, "", "    ")
	os.WriteFile(fileDataName, updatedData, 0644)

	fmt.Printf("Expense added successfully (ID: %d)\n", newID)
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "extence-tracker",
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new expance",
		Run:   runAddCommand,
	}

	// add flags for 'add' command
	addCmd.Flags().String("description", "", "description of expance")
	addCmd.Flags().Int("amount", 0, "amount of expance")
	rootCmd.AddCommand(addCmd)
	rootCmd.Execute()
}
