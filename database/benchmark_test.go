package database

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Simple performance test to get basic numbers
func TestSimplePerformance(t *testing.T) {
	fmt.Println("\n=== FiloDB Performance Test Results ===")

	// Test small batch insert performance
	testBatchInserts(t)

	// Test lookup performance
	testLookupPerformance(t)

	// Test database size characteristics
	testDatabaseSizes(t)

	fmt.Println("=========================================")
}

func testBatchInserts(t *testing.T) {
	dbPath := "perf_test.db"
	defer os.Remove(dbPath)

	db := setupBenchDB(t, dbPath)
	defer db.kv.Close()

	// Test inserting 50 records and measure time
	recordCount := 50
	start := time.Now()

	for i := 0; i < recordCount; i++ {
		var writer KVTX
		db.kv.Begin(&writer)

		rec := Record{
			Cols: []string{"id", "name", "data"},
			Vals: []Value{
				{Type: TYPE_INT64, I64: int64(i)},
				{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("User%d", i))},
				{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("Sample data for user %d", i))},
			},
		}

		_, err := db.Insert("bench_table", rec, &writer)
		if err != nil {
			t.Fatal(err)
		}

		db.kv.Commit(&writer)
	}

	elapsed := time.Since(start)
	opsPerSec := float64(recordCount) / elapsed.Seconds()

	fmt.Printf("Insert Performance: %.0f ops/sec (%d records in %v)\n",
		opsPerSec, recordCount, elapsed)
	fmt.Printf("Average per insert: %v\n", elapsed/time.Duration(recordCount))
}

func testLookupPerformance(t *testing.T) {
	dbPath := "lookup_test.db"
	defer os.Remove(dbPath)

	// Setup database with test data
	db := setupBenchDB(t, dbPath)
	defer db.kv.Close()

	// Insert test records
	for i := 0; i < 20; i++ {
		var writer KVTX
		db.kv.Begin(&writer)

		rec := Record{
			Cols: []string{"id", "name", "data"},
			Vals: []Value{
				{Type: TYPE_INT64, I64: int64(i)},
				{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("User%d", i))},
				{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("Data for user %d", i))},
			},
		}

		_, err := db.Insert("bench_table", rec, &writer)
		if err != nil {
			t.Fatal(err)
		}

		db.kv.Commit(&writer)
	}

	// Test lookup performance
	lookupCount := 100
	start := time.Now()

	for i := 0; i < lookupCount; i++ {
		var reader KVReader
		db.kv.BeginRead(&reader)

		rec := Record{
			Cols: []string{"id"},
			Vals: []Value{
				{Type: TYPE_INT64, I64: int64(i % 20)},
			},
		}

		_, err := db.Get("bench_table", &rec, &reader)
		if err != nil {
			t.Fatal(err)
		}

		db.kv.EndRead(&reader)
	}

	elapsed := time.Since(start)
	opsPerSec := float64(lookupCount) / elapsed.Seconds()

	fmt.Printf("Lookup Performance: %.0f ops/sec (%d lookups in %v)\n",
		opsPerSec, lookupCount, elapsed)
	fmt.Printf("Average per lookup: %v\n", elapsed/time.Duration(lookupCount))
}

func testDatabaseSizes(t *testing.T) {
	fmt.Println("\n=== Database Size Analysis ===")

	sizes := []int{10, 20, 30}

	for _, recordCount := range sizes {
		dbPath := fmt.Sprintf("size_test_%d.db", recordCount)

		db := setupBenchDB(t, dbPath)

		// Insert records
		for i := 0; i < recordCount; i++ {
			var writer KVTX
			db.kv.Begin(&writer)

			rec := Record{
				Cols: []string{"id", "name", "data"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: int64(i)},
					{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("User%d", i))},
					{Type: TYPE_BYTES, Str: []byte(fmt.Sprintf("Data content for user number %d with additional text", i))},
				},
			}

			_, err := db.Insert("bench_table", rec, &writer)
			if err != nil {
				t.Fatal(err)
			}

			db.kv.Commit(&writer)
		}

		// Check file size
		if stat, err := os.Stat(dbPath); err == nil {
			sizeKB := float64(stat.Size()) / 1024
			fmt.Printf("Records: %d -> Size: %.1f KB (%.2f KB per record)\n",
				recordCount, sizeKB, sizeKB/float64(recordCount))
		}

		db.kv.Close()
		os.Remove(dbPath)
	}
}

// Helper function to setup a database
func setupBenchDB(t *testing.T, path string) *DB {
	testDB := setupTestDB(t)
	testDB.Path = path
	testDB.kv.Path = path

	// Reopen with new path
	testDB.kv.Close()
	if err := testDB.kv.Open(); err != nil {
		t.Fatal("Failed to open:", err)
	}

	// Initialize internal tables
	err := initializeInternalTables(testDB)
	if err != nil {
		t.Fatal("Error while init table: ", err)
	}

	// Create benchmark table
	var writer KVTX
	testDB.kv.Begin(&writer)

	tdef := &TableDef{
		Name:  "bench_table",
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES},
		Cols:  []string{"id", "name", "data"},
		PKeys: 1,
	}

	if err := testDB.TableNew(tdef, &writer); err != nil {
		t.Fatal("failed to create bench table: ", err)
	}
	testDB.kv.Commit(&writer)

	return testDB
}
