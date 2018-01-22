package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/containous/structor/copy"
	"github.com/containous/structor/gh"
	"github.com/containous/structor/menu"
	"github.com/containous/structor/repository"
	"github.com/containous/structor/types"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/worktree"
	"github.com/pkg/errors"
)

const (
	baseRemote = "origin/"
)

type dockerfileInformation struct {
	name      string
	content   []byte
	imageName string
}

// Execute core process
func Execute(config *types.Configuration) error {
	workDir, err := ioutil.TempDir("", "structor")
	if err != nil {
		return err
	}

	defer func() {
		err = cleanAll(workDir, config.Debug)
		if err != nil {
			log.Println("Error during cleanning: ", err)
		}
	}()

	if config.Debug {
		log.Printf("Temp directory: %s", workDir)
	}

	menuTemplateContent, err := getMenuTemplateContent(config.MenuTemplateFile, config.MenuTemplateURL)
	if err != nil {
		return err
	}

	dockerFileContent, err := downloadFile(config.DockerfileURL)
	if err != nil {
		return errors.Wrap(err, "failed to download Dockerfile")
	}

	baseDockerfile := dockerfileInformation{
		name:      fmt.Sprintf("%v.Dockerfile", time.Now().UnixNano()),
		content:   dockerFileContent,
		imageName: config.DockerImageName,
	}

	repoID := types.RepoID{
		Owner:          config.Owner,
		RepositoryName: config.RepositoryName,
	}

	return process(workDir, repoID, baseDockerfile, menuTemplateContent, config.ExperimentalBranchName, config.Debug)
}

func process(workDir string, repoID types.RepoID, baseDockerfile dockerfileInformation, menuTemplateContent []byte, experimentalBranchName string, debug bool) error {
	latestTagName, err := gh.GetLatestReleaseTagName(repoID)
	if err != nil {
		return err
	}

	branches, err := getBranches(experimentalBranchName, debug)
	if err != nil {
		return err
	}

	siteDir, err := buildSiteDirectory()
	if err != nil {
		return err
	}

	for _, branchRef := range branches {
		versionName := strings.Replace(branchRef, baseRemote, "", 1)

		log.Printf("Generating doc for version %s", versionName)

		versionsInfo := types.VersionsInformation{
			Current:      versionName,
			Latest:       latestTagName,
			Experimental: experimentalBranchName,
			CurrentPath:  filepath.Join(workDir, versionName),
		}

		err = buildDocumentation(branches, branchRef, versionsInfo, baseDockerfile, menuTemplateContent, debug)
		if err != nil {
			return err
		}

		outputDir := siteDir
		if !strings.HasPrefix(latestTagName, versionName) {
			outputDir = filepath.Join(outputDir, versionName)
		}

		err = copy.Copy(filepath.Join(versionsInfo.CurrentPath, "site"), outputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildDocumentation(branches []string, branchRef string, versionsInfo types.VersionsInformation, baseDockerfile dockerfileInformation, menuTemplateContent []byte, debug bool) error {
	err := repository.CreateWorkTree(versionsInfo.CurrentPath, branchRef, debug)
	if err != nil {
		return err
	}

	err = menu.Build(versionsInfo, branches, menuTemplateContent)
	if err != nil {
		return err
	}

	dockerfileVersionPath := filepath.Join(versionsInfo.CurrentPath, baseDockerfile.name)
	err = ioutil.WriteFile(dockerfileVersionPath, baseDockerfile.content, os.ModePerm)
	if err != nil {
		return err
	}

	dockerTagName := baseDockerfile.imageName + ":" + versionsInfo.Current

	// Build image
	output, err := dockerCmd(debug, "build", "-t", dockerTagName, "-f", dockerfileVersionPath, versionsInfo.CurrentPath+"/")
	if err != nil {
		log.Println(output)
		return err
	}

	// Run image
	output, err = dockerCmd(debug, "run", "--rm", "-v", versionsInfo.CurrentPath+":/mkdocs", dockerTagName, "mkdocs", "build")
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
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

	return branches, nil
}

func buildSiteDirectory() (string, error) {
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

func getMenuTemplateContent(menuTemplateFile, menuTemplateURL string) ([]byte, error) {
	if len(menuTemplateFile) > 0 {
		content, err := ioutil.ReadFile(menuTemplateFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get template menu file content")
		}
		return content, nil
	}

	content, err := downloadFile(menuTemplateURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download menu template")
	}
	return content, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, resp.Body.Close()
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

func dockerCmd(debug bool, args ...string) (string, error) {
	name := "docker"

	if debug {
		log.Println(name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func createDirectory(directoryPath string) error {
	_, err := os.Stat(directoryPath)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return err
		}
	}

	return os.MkdirAll(directoryPath, os.ModePerm)
}
