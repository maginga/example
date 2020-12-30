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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Registers a Job to the flink cluster, and runs a Job to be registered stream job.",
	Long: `Registers a Job to the flink cluster, and runs a Job to be registered stream job.
For example: apm create stream [refiner, alarm, paramalarm, fdc, spc, mva, bae, current] [nest ID] [jar file]
`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		jobServer := fmt.Sprintf("%v", viper.Get("stream.jobServer"))
		log.Println("Job Server: " + jobServer)

		jarFile := args[2]

		// Your flink server HTTP API
		c, err := client.New(jobServer)
		if err != nil {
			panic(err)
		}

		d, err := c.Config()
		log.Println("Flink Version: " + d.FlinkVersion)

		k, err := c.Jars()
		existed := false
		jarID := ""
		for _, f := range k.Files {
			if f.Name == jarFile {
				existed = true
				jarID = f.ID
				log.Println("Existed Jars ID: " + jarID)
				break
			}
		}

		if !existed {
			u, err := c.UploadJar(jarFile)
			if err != nil {
				panic(err)
			}
			log.Println("file location: " + u.FileName)
			log.Println("This Jar was uploaded.")
			ids := strings.Split(u.FileName, "/")
			jarID = ids[len(ids)-1]
			log.Println("New Jars ID: " + jarID)
		}

		if len(jarID) > 0 {
			apmAddr := fmt.Sprintf("%v", viper.Get("grandview.url"))
			nestID := args[1]
			specURL := "http://" + apmAddr + "/api/spec/" + nestID

			opts := client.RunOpts{}
			opts.JarID = jarID

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
				opts.EntryClass = "com.bistel.apm.model.uv.oos.FDModel"
				opts.ProgramArg = []string{"--job-name", nestID + "-fdc",
					"--specification-url", specURL,
					"--program-identifier", "fdc",
					"--local-repository-location", "/var/tmp/flink/" + nestID + "/fd"}
				opts.Parallelism = 1
			} else if args[0] == "spc" {
				opts.EntryClass = "com.bistel.apm.model.uv.spc.SPCRulesModel"
				opts.ProgramArg = []string{"--job-name", nestID + "-spc",
					"--specification-url", specURL,
					"--program-identifier", "spc",
					"--reference-period", "120000",
					"--local-repository-location", "/var/tmp/flink/" + nestID + "/spc"}
				opts.Parallelism = 1
			} else if args[0] == "mva" {
				opts.EntryClass = "com.bistel.apm.model.mv.statistical.MVAModel"
				opts.ProgramArg = []string{"--job-name", nestID + "-mva",
					"--specification-url", specURL,
					"--program-identifier", "mva",
					"--local-repository-location", "/var/tmp/flink/" + nestID + "/mva"}
				opts.Parallelism = 1
			} else if args[0] == "bae" {
				opts.EntryClass = "com.bistel.apm.model.mv.unsupervised.AutoEncoderModel"
				opts.ProgramArg = []string{"--job-name", nestID + "-bae",
					"--specification-url", specURL,
					"--program-identifier", "unsupervised",
					"--local-repository-location", "/var/tmp/flink/" + nestID + "/bae"}
				opts.Parallelism = 1
			} else if args[0] == "current" {
				opts.EntryClass = "com.bistel.apm.model.domain.CurrentImbalanceModel"
				opts.ProgramArg = []string{"--job-name", nestID + "-current",
					"--specification-url", specURL,
					"--program-identifier", "current",
					"--local-repository-location", "/var/tmp/flink/" + nestID + "/current"}
				opts.Parallelism = 1
			} else {
				log.Println("Warning: There are no valid arguments.")
				return
			}

			r, err := c.RunJar(opts)
			if err != nil {
				panic(err)
			}
			log.Println("Running Job: " + r.ID)
		}
	},
}

func init() {
	tenantCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
