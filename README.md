# FiloDB

A lightweight relational database system written in Go, focusing on core database concepts and storage management.

## Table of Contents

- [Installation](#installation)
- [Features](#features)
- [Supported Commands](#supported-commands)
- [Usage](#usage)
- [License](#license)

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.17 or later)
- [Git](https://git-scm.com/downloads)
- Compatible with Linux and macOS

### Build

Build the project using Go:

```bash
go mod tidy
go build -o filodb
```

### Run

To start the FiloDB server:

```bash
./filodb
```

## Features

- **B+ Tree Storage Engine**: Fast data retrieval with indexing support
- **Transaction Support**: ACID-compliant transactions with rollback capabilities
- **Concurrent Reads**: Multiple simultaneous read operations
- **Free List Management**: Efficient storage space management and reuse
- **Memory-Mapped I/O**: Optimized file operations for better performance

## Supported Commands

- **CREATE** - Create a new table
- **INSERT** - Add records to a table
- **GET** - Retrieve records from a table
- **UPDATE** - Update existing records
- **DELETE** - Delete records from a table
- **BEGIN** - Start a transaction
- **COMMIT** - Commit a transaction
- **ABORT** - Rollback a transaction

## Usage

### Creating a Table
```
> create
Enter table name: users
Enter column names (comma-separated): id,name,email
Enter column types (comma-separated as numbers): 1,2,2
Enter indexes (format: col1+col2,col3, ... or leave empty): 
```

**Data Types:**
- `1` = INT64 (integers)
- `2` = BYTES (strings)

### Inserting Data
```
> insert
Enter table name: users
Enter value for id: 1
Enter value for name: John Doe
Enter value for email: john@example.com
```

### Querying Data
```
> get
Enter table name: users
Select query type:
1. Index lookup (primary/secondary index)
2. Range query
3. Column filter
Enter choice (1, 2 or 3): 1
```

## Testing

Run the test suite:

```bash
go test -v ./database/
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# FiloDB-
