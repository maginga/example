/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"example/apm/client"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// datasourceCmd represents the datasource command
var datasourceCmd = &cobra.Command{
	Use:   "datasource",
	Short: "(10) Create a datasource on DRUID from a file.",
	Long: `Create a datasource on DRUID from a file.
For example: 
	apm create datasource [direct] [Json File]
	apm create datasource [alarm, score, trace] [nest ID]

`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
		discoveryURL := fmt.Sprintf("%v", viper.Get("discovery.url"))

		c, err := client.NewDruidClient(discoveryURL)
		if err != nil {
			panic(err)
		}

		if args[0] == "direct" {
			c.Create(args[1])
			logger.Println("The datasource was created on Druid Cluster.")
		} else if args[0] == "alarm" {
			c.CreateAlarm()
			logger.Println("The ALARM datasource was created on Druid Cluster.")
		} else if args[0] == "score" {
			c.CreateScore(args[1])
			logger.Println("The SCORE datasource was created on Druid Cluster.")
		} else if args[0] == "trace" {
			c.CreateTrace(args[1])
			logger.Println("The TRACE datasource was created on Druid Cluster.")
		} else {
			logger.Println("The argument is not valid.")
		}
	},
}

func init() {
	createCmd.AddCommand(datasourceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// datasourceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// datasourceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
