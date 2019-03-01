package menu

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/containous/structor/file"
	"github.com/containous/structor/manifest"
	"github.com/containous/structor/types"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

const baseRemote = "origin/"

const (
	menuJsFileName  = "structor-menu.js"
	menuCSSFileName = "structor-menu.css"
)

// Content Content of menu files.
type Content struct {
	Js  []byte
	CSS []byte
}

type optionVersion struct {
	Path     string
	Text     string
	Selected bool
}

// GetTemplateContent Gets menu template content.
func GetTemplateContent(menu *types.MenuFiles) Content {
	var content Content

	if menu.HasJsFile() {
		jsContent, err := getMenuFileContent(menu.JsFile, menu.JsURL)
		if err != nil {
			return Content{}
		}
		content.Js = jsContent
	}

	if menu.HasCSSFile() {
		cssContent, err := getMenuFileContent(menu.CSSFile, menu.CSSURL)
		if err != nil {
			return Content{}
		}
		content.CSS = cssContent
	}

	return content
}

func getMenuFileContent(f string, u string) ([]byte, error) {
	if len(f) > 0 {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get template menu file content")
		}
		return content, nil
	}

	content, err := file.Download(u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download menu template")
	}
	return content, nil
}

// Build the menu.
func Build(versionsInfo types.VersionsInformation, branches []string, menuContent Content) error {
	manifestFile := filepath.Join(versionsInfo.CurrentPath, manifest.FileName)

	manif, err := manifest.Read(manifestFile)
	if err != nil {
		return err
	}

	manifestDocsDir, err := manifest.GetDocsDir(manifestFile, manif)
	if err != nil {
		return err
	}

	log.Printf("Using docs_dir from manifest: %s", manifestDocsDir)

	manifestJsFilePath, err := writeJsFile(manifestDocsDir, menuContent, versionsInfo, branches)
	if err != nil {
		return err
	}

	manifestCSSFilePath, err := writeCSSFile(manifestDocsDir, menuContent)
	if err != nil {
		return err
	}

	editManifest(manif, manifestJsFilePath, manifestCSSFilePath)

	err = manifest.Write(manifestFile, manif)
	if err != nil {
		return errors.Wrap(err, "error when edit MkDocs manifest")
	}

	return nil
}

func writeCSSFile(manifestDocsDir string, menuContent Content) (string, error) {
	if len(menuContent.CSS) == 0 {
		return "", nil
	}

	cssDir := filepath.Join(manifestDocsDir, "theme", "css")
	if _, errStat := os.Stat(cssDir); os.IsNotExist(errStat) {
		errDir := os.MkdirAll(cssDir, os.ModePerm)
		if errDir != nil {
			return "", errors.Wrap(errDir, "error when create CSS folder")
		}
	}

	err := ioutil.WriteFile(filepath.Join(cssDir, menuCSSFileName), menuContent.CSS, os.ModePerm)
	if err != nil {
		return "", errors.Wrap(err, "error when trying ro write CSS file")
	}

	return filepath.Join("theme", "css", menuCSSFileName), nil
}

func writeJsFile(manifestDocsDir string, menuContent Content, versionsInfo types.VersionsInformation, branches []string) (string, error) {
	if len(menuContent.Js) == 0 {
		return "", nil
	}

	jsDir := filepath.Join(manifestDocsDir, "theme", "js")
	if _, errStat := os.Stat(jsDir); os.IsNotExist(errStat) {
		errDir := os.MkdirAll(jsDir, os.ModePerm)
		if errDir != nil {
			return "", errors.Wrap(errDir, "error when create JS folder")
		}
	}

	menuFilePath := filepath.Join(jsDir, menuJsFileName)
	errBuild := buildJSFile(menuFilePath, versionsInfo, branches, string(menuContent.Js))
	if errBuild != nil {
		return "", errBuild
	}

	return filepath.Join("theme", "js", menuJsFileName), nil
}

func buildJSFile(filePath string, versionsInfo types.VersionsInformation, branches []string, menuTemplate string) error {
	temp := template.New("menu-js").Funcs(sprig.TxtFuncMap())

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
