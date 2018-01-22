package menu

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/containous/structor/types"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

const baseRemote = "origin/"

type optionVersion struct {
	Path     string
	Text     string
	Selected bool
}

const menuFileName = "version-menu.js"

// Build the menu
func Build(versionsInfo types.VersionsInformation, branches []string, menuTemplateContent []byte) error {
	manifestFile := filepath.Join(versionsInfo.CurrentPath, "mkdocs.yml")
	manifestMenuFilePath := filepath.Join("theme", "js", menuFileName)
	err := editManifest(manifestFile, manifestMenuFilePath)
	if err != nil {
		return errors.Wrap(err, "error when edit MkDocs manifest")
	}

	jsDir := filepath.Join(versionsInfo.CurrentPath, "docs", "theme", "js")
	_, err = os.Stat(jsDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(jsDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "error when create JS folder")
		}
	}

	menuFilePath := filepath.Join(jsDir, menuFileName)
	return buildMenuFile(menuFilePath, versionsInfo, branches, string(menuTemplateContent))
}

func buildMenuFile(filePath string, versionsInfo types.VersionsInformation, branches []string, menuTemplate string) error {
	temp := template.New("menu-js")

	_, err := temp.Parse(menuTemplate)
	if err != nil {
		return errors.Wrap(err, "error during parsing template")
	}

	versions, err := buildVersions(versionsInfo.Current, branches, versionsInfo.Latest, versionsInfo.Experimental)
	if err != nil {
		return errors.Wrap(err, "error when build versions")
	}

	model := struct {
		Versions []optionVersion
	}{
		Versions: versions,
	}

	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "error when create menu file")
	}

	return temp.Execute(f, model)
}

func buildVersions(currentVersion string, branches []string, latestTagName string, experimentalBranchName string) ([]optionVersion, error) {
	latestVersion, err := version.NewVersion(latestTagName)
	if err != nil {
		return nil, err
	}

	var versions []optionVersion
	for _, branch := range branches {
		versionName := strings.Replace(branch, baseRemote, "", 1)
		selected := currentVersion == versionName

		switch versionName {
		case latestTagName:
			// skip, because we must the branch instead of the tag
		case experimentalBranchName:
			versions = append(versions, optionVersion{
				Path:     experimentalBranchName,
				Text:     "Experimental",
				Selected: selected,
			})
		default:
			simpleVersion, err := version.NewVersion(versionName)
			if err != nil {
				return nil, err
			}

			if simpleVersion.GreaterThan(latestVersion) {
				versions = append(versions, optionVersion{
					Path:     versionName,
					Text:     versionName + " RC",
					Selected: selected,
				})
			} else if sameMinor(simpleVersion, latestVersion) {
				// latest version
				versions = append(versions, optionVersion{
					Text:     versionName + " Stable",
					Selected: selected,
				})
			} else {
				versions = append(versions, optionVersion{
					Path:     versionName,
					Text:     versionName,
					Selected: selected,
				})
			}
		}
	}

	return versions, nil
}

func sameMinor(v1 *version.Version, v2 *version.Version) bool {
	v1Parts := v1.Segments()
	v2Parts := v2.Segments()

	return v1Parts[0] == v2Parts[0] && v1Parts[1] == v2Parts[1]
}
