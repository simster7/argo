package main

import (
	"github.com/spf13/cobra/doc"

	"github.com/simster7/argo/v2/cmd/argo/commands"
)

func main() {
	println("generating docs/cli")
	cmd := commands.NewCommand()
	cmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTree(cmd, "docs/cli")
	if err != nil {
		panic(err)
	}
}
