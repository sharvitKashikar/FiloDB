package database

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.kv.Close()
}

func TestCreateTable(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	tests := []struct {
		name        string
		tableDef    *TableDef
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid table creation",
			tableDef: &TableDef{
				Name:  "users",
				Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES},
				Cols:  []string{"id", "name", "email"},
				PKeys: 1,
			},
			expectError: false,
		},
		{
			name: "duplicate table",
			tableDef: &TableDef{
				Name:  "users",
				Types: []uint32{TYPE_INT64, TYPE_BYTES},
				Cols:  []string{"id", "name"},
				PKeys: 1,
			},
			expectError: true,
			errorMsg:    "table already exists",
		},
		{
			name: "empty table name",
			tableDef: &TableDef{
				Name:  "",
				Types: []uint32{TYPE_INT64},
				Cols:  []string{"id"},
				PKeys: 1,
			},
			expectError: true,
			errorMsg:    "table name cannot be empty",
		},
		{
			name: "mismatched columns and types",
			tableDef: &TableDef{
				Name:  "test",
				Types: []uint32{TYPE_INT64},
				Cols:  []string{"id", "name"},
				PKeys: 1,
			},
			expectError: true,
			errorMsg:    "length of columns & types do not match",
		},
		{
			name: "invalid column type",
			tableDef: &TableDef{
				Name:  "test",
				Types: []uint32{99},
				Cols:  []string{"id"},
				PKeys: 1,
			},
			expectError: true,
			errorMsg:    "invalid data type",
		},
		{
			name: "multiple primary keys",
			tableDef: &TableDef{
				Name:  "test",
				Types: []uint32{TYPE_INT64, TYPE_INT64},
				Cols:  []string{"id1", "id2"},
				PKeys: 2,
			},
			expectError: true,
			errorMsg:    "only one primary key is allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer KVTX
			db.kv.Begin(&writer)
			err := db.TableNew(tt.tableDef, &writer)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !isEqual(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
				db.kv.Abort(&writer)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					db.kv.Abort(&writer)
				} else {
					db.kv.Commit(&writer)
				}
			}
		})
	}
}

func TestInsertRecord(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a test table first
	setupTestTable(t, db)

	tests := []struct {
		name        string
		record      Record
		mode        int
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid insert",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			mode:        MODE_INSERT_ONLY,
			expectError: false,
		},
		{
			name: "duplicate key",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			mode:        MODE_INSERT_ONLY,
			expectError: true,
			errorMsg:    "record already exists",
		},
		{
			name: "missing primary key",
			record: Record{
				Cols: []string{"name", "email"},
				Vals: []Value{
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			mode:        MODE_INSERT_ONLY,
			expectError: true,
			errorMsg:    "missing column",
		},
		{
			name: "type mismatch",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_BYTES, Str: []byte("invalid")},
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			mode:        MODE_INSERT_ONLY,
			expectError: true,
			errorMsg:    "invalid type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer KVTX
			db.kv.Begin(&writer)
			inserted, err := db.Insert("users", tt.record, &writer)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !isEqual(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
				db.kv.Abort(&writer)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					db.kv.Abort(&writer)
				} else if !inserted {
					t.Error("expected record to be inserted")
					db.kv.Abort(&writer)
				} else {
					db.kv.Commit(&writer)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a test table and insert a record first
	setupTestTable(t, db)
	insertTestRecord(t, db, 1)

	tests := []struct {
		name        string
		record      Record
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid delete",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			expectError: false,
		},
		{
			name: "non-existent record",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 999},
					{Type: TYPE_BYTES, Str: []byte("NonExistent")},
					{Type: TYPE_BYTES, Str: []byte("none@example.com")},
				},
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "missing primary key",
			record: Record{
				Cols: []string{"name", "email"},
				Vals: []Value{
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			expectError: true,
			errorMsg:    "missing primary key column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer KVTX
			db.kv.Begin(&writer)
			deleted, err := db.Delete("users", tt.record, &writer)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !isEqual(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
				db.kv.Abort(&writer)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					db.kv.Abort(&writer)
				} else if !deleted {
					t.Error("expected record to be deleted")
					db.kv.Abort(&writer)
				} else {
					db.kv.Commit(&writer)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a test table and insert a record first
	setupTestTable(t, db)
	insertTestRecord(t, db, 1)

	tests := []struct {
		name        string
		record      Record
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid update",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BYTES, Str: []byte("John Updated")},
					{Type: TYPE_BYTES, Str: []byte("updated@example.com")},
				},
			},
			expectError: false,
		},
		{
			name: "non-existent record",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 999},
					{Type: TYPE_BYTES, Str: []byte("NonExistent")},
					{Type: TYPE_BYTES, Str: []byte("none@example.com")},
				},
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "type mismatch",
			record: Record{
				Cols: []string{"id", "name", "email"},
				Vals: []Value{
					{Type: TYPE_BYTES, Str: []byte("invalid")},
					{Type: TYPE_BYTES, Str: []byte("John")},
					{Type: TYPE_BYTES, Str: []byte("john@example.com")},
				},
			},
			expectError: true,
			errorMsg:    "invalid type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer KVTX
			db.kv.Begin(&writer)
			updated, err := db.Update("users", tt.record, &writer)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !isEqual(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
				db.kv.Abort(&writer)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					db.kv.Abort(&writer)
				} else if !updated {
					t.Error("expected record to be updated")
					db.kv.Abort(&writer)
				} else {
					db.kv.Commit(&writer)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a test table and insert records
	setupTestTable(t, db)
	insertTestRecord(t, db, 1)
	insertTestRecord(t, db, 2)

	tests := []struct {
		name        string
		record      Record
		expectError bool
		errorMsg    string
		expectFound bool
	}{
		{
			name: "valid get",
			record: Record{
				Cols: []string{"id"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
				},
			},
			expectError: false,
			expectFound: true,
		},
		{
			name: "non-existent record",
			record: Record{
				Cols: []string{"id"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 999},
				},
			},
			expectError: false,
			expectFound: false,
		},
		{
			name: "invalid column",
			record: Record{
				Cols: []string{"invalid_col"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
				},
			},
			expectError: true,
			errorMsg:    "no index found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader KVReader
			db.kv.BeginRead(&reader)
			found, err := db.Get("users", &tt.record, &reader)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !isEqual(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if found != tt.expectFound {
					t.Errorf("expected found=%v, got %v", tt.expectFound, found)
				}
			}
			db.kv.EndRead(&reader)
		})
	}
}

func TestNewDataTypes(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Test FLOAT64, BOOLEAN, and DATETIME types
	tests := []struct {
		name     string
		tableDef *TableDef
		record   Record
	}{
		{
			name: "FLOAT64 table",
			tableDef: &TableDef{
				Name:  "prices",
				Types: []uint32{TYPE_INT64, TYPE_FLOAT64},
				Cols:  []string{"id", "price"},
				PKeys: 1,
			},
			record: Record{
				Cols: []string{"id", "price"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_FLOAT64, F64: 299.99},
				},
			},
		},
		{
			name: "BOOLEAN table",
			tableDef: &TableDef{
				Name:  "flags",
				Types: []uint32{TYPE_INT64, TYPE_BOOLEAN},
				Cols:  []string{"id", "active"},
				PKeys: 1,
			},
			record: Record{
				Cols: []string{"id", "active"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BOOLEAN, Bool: true},
				},
			},
		},
		{
			name: "DATETIME table",
			tableDef: &TableDef{
				Name:  "events",
				Types: []uint32{TYPE_INT64, TYPE_DATETIME},
				Cols:  []string{"id", "timestamp"},
				PKeys: 1,
			},
			record: Record{
				Cols: []string{"id", "timestamp"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_DATETIME, Time: time.Unix(1705320600, 0)},
				},
			},
		},
		{
			name: "Mixed types table",
			tableDef: &TableDef{
				Name:  "mixed",
				Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_FLOAT64, TYPE_BOOLEAN, TYPE_DATETIME},
				Cols:  []string{"id", "name", "score", "active", "created"},
				PKeys: 1,
			},
			record: Record{
				Cols: []string{"id", "name", "score", "active", "created"},
				Vals: []Value{
					{Type: TYPE_INT64, I64: 1},
					{Type: TYPE_BYTES, Str: []byte("Test Item")},
					{Type: TYPE_FLOAT64, F64: 95.5},
					{Type: TYPE_BOOLEAN, Bool: false},
					{Type: TYPE_DATETIME, Time: time.Unix(1705320600, 0)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer KVTX

			// Test table creation
			db.kv.Begin(&writer)
			err := db.TableNew(tt.tableDef, &writer)
			if err != nil {
				t.Fatalf("failed to create table: %v", err)
			}
			db.kv.Commit(&writer)

			// Test record insertion
			db.kv.Begin(&writer)
			inserted, err := db.Insert(tt.tableDef.Name, tt.record, &writer)
			if err != nil {
				t.Fatalf("failed to insert record: %v", err)
			}
			if !inserted {
				t.Fatal("record was not inserted")
			}
			db.kv.Commit(&writer)

			// Test record retrieval
			var reader KVReader
			db.kv.BeginRead(&reader)

			// Create a record with just the primary key for retrieval
			primaryKeyRecord := Record{
				Cols: []string{tt.tableDef.Cols[0]}, // Only primary key column
				Vals: []Value{tt.record.Vals[0]},    // Only primary key value
			}

			found, err := db.Get(tt.tableDef.Name, &primaryKeyRecord, &reader)
			db.kv.EndRead(&reader)

			if err != nil {
				t.Fatalf("failed to retrieve record: %v", err)
			}
			if !found {
				t.Fatal("record was not found")
			}

			// Verify all values match (primaryKeyRecord now contains the full record)
			for i, expectedVal := range tt.record.Vals {
				actualVal := primaryKeyRecord.Vals[i]
				if actualVal.Type != expectedVal.Type {
					t.Errorf("Type mismatch at index %d: expected %d, got %d", i, expectedVal.Type, actualVal.Type)
					continue
				}

				switch expectedVal.Type {
				case TYPE_INT64:
					if actualVal.I64 != expectedVal.I64 {
						t.Errorf("INT64 mismatch at index %d: expected %d, got %d", i, expectedVal.I64, actualVal.I64)
					}
				case TYPE_BYTES:
					if !bytes.Equal(actualVal.Str, expectedVal.Str) {
						t.Errorf("BYTES mismatch at index %d: expected %s, got %s", i, string(expectedVal.Str), string(actualVal.Str))
					}
				case TYPE_FLOAT64:
					if actualVal.F64 != expectedVal.F64 {
						t.Errorf("FLOAT64 mismatch at index %d: expected %f, got %f", i, expectedVal.F64, actualVal.F64)
					}
				case TYPE_BOOLEAN:
					if actualVal.Bool != expectedVal.Bool {
						t.Errorf("BOOLEAN mismatch at index %d: expected %t, got %t", i, expectedVal.Bool, actualVal.Bool)
					}
				case TYPE_DATETIME:
					if !actualVal.Time.Equal(expectedVal.Time) {
						t.Errorf("DATETIME mismatch at index %d: expected %v, got %v", i, expectedVal.Time, actualVal.Time)
					}
				}
			}
		})
	}
}

// Helper functions
func setupTestDB(t *testing.T) *DB {
	testPath := "test.db"

	testDB := &DB{
		Path:   testPath,
		kv:     *newKV(testPath),
		tables: make(map[string]*TableDef),
		pool:   NewPool(3),
	}

	if err := testDB.kv.Open(); err != nil {
		log.Fatalf("Failed to open  %v", err)
	}
	err := initializeInternalTables(testDB)
	if err != nil {
		if !errors.Is(err, ErrTableAlreadyExists) {
			fmt.Println("Error while init table: ", err)
			os.Exit(0)
		}
	}
	return testDB
}

func cleanupTestDB(t *testing.T, db *DB) {
	db.kv.Close()
	if err := os.Remove(db.Path); err != nil {
		t.Errorf("failed to cleanup test database: %v", err)
	}
}

func setupTestTable(t *testing.T, db *DB) {
	var writer KVTX
	db.kv.Begin(&writer)

	tdef := &TableDef{
		Name:  "users",
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES},
		Cols:  []string{"id", "name", "email"},
		PKeys: 1,
	}

	if err := db.TableNew(tdef, &writer); err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	db.kv.Commit(&writer)
}

func insertTestRecord(t *testing.T, db *DB, id int64) {
	var writer KVTX
	db.kv.Begin(&writer)

	record := Record{
		Cols: []string{"id", "name", "email"},
		Vals: []Value{
			{Type: TYPE_INT64, I64: id},
			{Type: TYPE_BYTES, Str: []byte("John")},
			{Type: TYPE_BYTES, Str: []byte("john@example.com")},
		},
	}

	inserted, err := db.Insert("users", record, &writer)
	if err != nil {
		t.Fatalf("failed to insert test record: %v", err)
	}
	if !inserted {
		t.Fatal("failed to insert test record")
	}
	db.kv.Commit(&writer)
}

func isEqual(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
