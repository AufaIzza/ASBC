package main

import (
	// "io"
	"os"
	// "bytes"
	"fmt"
	// "log"
	"net/http"
	"html/template"
	"database/sql"
	"github.com/AufaIzza/ASBC/api"
	"github.com/AufaIzza/ASBC/helper"
	// "github.com/mattn/go-sqlite3"
	// "modernc.org/sqlite"
	"flag"
	"os/signal"
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
	Tmpl *template.Template
	Port string
}

func MakeSite(port string, partialsPath string) (*Site, error) {
	site := Site {
		Port: port,
	}
	site.Tmpl, err = helper.MakeTemplateFromGlobAndCollect("./site_src/partials/*.gohtml")
	if err != nil {
		return nil, err
	}
	site.Mux = http.NewServeMux()
	return &site, nil
}

func (s *Site) PageRoute(path string, filepath string, props map[string]string) {
	err := helper.MakePageFromFile(s.Mux, s.Tmpl, path, filepath, props)
	if err != nil {
		panic(err)
	}
}

func (s *Site) Serve() error {
	fmt.Println("Listening on port ", s.Port)
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

func main() {
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

func db_test() {
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

func page_test() {
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
	site.PageRoute("GET /", "./site_src/pages/index.gohtml", map[string]string {
		"title": "ASBC - Main",
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
