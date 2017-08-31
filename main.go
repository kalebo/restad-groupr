package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-zoo/bone"
	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/cas.v2"
)

const schema = `
CREATE TABLE IF NOT EXISTS groupadmin (
	id INTEGER PRIMARY KEY,
	principle TEXT,
	type integer NOT NULL,
	targetgroup TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS scheduledactions (
	id INTEGER PRIMARY KEY,
	time integer NOT NULL,
	actionurl TEXT NOT NULL
);
`

type principleType int

const (
	groupPrinciple principleType = iota
	userPrinciple
)

var db *sql.DB

// !!! the following line is a preprocessor directive !!!
//go:generate go-bindata -debug -prefix "dist/" -pkg main -o bindata.go dist/...
//go:generate go-bindata-assetfs -debug -prefix "dist/..." -pkg main -o bindata.go dist/...

// TemplateBinding specifies the NetId username that the templates should be rendered with
type TemplateBinding struct {
	Username string
}

var (
	port   string
	casURL *url.URL
	apiURL *url.URL

	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	templates = template.New("").Funcs(templateMap)
)

func init() {
	// Init the DB connection and if need be create the table
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		glog.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	db.Ping()

	db.Exec(schema)

	// Get parameters
	var apiURLRaw string
	var casURLRaw string

	flag.StringVar(&port, "port", os.Getenv("PORT"), "the `port` to listen on.")
	flag.StringVar(&apiURLRaw, "api-url", os.Getenv("RESTAD_URL"), "URL of Addict API Service")
	flag.StringVar(&casURLRaw, "cas-url", os.Getenv("CAS_URL"), "URL of CAS service")
	flag.Parse()

	// Parse url parameters
	apiURL = mustParse(url.Parse(apiURLRaw))
	casURL = mustParse(url.Parse(casURLRaw))
	if port == "" {
		port = "8080"
	}
}

func main() {
	glog.Info("Starting...")

	r := bone.New()

	client := cas.NewClient(&cas.Options{
		URL: casURL,
	})

	// Backend API Routes
	r.HandleFunc("/api/group/*", APIEndpoints)
	r.HandleFunc("/api/user/managed", UserManagedGroups)

	// CAS Authentication Routes
	r.HandleFunc("/cas/login", cas.RedirectToLogin)
	r.HandleFunc("/cas/logout", cas.RedirectToLogout)

	// Main App Routing
	//r.HandleFunc("/app", App)
	r.HandleFunc("/", App)
	r.HandleFunc("/refresh", refreshCookie)

	// Static Asset Routing
	fs := http.FileServer(assetFS())
	r.Handle("/*", fs)

	//fs := http.FileServer(assetFS())
	//r.Handle("/", fs)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: client.Handle(r),
		// enforce timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	server.ListenAndServe()
}

// RenderTemplate loads an asset by the path passed in by `tmpl` and executes
// it with the values in `p`.
func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	//err := templates.ExecuteTemplate(w, tmpl, p)

	// development debug
	asset, err := Asset(tmpl)
	t, err := template.New("main").Parse(string(asset))
	t.Execute(w, p)
	// end

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func UserManagedGroups(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var managedGroups []string

	username := cas.Username(r)
	rows, err := db.Query("SELECT targetgroup FROM groupadmin WHERE type=? AND principle=?", userPrinciple, username)
	if err != nil {
		log.Fatalf("Could not get DB records: %s", err)
	}
	defer rows.Close()

	var _group string
	for rows.Next() {
		rows.Scan(&_group)
		managedGroups = append(managedGroups, _group)
	}

	_jsonGroup, _ := json.Marshal(managedGroups)
	w.Write(_jsonGroup)
}

func StaticAssetHandler(rw http.ResponseWriter, req *http.Request) {
	path := "dist" + req.URL.Path // path to static assets

	// We don't convert path "" to "index.html" as the index is handled and rendered by another handler

	fmt.Println("Path: " + path)
	if bs, err := Asset(path); err != nil {
		rw.WriteHeader(http.StatusNotFound)
	} else {
		var reader = bytes.NewBuffer(bs)
		io.Copy(rw, reader)
	}
}

// APIEndpoints acts as a reverse proxy to the RESTAD API backend
func APIEndpoints(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(apiURL)
	proxy.ServeHTTP(w, r)
}

// App renders and serves the front end index.html
func App(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}

	binding := &TemplateBinding{
		Username: cas.Username(r),
	}

	RenderTemplate(w, "dist/index.html", binding)
}

func refreshCookie(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

func mustParse(_url *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}

	return _url
}
