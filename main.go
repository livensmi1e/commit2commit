package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

func main() {
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

func help() string {
	return `Commands:
  next, n              move to the next commit
  prev, p              move to the previous commit
  goto, g <target>     navigate to a commit hash or index
  tree                 show commit history
  show                 show current commit
  quit, exit, q        exit`
}
