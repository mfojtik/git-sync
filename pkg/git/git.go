package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mfojtik/reposync/pkg/types"
)

func Checkout(r types.Repository, branch string) error {
	return runGit(r, "checkout", branch)
}

func Fetch(r types.Repository, remote string) error {
	return runGit(r, "fetch", remote, "--quiet")
}

func FetchTags(r types.Repository, remote string) error {
	return runGit(r, "fetch", remote, "--quiet", "--tags")
}

func Merge(r types.Repository, branch string) error {
	return runGit(r, "merge", branch, "--quiet")
}

func CurrentBranch(r types.Repository) (string, error) {
	return getGitOutput(r, "symbolic-ref", "--short", "HEAD")
}

func Cleanup(r types.Repository) error {
	if err := runGit(r, "remote", "prune", "origin"); err != nil {
		return err
	}
	if err := runGit(r, "gc", "--auto"); err != nil {
		return err
	}
	return nil
}

func Push(r types.Repository) error {
	return runGit(r, "push", "origin", "master", "--tags", "--force")
}

func getGitOutput(r types.Repository, options ...string) (string, error) {
	cmd, out := cmdForRepo(r, options...)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command %q failed: \n%s\n(%v)\n", strings.Join(cmd.Args, " "), out.String(), err)
	}
	return strings.TrimSpace(out.String()), nil
}

func runGit(r types.Repository, options ...string) error {
	cmd, out := cmdForRepo(r, options...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: \n%s\n(%v)\n", strings.Join(cmd.Args, " "), out.String(), err)
	}
	return nil
}

func cmdForRepo(r types.Repository, options ...string) (*exec.Cmd, bytes.Buffer) {
	cmd := exec.Command("git", options...)
	cmd.Dir = r.BaseDirectory
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return cmd, out
}
