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
	goTo(commits[0])
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
		next()
	case "prev", "p":
		prev()
	case "tree":
		fmt.Printf("run command: %s with no arg\n", cmd)
	case "show":
		fmt.Println("show")
	case "goto", "g":
		for _, commit := range commits {
			if args[0] == commit {
				goTo(args[0])
				return
			}
		}
		fmt.Println("commit not valid")
	case "help":
		help()
	case "quit", "exit", "q":
		restore()
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

func help() {
	cmd := `Commands:
  next, n              move to the next commit
  prev, p              move to the previous commit
  goto, g <target>     navigate to a commit hash or index
  tree                 show commit history
  show                 show current commit
  quit, exit, q        exit`
	fmt.Println(cmd)
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
	commits := strings.Fields(out)
	if len(commits) == 0 {
		fmt.Println("repository does not have any commits")
		os.Exit(1)
	}
	return commits
}

func goTo(commit string) {
	runGit("switch", "--detach", commit)
}

func restore() {
	if originalBranch != "" {
		runGit("switch", originalBranch)
	} else {
		runGit("switch", "--detach", originalCommit)
	}
}

func next() {
	if curr <= len(commits)-2 {
		curr++
		goTo(commits[curr])
	} else {
		fmt.Println("commit reach the end")
	}
}

func prev() {
	if curr >= 1 {
		curr--
		goTo(commits[curr])
	} else {
		fmt.Println("commit reach the start")
	}
}
