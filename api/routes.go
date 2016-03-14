package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danesparza/centralconfig/datastores"
	"github.com/gorilla/mux"
)

//	Bogus test route
func TestRoute(rw http.ResponseWriter, req *http.Request) {
	//	Parse the twitter name from the url
	twitterName := mux.Vars(req)["twitterName"]

	configItems := []datastores.ConfigItem{}

	configItems = append(configItems, datastores.ConfigItem{
		Application: "Testing",
		Name:        "Bogus",
		Value:       "Your mom's a config, " + twitterName})

	//	Set the content type header and return the JSON
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(configItems)
}

func GetConfig(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "This is where stuff would go, %s", "right here")
}
