/*
Copyright © 2019 Alex Kreidler

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
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var outfile string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get PACKAGE",
	Short: "Gets a data package by running its pipeline",
	Long: `Fetches data packages and performs processing needed.
	
Requires an output.yaml file to actually output anything, else the DAG won't be executed.
Only works for "type: data" packages.
	
PACKAGE refers to the Wiz package to fetch: either be a package identifier from an existing registry or the path to a local YAML manifest file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("get called with", args, " outfile:", outfile)
		packageName := args[0]
		fmt.Println("Installing ", packageName)
		f, err := ioutil.ReadFile(packageName)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(f)

	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringVarP(&outfile, "output", "o", "", "The output specification file (required)")
	//getCmd.MarkFlagRequired("output")
}
