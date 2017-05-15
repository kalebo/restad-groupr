package main

//"github.com/gorilla/mux"
//"github.com/go-zoo/bone"
//"github.com/dimfeld/httptreemux"

import (
	"flag"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	// finally settling on this muxer because:
	// it plays nice with the global map in cas
	// and it correctly serves static assets
	"github.com/szxp/mux"

	"github.com/golang/glog"

	"gopkg.in/cas.v1"
)

// TemplateBinding specifies the NetId username that the templates should be rendered with
type TemplateBinding struct {
	Username string
}

var (
	casURL = &url.URL{Scheme: "https", Host: "cas.byu.edu", Path: "/cas/"}
	apiURL = &url.URL{Scheme: "http", Host: "10.25.82.180:1234"}

	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	templates = template.New("").Funcs(templateMap)
)

func init() {
}

func main() {
	var port string

	flag.StringVar(&port, "port", "8080", "the `port` to listen on.")
	flag.Parse()

	glog.Info("Starting...")

	//r := http.NewServeMux()
	//r := mux.NewRouter() // gorilla mux isn't playing well with cas on go1.7
	//r := bone.New()
	//r := httptreemux.NewContextMux()
	r := mux.NewMuxer()

	client := cas.NewClient(&cas.Options{
		URL: casURL,
	})

	// Backend API Routes
	r.HandleFunc("/api/", APIEndpoints)
	//r.HandleFunc("/api/{path:.*}", ApiEndpoints) // requires gorilla mux

	// Main App Routes
	r.HandleFunc("/app", App)

	// CAS Authentication Routes
	r.HandleFunc("/cas/login", cas.RedirectToLogin)
	r.HandleFunc("/cas/logout", cas.RedirectToLogout)

	fs := http.FileServer(assetFS())
	r.Handle("/", fs)

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

// APIEndpoints acts as a reverse proxy to the RESTAD API backend
func APIEndpoints(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
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
