package dbutil

import (
	"fmt"
	"path"
	"testing"
)

func mustExecAndQuerySuccessfully(t *testing.T, db *DBUtil) {
	err := db.RunExec("CREATE TABLE A (ID TEXT)", nil, false)
	if err != nil {
		t.Errorf("Expected table creation, errored out")
	}
	err = db.RunExec("INSERT INTO A (ID) VALUES (?)", []any{"TestID"}, false)
	if err != nil {
		t.Errorf("Expected row creation, errored out")
	}
	rows, err := db.RunQuery("SELECT * FROM A", nil, false)
	if err != nil {
		t.Errorf("Expected to get rows, errored out")
	}
	if rows == nil {
		t.Errorf("Expected rows value not nil, found %v", rows)
	} else {
		defer rows.Close()
	}
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			t.Errorf("Expected scan of rows, errored out")
		}
		if id != "TestID" {
			t.Errorf("Expected id: TestID, found %s", id)
		}
	}
}

func mustDryRunExecSuccessfully(t *testing.T, db *DBUtil) {
	err := db.RunExec("INSERT INTO A (ID) VALUES (?)", []any{"TestID2"}, true)
	if err != nil {
		t.Errorf("Expected row creation, errored out")
	}
	rows, err := db.RunQuery("SELECT * FROM A", nil, false)
	if err != nil {
		t.Errorf("Expected to get rows, errored out")
	}
	if rows == nil {
		t.Errorf("Expected rows value not nil, found %v", rows)
	} else {
		defer rows.Close()
	}
	count := 0
	for rows.Next() {
		count += 1
		var id string
		if err = rows.Scan(&id); err != nil {
			t.Errorf("Expected scan of rows, errored out")
		}
		if id != "TestID" {
			t.Errorf("Expected id: TestID, found %s", id)
		}
	}
	if count > 1 {
		t.Errorf("Expected count value: 1, found %d", count)
	}
}

func mustDryRunQuerySuccessfully(t *testing.T, db *DBUtil) {
	rows, err := db.RunQuery("SELECT * FROM A", nil, true)
	if err != nil {
		t.Errorf("Expected to get rows, errored out")
	}
	if rows != nil {
		t.Errorf("Expected rows value nil, found %v", rows)
		defer rows.Close()
	}
}

func TestDBUtil(t *testing.T) {
	tempDir := t.TempDir()
	db, err := NewDBUtil(fmt.Sprintf("file:%s", path.Join(tempDir, "test.db")))
	if err != nil {
		t.Errorf("Expected database creation, errored out")
	}
	defer db.CloseDatabase()

	// test successful exec and query
	t.Run("SuccessfulExecAndQuery", func(t *testing.T) {
		mustExecAndQuerySuccessfully(t, db)
	})

	// test dry run on Exec
	t.Run("SuccessfulDryRunExec", func(t *testing.T) {
		mustDryRunExecSuccessfully(t, db)
	})

	// test dry run on query
	t.Run("SuccessfulDryRunQuery", func(t *testing.T) {
		mustDryRunQuerySuccessfully(t, db)
	})
}
