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

var (
	casUrl *url.URL = &url.URL{Scheme: "https", Host: "cas.byu.edu", Path: "/cas/"}
	apiUrl *url.URL = &url.URL{Scheme: "http", Host: "avari:1234"}
)

func init() {

}

func main() {
	fmt.Println("Starting...")

	glog.Info("Logging started!")

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
}
