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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Job on Flink Server.",
	Long: `Run a Job on Flink Server.
For example: apm flink run [Jar ID] [Entry Class] [Arguments]
`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		jobServer := fmt.Sprintf("%v", viper.Get("stream.jobServer"))
		log.Println("Job Server: " + jobServer)

		JarID := args[0]
		entryClass := args[1]
		programArgs := args[2]

		// Your flink server HTTP API
		c, err := client.New(jobServer)
		if err != nil {
			panic(err)
		}

		opts := client.RunOpts{}
		opts.JarID = JarID
		opts.EntryClass = entryClass
		opts.ProgramArg[0] = programArgs
		opts.Parallelism = 1

		r, err := c.RunJar(opts)

		log.Println("This Job was registered. [" + r.ID + "]")
	},
}

func init() {
	flinkCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
