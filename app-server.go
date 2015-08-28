package dokpi

import (
	"encoding/json"
	"log"
	"net/http"
)

// AppServer is server that run apps. The command are send through
// http to this server.
type AppServer struct {
	appsRunner *appsRunner
}

// NewAppServer create a new app server
func NewAppServer() *AppServer {
	ar, err := newAppsRunner()
	if err != nil {
		log.Fatal(err)
	}
	return &AppServer{appsRunner: ar}
}

func (as *AppServer) prepareResponse() (resp msg) {
	resp = msg{cmdSuccess: true}
	as.appsRunner.updateAvailableApps()
	resp.appsStatus = as.appsRunner.getAppsStatus()
	return
}

func sendResponse(w http.ResponseWriter, resp msg) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("send response error", err)
	}
}

// StartServer start the server
func (as *AppServer) StartServer() error {
	// see http://thenewstack.io/make-a-restful-json-api-go/
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		resp := as.prepareResponse()
		defer sendResponse(w, resp)

		name := r.FormValue("name")
		if name == "" {
			resp.setError("the url must contains an app name")
			return
		}

		err := as.appsRunner.startApp(name)
		if err != nil {
			resp.setError(err.Error())
			return
		}

		resp.appsStatus = as.appsRunner.getAppsStatus()
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		resp := as.prepareResponse()
		defer sendResponse(w, resp)

		name := r.FormValue("name")
		if name == "" {
			resp.setError("the url must contains an app name")
			return
		}

		err := as.appsRunner.stopApp(name)
		if err != nil {
			resp.setError(err.Error())
			return
		}

		resp.appsStatus = as.appsRunner.getAppsStatus()
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		resp := as.prepareResponse()
		defer sendResponse(w, resp)
	})

	log.Fatal(http.ListenAndServe(listenAndServePort, nil))
	return nil
}
