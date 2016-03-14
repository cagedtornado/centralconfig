package cmd

import (
	"log"
	"net/http"

	"github.com/danesparza/centralconfig/api"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Get configuration information
	log.Printf("[INFO] Using BoltDB database: %s", viper.GetString("boltdb.database"))

	//	Create a router and setup our REST endpoints...
	var Router = mux.NewRouter()

	//	Setup our routes
	Router.HandleFunc("/news/{twitterName}", api.TestRoute)
	Router.HandleFunc("/", api.GetConfig)

	log.Printf("[INFO] HTTP server info: %s:%s", viper.GetString("http.bind"), viper.GetString("http.port"))
	http.ListenAndServe(viper.GetString("http.bind")+":"+viper.GetString("http.port"), Router)
}

func init() {
	RootCmd.AddCommand(serveCmd)

	//	Setup our flags
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	serveCmd.Flags().StringVarP(&serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("http.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("http.bind", serveCmd.Flags().Lookup("bind"))
}
