package main

import (
	"fmt"
	"os"

	"github.com/Capsule7446/oh-my-markdown/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
		os.Exit(1)
	}
}
