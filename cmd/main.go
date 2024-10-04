package main

import (
	"flag"
	"fmt"
	"log"

	"github-workflow-explorer/internal"
	"github-workflow-explorer/internal/ui"
)

func main() {
	// comand line arguments
	var configFilepath string
	var authOverwrite string
	var ownerOverwrite string
	var repoOverwrite string

	flag.StringVar(&configFilepath, "config", "config.json", "Path to configuration json file")
	flag.StringVar(&configFilepath, "c", "config.json", "")
	flag.StringVar(&authOverwrite, "auth", "", "Token used to authenticate to GitHub, overwrites json configuration")
	flag.StringVar(&authOverwrite, "a", "", "")
	flag.StringVar(&repoOverwrite, "repo", "", "Repo name to view, overwrites json configuration")
	flag.StringVar(&repoOverwrite, "r", "", "")
	flag.StringVar(&ownerOverwrite, "owner", "", "Owner of repo, overwrites json configuration")
	flag.StringVar(&ownerOverwrite, "o", "", "")
	flag.Usage = func() {
		fmt.Print(`Usage of shipper: gh-we [-c | --config <path>] [-a | --auth <string>] [-r | --repo <path>]
	-c, --config  path to configuation json file
	-a, --auth    github Token used for authentication. Overwrites json configurations.
	-r, --repo    github repo to view. Overwrites json configurations.
	-o, --owner   owner of github repo. Overwrites json configurations.
	-h, --help    display this help message
`)
	}
	flag.Parse()

	// read script configuration file
	configs, err := internal.LoadConfigurations(configFilepath)
	if err != nil {
		log.Panic(err)
	}
	if authOverwrite != "" {
		configs.Auth = authOverwrite
	}
	if repoOverwrite != "" {
		configs.Name = repoOverwrite
	}
	if ownerOverwrite != "" {
		configs.Owner = ownerOverwrite
	}
	api := internal.GithubApi{
		Owner:   configs.Owner,
		Repo:    configs.Name,
		Auth:    fmt.Sprintf("Bearer %s", configs.Auth),
		Version: configs.GithubApiVersion,
	}

	appUI := ui.NewAppUI(api)

	// start graphical interface
	err = appUI.StartUI()
	if err != nil {
		log.Fatal(err)
	}
}
