package main

import (
	"log"
	"os"

	"github.com/containous/flaeg"
	"github.com/containous/structor/core"
	"github.com/containous/structor/meta"
	"github.com/containous/structor/types"
	"github.com/ogier/pflag"
)

const (
	defaultDockerImageName = "doc-site"
)

func main() {
	config := &types.Configuration{
		DockerImageName: defaultDockerImageName,
	}

	defaultPointersConfig := &types.Configuration{
		Menu: &types.MenuFiles{},
	}
	rootCmd := &flaeg.Command{
		Name:                  "structor",
		Description:           `Messor Structor: Manage multiple documentation versions with Mkdocs.`,
		DefaultPointersConfig: defaultPointersConfig,
		Config:                config,
		Run:                   runCommand(config),
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	// version
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           "Display the version.",
		Config:                &types.NoOption{},
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			meta.DisplayVersion()
			return nil
		},
	}

	flag.AddCommand(versionCmd)

	err := flag.Run()
	if err != nil && err != pflag.ErrHelp {
		log.Printf("Error: %v\n", err)
	}
}

func runCommand(config *types.Configuration) func() error {
	return func() error {
		if config.Debug {
			log.Printf("Run Structor command with config : %+v\n", config)
		}

		if len(config.DockerImageName) == 0 {
			log.Printf("'image-name' is undefined, fallback to %s.\n", defaultDockerImageName)
			config.DockerImageName = defaultDockerImageName
		}

		err := validateConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		err = core.Execute(config)
		if err != nil {
			log.Fatalf("Execute error: %+v", err)
		}
		return nil
	}
}

func validateConfig(config *types.Configuration) error {
	err := required(config.DockerfileURL, "dockerfile-url")
	if err != nil {
		return err
	}
	err = required(config.Owner, "owner")
	if err != nil {
		return err
	}
	return required(config.RepositoryName, "repo-name")
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
