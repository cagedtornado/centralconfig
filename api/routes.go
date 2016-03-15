package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/danesparza/centralconfig/datastores"
)

func ShowHelp(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "ShowHelp method")
}

func GetConfig(rw http.ResponseWriter, req *http.Request) {

	//	Read the request:
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))

	//	If we have any errors or problems getting the body, return an HTTP error:
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	if err := req.Body.Close(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	//	Unmarshal to a config item:
	request := &datastores.ConfigItem{}
	if err := json.Unmarshal(body, request); err != nil {

		//	If we have an error, return an HTTP error:
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(rw).Encode(err); err != nil {
			http.Error(rw, err.Error(), 422) // unprocessable entity
		}
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	response, _ := ds.Get(request)

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

func SetConfig(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "SetConfig method")
}

func GetAllConfig(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "GetAllConfig method")
}

func InitStore(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "InitStore method")
}
