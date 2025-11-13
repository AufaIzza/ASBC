package main

import (
	"io"
	"os"
	// "bytes"
	"fmt"
	// "log"
	"net/http"
	// "html/template"
	"database/sql"
	"github.com/AufaIzza/ASBC/api"
	"github.com/AufaIzza/ASBC/helper"
	// "github.com/mattn/go-sqlite3"
	// "modernc.org/sqlite"
	"flag"
	"os/signal"
	"github.com/joho/godotenv"
	"crypto/rand"
	"encoding/base32"
	// "strconv"
	// "math/big"
	// "os"
)

// func rootHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == http.MethodGet {
// 		io.WriteString(w, "Hello World\n")
// 	}
// }

// type Page struct {
// 	page string
// }

// func (p Page) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == http.MethodGet {
// 		io.WriteString(w, p.page)
// 	}
// }

var err error 

type Site struct {
	Mux *http.ServeMux
	// Tmpl *template.Template
	TmplPath string
	Port string
}

func MakeSite(port string, partialsPath string) (*Site, error) {
	site := Site {
		Port: port,
	}
	// site.Tmpl, err = helper.MakeTemplateFromGlobAndCollect("./site_src/partials/*.gohtml")
	// if err != nil {
	// 	return nil, err
	// }
	site.TmplPath = partialsPath
	site.Mux = http.NewServeMux()
	return &site, nil
}

func (s *Site) PageRouteDef(path string, filepath string, props map[string]string) {
	tmpl, err := helper.MakeTemplateFromGlobAndCollect(s.TmplPath)
	if err != nil {
		panic(err)
	}
	err = helper.MakePageFromFile(s.Mux, tmpl, path, filepath, props)
	if err != nil {
		panic(err)
	}
}

func (s *Site) Route(path string, handler http.Handler) {
	s.Mux.Handle(path, handler)
}

func (s *Site) RouteFunc(path string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	s.Mux.HandleFunc(path, handlerFunc)
}

func (s *Site) Serve() error {
	fmt.Println("Listening on port", s.Port)
	return http.ListenAndServe(s.Port, s.Mux)
}

func (s *Site) StaticPath(path string) {
	fs := http.FileServer(http.Dir(path))
	s.Mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
}

// func main() {
// 	page_test()
// 	// err = api.InitDB("./test.db")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer api.CloseDB()
// }



func getToken(length int) string {
    randomBytes := make([]byte, 32)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
    return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func main() {
	cleanDB := flag.Bool("clean-db", false, "Clears and initializes the database")
	seedDB := flag.Bool("seed-db", false, "Inserts dummy seed data into the database")
	debug := flag.Bool("debug", false, "Enable Debug Mode")
	startServer := flag.Bool("start-server", false, "Starts the server")
	genSessionKey := flag.Bool("gen-session-key", false, "Generate session key")	
	
	flag.Parse()

	if *genSessionKey {
		if err != nil {
			panic(err)
		}
		fmt.Println(getToken(32));
		os.Exit(0)
	} 

	if *cleanDB {
		fmt.Println("Cleaning DB")
		err = api.InitDB("./database.db", api.DebugMode(*debug))
		if err != nil {
			panic(err)
		}
		err = api.CleanDB()
		if err != nil {
			panic(err)
		}
		api.CloseDB()
		os.Exit(0)
	}

	if *seedDB {
		fmt.Println("Cleaning DB")
		err = api.InitDB("./database.db", api.DebugMode(*debug))
		if err != nil {
			panic(err)
		}
		err = api.InsertDummyData()
		if err != nil {
			panic(err)
		}
		api.CloseDB()
		os.Exit(0)
	}

	if *startServer {
		err := godotenv.Load()
		api.InitSession()
		if err != nil {
			panic(err)
		}
		server(api.DebugMode(*debug))
	}
}

func server(debug api.DebugMode ) {
	err = api.InitDB("./database.db", debug)
	if err != nil {
		panic(err)
	}
	defer api.CloseDB()

	site, err := MakeSite(":6969", "./site_src/partials/*.gohtml")

	site.StaticPath("./site_src/static/")

	site.PageRouteDef("GET /", "./site_src/pages/index.gohtml", nil)

	site.RouteFunc("GET /view_note/{noteID}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := helper.MakeTemplateFromGlobAndCollect(site.TmplPath)

		if err != nil {
			panic(err)
		}
		noteID := r.PathValue("noteID")
		note, err := api.QueryNoteID(noteID)
		if err != nil {
			panic(err)
		}
		if !note.IsPublic {
			_, id, ok := api.GetSessionUser(r)
			if !ok || id != note.UserID {
				io.WriteString(w, "Not Allowed")
				return
			}
		}
		page, err := helper.MakePageStr(tmpl, "./site_src/pages/view_note.gohtml", note)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, page)
	})

	site.RouteFunc("GET /new_note", func(w http.ResponseWriter, r *http.Request) {
		_, _, ok := api.GetSessionUser(r)
		if !ok {
			io.WriteString(w, "No User found")
			return
		}
		tmpl, err := helper.MakeTemplateFromGlobAndCollect(site.TmplPath)
		if err != nil {
			panic(err)
		}
		tags, err := api.QueryAllTags()
		if err != nil {
			panic(err)
		}
		page, err := helper.MakePageStr(tmpl, "./site_src/pages/new_notes.gohtml", tags)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, page)
	})

	site.PageRouteDef("GET /login", "./site_src/pages/login.gohtml", nil)

	site.PageRouteDef("GET /register", "./site_src/pages/register.gohtml", nil)

	site.PageRouteDef("GET /new_assignment", "./site_src/pages/new_assignment.gohtml", nil)

	site.PageRouteDef("GET /assignments", "./site_src/pages/view_all_assignment.gohtml", nil)

	
	// site.RouteFunc("GET /api/test", func(w http.ResponseWriter, r *http.Request) {
	// 	io.WriteString(w, "<p>This is from api</p>")
	// })
	site.RouteFunc("GET /logout", api.LogoutHandler)

	site.RouteFunc("POST /api/login", api.LoginHandler)
	site.RouteFunc("POST /api/register", api.RegisterHandler)
	site.RouteFunc("POST /api/new_note", api.NewNoteHandler)
	site.RouteFunc("POST /api/new_assignment", api.NewAssignmentHandler)

	site.RouteFunc("PATCH /api/check_assignment/{assignmentID}/{isDone}", api.CheckAssignmentHandler)

	site.RouteFunc("PUT /api/delete_assignment/{id}", api.DeleteAssignmentHandler)

	site.RouteFunc("GET /api/notes", api.AllPublicNotesHandlerTest)
	site.RouteFunc("GET /api/assignments", api.AllAssignmentHandler)
	site.RouteFunc("GET /api/navbaruser", api.NavBarUserHandler)
	go func() {
		err = site.Serve()
		if err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Stopping server")
}

func main2() {
	err = api.InitDB("./test.db", api.DebugModeEnabled)
	if err != nil {
		panic(err)
	}
	defer api.CloseDB()

	err = api.InitTestDummyDataDB()
	if err != nil {
		panic(err)
	}
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)
	// <-stop
	// fmt.Println("Closing app")
}

func flags_test() {
	is := flag.Bool("clear-db", false, "Clears the database")
	flag.Parse()
	if *is {
		fmt.Println("Cleaning Database")
	} else {
		fmt.Println("Doing nothing")		
	}
}

func main3() {
	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Db Connected")
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS posts;
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
INSERT INTO posts (userid, content) VALUES (1, 'HELLO WORLD I AM HERE');`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table created and inserted")
	res, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		panic(err)
	}
	columns, err := res.Columns()
	if err != nil {
		panic(err)
	}	
	fmt.Println(columns)
	for res.Next() {
		id := 0
		name := ""
		password := ""
		err = res.Scan(&id, &name, &password)
		if err != nil {
			panic(err)
		}	
		fmt.Println(id, name, password)
	}
}

func server_test() {
	// str := "<h1>this is not from partial</h1>{{template \"test\"}}"

	// site.Tmpl, err = helper.MakeTemplateFromGlobAndCollect("./site_src/partials/*.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	// tmpl, err := site.Tmpl.New("html").Parse(str)
	// if err != nil {
	// 	panic(err)
	// }

	// buf := bytes.NewBuffer(nil)
	// err = tmpl.Execute(buf, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// page := Page {
	// 	page: buf.String(),
	// }

	// site.Mux.Handle("/", page)

	site, err := MakeSite(":6969", "./site_src/partials/*.gohtml")

	site.StaticPath("./site_src/static/")
	
	site.PageRouteDef("GET /", "./site_src/pages/index.gohtml", nil)

	site.RouteFunc("GET /api/test", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<p>This is from api</p>")
	})

	go func() {
		err = site.Serve()
		if err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Stopping server")
}

