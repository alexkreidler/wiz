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
	"github.com/alexkreidler/wiz/api/packages"
	"github.com/alexkreidler/wiz/manager/local"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// packageCmd represents the pipeline command
var packageCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage your pipeline resources",
}

var installCmd = &cobra.Command{
	Use:   "install [PACKAGE]",
	Short: "Installs a package",
	Long:  "Right now, we only support local package spec files. In the future, it will be able to fetch from a registry. We only install data packages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		specFile := args[0]
		f, err := ioutil.ReadFile(specFile)
		if err != nil {
			log.Fatal("failed to read spec file", err)
		}
		p := packages.Package{}
		err = yaml.Unmarshal(f, p)
		if err != nil {
			log.Fatal("invalid spec file", err)
		}

		pipeline := p.Source

		pipeline.Spec.Sequential[0].Processor.Type = "input"
		//spew.Dump(p)
		file, err := homedir.Expand(cmd.Flag("manager").Value.String())
		if err != nil {
			log.Fatal(err)
		}
		m := local.NewManager(local.Options{StorageLocation: file, RestartExecutor: restart, PreserveRunIDs: debug, OverwritePipelines: debug})
		err = m.CreatePipeline(pipeline, "local")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	installCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Determines whether to restart the Wiz Executor each time")
	
	installCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enables debug features: don't overwrite RunIDs, allows the manager to overwrite existing pipelines")

	packageCmd.AddCommand(installCmd)

	rootCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
