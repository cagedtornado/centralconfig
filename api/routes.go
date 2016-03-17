package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danesparza/centralconfig/datastores"
)

func ShowHelp(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "ShowHelp method")
}

func GetConfig(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := &datastores.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		//	Do something with the err if we have one
		fmt.Println(err)
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	response, err := ds.Get(request)

	//	Do something with the err if we have one

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

func sendJSONError(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := datastores.ConfigResponse{
		Error:   err,
		Status:  code,
		Message: err.Error()}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
