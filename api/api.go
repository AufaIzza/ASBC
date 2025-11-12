package api

import (
	// "io"
	"fmt"
	// "net/http"
	"database/sql"
	_ "modernc.org/sqlite"
)

type DebugMode bool

var Db *sql.DB
var debug_enabled bool
var err error

const (
	DebugModeDisabled DebugMode = false
	DebugModeEnabled DebugMode = true
)

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

func DropTables() error {
	err = ExecDB(`
DROP TABLE IF EXISTS UserAssignments;
DROP TABLE IF EXISTS Assignments;
DROP TABLE IF EXISTS UsersNote;
DROP TABLE IF EXISTS NotesTags;
DROP TABLE IF EXISTS Tags;
DROP TABLE IF EXISTS Notes;
DROP TABLE IF EXISTS Users;
`)
	if err != nil {
		return err
	}
	return nil
}

func CreateTables() error {
	err = ExecDB(`
CREATE TABLE Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Username TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL
);
CREATE TABLE Notes (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    UserID INTEGER NOT NULL,
    Title TEXT NOT NULL,
    Content TEXT NOT NULL,
    IsPublic INTEGER NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
CREATE TABLE Tags (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Name TEXT NOT NULL
);
CREATE TABLE NotesTags (
    NoteID INTEGER NOT NULL,
    TagID INTEGER NOT NULL,
    FOREIGN KEY (NoteID) REFERENCES Notes(ID),
    FOREIGN KEY (TagID) REFERENCES Tags(ID)
);
CREATE TABLE UsersNote (
    NoteID INTEGER PRIMARY KEY NOT NULL,
    UserID INTEGER NOT NULL,
    FOREIGN KEY (NoteID) REFERENCES Notes(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
CREATE TABLE Assignments (
    ID INTEGER PRIMARY KEY NOT NULL,
    Title TEXT NOT NULL,
    Description TEXT NOT NULL,
    IsDone INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
CREATE TABLE UserAssignments (
    AssignmentID INTEGER PRIMARY KEY NOT NULL,
    UserID INTEGER NOT NULL,
    FOREIGN KEY (AssignmentID) REFERENCES Assignments(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
`)
	if err != nil {
		return err
	}

	return nil
}

func InsertDummyData() error {
	// Password is "password123"
	err = ExecDB(`
INSERT INTO Users (Username, Password) VALUES ('John', '$2a$12$Wnk5aXokULm/gi1zQgml0uk/zNOKAUUIwZasdDZm41VHt7PJi/7a6');
INSERT INTO Users (Username, Password) VALUES ('Smith', '$2a$12$cXSkUSaaKSZKW7YKOpOZL.WJFeDHzdfut2JpOhnQFL9iDdwb14YNG');
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO Notes (UserID, Title, Content, IsPublic) VALUES (1, 'PubNote', 'This is a public note 1, lorem ipsum', 1);
INSERT INTO Notes (UserID, Title, Content, IsPublic) VALUES (2, 'PubNote', 'This is a public note 2, lorem ipsum', 1);
INSERT INTO Notes (UserID, Title, Content, IsPublic) VALUES (1, 'PrivNote', 'This is a private note, lorem ipsum', 0);
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO Tags (Name) VALUES ('Programming');
INSERT INTO Tags (Name) VALUES ('Math');
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO NotesTags (NoteID, TagID) VALUES (1, 1);
INSERT INTO NotesTags (NoteID, TagID) VALUES (2, 2);
INSERT INTO NotesTags (NoteID, TagID) VALUES (3, 1);
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO UsersNote (NoteID, UserID) VALUES(1, 1);
INSERT INTO UsersNote (NoteID, UserID) VALUES(2, 2);
INSERT INTO UsersNote (NoteID, UserID) VALUES(3, 1);
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO Assignments (UserID, Title, Description, IsDone) VALUES (1, 'Task 1', 'This is a task 1', 0);
INSERT INTO Assignments (UserID, Title, Description, IsDone) VALUES (1, 'Task 2', 'This is a task 2', 0);
INSERT INTO Assignments (UserID, Title, Description, IsDone) VALUES (2, 'Task 3', 'This is a task 3', 0);
`)
	if err != nil {
		return err
	}

	err = ExecDB(`
INSERT INTO UserAssignments (AssignmentID, UserID) VALUES (1, 1);
INSERT INTO UserAssignments (AssignmentID, UserID) VALUES (2, 1);
INSERT INTO UserAssignments (AssignmentID, UserID) VALUES (3, 2);
`)
	if err != nil {
		return err
	}

	return nil
}

func CleanDB() error {
	// Drop Tables
	err = DropTables()
	if err != nil {
		return err
	}

	// Create Tables
	err = CreateTables()
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
