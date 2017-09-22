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

	"errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"io"

	"bufio"
	"github.com/tim15/wiz/api/pkg"
	"github.com/tim15/wiz/cli/backends/config"
	"os"
	"strings"
)

func check(e error, reason string) {
	if e != nil {
		fmt.Println(reason + e.Error())
	}
}

func backend(filename string) (string, error) {
	for _, backend := range config.GetBackends() {
		if strings.HasSuffix(filename, "."+backend.FileSuffix) {
			fmt.Println("Config backend: " + backend.Name)
			return backend.Name, nil
		}
	}
	return "", errors.New("No backends configured")
}

func e(err string) {
	fmt.Println("Error: " + err)
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [FILES ...]",
	Short: "Validate a spec file",
	Long: `Examples:
  wiz check ./wiz.json
  wiz check
  `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")
		fmt.Println(args)

		if len(args) > 0 {
			for _, filename := range args {
				// filename := args[0]

				if _, err := os.Stat(filename); os.IsNotExist(err) {
					e(filename + " does not exist")
					continue
				}

				// Wiz supports multiple config backends, and multiple spec file types.
				// Determine the backend to use based on filename
				b, err := backend(filename)
				if err != nil {
					e("No config backend available for " + filename)
					continue
				}
				fmt.Println(b)

				f, err := os.Open(args[0])
				defer f.Close()
				check(err, "")
				pb := pkg.Package{Name: "default"}
				// fmt.Printf("%+v\n", pb)
				fmt.Printf("package: %+v\n", pb)
				err = jsonpb.Unmarshal(f, &pb)
				check(err, "Invalid spec file: ")
				marshaller := jsonpb.Marshaler{}
				str, err := marshaller.MarshalToString(&pb)
				fmt.Printf("Package spec: %+v\n", pb)
				fmt.Println("Hello", pb.Type)
				for key, value := range pb.Dependencies {
					fmt.Println("Key:", key, "Value:", value)
				}
				fmt.Printf("Output: %s\n", str)
				fmt.Printf("package: %+v\n", pb)
				str = proto.MarshalTextString(&pb)
				fmt.Println(str)
				npb := &pkg.Package{}
				proto.UnmarshalText(str, npb)
				fmt.Printf("package: %+v\n", pb)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
