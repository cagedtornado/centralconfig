package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/cagedtornado/centralconfig/datastores"
)

var (
	WsHub = NewHub()
)

func ShowUI(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/ui/", 301)
}

//	Gets a specfic config item based on application and config item name
func GetConfig(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := datastores.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	response, err := ds.Get(request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	If we found an item, return it (otherwise, return an empty item):
	configItem := datastores.ConfigItem{}
	if response.Name != "" {
		configItem = response
		sendDataResponse(rw, "Config item found", configItem)
		return
	}

	sendDataResponse(rw, "No config item found with that application and name", configItem)
}

//	Set a specific config item
func SetConfig(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := datastores.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	response, err := ds.Set(request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
	} else {
		WsHub.Broadcast <- []byte(getWSResponse("Updated", response))
		sendDataResponse(rw, "Config item updated", response)
	}
}

//	Removes a specific config item
func RemoveConfig(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := datastores.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	err = ds.Remove(request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
	} else {
		WsHub.Broadcast <- []byte(getWSResponse("Removed", request))
		sendDataResponse(rw, "Config item removed", request)
	}
}

//	Gets all config information for a given application
func GetAllConfigForApp(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := &datastores.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	configItems, err := ds.GetAllForApplication(request.Application)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	If we found an item, return it (otherwise, return an empty array):
	if len(configItems) > 0 {
		sendDataResponse(rw, "Config items found", configItems)
		return
	}

	sendDataResponse(rw, "No config items found with that application", configItems)
}

//	Gets all config information
func GetAllConfig(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	configItems, err := ds.GetAll()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	If we found an item, return it (otherwise, return an empty array):
	if len(configItems) > 0 {
		sendDataResponse(rw, "Config items found", configItems)
		return
	}

	sendDataResponse(rw, "No config items found", configItems)
}

//	Gets all applications
func GetAllApplications(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get the current datastore:
	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	applications, err := ds.GetAllApplications()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	If we found an item, return it (otherwise, return an empty array):
	if len(applications) > 0 {
		sendDataResponse(rw, "Applications found", applications)
		return
	}

	sendDataResponse(rw, "No config items found", applications)
}

//	Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := datastores.ConfigResponse{
		Status:  code,
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

//	Used to send back a response with data
func sendDataResponse(rw http.ResponseWriter, message string, dataItems interface{}) {
	//	Our return value
	response := datastores.ConfigResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    dataItems}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

//	Gets a JSON formatted WebSocket event response
func getWSResponse(messageType string, item datastores.ConfigItem) string {
	//	Our default return value:
	retval := ""

	//	Our WebSocket return value
	response := datastores.WebSocketResponse{
		Data: item,
		Type: messageType}

	//	Serialize to JSON and return as a string:
	responseBytes := new(bytes.Buffer)
	if err := json.NewEncoder(responseBytes).Encode(&response); err == nil {
		retval = responseBytes.String()
	}

	return retval
}
