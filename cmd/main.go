package main

import (
	"flag"
	"fmt"
	"log"

	"github-workflow-explorer/internal"
)

func main() {
	// comand line arguments
	var configFilepath string
	var authOverwrite string

	flag.StringVar(&configFilepath, "config", "config.json", "Path to configuration json file")
	flag.StringVar(&configFilepath, "c", "config.json", "")
	flag.StringVar(&authOverwrite, "auth", "", "Token used to authenticate to GitHub, overwrites json configuration")
	flag.StringVar(&authOverwrite, "a", "", "")
	flag.Usage = func() {
		fmt.Print(`Usage of shipper:
	  -c, --config  path to configuation yaml
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
	api := internal.GithubApi{
		Owner:   configs.Owner,
		Repo:    configs.Name,
		Auth:    fmt.Sprintf("Bearer %s", configs.Auth),
		Version: configs.GithubApiVersion,
	}
	runs, err := api.ListWorkflows()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range runs {
		fmt.Println(v.ID, v.Name, v.Title)
	}
	fmt.Println(len(runs))

	//err = api.DeleteAPI(7834160683)
	//err = api.DeleteAPI(3)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
