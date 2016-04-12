package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danesparza/centralconfig/api"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverInterface   string
	serverPort        int
	serverUIDirectory string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the config server",
	Long: `Centralconfig provides its own webserver which can serve both the 
	API and the UI for the app.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {

	if ProblemWithConfigFile == true {
		fmt.Println(`
	There was a problem reading the server configuration file.  
	
	If you need help creating a configuration file, you can use 
	the 'defaults' command to generate a new server configuration file.  
	Use "centralconfig defaults --help" if you need help.

	=== Quick start === 
	To generate a server configuration file, run the following command: 
	
	centralconfig defaults > centralconfig.yaml
			`)

		//	We really shouldn't proceed.
		//	Use non-zero status to indicate failure.
		//	from https://golang.org/pkg/os/#Exit
		os.Exit(1)
	}

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Get configuration information
	log.Printf("[INFO] Using BoltDB database: %s", viper.GetString("boltdb.database"))

	//	Create a router and setup our REST endpoints...
	var Router = mux.NewRouter()

	//	Setup our routes
	Router.HandleFunc("/", api.ShowUI)
	Router.HandleFunc("/config/get", api.GetConfig)
	Router.HandleFunc("/config/set", api.SetConfig)
	Router.HandleFunc("/config/remove", api.RemoveConfig)
	Router.HandleFunc("/config/getall", api.GetAllConfig)
	Router.HandleFunc("/config/getallforapp", api.GetAllConfigForApp)
	Router.HandleFunc("/config/init", api.InitStore)
	Router.HandleFunc("/applications/getall", api.GetAllApplications)

	//	Use the static assets file generated with
	//	https://github.com/elazarl/go-bindata-assetfs using the centralconfig-ui from
	//	https://github.com/danesparza/centralconfig-ui.
	//
	//	To generate this file, place the 'ui'
	//	directory under the main centralconfig directory and run the commands:
	//	go-bindata-assetfs.exe -pkg cmd ./ui/...
	//	mv bindata_assetfs.go cmd
	//	go install ./...
	Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(assetFS())))

	log.Printf("[INFO] Starting HTTP server: %s:%s", viper.GetString("http.bind"), viper.GetString("http.port"))
	http.ListenAndServe(viper.GetString("http.bind")+":"+viper.GetString("http.port"), Router)
}

func init() {
	RootCmd.AddCommand(serveCmd)

	//	Setup our flags
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	serveCmd.Flags().StringVarP(&serverInterface, "bind", "i", "127.0.0.1", "interface to which the server will bind")
	serveCmd.Flags().StringVarP(&serverUIDirectory, "ui-dir", "u", "./ui/", "directory for the UI")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("http.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("http.bind", serveCmd.Flags().Lookup("bind"))
	viper.BindPFlag("http.ui-dir", serveCmd.Flags().Lookup("ui-dir"))
}
