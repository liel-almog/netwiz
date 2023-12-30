/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	url     string
	method  string
	data    string
	headers []string
	retry   int
)

const maxRetries = 5

// httpCmd represents the http command
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if retry > maxRetries {
			return fmt.Errorf("retry value cannot exceed %d", maxRetries)
		}

		upperCaseMethod := strings.ToUpper(method)
		performHttpRequest(url, upperCaseMethod, data, headers)
		return nil
	},
}

func init() {
	HttpCmd.Flags().StringVarP(&url, "url", "u", "", "URL to make the request to")
	HttpCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")
	HttpCmd.Flags().StringVarP(&data, "data", "d", "", "Data to send with the request")
	HttpCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "Headers to include in the request")
	HttpCmd.Flags().IntVarP(&retry, "retry", "r", 0, "Number of times to retry the request")
	HttpCmd.MarkFlagRequired("url")
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func drainRespBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}

	return false
}

func performHttpRequest(url, method, data string, headers []string) error {
	var req *http.Request
	var err error
	bytesData := []byte(data)

	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		nopCloserBody := io.NopCloser(bytes.NewBuffer(bytesData))
		req, err = http.NewRequest(method, url, nopCloserBody)
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

	for i := 0; shouldRetry(err, resp) && i < retry; i++ {
		delay := backoff(i)

		if err == nil {
			drainRespBody(resp)
			resp.Body.Close()
		}

		req.Body = io.NopCloser(bytes.NewBuffer(bytesData))
		time.Sleep(delay)
		resp, err = client.Do(req)
	}

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	return nil
}
