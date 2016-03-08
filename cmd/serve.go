package cmd

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {

	r := mux.NewRouter()
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
