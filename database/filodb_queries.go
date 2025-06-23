package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	CMP_GE = +3 // >=
	CMP_GT = +2 // >
	CMP_LT = -2 // <
	CMP_LE = -3 // <=
)

// the iterator for range queries
type Scanner struct {
	// the range, from Key1 to Key2
	db      *DB
	indexNo int // -1: use primary key; >= 0: use an index
	Cmp1    int
	Cmp2    int
	Key1    Record
	Key2    Record
	// internal
	tdef     *TableDef
	iter     *BIter // underlying BTree iterator
	keyEnd   []byte // the encoded Key2
	keyStart []byte // the encoded Key2
}

func (db *DB) Scan(table string, req *Scanner, tree *BTree) error {
	tdef := GetTableDef(db, table, tree)
	if tdef == nil {
		return fmt.Errorf("table not found: %s", table)
	}
	return dbScan(db, tdef, req, tree)
}

func dbScan(db *DB, tdef *TableDef, req *Scanner, tree *BTree) error {
	// sanity checks
	switch {
	case req.Cmp1 > 0 && req.Cmp2 < 0:
	case req.Cmp2 > 0 && req.Cmp1 < 0:
	default:
		return fmt.Errorf("bad range")
	}
	indexNo, err := findIndex(tdef, req.Key1.Cols)
	if err != nil {
		return err
	}
	index, prefix := tdef.Cols[:tdef.PKeys], tdef.Prefix
	if indexNo >= 0 {
		index, prefix = tdef.Indexes[indexNo], tdef.IndexPrefix[indexNo]
	}

	req.db = db

	req.tdef = tdef
	req.indexNo = indexNo
	// seek to the start key
	req.keyStart = encodeKeyPartial(nil, prefix, req.Key1.Vals, tdef, index, req.Cmp1)
	req.keyEnd = encodeKeyPartial(nil, prefix, req.Key2.Vals, tdef, index, req.Cmp2)
	req.iter = tree.Seek(req.keyStart, req.Cmp1)
	return nil
}

func (sc *Scanner) Valid() bool {
	if !sc.iter.Valid() {
		return false
	}
	key, _ := sc.iter.Deref()

	// First check if we've reached the end of valid keys
	if bytes.Compare(key, sc.keyEnd) > 0 {
		return false
	}

	// Check if we're still within range
	if bytes.Compare(key, sc.keyStart) < 0 {
		return false
	}

	return true
}

func (sc *Scanner) Next() {
	if !sc.iter.Valid() {
		return
	}

	currentKey, _ := sc.iter.Deref()
	sc.iter.Next()

	// If after moving Next(), we get the same key or invalid iterator,
	// we've reached the end of valid data
	if !sc.iter.Valid() {
		return
	}

	nextKey, _ := sc.iter.Deref()
	if bytes.Equal(currentKey, nextKey) {
		// If we get the same key, invalidate iterator to stop
		sc.iter = &BIter{}
		return
	}
}

// fetch the current row
func (sc *Scanner) Deref(rec *Record, tree *BTree) {
	if !sc.Valid() {
		return
	}
	tdef := sc.tdef
	rec.Cols = tdef.Cols
	rec.Vals = rec.Vals[:0]
	key, val := sc.iter.Deref()
	if sc.indexNo < 0 {
		values := make([]Value, len(rec.Cols))
		for i := range rec.Cols {
			values[i].Type = tdef.Types[i]
		}
		decodeValues(key[4:], values[:tdef.PKeys])
		decodeValues(val, values[tdef.PKeys:])
		rec.Vals = append(rec.Vals, values...)
	} else {
		index := tdef.Indexes[sc.indexNo]
		ival := make([]Value, len(index))
		for i, col := range index {
			ival[i].Type = tdef.Types[ColIndex(tdef, col)]
		}
		decodeValues(key[4:], ival)
		icol := Record{index, ival}

		rec.Cols = rec.Cols[:tdef.PKeys]
		for _, col := range rec.Cols {
			rec.Vals = append(rec.Vals, *icol.Get(col))
		}

		ok, err := dbGet(sc.db, tdef, rec, tree)
		if !ok && err != nil {
			fmt.Println("Error getting record from DB")
		}
	}
}

// B-Tree Iterator
type BIter struct {
	tree *BTree
	path []BNode  // from root to leaf
	pos  []uint16 // indexes into nodes
}

// get current KV pair
func (iter *BIter) Deref() (key []byte, val []byte) {
	currentNode := iter.path[len(iter.path)-1]
	idx := iter.pos[len(iter.pos)-1]
	key = currentNode.getKey(idx)
	val = currentNode.getVal(idx)
	return
}

// precondition of the Deref()
func (iter *BIter) Valid() bool {
	if len(iter.path) == 0 {
		return false
	}
	lastNode := iter.path[len(iter.path)-1]
	return lastNode.data != nil && iter.pos[len(iter.pos)-1] < lastNode.nKeys()
}

// moving backward and forward
func (iter *BIter) Prev() {
	iterPrev(iter, len(iter.path)-1)
}

func (iter *BIter) Next() {
	iterNext(iter, 0)
}

func (tree *BTree) Seek(key []byte, cmp int) *BIter {
	iter := tree.SeekLE(key)
	if cmp != CMP_LE && iter.Valid() {
		cur, _ := iter.Deref()
		if !cmpOK(cur, cmp, key) {
			if cmp > 0 {
				iter.Next()
			} else {
				iter.Prev()
			}
		}
	}
	return iter
}

func (tree *BTree) SeekLE(key []byte) *BIter {
	iter := &BIter{tree: tree}
	for ptr := tree.root; ptr != 0; {
		node := tree.get(ptr)
		idx := nodeLookupLE(node, key)
		iter.path = append(iter.path, node)
		iter.pos = append(iter.pos, idx)
		if node.bNodeType() == BNODE_INODE {
			ptr = node.getPtr(idx)
		} else {
			ptr = 0
		}
	}
	return iter
}

// compares current key & ref key & checks if cmp is valid
func cmpOK(key []byte, cmp int, ref []byte) bool {
	r := bytes.Compare(key, ref)
	switch cmp {
	case CMP_GE:
		return r >= 0
	case CMP_GT:
		return r > 0
	case CMP_LT:
		return r < 0
	case CMP_LE:
		return r <= 0
	default:
		panic("wrong comparison")
	}
}

func iterPrev(iter *BIter, level int) {
	if iter.pos[level] > 0 {
		iter.pos[level]-- // move within this node
	} else if level > 0 { // make sure the level is not less than the `root`
		iterPrev(iter, level-1)
	} else {
		return
	}
	if level+1 < len(iter.pos) {
		// update the kid prevNode
		prevNode := iter.path[level]
		kid := iter.tree.get(prevNode.getPtr(iter.pos[level]))
		iter.path[level+1] = kid
		iter.pos[level+1] = kid.nKeys() - 1
	}
}

func iterNext(iter *BIter, level int) {
	currentNode := iter.path[level]
	if iter.pos[level] < uint16(currentNode.nKeys())-1 {
		iter.pos[level]++ // move within this node
	} else if level < len(iter.path)-1 {
		iterNext(iter, level+1)
	} else {
		return
	}
	if level+1 < len(iter.pos) {
		// update the kid nextNode
		nextNode := iter.path[level]
		kid := iter.tree.get(nextNode.getPtr(iter.pos[level]))
		iter.path[level+1] = kid
		iter.pos[level+1] = 0
	}
}

// JSONQuery represents a JSON-style query
type JSONQuery map[string]interface{}

// ParseJSONQuery parses JSON query syntax like {"id": 1, "age": {">": 25}}
func ParseJSONQuery(queryStr string) (JSONQuery, error) {
	var query JSONQuery
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return nil, fmt.Errorf("invalid JSON query: %w", err)
	}
	return query, nil
}

// ExecuteJSONQuery executes a JSON-style query on a table
func (db *DB) ExecuteJSONQuery(table string, queryStr string, kvReader *KVReader) ([]*Record, error) {
	query, err := ParseJSONQuery(queryStr)
	if err != nil {
		return nil, err
	}

	tdef := GetTableDef(db, table, &kvReader.Tree)
	if tdef == nil {
		return nil, fmt.Errorf("table not found: %s", table)
	}

	// For now, implement simple equality checks
	// This can be extended to support operators like >, <, etc.
	var results []*Record

	// Create a scanner for full table scan (can be optimized later)
	scanner := Scanner{
		Cmp1: CMP_GE,
		Cmp2: CMP_LE,
	}

	if err := dbScan(db, tdef, &scanner, &kvReader.Tree); err != nil {
		return nil, err
	}

	// Filter results based on JSON query
	for scanner.Valid() {
		rec := &Record{
			Cols: make([]string, len(tdef.Cols)),
			Vals: make([]Value, len(tdef.Cols)),
		}
		copy(rec.Cols, tdef.Cols)
		scanner.Deref(rec, &kvReader.Tree)

		if matchesJSONQuery(rec, query, tdef) {
			results = append(results, rec)
		}

		scanner.Next()
		if len(results) > 1000 { // Safety limit
			break
		}
	}

	return results, nil
}

// matchesJSONQuery checks if a record matches the JSON query conditions
func matchesJSONQuery(rec *Record, query JSONQuery, tdef *TableDef) bool {
	for field, condition := range query {
		// Find the column index
		colIndex := -1
		for i, col := range tdef.Cols {
			if col == field {
				colIndex = i
				break
			}
		}

		if colIndex == -1 {
			continue // Skip unknown columns
		}

		if !matchesCondition(rec.Vals[colIndex], condition) {
			return false
		}
	}
	return true
}

// matchesCondition checks if a value matches a specific condition
func matchesCondition(val Value, condition interface{}) bool {
	switch cond := condition.(type) {
	case string:
		// Simple string equality
		return string(val.Str) == cond
	case float64:
		// Numeric equality
		if val.Type == TYPE_INT64 {
			return float64(val.I64) == cond
		}
	case map[string]interface{}:
		// Complex conditions like {">": 25, "<": 40}
		for op, value := range cond {
			if !evaluateOperator(val, op, value) {
				return false
			}
		}
		return true
	}

	// Default: simple equality
	return reflect.DeepEqual(val, condition)
}

// evaluateOperator evaluates comparison operators
func evaluateOperator(val Value, operator string, compareValue interface{}) bool {
	if val.Type == TYPE_INT64 {
		if compareVal, ok := compareValue.(float64); ok {
			switch operator {
			case ">":
				return float64(val.I64) > compareVal
			case "<":
				return float64(val.I64) < compareVal
			case ">=":
				return float64(val.I64) >= compareVal
			case "<=":
				return float64(val.I64) <= compareVal
			case "=", "==":
				return float64(val.I64) == compareVal
			}
		}
	}
	return false
}
