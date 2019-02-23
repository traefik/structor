package core

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/containous/structor/docker"
	"github.com/containous/structor/file"
	"github.com/containous/structor/gh"
	"github.com/containous/structor/manifest"
	"github.com/containous/structor/menu"
	"github.com/containous/structor/repository"
	"github.com/containous/structor/requirements"
	"github.com/containous/structor/types"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/worktree"
	"github.com/pkg/errors"
)

const (
	baseRemote      = "origin/"
	envVarLatestTag = "STRUCTOR_LATEST_TAG"
)

// Execute core process
func Execute(config *types.Configuration) error {
	workDir, err := ioutil.TempDir("", "structor")
	if err != nil {
		return err
	}

	defer func() {
		if err = cleanAll(workDir, config.Debug); err != nil {
			log.Println("Error during cleaning: ", err)
		}
	}()

	if config.Debug {
		log.Printf("Temp directory: %s", workDir)
	}

	return process(workDir, config)
}

func process(workDir string, config *types.Configuration) error {
	menuContent := menu.GetTemplateContent(config.Menu)

	fallbackDockerfile, err := docker.GetDockerfileFallback(config.DockerfileURL, config.DockerImageName)
	if err != nil {
		return err
	}

	requirementsContent, err := requirements.GetContent(config.RequirementsURL)
	if err != nil {
		return err
	}

	latestTagName, err := getLatestReleaseTagName(config.Owner, config.RepositoryName)
	if err != nil {
		return err
	}

	log.Printf("Latest tag: %s", latestTagName)

	branches, err := getBranches(config.ExperimentalBranchName, config.Debug)
	if err != nil {
		return err
	}

	siteDir, err := createSiteDirectory()
	if err != nil {
		return err
	}

	for _, branchRef := range branches {
		versionName := strings.Replace(branchRef, baseRemote, "", 1)
		versionCurrentPath := filepath.Join(workDir, versionName)

		log.Printf("Generating doc for version %s", versionName)

		err := repository.CreateWorkTree(versionCurrentPath, branchRef, config.Debug)
		if err != nil {
			return err
		}

		versionDocsRoot, err := getDocumentationRoot(versionCurrentPath)
		if err != nil {
			return err
		}

		err = requirements.Check(versionDocsRoot)
		if err != nil {
			return err
		}

		versionsInfo := types.VersionsInformation{
			Current:      versionName,
			Latest:       latestTagName,
			Experimental: config.ExperimentalBranchName,
			CurrentPath:  versionDocsRoot,
		}
		fallbackDockerfile.Path = filepath.Join(versionsInfo.CurrentPath, fallbackDockerfile.Name)

		err = buildDocumentation(branches, versionsInfo, fallbackDockerfile, menuContent, requirementsContent, config)
		if err != nil {
			return err
		}

		err = copyVersionSiteToOutputSite(versionsInfo, siteDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// getDocumentationRoot returns the path to the documentation's root by searching for "${menu.ManifestFileName}".
// Search is done from the docsRootSearchPath, relatively to the provided repository path.
// An additional sanity checking is done on the file named "requirements.txt" which must be located in the same directory.
func getDocumentationRoot(repositoryRoot string) (string, error) {
	var docsRootSearchPaths = []string{"/", "docs/"}

	for _, docsRootSearchPath := range docsRootSearchPaths {
		candidateDocsRootPath := filepath.Join(repositoryRoot, docsRootSearchPath)

		if _, err := os.Stat(filepath.Join(candidateDocsRootPath, manifest.FileName)); !os.IsNotExist(err) {
			log.Printf("Found %s for building documentation in %s.", manifest.FileName, candidateDocsRootPath)
			return candidateDocsRootPath, nil
		}
	}

	return "", errors.Errorf("no file %s found in %s (search path was: %s)", manifest.FileName, repositoryRoot, strings.Join(docsRootSearchPaths, ","))
}

// copyVersionSiteToOutputSite adds the generated documentation for the version described in ${versionsInfo} to the output directory.
// If the current version (branch) name is related to the latest tag, then it's copied at the root of the output directory.
// Else it is copied under a directory named after the version, at the root of the output directory.
func copyVersionSiteToOutputSite(versionsInfo types.VersionsInformation, siteDir string) error {
	currentSiteDir, err := getDocumentationRoot(versionsInfo.CurrentPath)
	if err != nil {
		return err
	}

	outputDir := siteDir
	if !strings.HasPrefix(versionsInfo.Latest, versionsInfo.Current) {
		outputDir = filepath.Join(siteDir, versionsInfo.Current)
	}

	return file.Copy(filepath.Join(currentSiteDir, "site"), outputDir)
}

func buildDocumentation(branches []string, versionsInfo types.VersionsInformation,
	fallbackDockerfile docker.DockerfileInformation, menuTemplateContent menu.Content, requirementsContent []byte,
	config *types.Configuration) error {

	err := menu.Build(versionsInfo, branches, menuTemplateContent)
	if err != nil {
		return err
	}

	err = requirements.Build(versionsInfo, requirementsContent)
	if err != nil {
		return err
	}

	baseDockerfile, err := docker.GetDockerfile(versionsInfo.CurrentPath, fallbackDockerfile, config.DockerfileName)
	if err != nil {
		return err
	}

	dockerImageFullName, err := baseDockerfile.BuildImage(versionsInfo, config.NoCache, config.Debug)
	if err != nil {
		return err
	}

	// Run image
	output, err := docker.Exec(config.Debug, "run", "--rm", "-v", versionsInfo.CurrentPath+":/mkdocs", dockerImageFullName, "mkdocs", "build")
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}

func getLatestReleaseTagName(owner, repositoryName string) (string, error) {
	latest := os.Getenv(envVarLatestTag)
	if len(latest) > 0 {
		return latest, nil
	}

	return gh.GetLatestReleaseTagName(owner, repositoryName)
}

func getBranches(experimentalBranchName string, debug bool) ([]string, error) {
	var branches []string

	if len(experimentalBranchName) > 0 {
		branches = append(branches, baseRemote+experimentalBranchName)
	}

	gitBranches, err := repository.ListBranches(debug)
	if err != nil {
		return nil, err
	}
	branches = append(branches, gitBranches...)

	if len(branches) == 0 {
		log.Println("[WARN] no branch.")
	}

	return branches, nil
}

func createSiteDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	siteDir := filepath.Join(currentDir, "site")
	err = createDirectory(siteDir)
	if err != nil {
		return "", err
	}

	return siteDir, nil
}

func createDirectory(directoryPath string) error {
	if _, err := os.Stat(directoryPath); !os.IsNotExist(err) {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return err
		}
	}

	return os.MkdirAll(directoryPath, os.ModePerm)
}

func cleanAll(workDir string, debug bool) error {
	err := os.RemoveAll(workDir)
	if err != nil {
		return err
	}

	output, err := git.Worktree(worktree.Prune, git.Debugger(debug))
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}
