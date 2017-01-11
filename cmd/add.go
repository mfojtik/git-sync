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
	"os"

	"github.com/mfojtik/reposync/pkg/types"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [REPO1, REPO2...]",
	Short: "Add GIT repository to synchronize",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := types.GetUserConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to read configuration %q: %v\n", types.ConfigDefaultLocation, err)
			return
		}
		needSave := false
		for _, name := range args {
			exists := false
			for _, repo := range config.Repositories {
				if repo.BaseDirectory == name {
					exists = true
					break
				}
			}
			if exists {
				fmt.Fprintf(os.Stderr, "%q: directory already registered as repository (skipping)\n", name)
				continue
			}
			config.Repositories = append(config.Repositories, types.Repository{BaseDirectory: name})
			needSave = true
			fmt.Fprintf(os.Stdout, "%q: succesfully added\n", name)
		}
		if needSave {
			if err := types.SaveUserConfig(config); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: unable to write configuration %q: %v\n", types.ConfigDefaultLocation, err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
