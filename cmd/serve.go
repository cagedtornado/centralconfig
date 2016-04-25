package cmd

import (
	"log"
	"net/http"

	"github.com/danesparza/centralconfig/api"
	"github.com/danesparza/centralconfig/datastores"

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

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Log the datastore information we have:
	logDatastoreInfo()

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

	//	If we don't have a UI directory specified...
	if viper.GetString("http.ui-dir") == "" {
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
	} else {
		//	Use the supplied directory:
		log.Printf("[INFO] Using UI directory: %s", viper.GetString("http.ui-dir"))
		Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("http.ui-dir")))))
	}

	log.Printf("[INFO] Starting HTTP server: %s:%s", viper.GetString("http.bind"), viper.GetString("http.port"))
	http.ListenAndServe(viper.GetString("http.bind")+":"+viper.GetString("http.port"), Router)
}

func init() {
	RootCmd.AddCommand(serveCmd)

	//	Setup our flags
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	serveCmd.Flags().StringVarP(&serverInterface, "bind", "i", "127.0.0.1", "interface to which the server will bind")
	serveCmd.Flags().StringVarP(&serverUIDirectory, "ui-dir", "u", "", "directory for the UI")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("http.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("http.bind", serveCmd.Flags().Lookup("bind"))
	viper.BindPFlag("http.ui-dir", serveCmd.Flags().Lookup("ui-dir"))
}

func logDatastoreInfo() {
	//	Get configuration information
	ds := datastores.GetConfigDatastore()

	switch t := ds.(type) {
	case datastores.MySqlDB:
		log.Printf("[INFO] Using MySQL server: %s", ds.(datastores.MySqlDB).Address)
		log.Printf("[INFO] Using MySQL database: %s", ds.(datastores.MySqlDB).Database)
	case datastores.MSSqlDB:
		log.Printf("[INFO] Using MSSQL server: %s", ds.(datastores.MSSqlDB).Address)
		log.Printf("[INFO] Using MSSQL database: %s", ds.(datastores.MSSqlDB).Database)
	case datastores.BoltDB:
		log.Printf("[INFO] Using BoltDB database: %s", ds.(datastores.BoltDB).Database)
	default:
		_ = t
		log.Println("[ERROR] Can't determine datastore type")
	}
}
