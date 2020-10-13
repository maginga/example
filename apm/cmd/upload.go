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

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload Jar file on Flink Server.",
	Long: `Upload Jar file on Flink Server.
For example: apm flink upload [Jar File Path]
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
		jobServer := fmt.Sprintf("%v", viper.Get("stream.jobServer"))
		logger.Println("Job Server: " + jobServer)

		jarFile := args[0]

		// Your flink server HTTP API
		c, err := client.New(jobServer)
		if err != nil {
			panic(err)
		}

		d, err := c.Config()
		logger.Println("Flink Version: " + d.FlinkVersion)

		k, err := c.Jars()
		for _, f := range k.Files {
			logger.Println("Existed Jars: " + f.ID)
		}

		u, err := c.UploadJar(jarFile)
		logger.Println("Jar ID: " + u.FileName)
		logger.Println("This Jar was uploaded.")

	},
}

func init() {
	flinkCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
