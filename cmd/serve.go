package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	serverInterface string
	serverPort      int
)

type Tweet struct {
	Id         int64  `json:"id"`
	CreateTime int64  `json:"createtime"`
	Text       string `json:text`
	MediaUrl   string `json:url`
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the config server",
	Long: `Centralconfig provides its own webserver which can serve both the 
	API and the UI for the app.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	//	Get configuration information
	fmt.Println(viper.GetString("boltdb.database"))

	//	Create a router and setup our REST endpoints...
	r := mux.NewRouter()

	//	Handle config get
	r.HandleFunc("/news/{twitterName}", func(w http.ResponseWriter, r *http.Request) {

		//	Parse the twitter name from the url
		twitterName := mux.Vars(r)["twitterName"]

		//	Our return values:
		tweets := []Tweet{}

		tweets = append(tweets, Tweet{
			Id:       42,
			Text:     "Some bogus text for " + twitterName,
			MediaUrl: "http://yourmomsaurl.com"})

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(tweets)
	})

	portString := strconv.Itoa(serverPort)
	http.ListenAndServe(serverInterface+":"+portString, r)
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	serveCmd.Flags().StringVarP(&serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")

}
