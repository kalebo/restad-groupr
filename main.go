package main

import (
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gopkg.in/cas.v1"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"flag"
	"strings"
)

type Empty struct{}

var (
	casUrl *url.URL = &url.URL{Scheme: "https", Host: "cas.byu.edu", Path: "/cas/"}
	apiUrl *url.URL = &url.URL{Scheme: "http", Host: "avari:1234"}

	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	templates = template.New("").Funcs(templateMap)
)

func init() {
	for _, path := range AssetNames() {
		bytes, err := Asset(path)
		if err != nil {
			glog.Warningf("Unable to parse: path=%s, err=%s", path, err)

		}
		templates.New(path).Parse(string(bytes))
	}
}

func main() {
	flag.Parse()
	glog.Info("Starting...")

	r := mux.NewRouter()

	client := cas.NewClient(&cas.Options{
		URL: casUrl,
	})

	// Backend API Routes
	r.HandleFunc("/api/{path:.*}", ApiEndpoints)

	// Main App Routes
	r.HandleFunc("/app", MainApp)

	// CAS Authentication Routes
	r.HandleFunc("/cas/login", cas.RedirectToLogin)
	r.HandleFunc("/cas/logout", cas.RedirectToLogout)

	server := &http.Server{
		Addr:    ":8080",
		Handler: client.Handle(r),
	}

	server.ListenAndServe()
}

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

func ApiEndpoints(w http.ResponseWriter, r *http.Request) {
	// Passes incoming requests to RESTAD
	if !cas.IsAuthenticated(r) {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(apiUrl)
	proxy.ServeHTTP(w, r)
}

func MainApp(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}

	RenderTemplate(w, "data/test.html", &Empty{})
}
