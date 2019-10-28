package main

import (
	"fmt"
	"log"
	"os"

	"github.com/containous/structor/core"
	"github.com/containous/structor/types"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	defaultDockerImageName = "doc-site"
	defaultDockerfileName  = "docs.Dockerfile"
)

func main() {
	cfg := &types.Configuration{
		DockerImageName: defaultDockerImageName,
		DockerfileName:  defaultDockerfileName,
		NoCache:         false,
		Menu:            &types.MenuFiles{},
	}

	rootCmd := &cobra.Command{
		Use:     "structor",
		Short:   "Messor Structor: Manage multiple documentation versions with Mkdocs.",
		Long:    `Messor Structor: Manage multiple documentation versions with Mkdocs.`,
		Version: version,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if cfg.Debug {
				log.Printf("Run Structor command with config : %+v", cfg)
			}

			if len(cfg.DockerImageName) == 0 {
				log.Printf("'image-name' is undefined, fallback to %s.", defaultDockerImageName)
				cfg.DockerImageName = defaultDockerImageName
			}

			return validateConfig(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return core.Execute(cfg)
		},
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&cfg.Owner, "owner", "o", "", "Repository owner. [required]")
	flags.StringVarP(&cfg.RepositoryName, "repo-name", "r", "", "Repository name. [required]")

	flags.BoolVar(&cfg.Debug, "debug", false, "Debug mode.")

	flags.StringVarP(&cfg.DockerfileURL, "dockerfile-url", "d", "", "Use this Dockerfile when --dockerfile-name is not found. Can be a file path. [required]")
	flags.StringVar(&cfg.DockerfileURL, "dockerfile-name", defaultDockerfileName, "Search and use this Dockerfile in the repository (in './docs/' or in './') for building documentation.")
	flags.StringVar(&cfg.DockerImageName, "image-name", defaultDockerImageName, "Docker image name.")
	flags.BoolVar(&cfg.NoCache, "no-cache", false, "Set to 'true' to disable the Docker build cache.")

	flags.StringVar(&cfg.ExperimentalBranchName, "exp-branch", "", "Build a branch as experimental.")
	flags.StringSliceVar(&cfg.ExcludedBranches, "exclude", nil, "Exclude branches from the documentation generation.")

	flags.BoolVar(&cfg.ForceEditionURI, "force-edit-url", false, "Add a dedicated edition URL for each version.")
	flags.StringVar(&cfg.RequirementsURL, "rqts-url", "", "Use this requirements.txt to merge with the current requirements.txt. Can be a file path.")

	flags.StringVar(&cfg.Menu.JsURL, "menu.js-url", "", "URL of the template of the JS file use for the multi version menu.")
	flags.StringVar(&cfg.Menu.JsFile, "menu.js-file", "", "File path of the template of the JS file use for the multi version menu.")
	flags.StringVar(&cfg.Menu.CSSURL, "menu.css-url", "", "URL of the template of the CSS file use for the multi version menu.")
	flags.StringVar(&cfg.Menu.CSSFile, "menu.css-file", "", "File path of the template of the CSS file use for the multi version menu.")

	docCmd := &cobra.Command{
		Use:    "doc",
		Short:  "Generate documentation",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, "./docs")
		},
	}

	rootCmd.AddCommand(docCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Run: func(_ *cobra.Command, _ []string) {
			displayVersion(rootCmd.Name())
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		return fmt.Errorf("%s is mandatory", fieldName)
	}
	return nil
}
