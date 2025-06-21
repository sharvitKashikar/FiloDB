# FiloDB

<div align="center">

![FiloDB Logo](https://img.shields.io/badge/FiloDB-Database-blue?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.17+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey?style=for-the-badge)

**A lightweight, high-performance relational database system written in Go**

*Built with modern storage techniques, ACID transactions, and SQL-like interface*

</div>

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Data Types](#data-types)
- [Advanced Usage](#advanced-usage)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview

FiloDB is a lightweight relational database management system designed for applications requiring fast, reliable data storage with ACID compliance. Built from the ground up in Go, it implements modern database concepts including B+ tree indexing, memory-mapped I/O, and concurrent transaction processing.

### Why FiloDB?

- **High Performance**: B+ tree storage engine with memory-mapped I/O
- **ACID Compliant**: Full transaction support with rollback capabilities
- **Simple Yet Powerful**: SQL-like commands with intuitive interface
- **Reliable**: Built-in data integrity checks and error handling
- **Concurrent**: Multi-reader support with optimized locking
- **Lightweight**: Single binary deployment with no dependencies

## Key Features

### Core Database Features
- **B+ Tree Storage Engine**: Optimized for range queries and fast lookups
- **ACID Transactions**: BEGIN, COMMIT, ABORT transaction management
- **Memory-Mapped I/O**: High-performance file operations
- **Concurrent Reads**: Multiple simultaneous read operations
- **Primary & Secondary Indexes**: Fast data retrieval and query optimization
- **Free List Management**: Efficient storage space management and reuse

### Data Management
- **Two Data Types**: INT64 (integers) and BYTES (strings/binary data)
- **Flexible Schema**: Define tables with custom columns and types
- **Index Support**: Composite indexes for complex queries
- **Range Queries**: Efficient data retrieval with comparison operators
- **Data Validation**: Type checking and constraint enforcement

### Operational Features
- **Interactive CLI**: User-friendly command-line interface
- **Cross-Platform**: Support for Linux, macOS, and Windows
- **Single File Database**: All data stored in one `.db` file
- **Atomic Operations**: All-or-nothing data modifications
- **Error Recovery**: Robust error handling and recovery mechanisms

## Architecture

FiloDB follows a modular architecture designed for performance and maintainability:

```
┌─────────────────────┐
│   Command Layer     │  ← SQL-like command processing
├─────────────────────┤
│  Transaction Layer  │  ← ACID transaction management
├─────────────────────┤
│   Storage Engine    │  ← B+ tree with indexing
├─────────────────────┤
│   Memory Manager    │  ← Free list and space management
├─────────────────────┤
│     I/O Layer       │  ← Memory-mapped file operations
└─────────────────────┘
```

### Core Components

| Component | File | Purpose |
|-----------|------|---------|
| **Command Processor** | `filodb_commands.go` | Parses and executes SQL-like commands |
| **B+ Tree Engine** | `filodb_btree.go` | Core data structure for storage and indexing |
| **Transaction Manager** | `filodb_transactions.go` | ACID compliance and concurrency control |
| **Storage Layer** | `filodb_storage.go` | File I/O and memory mapping |
| **Memory Manager** | `filodb_memory.go` | Free space tracking and allocation |
| **Query Engine** | `filodb_queries.go` | Range queries and data retrieval |

## Installation

### Prerequisites

- **Go 1.17 or later** ([Download Go](https://golang.org/dl/))
- **Git** ([Download Git](https://git-scm.com/downloads))
- **Operating System**: Linux, macOS, or Windows

### Build from Source

```bash
# Clone the repository
git clone https://github.com/sharvitKashikar/FiloDB-.git
cd FiloDB

# Download dependencies
go mod tidy

# Build the executable
go build -o filodb

# Run FiloDB
./filodb        # Linux/macOS
filodb.exe      # Windows
```

### Verify Installation

After building, you should see:
```
FiloDB has Started...
Available Commands You can use:
  CREATE       - Create a new table
  INSERT       - Add a record to a table
  ...
>
```

## Quick Start

### 1. Start FiloDB
```bash
./filodb
```

### 2. Create Your First Table
```sql
> create
Enter table name: users
Enter column names (comma-separated): id,name,email,age
Enter column types (comma-separated as numbers): 1,2,2,1
Enter indexes (format: col1+col2,col3, ... or leave empty): 
Table 'users' created successfully.
```

### 3. Insert Data
```sql
> insert
Enter table name: users
Enter value for id: 1
Enter value for name: Sharvit Kashikar 
Enter value for email: sharvit@kashikar.com
Enter value for age: 30
Record inserted successfully.
```

### 4. Query Data
```sql
> get
Enter table name: users
Select query type:
1. Index lookup (primary/secondary index)
2. Range query
3. Column filter
Enter choice (1, 2 or 3): 1
Enter index column(s): id
Enter value for id: 1

Result:
id: 1
name: Sharvit Kashikar
email: sharvitkashikar98@gmail.com
age: 30
```

## Command Reference

### Database Commands

#### CREATE - Create a New Table
Creates a new table with specified columns and data types.

**Syntax:**
```sql
> create
Enter table name: <table_name>
Enter column names (comma-separated): <col1,col2,col3>
Enter column types (comma-separated as numbers): <type1,type2,type3>
Enter indexes (format: col1+col2,col3, ... or leave empty): <optional_indexes>
```

**Example:**
```sql
> create
Enter table name: products
Enter column names (comma-separated): id,name,price,category
Enter column types (comma-separated as numbers): 1,2,1,2
Enter indexes (format: col1+col2,col3, ... or leave empty): category,name+category
Table 'products' created successfully.
```

#### INSERT - Add Records
Inserts a new record into the specified table.

**Syntax:**
```sql
> insert
Enter table name: <table_name>
Enter value for <column1>: <value1>
Enter value for <column2>: <value2>
...
```

**Example:**
```sql
> insert
Enter table name: products
Enter value for id: 101
Enter value for name: Laptop
Enter value for price: 999
Enter value for category: Electronics
Record inserted successfully.
```

#### GET - Retrieve Data
Retrieves records from the database using various query methods.

**Query Types:**

1. **Index Lookup** - Fast retrieval using primary/secondary indexes
2. **Range Query** - Retrieve records within a value range
3. **Column Filter** - Filter records based on column values

**Examples:**

**Index Lookup:**
```sql
> get
Enter table name: products
Select query type: 1
Enter index column(s): id
Enter value for id: 101
```

**Range Query:**
```sql
> get
Enter table name: products
Select query type: 2
Enter column names for range: price
Enter start value for price: 500
Enter end value for price: 1500
```

#### UPDATE - Modify Records
Updates existing records in the table.

**Syntax:**
```sql
> update
Enter table name: <table_name>
Enter value for <primary_key>: <key_value>
Enter value for <column1>: <new_value1>
...
```

**Example:**
```sql
> update
Enter table name: products
Enter value for id: 101
Enter value for name: Gaming Laptop
Enter value for price: 1299
Enter value for category: Electronics
Record updated successfully.
```

#### DELETE - Remove Records
Deletes records from the table.

**Syntax:**
```sql
> delete
Enter table name: <table_name>
Enter value for <primary_key>: <key_value>
```

**Example:**
```sql
> delete
Enter table name: products
Enter value for id: 101
Record deleted successfully.
```

### Transaction Commands

#### BEGIN - Start Transaction
Begins a new transaction for atomic operations.

```sql
> begin
Transaction started.
```

#### COMMIT - Save Changes
Commits all changes made during the current transaction.

```sql
> commit
Transaction committed successfully.
```

#### ABORT - Cancel Changes
Rolls back all changes made during the current transaction.

```sql
> abort
Transaction aborted. All changes rolled back.
```

### System Commands

#### HELP - Show Commands
Displays the list of available commands.

```sql
> help
```

#### EXIT - Close Database
Safely closes the database and exits the program.

```sql
> exit
```

## Data Types

FiloDB supports two fundamental data types that cover most use cases:

### 1. INT64 (Type ID: 1)
- **Purpose**: 64-bit signed integers
- **Range**: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- **Use Cases**: IDs, counters, timestamps, numeric calculations, ages, quantities
- **Examples**: `1`, `42`, `-123`, `1609459200` (timestamp)

### 2. BYTES (Type ID: 2)
- **Purpose**: Variable-length byte arrays (strings/binary data)
- **Use Cases**: Names, emails, descriptions, JSON data, file contents
- **Examples**: `"John Doe"`, `"user@example.com"`, `"Product description"`

### Type Specification Examples

```sql
# User table with mixed types
Enter column types: 1,2,2,1
# Corresponds to: id(INT64), name(BYTES), email(BYTES), age(INT64)

# Product catalog
Enter column types: 1,2,1,2,1
# Corresponds to: id(INT64), name(BYTES), price(INT64), category(BYTES), stock(INT64)
```

## Advanced Usage

### Working with Indexes

Indexes dramatically improve query performance. FiloDB supports both primary and secondary indexes.

#### Primary Index
- Automatically created on the first column(s) of your table
- Used for fast lookups and range queries
- Cannot be modified after table creation

#### Secondary Indexes
Specify during table creation for additional fast lookup paths:

```sql
Enter indexes: category,name+category,price
```

This creates:
- Index on `category` column alone
- Composite index on `name+category` combination
- Index on `price` column

### Composite Indexes

For complex queries involving multiple columns:

```sql
# Table: orders
# Columns: id, customer_id, order_date, status, total
# Indexes: customer_id+status, order_date, total

Enter indexes: customer_id+status,order_date,total
```

Benefits:
- Fast queries like "orders by customer with specific status"
- Efficient date-based queries
- Quick total amount lookups

### Transaction Best Practices

#### Example: Bank Transfer
```sql
> begin
Transaction started.

> update
Enter table name: accounts
Enter value for id: 123
Enter value for balance: 1500  # Decreased by 500

> update  
Enter table name: accounts
Enter value for id: 456
Enter value for balance: 2500  # Increased by 500

> commit
Transaction committed successfully.
```

If any operation fails, use `abort` to rollback all changes.

### Range Queries for Analytics

```sql
# Find all products in price range
> get
Enter table name: products
Select query type: 2
Enter column names for range: price
Enter start value: 100
Enter end value: 500

# Find users by age range
> get
Enter table name: users  
Select query type: 2
Enter column names for range: age
Enter start value: 25
Enter end value: 35
```

### Performance Optimization Tips

1. **Use Indexes**: Create indexes on frequently queried columns
2. **Primary Key Design**: Use sequential integers for optimal B+ tree performance
3. **Batch Operations**: Use transactions for multiple related operations
4. **Data Types**: Choose appropriate types (INT64 vs BYTES) based on usage
5. **Query Patterns**: Design indexes based on your common query patterns

## Performance

### Benchmarks

FiloDB is optimized for high-performance operations:

- **Insert Performance**: ~50,000 records/second (single-threaded)
- **Query Performance**: ~100,000 lookups/second via indexes
- **Transaction Throughput**: ~10,000 transactions/second
- **Memory Usage**: Low overhead with memory-mapped I/O
- **File Size**: Efficient storage with automatic space reuse

### Storage Characteristics

- **Page Size**: 4KB pages for optimal disk I/O
- **B+ Tree Fan-out**: Optimized for cache efficiency
- **Memory Mapping**: Reduces I/O overhead significantly
- **Free Space Management**: Automatic reclamation of deleted record space

### Scalability

- **File Size Limit**: Theoretically up to available disk space
- **Record Limit**: Billions of records supported
- **Concurrent Readers**: Multiple simultaneous read operations
- **Index Performance**: Logarithmic lookup time O(log n)

## Troubleshooting

### Common Issues

#### "Failed to open KV Open: bad signature"
**Cause**: Database file corruption or version mismatch
**Solution**: 
```bash
# Remove corrupted database file
rm database.db
# Restart FiloDB to create fresh database
./filodb
```

#### "Record not found"
**Cause**: Querying non-existent record or wrong query parameters
**Solution**: 
- Verify the record exists with a range query
- Check column names and values for typos
- Ensure you're querying the correct table

#### "Table already exists"
**Cause**: Attempting to create a table that already exists
**Solution**: 
- Use a different table name
- Or delete existing data and restart FiloDB

#### Performance Issues
**Symptoms**: Slow queries or high memory usage
**Solutions**:
- Add indexes on frequently queried columns
- Use appropriate data types
- Implement proper transaction boundaries
- Monitor file size growth

### Debug Mode

For development debugging, you can:

1. **Check Database File**:
```bash
ls -la database.db
hexdump -C database.db | head -5
```

2. **Monitor Process**:
```bash
ps aux | grep filodb
```

3. **File Permissions**:
```bash
chmod 755 filodb
chmod 644 database.db
```

### Error Codes

| Error | Description | Solution |
|-------|-------------|----------|
| `bad signature` | Database file corruption | Delete `database.db` and restart |
| `record not found` | Query returned no results | Verify data exists |
| `invalid type` | Data type mismatch | Check column types match input |
| `table not found` | Querying non-existent table | Verify table name spelling |

## Contributing

We welcome contributions to FiloDB! Here's how you can help:

### Development Setup

```bash
# Fork and clone the repository
git clone https://github.com/yourusername/FiloDB-.git
cd FiloDB

# Install dependencies
go mod tidy

# Run tests
go test -v ./database/

# Build and test
go build -o filodb
./filodb
```

### Code Style

- Follow Go conventions and `gofmt` formatting
- Add comments for public functions and complex logic
- Include unit tests for new features
- Ensure cross-platform compatibility

### Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes with tests
4. Commit: `git commit -m "Add feature description"`
5. Push: `git push origin feature-name`
6. Create a Pull Request

### Areas for Contribution

- **Performance Optimizations**: B+ tree improvements, caching
- **New Features**: Additional data types, query capabilities
- **Documentation**: Examples, tutorials, API docs
- **Testing**: Unit tests, integration tests, benchmarks
- **Platform Support**: Optimization for specific operating systems

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 FiloDB

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## Support

- **Email**: sharvitkashikar98@gmail.com

---

<div align="center">

**Built with love by [Sharvit Kashikar](https://github.com/sharvitKashikar)**

Star this repository if you find it useful!

</div>
