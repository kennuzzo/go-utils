package sqlite3Util

import (
	"database/sql"
	"fmt"
	//for sqlite3 support
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

//Results represents the parsed query results.
type Results []map[string]interface{}

//Sqlite3Wrapper wraps a sql.Db with methods of convenience.

type rowStack struct {
	rowStack []*sql.Rows
}

//Sqlite3 is a wrapper to go-sqlite-3 lib.
type Sqlite3 struct {
	rowStack
	dbHandle *sql.DB // db handler
}

//Close the database connection
func (dbw *Sqlite3) Close() {
	dbw.dbHandle.Close()
	dbw.rowStack.Close()
}

//Connect to a database.
func (dbw *Sqlite3) Connect() (err error) {
	dbw.dbHandle, err = sql.Open("sqlite3", "./development.sqlite3")
	return
}

//Execute a query
func (dbw *Sqlite3) Execute(query string) (rows *sql.Rows, err error) {
	rows, err = dbw.dbHandle.Query(query)
	dbw.rowStack.rowStack = append(dbw.rowStack.rowStack, rows)
	// defer rows.Close()
	return
}

//Retrive execute and parse a query.
func (dbw *Sqlite3) Retrive(query string) (rows Results, err error) {
	var tempRows *sql.Rows
	tempRows, err = dbw.Execute(query)
	// defer tempRows.Close()
	rows, err = dbw.Parse(tempRows)
	checkErrors(err)
	return
}

//Parse an executed query
func (dbw *Sqlite3) Parse(rows *sql.Rows) (results Results, err error) {
	// here im not sure what im getting from the db as it can be anything from the results
	//therefore i wish i can just put everything in a bucket, and think later.
	tc, tn := dbw.TotalCount(rows)
	// multiFields := make([]interface{}, tc)
	multiFieldsPtrs := make([]interface{}, tc)
	for i := 0; i < tc; i++ {
		var multiFields interface{}
		multiFieldsPtrs[i] = &multiFields
	}

	headers := make(map[string]interface{})
	headers["fields"] = tn
	results = append(results, headers)
	for rows.Next() {
		err = rows.Scan(multiFieldsPtrs...)
		checkErrors(err)
		tempMap := make(map[string]interface{})
		for idx, label := range tn {
			row := *(multiFieldsPtrs[idx].(*interface{}))
			// if sizeOf(row) < 1 {
			// 	skip = true
			// 	break
			// }
			tempMap[label] = row
		}
		results = append(results, tempMap)
	}
	return
}

func sizeOf(data interface{}) int {

	switch data.(type) {
	case string:
		return len(data.(string))
	case int32, int64:
		return len(strconv.Itoa(data.(int)))
	case []uint8:
		return len(toString(data))
	case time.Time:
		return 1
	}
	return 0
}

//TotalCount returns the db fields count and their fields name.
func (dbw *Sqlite3) TotalCount(rows *sql.Rows) (int, []string) {
	tmpArrayString, _ := rows.Columns()
	return len(tmpArrayString), tmpArrayString
}

//Fields returns the fields found in the results.
func (r Results) Fields() interface{} {
	return r[0]["fields"]
}

//Close the rows handles
func (rs rowStack) Close() {
	for _, s := range rs.rowStack {
		s.Close()
	}
}

func toString(a interface{}) string {
	if a == nil {
		return ""
	}
	return string(a.([]byte))
}

func checkErrors(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
