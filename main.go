package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
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
const fileExportName = "expenses.csv"

// loadData loads data from expenses.json
func loadData() []Expense {
	var expenses []Expense
	fileData, err := os.ReadFile(fileDataName)
	if err == nil && len(fileData) > 0 {
		json.Unmarshal(fileData, &expenses)
	}

	return expenses
}

func saveData(expenses []Expense) {
	updatedData, _ := json.MarshalIndent(expenses, "", "    ")
	os.WriteFile(fileDataName, updatedData, 0644)
}

func runAddCommand(cmd *cobra.Command, args []string) {
	description, _ := cmd.Flags().GetString("description")
	amount, _ := cmd.Flags().GetInt("amount")

	var expenses []Expense
	expenses = loadData()
	newID := 1
	if len(expenses) > 0 {
		newID = expenses[len(expenses)-1].ID + 1
	}

	newExpense := Expense{
		ID:          newID,
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}

	expenses = append(expenses, newExpense)
	saveData(expenses)

	fmt.Printf("Expense added successfully (ID: %d)\n", newID)
}

func runListCommand(cmd *cobra.Command, args []string) {
	var expenses []Expense
	expenses = loadData()

	const padding = 3
	const symbol = ' '

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, symbol, 0)
	fmt.Fprint(w, "ID\tDate\tDescription\tAmount\n")

	for _, e := range expenses {
		fmt.Fprintf(w, "%d\t%s\t%s\t$%d\n", e.ID, e.Date.Format("2006-01-02"), e.Description, e.Amount)
	}

	w.Flush()
}

func runDeleteCommand(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetInt("id")

	var expenses, NewExpenses []Expense
	expenses = loadData()

	for _, e := range expenses {
		if e.ID != id {
			NewExpenses = append(NewExpenses, e)
		}
	}

	saveData(NewExpenses)

	fmt.Printf("Expense deleted successfully (ID=%d)", id)
}

func runSummaryCommand(cmd *cobra.Command, args []string) {
	month, _ := cmd.Flags().GetInt("month")

	var expenses []Expense
	expenses = loadData()

	var totalExpenseAmount int
	if month == 0 {
		for _, e := range expenses {
			totalExpenseAmount += e.Amount
		}

		fmt.Printf("Total expenses: $%d", totalExpenseAmount)
	} else {
		month := time.Month(month)
		for _, e := range expenses {
			if e.Date.Month() == month {
				totalExpenseAmount += e.Amount
			}
		}

		fmt.Printf("Total expenses for %s: %d", month, totalExpenseAmount)
	}
}

func runExportCommand(cmd *cobra.Command, args []string) {
	var expenses []Expense
	expenses = loadData()

	file, err := os.Create(fileExportName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headingRow := []string{
		"ID",
		"Date",
		"Description",
		"Amount",
	}
	writer.Write(headingRow)
	
	for _, e := range expenses {
		row := []string{
			strconv.Itoa(e.ID),
			e.Date.Format("2006-01-02"),
			e.Description,
			strconv.Itoa(e.Amount),
		}
		writer.Write(row)
	}

	fmt.Println("Data was successully exported to csv file")
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "expense-tracker",
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new expense",
		Run:   runAddCommand,
	}

	// add flags for 'add' command
	addCmd.Flags().String("description", "", "description of expense")
	addCmd.Flags().Int("amount", 0, "amount of expense")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show all expenses in your list",
		Run:   runListCommand,
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete expense with corresponding ID",
		Run:   runDeleteCommand,
	}

	// add flags for 'delete' command
	deleteCmd.Flags().Int("id", 0, "id of expense to delete")

	var summaryCmd = &cobra.Command{
		Use:   "summary",
		Short: "Print out total expenses",
		Run:   runSummaryCommand,
	}

	// add flags for 'summary' command
	summaryCmd.Flags().Int("month", 0, "month expenses")

	var exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export all data in CSV file",
		Run:   runExportCommand,
	}

	// add commands to the root
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.Execute()
}
