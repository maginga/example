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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "(11) Registers a Job to the flink cluster, and runs a Job to be registered stream job.",
	Long: `Registers a Job to the flink cluster, and runs a Job to be registered stream job.
For example: apm create stream [refiner, alarm, paramalarm, fdc, spc, mva, bae, current] [nest ID] [jar file]
`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
		jobServer := fmt.Sprintf("%v", viper.Get("stream.jobServer"))
		logger.Println("Job Server: " + jobServer)

		jarFile := args[2]

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
		if err != nil {
			panic(err)
		}

		if len(u.FileName) > 0 {
			logger.Println("file location: " + u.FileName)
			logger.Println("This Jar was uploaded.")

			apmAddr := fmt.Sprintf("%v", viper.Get("grandview.url"))
			nestID := args[1]
			specURL := "http://" + apmAddr + "/api/spec/" + nestID

			jarID := strings.Split(u.FileName, "/")
			logger.Println("jar ID: " + jarID[len(jarID)-1])

			opts := client.RunOpts{}
			opts.JarID = jarID[len(jarID)-1]

			if args[0] == "refiner" {
				opts.EntryClass = "com.skt.apm.refinement.ParameterRefiner"
				opts.ProgramArg = []string{"--job-name", nestID + "-Refiner",
					"--specification-url", specURL,
					"--program-identifier", "parameter_refiner"}
				opts.Parallelism = 1
			} else if args[0] == "alarm" {
				opts.EntryClass = "com.skt.apm.alarm.asset.AssetAlarm"
				opts.ProgramArg = []string{"--job-name", nestID + "-AssetAlarm",
					"--specification-url", specURL,
					"--program-identifier", "asset"}
				opts.Parallelism = 1
			} else if args[0] == "paramalarm" {
				opts.EntryClass = "com.skt.apm.alarm.parameter.ParameterAlarm"
				opts.ProgramArg = []string{"--job-name", nestID + "-ParamAlarm",
					"--specification-url", specURL,
					"--program-identifier", "parameter"}
				opts.Parallelism = 1
			} else if args[0] == "fdc" {
				opts.EntryClass = ""
				opts.ProgramArg = []string{"--job-name", nestID + "-fdc",
					"--specification-url", specURL,
					"--program-identifier", "fdc",
					"--local-repository-location", "/var/tmp/flink/fd"}
				opts.Parallelism = 1
			} else if args[0] == "spc" {
				opts.EntryClass = ""
				opts.ProgramArg = []string{"--job-name", nestID + "-spc",
					"--specification-url", specURL,
					"--program-identifier", "spc",
					"--reference-period", "120000",
					"--local-repository-location", "/var/tmp/flink/spc"}
				opts.Parallelism = 1
			} else if args[0] == "mva" {
				opts.EntryClass = ""
				opts.ProgramArg = []string{"--job-name", nestID + "-mva",
					"--specification-url", specURL,
					"--program-identifier", "mva",
					"--local-repository-location", "/var/tmp/flink/mva"}
				opts.Parallelism = 1
			} else if args[0] == "bae" {
				opts.EntryClass = ""
				opts.ProgramArg = []string{"--job-name", nestID + "-bae",
					"--specification-url", specURL,
					"--program-identifier", "unsupervised",
					"--local-repository-location", "/var/tmp/flink/bae"}
				opts.Parallelism = 1
			} else if args[0] == "current" {
				opts.EntryClass = ""
				opts.ProgramArg = []string{"--job-name", nestID + "-current",
					"--specification-url", specURL,
					"--program-identifier", "current",
					"--local-repository-location", "/var/tmp/flink/current"}
				opts.Parallelism = 1
			} else {
				logger.Println("Warning: There are no valid arguments.")
				return
			}

			r, err := c.RunJar(opts)
			if err != nil {
				panic(err)
			}
			logger.Println("Running Job: " + r.ID)
		}
	},
}

func init() {
	createCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
