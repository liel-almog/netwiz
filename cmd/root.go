/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/lielalmog/netwiz/cmd/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netwiz",
	Short: "Network Wizard: A versatile network toolkit",
	Long: `Network Wizard (netwiz) is a versatile network toolkit built in Go, offering a range of utilities for network exploration and diagnostics. It includes a simple HTTP client, a ping tool, and a port scanner, making it a handy toolset for network administrators, developers, and IT professionals.

The nw toolkit is designed to be intuitive and user-friendly, with commands and flags that follow conventional CLI patterns. 

Examples of use:

1. HTTP Client:
   Send a GET request:
   $ netwiz http --url https://example.com
   
   Send a POST request with data:
   $ netwiz http --url https://example.com/api --method post --data '{"key":"value"}'

2. Ping Tool:
   Ping a host:
   $ netwiz ping google.com

   Ping a host with a specific number of echo requests:
   $ netwiz ping --number 5 google.com

3. Port:
   Scan a host for open ports:
   $ netwiz port --host 192.168.1.1

   Scan a host within a specific port range:
   $ netwiz port --host 192.168.1.1 --range 80-100

   Kill a process running on a specific port:
   $ netwiz port --kill 8080

For more information and detailed usage of each command, use the help command followed by the command name, e.g., 'nw help httpclient'.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netwiz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(http.HttpCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".netwiz" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".netwiz")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
