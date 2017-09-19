// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tim15/wiz/api/daemon"
	"github.com/tim15/wiz/cli/client"
	"golang.org/x/net/context"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version and exits",
	Long:  `Prints wiz CLI and daemon version`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.GetClient()
		if err != nil {
			panic(err)
		}
		fmt.Println("good")
		fmt.Println(c.GetVersion(context.Background(), &daemon.Empty{}))
		// Contact the server and print out its response.
		// name := defaultName
		// if len(os.Args) > 1 {
		//   name = os.Args[1]
		// }
		// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		// if err != nil {
		//   log.Fatalf("could not greet: %v", err)
		// }
		// log.Printf("Greeting: %s", r.Message)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
