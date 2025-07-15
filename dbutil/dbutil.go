package dbutil

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DBUtil struct {
	connectionString string
	dbSession        *sql.DB
}

func NewDBUtil(connectionString string) (*DBUtil, error) {
	dbUtil := &DBUtil{}
	dbUtil.connectionString = connectionString
	dbSession, err := sql.Open("sqlite3", dbUtil.connectionString)
	if err != nil {
		return nil, err
	}
	dbUtil.dbSession = dbSession
	return dbUtil, nil
}

func (dbUtil *DBUtil) CloseDatabase() error {
	err := dbUtil.dbSession.Close()
	if err != nil {
		return err
	}
	return nil
}

func (dbUtil *DBUtil) CreateDatabase() error {
	DB, err := sql.Open("sqlite3", dbUtil.connectionString)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (dbUtil *DBUtil) RunExec(queryString string, args []any, dryRun bool) error {
	if dryRun {
		fmt.Printf("\n%s\n", queryString)
		return nil
	} else {
		txn, err := dbUtil.dbSession.Begin()
		if err != nil {
			return err
		}
		_, err = txn.Exec(queryString, args...)
		if err != nil {
			txn.Rollback()
			return err
		}
		err = txn.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (dbUtil *DBUtil) RunQuery(queryString string, args []any, isDryRun bool) (*sql.Rows, error) {
	if isDryRun {
		fmt.Printf("\n%s\n", queryString)
		return nil, nil
	}
	rows, err := dbUtil.dbSession.Query(queryString, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
