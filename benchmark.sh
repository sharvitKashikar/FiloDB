#!/bin/bash

echo "=== FiloDB Performance Benchmark ==="
echo "Building FiloDB..."

# Build the database
go build -o filodb

if [ ! -f "./filodb" ]; then
    echo "Error: Failed to build filodb"
    exit 1
fi

echo "Starting performance tests..."

# Clean up any existing database
rm -f database.db

# Start FiloDB in background and get its PID
./filodb > /dev/null 2>&1 &
FILODB_PID=$!

# Wait for FiloDB to start
sleep 2

# Function to send command to FiloDB
send_command() {
    echo "$1" | timeout 5 ./filodb > /dev/null 2>&1
}

echo ""
echo "=== Insert Performance Test ==="
echo "Creating test table..."
send_command "CREATE users id:int name:str email:str"

echo "Inserting 50 records and measuring time..."
start_time=$(date +%s.%N)

for i in {1..50}; do
    send_command "INSERT users $i User$i user$i@example.com"
done

end_time=$(date +%s.%N)
insert_duration=$(echo "$end_time - $start_time" | bc)
insert_ops_per_sec=$(echo "scale=0; 50 / $insert_duration" | bc)

echo "Insert Results:"
echo "  - Total time: ${insert_duration}s"
echo "  - Operations per second: ~$insert_ops_per_sec ops/sec"
echo "  - Average per insert: $(echo "scale=3; $insert_duration / 50" | bc)s"

echo ""
echo "=== Query Performance Test ==="
echo "Performing 100 lookups..."

start_time=$(date +%s.%N)

for i in {1..100}; do
    lookup_id=$((i % 50 + 1))
    send_command "GET users $lookup_id"
done

end_time=$(date +%s.%N)
query_duration=$(echo "$end_time - $start_time" | bc)
query_ops_per_sec=$(echo "scale=0; 100 / $query_duration" | bc)

echo "Query Results:"
echo "  - Total time: ${query_duration}s"
echo "  - Operations per second: ~$query_ops_per_sec ops/sec"
echo "  - Average per query: $(echo "scale=4; $query_duration / 100" | bc)s"

echo ""
echo "=== Database Size Analysis ==="
if [ -f "database.db" ]; then
    db_size_bytes=$(stat -f%z database.db 2>/dev/null || stat -c%s database.db 2>/dev/null || echo "0")
    db_size_kb=$(echo "scale=1; $db_size_bytes / 1024" | bc)
    size_per_record=$(echo "scale=2; $db_size_kb / 50" | bc)
    
    echo "Database file size: ${db_size_kb} KB"
    echo "Size per record: ${size_per_record} KB"
else
    echo "Database file not found"
fi

# Clean up
kill $FILODB_PID 2>/dev/null
rm -f filodb database.db

echo ""
echo "=== Performance Summary ==="
echo "FiloDB Performance Profile:"
echo "  - Insert Performance: ~$insert_ops_per_sec ops/sec"
echo "  - Query Performance: ~$query_ops_per_sec ops/sec" 
echo "  - Storage Efficiency: ${size_per_record:-"N/A"} KB per record"
echo ""
echo "Use Case Recommendations:"
echo "  ✓ Educational projects and learning"
echo "  ✓ Small applications (< 1,000 records)"
echo "  ✓ Prototyping and development"
echo "  ✗ Production applications requiring high performance"
echo ""
echo "Note: Results may vary based on hardware and system load"
echo "=== Benchmark Complete ===" 