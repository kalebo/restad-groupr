package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gopkg.in/cas.v1"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Globals

type ApiHandler struct{}
type FrontEndHandler struct{}

var (
	casURLstr string = "https://cas.byu.edu/cas/"
	apiURLstr string = "http://avari:1234"
	//apiURLstr string = "http://localhost:1111"
)

func init() {

}

func main() {
	glog.Info("starting..")

	r := mux.NewRouter()

	casUrl, _ := url.Parse(casURLstr)
	apiUrl, _ := url.Parse(apiURLstr)

	client := cas.NewClient(&cas.Options{
		URL: casUrl,
	})

	// Backend API Routes
	r.Handle("/api/{path:.*}", httputil.NewSingleHostReverseProxy(apiUrl))

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

func ApiEndpoints(w http.ResponseWriter, r *http.Request) {
	// Passes incoming requests to RESTAD
	//if !cas.IsAuthenticated(r) {
	//       return
	//}

	vars := mux.Vars(r)
	apiURLstr := vars["path"]

	fmt.Println(apiURLstr + r.URL.Path)
	// apiEndpoint, _ := url.Parse(apiURLstr + r.URL.Path)
	// proxy := httputil.NewSingleHostReverseProxy(apiEndpoint)

}

func MainApp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Is hit")
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}
	fmt.Println("Is authd")
}
