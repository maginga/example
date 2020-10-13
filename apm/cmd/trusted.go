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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// trustedCmd represents the trusted command
var trustedCmd = &cobra.Command{
	Use:   "trusted",
	Short: "trusted",
	Long: `
For example: apm trusted
`,
	Run: func(cmd *cobra.Command, args []string) {
		discoveryURL := fmt.Sprintf("%v", viper.Get("discovery.url"))
		log.Println("Discovery URL: " + discoveryURL)
		c, err := client.NewDruidClient(discoveryURL)
		if err != nil {
			panic(err)
		}

		// t, err := c.InitOrganization()
		// if err != nil {
		// 	panic(err)
		// }
		// logger.Println("The organization was initialized.: " + t)

		q, err := c.ConfigureUser()
		if err != nil {
			panic(err)
		}

		log.Println("Login Setup Completed. " + q)
	},
}

func init() {
	rootCmd.AddCommand(trustedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trustedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trustedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
