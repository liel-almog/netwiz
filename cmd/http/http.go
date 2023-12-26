/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	url     string
	method  string
	data    string
	headers []string
)

// httpCmd represents the http command
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create an http request
		performHttpRequest(url, method, data, headers)
	},
}

func init() {
	HttpCmd.Flags().StringVarP(&url, "url", "u", "", "URL to make the request to")
	HttpCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")
	HttpCmd.Flags().StringVarP(&data, "data", "d", "", "Data to send with the request")
	HttpCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "Headers to include in the request")
	HttpCmd.MarkFlagRequired("url")
}

func performHttpRequest(url, method, data string, headers []string) {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		var jsonData = []byte(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	}

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add headers if provided
	if len(headers) != 0 {
		for _, header := range headers {
			var headerSplit = strings.Split(header, ":")
			var headerKey = headerSplit[0]
			var headerValue = headerSplit[1]
			req.Header.Add(headerKey, headerValue)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))

}
