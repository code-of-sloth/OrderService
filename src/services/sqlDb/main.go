package sqlDb

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func initializeDB() (db *sql.DB, err error) {
	// Implement the code for initializing the database connection
	// add sql credendial
	connStr := "user=postgres password=password dbname=mydb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func sqlHandler(sqlString string, result chan<- []map[string]interface{}, db *sql.DB) {
	if strings.ToLower(sqlString[:6]) == "select" {
		stmt, err := db.Prepare(sqlString)
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		tableData := make([]map[string]interface{}, 0)
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			entry := make(map[string]interface{})
			if err != nil {
				result <- []map[string]interface{}{
					{"error": err.Error()},
				}
				return
			}

			for i, col := range values {
				if col == nil {
					entry[columns[i]] = nil
				} else if numb, ok := strconv.ParseFloat(string(col), 64); ok == nil {
					entry[columns[i]] = numb
				} else if booln, ok := strconv.ParseBool(string(col)); ok == nil {
					entry[columns[i]] = booln
				} else {
					entry[columns[i]] = string(col)
				}
			}
			tableData = append(tableData, entry)
		}
		result <- tableData
	} else {
		entry := make(map[string]interface{})
		stmtIns, err := db.Prepare(sqlString)
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		defer stmtIns.Close()
		res, err := stmtIns.Exec()
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		entry["lastInsertId"] = lastId
		rowCnt, err := res.RowsAffected()
		if err != nil {
			result <- []map[string]interface{}{
				{"error": err.Error()},
			}
			return
		}

		entry["rowsAffected"] = rowCnt

		result <- []map[string]interface{}{entry}
	}
}

func RunQuery(querys ...string) (rows [][]map[string]interface{}, err error) {
	db, err := initializeDB()
	defer db.Close()
	channelArr := make([]chan []map[string]interface{}, len(querys))
	for _, query := range querys {
		sqlResult := make(chan []map[string]interface{})
		go sqlHandler(strings.TrimSpace(query), sqlResult, db)
		channelArr = append(channelArr, sqlResult)
	}

	for _, v := range channelArr {
		data := <-v
		rows = append(rows, data)
	}

	return
}
