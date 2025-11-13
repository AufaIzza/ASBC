package api

import (
	"io"
	"net/http"
	"strings"
	"fmt"
	"strconv"
)

type LoginBody struct {
	Username string 
	Password string
}

type AssignmentBody struct {
	Title string
	Description string
}

type NoteBody struct {
	Title string
	Content string
	Public int
	Tag string
}

func assignmentTmpl(ID int, title string, description string, check bool) string {
	var checkBtn string
	// TODO: put hx-put properly
	if check {
		checkBtn = fmt.Sprintf(`<input type="checkbox" hx-swap="outerHTML" checked hx-trigger="change" hx-patch="/api/check_assignment/%d/0">`, ID)
	} else {
		checkBtn = fmt.Sprintf(`<input type="checkbox" hx-swap="outerHTML" hx-trigger="change" hx-patch="/api/check_assignment/%d/1">`, ID)
	}
	return fmt.Sprintf(`<div class="note">
		<h3>%s</h3>
		<p>%s</p>
		<div class="note-actions">
		    %s
		    <button hx-trigger="click" hx-put="/api/delete_assignment/%d" class="delete-btn">🗑 Delete</button>
		</div>
	</div>`, title, description, checkBtn, ID)
}


func DeleteAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	_, _, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}
	ExecDeleteAssignment(id)
	w.Header().Add("HX-Refresh", "true")
}

func NewNoteHandler(w http.ResponseWriter, r *http.Request) {
	_, id, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}
	res, err := strconv.Atoi(r.FormValue("public"))
	if err != nil {
		panic(err)
	}
	b := NoteBody{
		Title: r.FormValue("title"),
		Content: r.FormValue("content"),
		Public: res,
		Tag: r.FormValue("tag"),
	}

	err = ExecNewNote(id, b.Title, b.Content, b.Public, b.Tag)
	if err != nil {
		panic(err)
	}

	ok, noteid, err := QueryNoteGetID(id, b.Title, b.Content, b.Public, b.Tag)
	if err != nil || !ok {
		panic(err)
	}
	
	w.Header().Add("HX-Redirect", fmt.Sprintf("/view_note/%d", noteid))
	w.WriteHeader(http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/view_note/%d", noteid), http.StatusOK)
}

func CheckAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	assignmentID, err := strconv.Atoi(r.PathValue("assignmentID"))
	if err != nil {
		panic(err)
	}
	isDone, err := strconv.Atoi(r.PathValue("isDone"))
	if err != nil {
		panic(err)
	}
	ExecUpdateAssignment(assignmentID, isDone)

	if isDone == 1 {
		io.WriteString(w, fmt.Sprintf(`<input type="checkbox" hx-swap="outerHTML" checked hx-trigger="change" hx-patch="/api/check_assignment/%d/0">`, assignmentID))
	} else {
		io.WriteString(w, fmt.Sprintf(`<input type="checkbox" hx-swap="outerHTML" hx-trigger="change" hx-patch="/api/check_assignment/%d/1">`, assignmentID))
	}
}

func AllAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	_, id, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}
	assignments, err := QueryAllUserAssignment(id)
	if err != nil {
		panic(err)
	}
	
	for _, a := range assignments {
		sb.WriteString(assignmentTmpl(a.ID, a.Title, a.Description, a.IsDone))
	}
	sb.WriteString(`
            <div class="add-note">
                <a href="/new_assignment">➕ Add New Assignment</a>
            </div>
`)
	io.WriteString(w, sb.String())
}

func createNoteView(id int, title string, isPublic bool, username string, tagname string, content string, userID int, r *http.Request) string{
	_, userSessionID, ok := GetSessionUser(r)
	var visibility string
	if isPublic {
		visibility = "Public"
	} else {
		visibility = "Private"
	}
	var deleteBTN string
	if ok && userID == userSessionID {
		deleteBTN = fmt.Sprintf(`<button hx-trigger="click" hx-put="/api/delete_note/%d" class="delete-btn">🗑 Delete</button>`, id)
	} else {
		deleteBTN = ""
	}
	return fmt.Sprintf(`
<a href="/view_note/%d">
<div class="view-note-container">
	<div class="note-header">
	<h1 class="note-title">
	%s
<span class="visibility-badge">%s</span>
	</h1>
	<div class="note-meta">
	<span class="tag">%s</span>
	<span class="tag">%s</span>
	</div>
	</div>

	<div class="note-content">
	<span>%s</span>
	</div>
    %s
	</div>
`, id, title, visibility, username, tagname, content, deleteBTN)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	_, userID, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}	
	deleteID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}

	err = ExecDeleteNote(deleteID, userID)
	if err != nil {
		panic(err)
	}
	w.Header().Add("HX-Refresh", "true")
}

func AllPrivateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	_, id, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}
	notes, err := QueryAllPrivateNotes(id)
	if err != nil {
		panic(err)
	}
	var public bool

	for _, note := range notes {
		if note.IsPublic == 1 {
			public = true
		} else {
			public = false
		}
		sb.WriteString(createNoteView(note.ID, note.Title, public, note.Username, note.TagName, note.Content, note.UserID, r))
	}

	io.WriteString(w, sb.String())
}

func AllPublicNotesHandler(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder

	notes, err := QueryAllPublicNotes()
	if err != nil {
		panic(err)
	}

	for _, note := range notes {
		sb.WriteString(createNoteView(note.ID, note.Title, true, note.Username, note.TagName, note.Content, note.UserID, r))
	}

	io.WriteString(w, sb.String())
}

func AllPublicNotesHandlerTest(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	
	notes, err := QueryAllPublicNotes()
	if err != nil {
		panic(err)
	}

	for _, note := range notes {
		sb.WriteString(fmt.Sprintf("%v", note))
	}

	io.WriteString(w, sb.String())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}	

	body := LoginBody{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	ok, _, id := QueryLogin(body.Username, body.Password)
	if ok {
		SetSessionUser(w, r, body.Username, id)
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
		http.Redirect(w, r, "/", http.StatusOK)

	} else {
		io.WriteString(w, "Username/Password Incorrect")
	}
}

func NewAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	_, id, ok := GetSessionUser(r)
	if !ok {
		io.WriteString(w, "No User found")
		return
	}

	b := AssignmentBody{
		Title: r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	err = ExecNewAssignment(id, b.Title, b.Description)
	if err != nil {
		panic(err)
	}
	
	w.Header().Add("HX-Redirect", "/assignments")
	w.WriteHeader(http.StatusSeeOther)
	http.Redirect(w, r, "/assignments", http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}	

	body := LoginBody{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	ok, err := QueryRegister(body.Username, body.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
		http.Redirect(w, r, "/", http.StatusOK)

	} else {
		io.WriteString(w, "Cannot register, Username might be unavailable")
	}
}

func NavBarUserHandler(w http.ResponseWriter, r *http.Request) {
	logged_in := `
<a href="/private_notes">Personal Notes</a>
<a href="/assignments">Assignments</a>
<a href="/logout">Logout</a>
`
	logged_out := `
<a href="/login">Log in</a>
<a href="/register">Register</a>
`

	_, _, ok := GetSessionUser(r)
	
	if ok {
		io.WriteString(w, logged_in)
	} else {
		io.WriteString(w, logged_out)		
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := Store.Get(r, "session-name")
    // Clear the session
    delete(session.Values, "username")
    delete(session.Values, "userID")
    session.Save(r, w)

    // Redirect to the homepage or login page
	// w.Header().Add("HX-Redirect", "/")
	// w.WriteHeader(http.StatusSeeOther)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func LoginFailedHandler(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Login Failed")	
}

func TestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, "I AM Testing")
	}
}
