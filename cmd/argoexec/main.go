package main

import (
	"fmt"
	"os"

	"github.com/simster7/argo/v2/cmd/argoexec/commands"
	// load authentication plugin for obtaining credentials from cloud providers.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
