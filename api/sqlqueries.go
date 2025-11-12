package api

type NoteQuery struct {
	ID int
	Title string
	Content string
	UserID int
	Username string
}

func queryAllPublicNotes() ([]NoteQuery, error) {
	var notes []NoteQuery
	query := `
SELECT n.ID, n.Title, n.Content, n.UserID, a.Username
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
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
