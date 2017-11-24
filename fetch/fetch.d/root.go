// Copyright © 2017 Slotix s.r.o. <dm@slotix.sk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"os"

	"github.com/slotix/dataflowkit/fetch"
	"github.com/slotix/dataflowkit/healthcheck"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//VERSION               string // VERSION is set during build
	DFKFetch string //API_GATEWAY_ADDRESS

	splashHost            string
	splashTimeout         int
	splashResourceTimeout int
	splashWait            float64

	storageType     string
	storageExpires  int64 //how long in seconds object stay in a cache before expiration.
	diskvBaseDir    string
	fetchBucket     string
	redisHost       string
	redisExpire     int
	redisNetwork    string
	redisPassword   string
	redisDB         int
	redisSocketPath string

	//sqsQueueFetchURLIn  string
	//sqsQueueFetchURLOut string
	//sqsAWSRegion        string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dataflowkit",
	Short: "Dataflow Kit html fetcher",
	Long:  `Dataflow Kit fetch service retrieves html pages from websites and passes content to DFK parser service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking services ... ")

		services := []healthcheck.Checker{

			healthcheck.SplashConn{
				Host: splashHost,
			},
		}
		if storageType == "Redis" {
			services = append(services, healthcheck.RedisConn{
				Network: redisNetwork,
				Host:    redisHost})
		}
		status := healthcheck.CheckServices(services...)
		allAlive := true

		for k, v := range status {
			fmt.Printf("%s: %s\n", k, v)
			if v != "Ok" {
				allAlive = false
			}
		}
		if allAlive {
			fmt.Printf("Starting Server %s\n", DFKFetch)
			fetch.Start(DFKFetch)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	//flags and configuration settings. They are global for the application.

	RootCmd.Flags().StringVarP(&DFKFetch, "DFK_FETCH", "a", "127.0.0.1:8000", "HTTP listen address")

	RootCmd.Flags().StringVarP(&splashHost, "SPLASH", "s", "127.0.0.1:8050", "Splash host address")
	RootCmd.Flags().IntVarP(&splashTimeout, "SPLASH_TIMEOUT", "", 20, "A timeout (in seconds) for the render.")
	RootCmd.Flags().IntVarP(&splashResourceTimeout, "SPLASH_RESOURCE_TIMEOUT", "", 30, "A timeout (in seconds) for individual network requests.")
	RootCmd.Flags().Float64VarP(&splashWait, "SPLASH_WAIT", "", 0.5, "Time in seconds to wait until js scripts loaded.")

	//set here default type of storage
	RootCmd.Flags().StringVarP(&storageType, "STORAGE_TYPE", "", "Diskv", "Storage backend for intermediary data passed to html parser. Types: S3, Redis, Diskv")
	RootCmd.Flags().Int64VarP(&storageExpires, "STORAGE_EXPIRE", "", 3600, "Default Storage expire value in seconds")
	RootCmd.Flags().StringVarP(&diskvBaseDir, "DISKV_BASE_DIR", "", "diskv", "diskv base directory for storing fetch results")
	RootCmd.Flags().StringVarP(&fetchBucket, "FETCH_BUCKET", "", "fetch-bucket", "S3 bucket name for storing fetch results")

	RootCmd.Flags().StringVarP(&redisHost, "REDIS", "r", "127.0.0.1:6379", "Redis host address")
	RootCmd.Flags().IntVarP(&redisExpire, "REDIS_EXPIRE", "", 3600, "Default Redis expire value in seconds")
	RootCmd.Flags().StringVarP(&redisNetwork, "REDIS_NETWORK", "", "tcp", "Redis Network")
	RootCmd.Flags().StringVarP(&redisPassword, "REDIS_PASSWORD", "", "", "Redis Password")
	RootCmd.Flags().IntVarP(&redisDB, "REDIS_DB", "", 0, "Redis DB")
	RootCmd.Flags().StringVarP(&redisSocketPath, "REDIS_SOCKET_PATH", "", "", "Redis Socket Path")

	//RootCmd.Flags().StringVarP(&sqsQueueFetchURLIn, "SQS_QUEUE_FETCH_URL_IN", "", "https://sqs.us-east-1.amazonaws.com/060679207441/fetch-in", "SQS Queue Fetch URL In")
	//RootCmd.Flags().StringVarP(&sqsQueueFetchURLOut, "SQS_QUEUE_FETCH_URL_OUT", "", "https://sqs.us-east-1.amazonaws.com/060679207441/fetch-out", "SQS Queue Fetch URL Out")
	//RootCmd.Flags().StringVarP(&sqsAWSRegion, "SQS_AWS_REGION", "", "us-east-1", "SQS AWS Region")

	viper.AutomaticEnv() // read in environment variables that match
	viper.BindPFlag("DFK_FETCH", RootCmd.Flags().Lookup("DFK_FETCH"))

	viper.BindPFlag("SPLASH", RootCmd.Flags().Lookup("SPLASH"))
	viper.BindPFlag("SPLASH_TIMEOUT", RootCmd.Flags().Lookup("SPLASH_TIMEOUT"))
	viper.BindPFlag("SPLASH_RESOURCE_TIMEOUT", RootCmd.Flags().Lookup("SPLASH_RESOURCE_TIMEOUT"))
	viper.BindPFlag("SPLASH_WAIT", RootCmd.Flags().Lookup("SPLASH_WAIT"))

	viper.BindPFlag("STORAGE_TYPE", RootCmd.Flags().Lookup("STORAGE_TYPE"))
	viper.BindPFlag("STORAGE_EXPIRE", RootCmd.Flags().Lookup("STORAGE_EXPIRE"))
	viper.BindPFlag("DISKV_BASE_DIR", RootCmd.Flags().Lookup("DISKV_BASE_DIR"))
	viper.BindPFlag("FETCH_BUCKET", RootCmd.Flags().Lookup("FETCH_BUCKET"))
	viper.BindPFlag("REDIS", RootCmd.Flags().Lookup("REDIS"))
	viper.BindPFlag("REDIS_EXPIRE", RootCmd.Flags().Lookup("REDIS_EXPIRE"))
	viper.BindPFlag("REDIS_NETWORK", RootCmd.Flags().Lookup("REDIS_NETWORK"))
	viper.BindPFlag("REDIS_PASSWORD", RootCmd.Flags().Lookup("REDIS_PASSWORD"))
	viper.BindPFlag("REDIS_DB", RootCmd.Flags().Lookup("REDIS_DB"))
	viper.BindPFlag("REDIS_SOCKET_PATH", RootCmd.Flags().Lookup("REDIS_SOCKET_PATH"))

	//	viper.BindPFlag("SQS_QUEUE_FETCH_URL_IN", RootCmd.Flags().Lookup("SQS_QUEUE_FETCH_URL_IN"))
	//	viper.BindPFlag("SQS_QUEUE_FETCH_URL_OUT", RootCmd.Flags().Lookup("SQS_QUEUE_FETCH_URL_OUT"))
	//	viper.BindPFlag("SQS_AWS_REGION", RootCmd.Flags().Lookup("SQS_AWS_REGION"))
}