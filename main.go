package main

import (
	"os"

	"brewjenkins/go-ai-insight-advisor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
