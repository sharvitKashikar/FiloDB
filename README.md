<p align="center">
  <img src="Assets/FiloDBLogo.png" alt="FiloDB Logo" />
</p>

# FiloDB

<div align="center">

![FiloDB Logo](https://img.shields.io/badge/FiloDB-Database-blue?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)
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

## **What Makes FiloDB Different**

| Feature | FiloDB | Other Go Databases |
|---------|--------|-------------------|
| **Dependencies** | Only `golang.org/x/sys` | Multiple external deps |
| **CLI Experience** | Interactive with aggregate functions | Basic command interface |
| **Query Features** | SQL-like with range queries | Standard SQL only |
| **Performance Metrics** | Built-in stats and monitoring | Limited visibility |
| **Cross-Platform** | Optimized memory mapping per OS | Basic compatibility |
| **Worker Pools** | Background processing support | Basic threading |

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
- **Five Data Types**: INT64 (integers), BYTES (strings/binary), FLOAT64 (decimals), BOOLEAN (true/false), DATETIME (timestamps)
- **Flexible Schema**: Define tables with custom columns and types
- **Index Support**: Composite indexes for complex queries
- **Range Queries**: Efficient data retrieval with comparison operators
- **Data Validation**: Type checking and constraint enforcement

### Operational Features
- **Interactive CLI**: User-friendly command-line interface
- **Aggregate Functions**: Built-in COUNT, SUM, AVG, MIN, MAX operations
- **Data Analysis Tools**: SCAN and DEBUG commands for table inspection
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

- **Go 1.23 or later** ([Download Go](https://golang.org/dl/))
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
Enter value for name: Arjun Singh
Enter value for email: arjun.singh@techcorp.in
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
Enter index column(s) (comma-separated for composite index): id
Enter value for id: 1

Result:
id: 1
name: Arjun Singh
email: arjun.singh@techcorp.in
age: 30
```

### 5. Advanced Example with New Data Types
```sql
# Create a products table with all data types
> create
Enter table name: products
Enter column names (comma-separated): id,name,price,active,created_at
Enter column types (comma-separated as numbers): 1,2,3,4,5
Enter indexes (format: col1+col2,col3, ... or leave empty): active,created_at
Table 'products' created successfully.

# Insert a product with new data types
> insert
Enter table name: products
Enter value for id: 101
Enter value for name: Wireless Headphones
Enter value for price: 299.99
Enter value for active: true
Enter value for created_at: 2024-03-15 09:30:00
Record inserted successfully.

# Query products by date range
> get
Enter table name: products
Select query type: 2
Enter column name for range lookup(index col): created_at
Enter start range value: 2024-03-01
Enter end range value: 2024-03-31

Result:
id: 101
name: Wireless Headphones
price: 299.990000
active: true
created_at: 2024-03-15 09:30:00
```

## Performance

FiloDB includes comprehensive benchmarking tools and real-world performance metrics:

### Running Benchmarks

```bash
# Run automated performance benchmark
./benchmark.sh
```

### Performance Documentation

For detailed performance analysis and optimization guides, see [PERFORMANCE.md](./PERFORMANCE.md)

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
Enter column types (comma-separated as numbers): 1,2,3,2
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
Enter value for name: Gaming Laptop
Enter value for price: 75999.50
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
Enter index column(s) (comma-separated for composite index): id
Enter value for id: 101
```

**Range Query:**
```sql
> get
Enter table name: products
Select query type: 2
Enter column name for range lookup(index col): price
Enter start range value: 500
Enter end range value: 1500
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
Enter value for name: Premium Gaming Laptop
Enter value for price: 89999.00
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

### Aggregate Functions

FiloDB includes a powerful set of aggregate functions for data analysis and reporting. These functions let you perform calculations across multiple records in your tables without needing external tools.

#### COUNT - Count Records
Counts the total number of records in a table.

**Syntax:**
```sql
> count
Enter table name: <table_name>
```

**Example:**
```sql
> count
Enter table name: sales
Table: sales
Count: 1,250
```

This is particularly useful for getting quick insights into data volume and checking if records exist.

#### SUM - Calculate Totals
Adds up all numeric values in a specified column.

**Syntax:**
```sql
> sum
Enter table name: <table_name>
Enter column name for SUM: <numeric_column>
```

**Example:**
```sql
> sum
Enter table name: sales
Enter column name for SUM: amount
Table: sales
Column: amount
Records processed: 1,250
SUM(amount): 485,750
```

**Supported Types**: INT64, FLOAT64

Great for calculating totals like revenue, quantities, or any numeric aggregations.

#### AVG - Find Averages
Calculates the average value of a numeric column.

**Syntax:**
```sql
> avg
Enter table name: <table_name>
Enter column name for AVG: <numeric_column>
```

**Example:**
```sql
> avg
Enter table name: products
Enter column name for AVG: price
Table: products
Column: price
Records processed: 45
AVG(price): 324.500000
```

**Supported Types**: INT64, FLOAT64

This helps you understand typical values like average order amounts or product prices.

#### MIN - Find Minimum Values
Finds the smallest value in a column (works with all data types).

**Syntax:**
```sql
> min
Enter table name: <table_name>
Enter column name for MIN: <column_name>
```

**Example:**
```sql
> min
Enter table name: products
Enter column name for MIN: price
Table: products
Column: price
Records processed: 45
MIN(price): 12.500000
```

**Supported Types**: All types (INT64, BYTES, FLOAT64, BOOLEAN, DATETIME)

#### MAX - Find Maximum Values
Finds the largest value in a column (works with all data types).

**Syntax:**
```sql
> max
Enter table name: <table_name>
Enter column name for MAX: <column_name>
```

**Example:**
```sql
> max
Enter table name: sales
Enter column name for MAX: amount
Table: sales
Column: amount
Records processed: 1,250
MAX(amount): 2500.000000
```

**Supported Types**: All types (INT64, BYTES, FLOAT64, BOOLEAN, DATETIME)

### Utility Commands

#### SCAN - View All Records
Shows every record in a table with clean, formatted output.

**Syntax:**
```sql
> scan
Enter table name: <table_name>
```

**Example:**
```sql
> scan
Enter table name: users
Record 1: id=1, name=Arjun Patel, email=arjun@techsolutions.com, age=28
Record 2: id=2, name=Priya Sharma, email=priya@innovatetech.in, age=32
...
Total records: 5
```

This command is really useful for debugging and seeing all your data at once.

#### DEBUG - Table Information
Displays detailed information about a table's structure and contents.

**Syntax:**
```sql
> debug
Enter table name: <table_name>
```

**Example:**
```sql
> debug
Enter table name: products
Table: products
Columns: [id, name, price, category]
Types: [INT64, BYTES, INT64, BYTES]
Total records found: 45
Sample records displayed above
```

Helps you understand table schemas and verify data integrity.

### Practical Examples

**Monthly Sales Analysis for Mumbai Store:**
```sql
# Get total sales count
> count
Enter table name: mumbai_sales

# Calculate total revenue  
> sum
Enter table name: mumbai_sales
Enter column name for SUM: revenue

# Find average order value
> avg
Enter table name: mumbai_sales
Enter column name for AVG: order_value
```

**Delhi Warehouse Inventory Insights:**
```sql
# Check minimum stock levels
> min
Enter table name: delhi_inventory
Enter column name for MIN: stock_level

# Find most expensive item
> max
Enter table name: delhi_inventory
Enter column name for MAX: unit_price

# Review all products
> scan
Enter table name: delhi_inventory
```

These aggregate functions work efficiently with FiloDB's B+ tree storage engine, making data analysis fast even with thousands of records.

## Data Types

FiloDB supports five fundamental data types that cover most use cases:

### 1. INT64 (Type ID: 1)
- **Purpose**: 64-bit signed integers
- **Range**: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- **Use Cases**: IDs, counters, timestamps, numeric calculations, ages, quantities
- **Examples**: `1`, `42`, `-123`, `1609459200` (timestamp)

### 2. BYTES (Type ID: 2)
- **Purpose**: Variable-length byte arrays (strings/binary data)
- **Use Cases**: Names, emails, descriptions, JSON data, file contents
- **Examples**: `"Ananya Patel"`, `"ananya@startupindia.com"`, `"Premium smartphone with dual camera"`

### 3. FLOAT64 (Type ID: 3)
- **Purpose**: 64-bit floating-point numbers
- **Range**: IEEE 754 double precision floating point
- **Use Cases**: Prices, percentages, scientific calculations, measurements
- **Examples**: `3.14159`, `1299.50`, `45999.99`, `2.5e6` (for scientific calculations)

### 4. BOOLEAN (Type ID: 4)
- **Purpose**: True/false values
- **Use Cases**: Flags, status indicators, yes/no questions, active/inactive states
- **Input Formats**: `true`/`false`, `1`/`0`, `yes`/`no`, `y`/`n` (case insensitive)
- **Examples**: `true`, `false`, `1`, `0`

### 5. DATETIME (Type ID: 5)
- **Purpose**: Date and time values
- **Storage**: Unix timestamp (seconds since epoch)
- **Use Cases**: Created dates, timestamps, scheduling, logging
- **Input Formats**:
  - `2024-01-15 14:30:00` (YYYY-MM-DD HH:MM:SS)
  - `2024-01-15` (YYYY-MM-DD, defaults to 00:00:00)
  - `2024-01-15T14:30:00Z` (ISO 8601)
  - Unix timestamps as integers (e.g., `1705320600`)
- **Display Format**: `2024-01-15 14:30:00`

### Type Specification Examples

```sql
# Customer database for Indian e-commerce
Enter column types: 1,2,2,1,4
# Corresponds to: id(INT64), name(BYTES), email(BYTES), age(INT64), premium_member(BOOLEAN)

# Product catalog for electronics store
Enter column types: 1,2,3,2,1,5
# Corresponds to: id(INT64), name(BYTES), price(FLOAT64), category(BYTES), stock(INT64), created_at(DATETIME)

# Order tracking system
Enter column types: 1,2,5,4,3
# Corresponds to: order_id(INT64), customer_name(BYTES), order_date(DATETIME), delivered(BOOLEAN), amount(FLOAT64)
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

#### Example: Bank Transfer (Priya to Rahul)
```sql
> begin
Transaction started.

> update
Enter table name: accounts
Enter value for id: 123
Enter value for account_holder: Priya Sharma
Enter value for balance: 25000  # Decreased by 5000

> update  
Enter table name: accounts
Enter value for id: 456
Enter value for account_holder: Rahul Gupta
Enter value for balance: 35000  # Increased by 5000

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
Enter column name for range lookup(index col): price
Enter start range value: 1000
Enter end range value: 50000

# Find users by age range
> get
Enter table name: users  
Select query type: 2
Enter column name for range lookup(index col): age
Enter start range value: 25
Enter end range value: 35
```

### Performance Optimization Tips

1. **Use Indexes**: Create indexes on frequently queried columns
2. **Primary Key Design**: Use sequential integers for optimal B+ tree performance
3. **Batch Operations**: Use transactions for multiple related operations
4. **Data Types**: Choose appropriate types (INT64 vs BYTES) based on usage
5. **Query Patterns**: Design indexes based on your common query patterns

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

## References and Learning Resources

This project's implementation was inspired by and learned from various resources in the database development community:

- "Build Your Own Database From Scratch in Go" - A foundational resource that provided valuable insights into database implementation patterns
- Database Internals: A Deep Dive into How Distributed Data Systems Work
- Designing Data-Intensive Applications by Martin Kleppmann

While the core concepts and patterns are similar to standard database implementations, FiloDB includes significant additional features and enhancements:

**Extended Data Type System:**
- FLOAT64 support with IEEE 754 binary encoding for precise decimal calculations
- BOOLEAN data type with flexible input formats (true/false, 1/0, yes/no, y/n)
- DATETIME support with multiple input formats and UTC timezone consistency
- Enhanced serialization/deserialization for all five data types

**Advanced Query Capabilities:**
- Range queries on indexed columns with optimized B+ tree scanning
- Composite index support for complex multi-column queries
- Enhanced aggregate functions (SUM/AVG support for both INT64 and FLOAT64)
- MIN/MAX operations working across all data types including datetime comparisons

**Robust Command Interface:**
- Interactive CLI with intuitive menu-driven operations
- Enhanced error handling and user-friendly error messages
- Comprehensive table scanning and debugging utilities
- Safe extension approach maintaining backward compatibility

**Performance and Reliability:**
- Fixed scanner initialization preventing data corruption
- Consistent timezone handling for datetime operations
- Memory-safe range query implementation
- Comprehensive test coverage for all data types and operations

---

## Support

- **Email**: sharvitkashikar98@gmail.com

---

<div align="center">

**Built with love by [Sharvit Kashikar](https://github.com/sharvitKashikar)**

Star this repository if you find it useful!

</div>
