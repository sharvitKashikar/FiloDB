package helper

import (
	"bufio"
	"fmt"
	"strings"
)

type TableInput struct {
	Name    string
	Types   []uint32
	Cols    []string
	Indexes [][]string
}

func GetTableInput(scanner *bufio.Reader) TableInput {
	name := GetTableName(scanner)

	fmt.Print("Enter column names (comma-separated): ")
	colsInput, _ := scanner.ReadString('\n')
	colsInput = strings.TrimSpace(colsInput)
	cols := strings.Split(colsInput, ",")

	fmt.Print("Enter column types (comma-separated as numbers): ")
	typesInput, _ := scanner.ReadString('\n')
	typesInput = strings.TrimSpace(typesInput)
	typesStr := strings.Split(typesInput, ",")
	types := make([]uint32, len(typesStr))
	for i, t := range typesStr {
		var typeValue uint32
		fmt.Sscanf(t, "%d", &typeValue)
		types[i] = typeValue
	}

	fmt.Print("Enter indexes (format: col1+col2,col3, ... or leave empty): ")
	indexInput, _ := scanner.ReadString('\n')
	indexInput = strings.TrimSpace(indexInput)

	indexes := [][]string{}
	if indexInput != "" {
		indexList := strings.Split(indexInput, ",")
		for _, indexCols := range indexList {
			indexes = append(indexes, strings.Split(indexCols, "+"))
		}
	}
	tdef := TableInput{
		Name:    name,
		Cols:    cols,
		Types:   types,
		Indexes: indexes,
	}
	return tdef
}

func GetTableName(scanner *bufio.Reader) string {
	fmt.Print("Enter table name: ")
	name, _ := scanner.ReadString('\n')
	name = strings.TrimSpace(name)
	return name
}

func PrintWelcomeMessage(isWelcome bool) {
	if isWelcome {
		fmt.Println("FiloDB has Started...")
	}
	fmt.Println("Available Commands You can use:")
	fmt.Println("  CREATE       - Create a new table")
	fmt.Println("  INSERT       - Add a record to a table")
	fmt.Println("  DELETE       - Delete a record from a table")
	fmt.Println("  GET          - Retrieve a record from a table")
	fmt.Println("  UPDATE       - Update a record in a table")
	fmt.Println("  BEGIN        - Begin new transaction")
	fmt.Println("  COMMIT       - Commit transaction")
	fmt.Println("  ABORT        - Rollback transaction")
	fmt.Println("  STATS        - Show database statistics")
	fmt.Println("  HELP         - List all commands")
	fmt.Println("  EXIT         - Exit the program")
	fmt.Println()
	fmt.Println("AGGREGATE FUNCTIONS:")
	fmt.Println("  COUNT        - Count records in a table")
	fmt.Println("  SUM          - Sum numeric values")
	fmt.Println("  AVG          - Calculate averages")
	fmt.Println("  MIN          - Find minimum values")
	fmt.Println("  MAX          - Find maximum values")
	fmt.Println()
	fmt.Println("UTILITY COMMANDS:")
	fmt.Println("  SCAN         - Show all records in a table")
	fmt.Println("  DEBUG        - Show table structure and info")
	fmt.Println()
}
