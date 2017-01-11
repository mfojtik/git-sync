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
	"text/tabwriter"

	humanize "github.com/dustin/go-humanize"
	"github.com/mfojtik/git-sync/pkg/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List GIT repositories that we track",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := types.GetUserConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to read configuration %q: %v\n", types.ConfigDefaultLocation, err)
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
		for _, r := range config.Repositories {
			fmt.Fprintf(w, "%s\t%s\t\n", repoName(r), "last synced "+humanize.Time(r.LastSync))
		}
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
