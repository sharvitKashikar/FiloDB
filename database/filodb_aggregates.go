// FiloDB Aggregates - Fixed Aggregate Functions using Table Scanning
// Author: Sharvit Kashikar
// Component: COUNT, SUM, AVG, MIN, MAX functions using proven table scanning method

package database

import (
	"bufio"
	"filodb/database/helper"
	"fmt"
	"strings"
	"time"
)

// RegisterAggregateCommands adds working aggregate commands to the command map
func RegisterAggregateCommands(commands map[string]Command) {
	commands["count"] = HandleCount
	commands["sum"] = HandleSum
	commands["avg"] = HandleAvg
	commands["min"] = HandleMin
	commands["max"] = HandleMax
	commands["debug"] = HandleDebugTable // Debug command to check table structure
	commands["scan"] = HandleTableScan   // New working scan command
}

// Helper function to get all records from a table using the same method as GET command
func getAllRecords(db *DB, tableName string) ([]*Record, error) {
	var kvReader KVReader
	db.kv.BeginRead(&kvReader)
	defer db.kv.EndRead(&kvReader)

	tdef := GetTableDef(db, tableName, &kvReader.Tree)
	if tdef == nil {
		return nil, fmt.Errorf("table '%s' not found", tableName)
	}

	// Use fullTableScan directly - this is the actual function that QueryWithFilter calls internally
	// This avoids the empty filter issue since fullTableScan doesn't require any filter
	results, err := fullTableScan(db, tableName, tdef)
	return results, err
}

// HandleCount - Count records in a table
func HandleCount(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	count := int64(len(results))

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("COUNT RESULT\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Count: %d\n", count)
	fmt.Println(strings.Repeat("=", 50))
}

// HandleSum - Sum numeric values in a column
func HandleSum(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	fmt.Print("Enter column name for SUM: ")
	columnInput, _ := scanner.ReadString('\n')
	columnName := strings.TrimSpace(columnInput)

	// First check if table exists and get its definition
	var reader KVReader
	db.kv.BeginRead(&reader)
	defer db.kv.EndRead(&reader)

	tdef := GetTableDef(db, tableName, &reader.Tree)
	if tdef == nil {
		fmt.Printf("Table '%s' not found.\n", tableName)
		return
	}

	colIndex := ColIndex(tdef, columnName)
	if colIndex < 0 {
		fmt.Printf("Column '%s' not found.\n", columnName)
		return
	}

	if tdef.Types[colIndex] != TYPE_INT64 && tdef.Types[colIndex] != TYPE_FLOAT64 {
		fmt.Printf("SUM only works with numeric (INT64 and FLOAT64) columns.\n")
		return
	}

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var sum float64 = 0 // Use float64 to handle both INT64 and FLOAT64
	count := int64(len(results))

	for _, record := range results {
		if tdef.Types[colIndex] == TYPE_INT64 {
			sum += float64(record.Vals[colIndex].I64)
		} else if tdef.Types[colIndex] == TYPE_FLOAT64 {
			sum += record.Vals[colIndex].F64
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("SUM RESULT\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Column: %s\n", columnName)
	fmt.Printf("Records processed: %d\n", count)
	if tdef.Types[colIndex] == TYPE_INT64 {
		fmt.Printf("SUM(%s): %.0f\n", columnName, sum) // Display as integer for INT64 columns
	} else {
		fmt.Printf("SUM(%s): %.6f\n", columnName, sum) // Display with decimals for FLOAT64 columns
	}
	fmt.Println(strings.Repeat("=", 50))
}

// HandleAvg - Calculate average of numeric values
func HandleAvg(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	fmt.Print("Enter column name for AVG: ")
	columnInput, _ := scanner.ReadString('\n')
	columnName := strings.TrimSpace(columnInput)

	var reader KVReader
	db.kv.BeginRead(&reader)
	defer db.kv.EndRead(&reader)

	tdef := GetTableDef(db, tableName, &reader.Tree)
	if tdef == nil {
		fmt.Printf("Table '%s' not found.\n", tableName)
		return
	}

	colIndex := ColIndex(tdef, columnName)
	if colIndex < 0 {
		fmt.Printf("Column '%s' not found.\n", columnName)
		return
	}

	if tdef.Types[colIndex] != TYPE_INT64 && tdef.Types[colIndex] != TYPE_FLOAT64 {
		fmt.Printf("AVG only works with numeric (INT64 and FLOAT64) columns.\n")
		return
	}

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	count := int64(len(results))
	if count == 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Printf("AVG RESULT\n")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Table: %s\n", tableName)
		fmt.Printf("Column: %s\n", columnName)
		fmt.Printf("AVG(%s): No records found\n", columnName)
		fmt.Println(strings.Repeat("=", 50))
		return
	}

	var sum float64 = 0 // Use float64 to handle both INT64 and FLOAT64
	for _, record := range results {
		if tdef.Types[colIndex] == TYPE_INT64 {
			sum += float64(record.Vals[colIndex].I64)
		} else if tdef.Types[colIndex] == TYPE_FLOAT64 {
			sum += record.Vals[colIndex].F64
		}
	}

	average := sum / float64(count)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("AVG RESULT\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Column: %s\n", columnName)
	fmt.Printf("Records processed: %d\n", count)
	fmt.Printf("AVG(%s): %.6f\n", columnName, average)
	fmt.Println(strings.Repeat("=", 50))
}

// HandleMin - Find minimum value in a column
func HandleMin(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	fmt.Print("Enter column name for MIN: ")
	columnInput, _ := scanner.ReadString('\n')
	columnName := strings.TrimSpace(columnInput)

	var reader KVReader
	db.kv.BeginRead(&reader)
	defer db.kv.EndRead(&reader)

	tdef := GetTableDef(db, tableName, &reader.Tree)
	if tdef == nil {
		fmt.Printf("Table '%s' not found.\n", tableName)
		return
	}

	colIndex := ColIndex(tdef, columnName)
	if colIndex < 0 {
		fmt.Printf("Column '%s' not found.\n", columnName)
		return
	}

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error scanning table: %v\n", err)
		return
	}

	count := int64(len(results))
	if count == 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Printf("MIN RESULT\n")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Table: %s\n", tableName)
		fmt.Printf("Column: %s\n", columnName)
		fmt.Printf("MIN(%s): No records found\n", columnName)
		fmt.Println(strings.Repeat("=", 50))
		return
	}

	var minVal int64
	var minStr string
	var minFloat float64
	var minBool bool
	var minTime time.Time

	for i, record := range results {
		if i == 0 {
			// Initialize with first record
			switch tdef.Types[colIndex] {
			case TYPE_INT64:
				minVal = record.Vals[colIndex].I64
			case TYPE_BYTES:
				minStr = string(record.Vals[colIndex].Str)
			case TYPE_FLOAT64:
				minFloat = record.Vals[colIndex].F64
			case TYPE_BOOLEAN:
				minBool = record.Vals[colIndex].Bool
			case TYPE_DATETIME:
				minTime = record.Vals[colIndex].Time
			}
		} else {
			// Compare and update minimum
			switch tdef.Types[colIndex] {
			case TYPE_INT64:
				val := record.Vals[colIndex].I64
				if val < minVal {
					minVal = val
				}
			case TYPE_BYTES:
				val := string(record.Vals[colIndex].Str)
				if val < minStr {
					minStr = val
				}
			case TYPE_FLOAT64:
				val := record.Vals[colIndex].F64
				if val < minFloat {
					minFloat = val
				}
			case TYPE_BOOLEAN:
				val := record.Vals[colIndex].Bool
				if !val && minBool { // false is less than true
					minBool = val
				}
			case TYPE_DATETIME:
				val := record.Vals[colIndex].Time
				if val.Before(minTime) {
					minTime = val
				}
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("MIN RESULT\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Column: %s\n", columnName)
	fmt.Printf("Records processed: %d\n", count)
	switch tdef.Types[colIndex] {
	case TYPE_INT64:
		fmt.Printf("MIN(%s): %d\n", columnName, minVal)
	case TYPE_BYTES:
		fmt.Printf("MIN(%s): %s\n", columnName, minStr)
	case TYPE_FLOAT64:
		fmt.Printf("MIN(%s): %.6f\n", columnName, minFloat)
	case TYPE_BOOLEAN:
		fmt.Printf("MIN(%s): %t\n", columnName, minBool)
	case TYPE_DATETIME:
		fmt.Printf("MIN(%s): %s\n", columnName, minTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println(strings.Repeat("=", 50))
}

// HandleMax - Find maximum value in a column
func HandleMax(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	fmt.Print("Enter column name for MAX: ")
	columnInput, _ := scanner.ReadString('\n')
	columnName := strings.TrimSpace(columnInput)

	var reader KVReader
	db.kv.BeginRead(&reader)
	defer db.kv.EndRead(&reader)

	tdef := GetTableDef(db, tableName, &reader.Tree)
	if tdef == nil {
		fmt.Printf("Table '%s' not found.\n", tableName)
		return
	}

	colIndex := ColIndex(tdef, columnName)
	if colIndex < 0 {
		fmt.Printf("Column '%s' not found.\n", columnName)
		return
	}

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error scanning table: %v\n", err)
		return
	}

	count := int64(len(results))
	if count == 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Printf("MAX RESULT\n")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Table: %s\n", tableName)
		fmt.Printf("Column: %s\n", columnName)
		fmt.Printf("MAX(%s): No records found\n", columnName)
		fmt.Println(strings.Repeat("=", 50))
		return
	}

	var maxVal int64
	var maxStr string
	var maxFloat float64
	var maxBool bool
	var maxTime time.Time

	for i, record := range results {
		if i == 0 {
			// Initialize with first record
			switch tdef.Types[colIndex] {
			case TYPE_INT64:
				maxVal = record.Vals[colIndex].I64
			case TYPE_BYTES:
				maxStr = string(record.Vals[colIndex].Str)
			case TYPE_FLOAT64:
				maxFloat = record.Vals[colIndex].F64
			case TYPE_BOOLEAN:
				maxBool = record.Vals[colIndex].Bool
			case TYPE_DATETIME:
				maxTime = record.Vals[colIndex].Time
			}
		} else {
			// Compare and update maximum
			switch tdef.Types[colIndex] {
			case TYPE_INT64:
				val := record.Vals[colIndex].I64
				if val > maxVal {
					maxVal = val
				}
			case TYPE_BYTES:
				val := string(record.Vals[colIndex].Str)
				if val > maxStr {
					maxStr = val
				}
			case TYPE_FLOAT64:
				val := record.Vals[colIndex].F64
				if val > maxFloat {
					maxFloat = val
				}
			case TYPE_BOOLEAN:
				val := record.Vals[colIndex].Bool
				if val && !maxBool { // true is greater than false
					maxBool = val
				}
			case TYPE_DATETIME:
				val := record.Vals[colIndex].Time
				if val.After(maxTime) {
					maxTime = val
				}
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("MAX RESULT\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Column: %s\n", columnName)
	fmt.Printf("Records processed: %d\n", count)
	switch tdef.Types[colIndex] {
	case TYPE_INT64:
		fmt.Printf("MAX(%s): %d\n", columnName, maxVal)
	case TYPE_BYTES:
		fmt.Printf("MAX(%s): %s\n", columnName, maxStr)
	case TYPE_FLOAT64:
		fmt.Printf("MAX(%s): %.6f\n", columnName, maxFloat)
	case TYPE_BOOLEAN:
		fmt.Printf("MAX(%s): %t\n", columnName, maxBool)
	case TYPE_DATETIME:
		fmt.Printf("MAX(%s): %s\n", columnName, maxTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println(strings.Repeat("=", 50))
}

// HandleTableScan - Shows all records in a table (debugging/verification)
func HandleTableScan(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	// Use the SAME table scanning approach as the working GET command
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("TABLE SCAN RESULT: %s\n", tableName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total records found: %d\n", len(results))
	fmt.Println(strings.Repeat("=", 50))

	if len(results) > 0 {
		// Print records in the same format as GET command
		printRecords(results)
	} else {
		fmt.Println("No records found in table")
	}
}

// HandleDebugTable - Debug function to check table structure
func HandleDebugTable(scanner *bufio.Reader, db *DB, currentTX *DBTX) {
	tableName := helper.GetTableName(scanner)

	// Get table definition first to show structure
	var reader KVReader
	db.kv.BeginRead(&reader)
	defer db.kv.EndRead(&reader)

	tdef := GetTableDef(db, tableName, &reader.Tree)
	if tdef == nil {
		fmt.Printf("Table '%s' not found.\n", tableName)
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("DEBUG TABLE INFO: %s\n", tableName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Columns: %v\n", tdef.Cols)
	fmt.Printf("Types: %v (1=INT64, 2=BYTES)\n", tdef.Types)

	// Now test the new scanning approach
	results, err := getAllRecords(db, tableName)
	if err != nil {
		fmt.Printf("Error with table scan: %v\n", err)
	} else {
		fmt.Printf("Records found with new method: %d\n", len(results))
		if len(results) > 0 {
			// Display sample record in human-readable format
			fmt.Print("Sample record: [")
			for i, val := range results[0].Vals {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s: \"%s\"", tdef.Cols[i], formatValue(val))
			}
			fmt.Println("]")
		}
	}
	fmt.Println(strings.Repeat("=", 50))
}

// Note: formatValue and printRecords functions are defined in filodb_commands.go
