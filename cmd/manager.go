/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/alexkreidler/wiz/manager/local"
	"github.com/mitchellh/go-homedir"
	"log"

	"github.com/spf13/cobra"
)

// managerCmd represents the manager command
var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Utilities to manage the local Wiz Manager",
}

var resetCmd = &cobra.Command{Use: "reset", Short: "Reset the manager's state", Run: func(cmd *cobra.Command, args []string) {
	//TODO: maybe use viper for this
	file, err := homedir.Expand(cmd.Flag("manager").Value.String())
	if err != nil {
		log.Fatal(err)
	}

	m := local.NewManager(local.Options{StorageLocation: file})
	m.ResetState()
}}

func init() {
	managerCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(managerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// managerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// managerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
