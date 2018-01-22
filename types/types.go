package types

// Configuration task configuration.
type Configuration struct {
	Owner                  string `short:"o" description:"Repository owner. [required]"`
	RepositoryName         string `short:"r" long:"repo-name" description:"Repository name. [required]"`
	Debug                  bool   `long:"debug" description:"Debug mode."`
	DockerfileURL          string `short:"d" long:"dockerfile-url" description:"Dockerfile URL. [required]"`
	ExperimentalBranchName string `long:"exp-branch" description:"Build a branch as experimental."`
	DockerImageName        string `long:"image-name" description:"Docker image name."`
	MenuTemplateURL        string `short:"m" long:"menu-url" description:"URL of the template of the JS file use for the multi version menu."`
	MenuTemplateFile       string `long:"menu-file" description:"File path of the template of the JS file use for the multi version menu."`
}

// VersionsInformation versions information
type VersionsInformation struct {
	Current      string
	Latest       string
	Experimental string
	CurrentPath  string
}

// RepoID Repository identifier.
type RepoID struct {
	Owner          string
	RepositoryName string
}

// NoOption empty struct.
type NoOption struct{}
