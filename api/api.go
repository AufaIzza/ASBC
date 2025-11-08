package api

import (
	"io"
	"fmt"
	"net/http"
	"database/sql"
	_ "modernc.org/sqlite"
)

type DebugMode int

var Db *sql.DB
var debug_enabled bool
var err error

const (
	DebugModeDisabled DebugMode = iota
	DebugModeEnabled
)

func TestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, "I AM Testing")
	}
}

func ExecDB(command string) error {
	_, err = Db.Exec(command)
	if err != nil {
		return err
	}
	if debug_enabled {
		fmt.Println("DATABASE:", command)
	}
	return nil
}

func QueryDB(command string) (*sql.Rows, error) {
	rows, err := Db.Query(command)
	if err != nil {
		return nil, err
	}
	if debug_enabled {
		fmt.Println("DATABASE:", command)
	}
	return rows, nil
}

func InitDB(path string, debug_mode DebugMode) error {
	if debug_mode == DebugModeEnabled {
		debug_enabled = true
	} else if debug_mode == DebugModeDisabled {
		debug_enabled = false
	}

	Db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	fmt.Println("DATABASE: Opening", path)

	err = ExecDB(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	Db.Close()
	fmt.Println("DATABASE: Closing")
}

func DropTables() err {
	err = ExecDB(
		``)
	if err != nil {
		return err
	}
	return nil
}

func CreateTables() err {
	err = ExecDB(
		``)
	if err != nil {
		return err
	}
	return nil
}

func InitTestDummyDataDB() error {
	err = ExecDB(
		`DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
CREATE TABLE users(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    password TEXT NOT NULL
);
CREATE TABLE posts(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    userid INTEGER NOT NULL,
    content TEXT,
    FOREIGN KEY (userid) REFERENCES users(id)
);
INSERT INTO users (name, password) VALUES ('Bran', 'Lain');
INSERT INTO users (name, password) VALUES ('Lain', 'Bran');
INSERT INTO posts (userid, content) VALUES (1, 'HELLO WORLD I AM HERE');
INSERT INTO posts (userid, content) VALUES (2, 'More test data')`)
	if err != nil {
		return err
	}
	return nil
}
