package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func Connect(fp string) {
	var e error
	db, e = sql.Open("sqlite3", "file:"+fp+"?cache=shared")
	if e != nil {
		logrus.Fatal(e)
	}
}

func Exec(db *sql.DB, query string, parameters ...interface{}) int {
	result, e := db.Exec(query, parameters...)
	if e != nil {
		logrus.WithFields(
			logrus.Fields{
				"file":
				"line":
				"function": "
			}).Error("dbExec error: %v, query: %v", e, query)
		return -1
	}

	affected, e := result.RowsAffected()
	if e != nil {
		logrus.Error("dbExec error: %v, query: %v", e, query)
		return -1
	}

	return int(affected)
}

func Query(db *sql.DB, query string, parameters ...interface{}) [][]string {
	result := [][]string{}

	rows, e := db.Query(query, parameters...)
	if e != nil {
		logrus.Error("dbQuery error: %v, query: %v", e, query)
		return nil
	}
	defer rows.Close()

	columns, e := rows.Columns()
	if e != nil {
		logrus.Error("dbQuery error: %v, query: %v", e, query)
		return nil
	}

	tempBytePtr := make([]interface{}, len(columns))
	tempByte := make([][]byte, len(columns))
	tempString := make([]string, len(columns))
	for i := range tempByte {
		tempBytePtr[i] = &tempByte[i]
	}

	for rows.Next() {
		if e := rows.Scan(tempBytePtr...); e != nil {
			logrus.Error("dbQuery error: %v, query: %v", e, query)
			return nil
		}

		for i, rawByte := range tempByte {
			if rawByte == nil {
				tempString[i] = "\\N"
			} else {
				tempString[i] = string(rawByte)
			}
		}

		result = append(result, make([]string, len(columns)))
		copy(result[len(result)-1], tempString)
	}

	return result
}

func Tx(db *sql.DB, procedure func(*sql.Tx) bool) bool {
	tx, e := db.Begin()
	if e != nil {
		logrus.Error("dbTx Begin error: %v", e)
		return false
	}

	if procedure(tx) {
		e := tx.Commit()
		if e != nil {
			logrus.Error("dbTx Commit error: %v", e)
			return false
		}
		return true
	} else {
		e := tx.Rollback()
		if e != nil {
			logrus.Error("dbTx Rollback error: %v", e)
			return false
		}
		return false
	}
}

func TxExec(tx *sql.Tx, query string, parameters ...interface{}) int {
	result, e := tx.Exec(query, parameters...)
	if e != nil {
		logrus.Error("dbTxExec error: %v, query: %v", e, query)
		return -1
	}

	affected, e := result.RowsAffected()
	if e != nil {
		logrus.Error("dbTxExec error: %v, query: %v", e, query)
		return -1
	}

	return int(affected)
}

func TxQuery(tx *sql.Tx, query string, parameters ...interface{}) [][]string {
	result := [][]string{}

	rows, e := tx.Query(query, parameters...)
	if e != nil {
		logrus.Error("dbTxQuery error: %v, query: %v", e, query)
		return nil
	}
	defer rows.Close()

	columns, e := rows.Columns()
	if e != nil {
		logrus.Error("dbTxQuery error: %v, query: %v", e, query)
		return nil
	}

	tempBytePtr := make([]interface{}, len(columns))
	tempByte := make([][]byte, len(columns))
	tempString := make([]string, len(columns))
	for i := range tempByte {
		tempBytePtr[i] = &tempByte[i]
	}

	for rows.Next() {
		if e := rows.Scan(tempBytePtr...); e != nil {
			logrus.Error("dbTxQuery error: %v, query: %v", e, query)
			return nil
		}

		for i, rawByte := range tempByte {
			if rawByte == nil {
				tempString[i] = "\\N"
			} else {
				tempString[i] = string(rawByte)
			}
		}

		result = append(result, make([]string, len(columns)))
		copy(result[len(result)-1], tempString)
	}

	return result
}
