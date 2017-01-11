package sync

import (
	"github.com/mfojtik/reposync/pkg/git"
	"github.com/mfojtik/reposync/pkg/types"
)

func Repository(r types.Repository, progress chan<- int) error {
	defer func() {
		close(progress)
	}()

	// Basic synchronization:
	currentBranch, err := git.CurrentBranch(r)
	if err != nil {
		return err
	}

	progress <- 10
	if currentBranch != "master" {
		if err := git.Checkout(r, "master"); err != nil {
			return err
		}
	}

	progress <- 20
	if err := git.Fetch(r, "upstream"); err != nil {
		return err
	}

	progress <- 30
	if err := git.FetchTags(r, "upstream"); err != nil {
		return err
	}

	progress <- 50
	if err := git.Merge(r, "upstream/master"); err != nil {
		return err
	}

	progress <- 60
	if err := git.Push(r); err != nil {
		return err
	}

	progress <- 80
	if currentBranch != "master" {
		return git.Checkout(r, currentBranch)
	}

	progress <- 100
	return nil
}
