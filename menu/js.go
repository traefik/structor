package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/go-version"
	"github.com/traefik/structor/types"
)

const menuJsFileName = "structor-menu.js"

const (
	stateLatest          = "LATEST"
	stateExperimental    = "EXPERIMENTAL"
	statePreFinalRelease = "PRE_FINAL_RELEASE"
	stateObsolete        = "OBSOLETE"
)

type optionVersion struct {
	Path     string
	Text     string
	Name     string
	State    string
	Selected bool
}

func writeJsFile(manifestDocsDir string, menuContent Content, versionsInfo types.VersionsInformation, branches []string) (string, error) {
	if len(menuContent.Js) == 0 {
		return "", nil
	}

	jsDir := filepath.Join(manifestDocsDir, "theme", "js")
	if _, errStat := os.Stat(jsDir); os.IsNotExist(errStat) {
		errDir := os.MkdirAll(jsDir, os.ModePerm)
		if errDir != nil {
			return "", fmt.Errorf("error when create JS folder: %w", errDir)
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
		return fmt.Errorf("error during parsing template: %w", err)
	}

	versions, err := buildVersions(versionsInfo.Current, branches, versionsInfo.Latest, versionsInfo.Experimental)
	if err != nil {
		return fmt.Errorf("error when build versions: %w", err)
	}

	model := struct {
		Latest   string
		Current  string
		Versions []optionVersion
	}{
		Latest:   versionsInfo.Latest,
		Current:  versionsInfo.Current,
		Versions: versions,
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error when create menu file: %w", err)
	}

	return temp.Execute(f, model)
}

func buildVersions(currentVersion string, branches []string, latestTagName, experimentalBranchName string) ([]optionVersion, error) {
	latestVersion, err := version.NewVersion(latestTagName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse latest tag version %s: %w", latestTagName, err)
	}

	rawVersions, heads := parseBranches(branches)

	var versions []optionVersion
	for _, versionName := range rawVersions {
		selected := currentVersion == versionName

		switch versionName {
		case latestTagName:
			// skip, because we must use the branch instead of the tag
		case experimentalBranchName:
			versions = append(versions, optionVersion{
				Path:     experimentalBranchName,
				Text:     "Experimental",
				Name:     experimentalBranchName,
				State:    stateExperimental,
				Selected: selected,
			})

		default:
			simpleVersion, err := version.NewVersion(versionName)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version %s: %w", versionName, err)
			}

			v := optionVersion{
				Name:     versionName,
				Selected: selected,
			}

			switch {
			case simpleVersion.GreaterThan(latestVersion):
				v.Path = versionName
				v.Text = versionName + " RC"
				v.State = statePreFinalRelease
			case sameMinor(simpleVersion, latestVersion):
				// latest version
				v.Text = versionName + " Latest"
				v.State = stateLatest
			default:
				v.Path = versionName
				v.Text = versionName
				if !isHeads(heads, simpleVersion) {
					v.State = stateObsolete
				}
			}

			versions = append(versions, v)
		}
	}

	return versions, nil
}

func parseBranches(branches []string) ([]string, map[int]*version.Version) {
	heads := map[int]*version.Version{}

	var rawVersions []string
	for _, branch := range branches {
		versionName := strings.Replace(branch, baseRemote, "", 1)
		rawVersions = append(rawVersions, versionName)

		v, err := version.NewVersion(versionName)
		if err != nil {
			continue
		}

		if p, ok := heads[v.Segments()[0]]; ok {
			if v.GreaterThan(p) {
				heads[v.Segments()[0]] = v
			}
		} else {
			heads[v.Segments()[0]] = v
		}
	}
	return rawVersions, heads
}

func sameMinor(v1, v2 *version.Version) bool {
	v1Parts := v1.Segments()
	v2Parts := v2.Segments()

	return v1Parts[0] == v2Parts[0] && v1Parts[1] == v2Parts[1]
}

func isHeads(heads map[int]*version.Version, v *version.Version) bool {
	for _, head := range heads {
		if v.Equal(head) {
			return true
		}
	}

	return false
}
