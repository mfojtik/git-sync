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
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	reposync "github.com/mfojtik/git-sync/pkg/sync"
	"github.com/mfojtik/git-sync/pkg/types"
	"github.com/sethgrid/multibar"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run synchronization",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := types.GetUserConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to read configuration %q: %v\n", types.ConfigDefaultLocation, err)
			return
		}
		bar, _ := multibar.New()
		wg := &sync.WaitGroup{}
		wg.Add(len(config.Repositories))
		type repoBarType struct {
			ticker     chan<- int
			repository types.Repository
			bar        *pb.ProgressBar
			finish     chan struct{}
		}
		repoProgressBars := []repoBarType{}
		pool, _ := pb.StartPool()
		// Create progress bars for all repositories first
		for _, repo := range config.Repositories {
			go func(r types.Repository) {
				syncProgresChan := make(chan int)
				finish := make(chan struct{})
				bar := pb.New(100).Prefix(repoName(r) + " ")
				bar.ShowCounters = false
				pool.Add(bar)
				repoProgressBars = append(repoProgressBars, repoBarType{
					ticker:     syncProgresChan,
					repository: r,
					bar:        bar,
					finish:     finish,
				})
				wg.Done()
				for {
					select {
					case p := <-syncProgresChan:
						if p > 0 {
							bar.Set(p)
						}
					case <-finish:
						bar.Finish()
						break
					}
				}
			}(repo)
		}
		wg.Wait()

		errors := []error{}

		wg.Add(len(repoProgressBars))
		go bar.Listen()
		for i, repo := range repoProgressBars {
			go func(r repoBarType, repoIndex int) {
				defer func() {
					close(r.finish)
					wg.Done()
				}()
				if err := reposync.Repository(r.repository, r.ticker); err != nil {
					errors = append(errors, fmt.Errorf("ERROR: %q: %s", repoName(r.repository), strings.TrimSpace(err.Error())))
					return
				}
				config.Repositories[repoIndex].LastSync = time.Now()
			}(repo, i)
		}

		wg.Wait()
		pool.Stop()
		types.SaveUserConfig(config)
		if len(errors) > 0 {
			for _, e := range errors {
				fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			}
			os.Exit(1)
		}
	},
}

func repoName(r types.Repository) string {
	parts := strings.Split(r.BaseDirectory, "/")
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

func init() {
	RootCmd.AddCommand(runCmd)
}
