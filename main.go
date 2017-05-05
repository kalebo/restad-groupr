package main

import (
	"github.com/dimfeld/httptreemux"
	"github.com/golang/glog"
	"gopkg.in/cas.v1"
	"net/http"
	"net/url"
)

// Globals

var casURLstr string = "https://cas.byu.edu/cas/"
var apiURLstr string = "http://avari:1234/api/"

type apiHandler struct{}
type frontEndHandler struct{}

func init() {

}

func main() {
	glog.Info("starting..")

	router := httptreemux.New()
	api := router.NewGroup("/api")
	app := router.NewGroup("/app")

	casUrl, _ := url.Parse(casURLstr)
	casClient := cas.NewClient(&cas.Options{
		URL: casUrl,
	})

	api.GET("/*", ProxyAPI)

	http.ListenAndServe(":8080", router)

}

func (h *apiHandler) ProxyAPI(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// Passes incoming requests to RESTAD
	if !cas.IsAuthenticated(r) {
		//cas.RedirectToLogin(w, r)
		return
	}

}

func (h *frontEndHandler) MainApp(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}

}

func (h *frontEndHandler) CasLogout(w http.ResponseWriter, r *http.Request, params map[string]string)  {
	cas.RedirectToLogout(w, r)
	return
}
