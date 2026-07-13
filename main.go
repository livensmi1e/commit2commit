package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

var (
	repoDir        string   // repo root directory
	originalBranch string   // HEAD can be restored and point to original branch as is
	originalCommit string   // original branch can be restored and point to original commit as is
	commits        []string // short hash commit history to navigate
	curr           int      // current commit index 0 as first commit
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: commit2commit [repo-directory-path]")
		os.Exit(1)
	}
	repoDir = findRepoRootDir(os.Args...)
	originalBranch, originalCommit = findOriginals()
	commits = findCommitHistoryFromStart()
	debug()
	p := prompt.New(executor, completer, prompt.OptionPrefix("commit2commit> "), prompt.OptionTitle("commit2commit"))
	p.Run()
}

func executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	parts := strings.Fields(s)
	cmd := parts[0]
	args := parts[1:]
	switch cmd {
	case "next", "n":
	case "prev", "p":
	case "tree":
		fmt.Printf("run command: %s with no arg\n", cmd)
	case "show":
		fmt.Println("show")
	case "goto", "g":
		fmt.Printf("run command: %s with with arg %s\n", cmd, args[0])
	case "help":
		fmt.Println(help())
	case "quit", "exit", "q":
		fmt.Println("bye!")
		os.Exit(0)
	default:
		fmt.Printf("unknown command: %s\n", cmd)
	}
}

func completer(document prompt.Document) []prompt.Suggest {
	return nil
}

func debug() {
	fmt.Printf("repoDir: %s\n", repoDir)
	fmt.Printf("original branch: %s\n", originalBranch)
	fmt.Printf("original commit: %s\n", originalCommit)
	fmt.Printf("current index: %d\n", curr)
	fmt.Printf("commits: %v\n", commits)
}

func help() string {
	return `Commands:
  next, n              move to the next commit
  prev, p              move to the previous commit
  goto, g <target>     navigate to a commit hash or index
  tree                 show commit history
  show                 show current commit
  quit, exit, q        exit`
}

// shared function to run git command
func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf(
			"git command failed: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}
	return strings.TrimSpace(string(output)), nil
}

func findRepoRootDir(args ...string) string {
	dir := "."
	if len(args) == 2 {
		dir = args[1]
	}
	dir, _ = runGit("-C", dir, "rev-parse", "--show-toplevel")
	return dir
}

func findOriginals() (branch, commit string) {
	branch, _ = runGit("branch", "--show-current")
	commit, _ = runGit("rev-parse", "--short", "HEAD")
	return
}

func findCommitHistoryFromStart() []string {
	out, _ := runGit("rev-list", "--reverse", "--first-parent", "--abbrev-commit", "HEAD")
	return strings.Fields(out)
}
