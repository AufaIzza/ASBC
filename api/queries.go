package api

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserQuery struct {
	ID int
	Username string
	Password string
}

type NoteQuery struct {
	ID int
	Title string
	Content string
	TagName string
	IsPublic int
	UserID int
	Username string
}

type AssignmentQuery struct {
	ID int
	Title string
	Description string
	IsDone bool
}

func ExecDeleteAssignment(id int) error {
	query := fmt.Sprintf(`
DELETE FROM Assignments
WHERE ID = %d
`, id)
	return ExecDB(query)
}

func ExecUpdateAssignment(id int, isDone int) error {
	query := fmt.Sprintf(`
UPDATE Assignments
SET IsDone = %d
WHERE ID = %d;
`, isDone, id)
	return ExecDB(query)
}

func ExecNewAssignment(userID int, title string, description string) error {
	query := fmt.Sprintf(`
INSERT INTO Assignments (UserID, Title, Description, IsDone) VALUES (%d, '%s', '%s', 0);
`, userID, title, description)

	return ExecDB(query)
}

func ExecNewNote(userid int, title string, content string, public int, tag string) error {
	
	query := fmt.Sprintf(`
INSERT INTO Notes (UserID, Title, Content, IsPublic, TagName) VALUES (%d, '%s', '%s', %d, '%s');
`, userid, title, content, public, tag)

	return ExecDB(query)
}

func QueryNoteGetID(userid int, title string, content string, public int, tag string) (bool, int, error) {
	
	query := fmt.Sprintf(`
SELECT a.ID
FROM Notes as a
WHERE a.UserID = %d AND a.Title = '%s' AND a.Content = '%s' AND a.IsPublic = %d AND a.TagName = '%s';
`, userid, title, content, public, tag)
	
	rows, err := QueryDB(query)
	if err != nil {
		// fmt.Println(err)
		return false, -1, err
	}	
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return false, -1, err
		}
		break;
	}
	return true, id, nil
}

func QueryRegister(username string, password string) (bool, error) {
	query := fmt.Sprintf(`
SELECT u.ID, u.Username, u.Password
FROM Users as u
WHERE u.Username = '%s';
`, username)
	rows, err := QueryDB(query)
	if err != nil {
		return false, err
	}	
	if rows.Next() {
		return false, nil
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return false, err
	}
	
	queryCreate := fmt.Sprintf(`
INSERT INTO Users (Username, Password) VALUES ('%s', '%s');
`, username, string(pass))

	err = ExecDB(queryCreate)
	if err != nil {
		return false, err
	}

	return true, nil
}

func QueryAllTags() ([]string, error) {
	rows, err := QueryDB(`
SELECT t.Name
FROM Tags as t;
`) 
	if err != nil {
		return []string{}, err
	}

	var tags []string

	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return []string{}, err
		}
		tags = append(tags, str)
	}
	return tags, nil
}

func QueryLogin(username string, password string) (bool, error, int) {
	query := fmt.Sprintf(`
SELECT u.ID, u.Username, u.Password
FROM Users as u
WHERE u.Username = '%s';
`, username)
	rows, err := QueryDB(query)
	if err != nil {
		// fmt.Println(err)
		return false, err, -1
	}	
	var user UserQuery
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
		)
		if err != nil {
			return false, err, -1
		}
		break;
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err, -1
	}
	return true, nil, user.ID
}

func QueryNoteID(noteID string) (NoteQuery, error) {
	query := fmt.Sprintf(`
SELECT a.ID, a.Title, a.Content, b.Username, a.TagName, a.IsPublic, a.UserID
FROM Notes AS a
JOIN Users AS b ON a.UserID = b.ID
WHERE a.ID = %s;
`, noteID)	
	rows, err := QueryDB(query)
	if err != nil {
		return NoteQuery{}, err
	}

	for rows.Next() {
		a := NoteQuery{}
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.Username,
			&a.TagName,
			&a.IsPublic,
			&a.UserID,
		)
		if err != nil {
			return NoteQuery{}, err
		}
		
		return a, nil
	}

	return NoteQuery{}, nil
}

func QueryAllUserAssignment(userID int) ([]AssignmentQuery, error) {
	var assignments []AssignmentQuery
	query := fmt.Sprintf(`
SELECT a.ID, a.Title, a.Description, a.IsDone
From Assignments as a
WHERE a.UserID = %d;
`, userID)
	rows, err := QueryDB(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a := AssignmentQuery{}
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.IsDone,
		)
		if err != nil {
			return nil, err
		}
		
		assignments = append(assignments, a)
	}

	return assignments, nil
}

func QueryAllPrivateNotes(id int) ([]NoteQuery, error) {
	var notes []NoteQuery
	query := fmt.Sprintf(`
SELECT n.ID, n.Title, n.Content, n.UserID, a.Username, n.IsPublic, n.TagName
FROM Notes as n
JOIN Users as a ON n.UserID = a.ID
WHERE n.UserID = %d;
`, id)
	rows, err := QueryDB(query)
	if err != nil {
		return nil, err
	}


	for rows.Next() {
		note := NoteQuery{}
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.UserID,
			&note.Username,
			&note.IsPublic,
			&note.TagName,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func ExecDeleteNote(noteID int, userID int) error {
	query := fmt.Sprintf(`
DELETE FROM Notes
WHERE ID = %d AND UserID = %d
`, noteID, userID)
	return ExecDB(query)
}


func QueryAllPublicNotes() ([]NoteQuery, error) {
	var notes []NoteQuery
	query := `
SELECT n.ID, n.Title, n.Content, n.UserID, a.Username, n.TagName
FROM Notes as n
JOIN Users as a ON n.UserID = a.ID
WHERE n.IsPublic = 1;
`
	rows, err := QueryDB(query)
	if err != nil {
		return nil, err
	}


	for rows.Next() {
		note := NoteQuery{}
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.UserID,
			&note.Username,
			&note.TagName,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
