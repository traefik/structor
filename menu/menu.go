package menu

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/containous/structor/manifest"
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

const menuJsFileName = "structor-menu.js"
const menuCSSFileName = "structor-menu.css"

// Build the menu
func Build(versionsInfo types.VersionsInformation, branches []string, menuContent types.MenuContent) error {
	manifestFile := filepath.Join(versionsInfo.CurrentPath, manifest.FileName)
	manifestDocsDir, err := manifest.GetManifestDocsDir(manifestFile)
	if err != nil {
		return err
	}
	log.Printf("Using docs_dir from manifest: %s", manifestDocsDir)

	var manifestJsFilePath string
	if len(menuContent.Js) > 0 {

		manifestJsFilePath = filepath.Join("theme", "js", menuJsFileName)

		jsDir := filepath.Join(manifestDocsDir, "theme", "js")
		_, errStat := os.Stat(jsDir)
		if os.IsNotExist(errStat) {
			errDir := os.MkdirAll(jsDir, os.ModePerm)
			if errDir != nil {
				return errors.Wrap(errDir, "error when create JS folder")
			}
		}

		menuFilePath := filepath.Join(jsDir, menuJsFileName)
		errBuild := buildJSFile(menuFilePath, versionsInfo, branches, string(menuContent.Js))
		if errBuild != nil {
			return errBuild
		}
	}

	var manifestCSSFilePath string
	if len(menuContent.CSS) > 0 {
		manifestCSSFilePath = filepath.Join("theme", "css", menuCSSFileName)

		cssDir := filepath.Join(manifestDocsDir, "theme", "css")
		_, errStat := os.Stat(cssDir)
		if os.IsNotExist(errStat) {
			errDir := os.MkdirAll(cssDir, os.ModePerm)
			if errDir != nil {
				return errors.Wrap(errDir, "error when create CSS folder")
			}
		}

		errWrite := ioutil.WriteFile(cssDir, menuContent.CSS, os.ModePerm)
		if errWrite != nil {
			return errors.Wrap(errWrite, "error when trying ro write CSS file")
		}
	}

	err = editManifest(manifestFile, manifestJsFilePath, manifestCSSFilePath)
	if err != nil {
		return errors.Wrap(err, "error when edit MkDocs manifest")
	}

	return nil
}

func buildJSFile(filePath string, versionsInfo types.VersionsInformation, branches []string, menuTemplate string) error {
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

			var v optionVersion
			switch {
			case simpleVersion.GreaterThan(latestVersion):
				v = optionVersion{
					Path:     versionName,
					Text:     versionName + " RC",
					Selected: selected,
				}
			case sameMinor(simpleVersion, latestVersion):
				// latest version
				v = optionVersion{
					Text:     versionName + " Latest",
					Selected: selected,
				}
			default:
				v = optionVersion{
					Path:     versionName,
					Text:     versionName,
					Selected: selected,
				}
			}
			versions = append(versions, v)
		}
	}

	return versions, nil
}

func sameMinor(v1 *version.Version, v2 *version.Version) bool {
	v1Parts := v1.Segments()
	v2Parts := v2.Segments()

	return v1Parts[0] == v2Parts[0] && v1Parts[1] == v2Parts[1]
}
