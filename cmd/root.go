package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "centralconfig",
	Short: "A simple REST service and UI for application configuration",
	Long: `Centralconfig is a REST based service for managing application configuration.
It's designed to be used with one of many different SQL (or nosql) backends.  You can
use both an API and a web UI to manage configuration information.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/centralconfig.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("centralconfig") // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // adding home directory as first search path
	viper.AddConfigPath(".")             // also look in the working directory
	viper.AutomaticEnv()                 // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[ERROR] Problem reading config file: %s\n", err)
		fmt.Println(`
	There was a problem reading the server configuration file.  
	If you need help creating a configuration file, you can use 
	the 'defaults' command to generate a new server configuration file.  
	Use "centralconfig defaults --help" if you need help.

	Quick start:
	To generate a server configuration file, run the following command: 

	centralconfig defaults > centralconfig.yaml
			`)

		//	We really shouldn't proceed.
		//	Use non-zero status to indicate failure.
		//	from https://golang.org/pkg/os/#Exit
		os.Exit(1)
	}
}
