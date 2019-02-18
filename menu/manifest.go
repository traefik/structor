package menu

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ManifestFileName file name of the mkdocs manifest file
const ManifestFileName = "mkdocs.yml"

func editManifest(mkdocsFilePath string, versionJsFile string, versionCSSFile string) error {
	bytes, err := ioutil.ReadFile(mkdocsFilePath)
	if err != nil {
		return errors.Wrap(err, "error when reading MkDocs Manifest.")
	}

	manif := make(map[string]interface{})

	err = yaml.Unmarshal(bytes, manif)
	if err != nil {
		return errors.Wrap(err, "error when unmarshal MkDocs Manifest.")
	}

	// Append menu JS file
	if len(versionJsFile) > 0 {
		var extraJs []interface{}
		if val, ok := manif["extra_javascript"]; ok {
			extraJs = val.([]interface{})
		}
		extraJs = append(extraJs, versionJsFile)
		manif["extra_javascript"] = extraJs
	}

	// Append menu CSS file
	if len(versionCSSFile) > 0 {
		var extraCSS []interface{}
		if val, ok := manif["extra_css"]; ok {
			extraCSS = val.([]interface{})
		}
		extraCSS = append(extraCSS, versionCSSFile)
		manif["extra_css"] = extraCSS
	}

	// reset site URL
	manif["site_url"] = ""

	out, err := yaml.Marshal(manif)
	if err != nil {
		return errors.Wrap(err, "error when marshal MkDocs Manifest.")
	}

	return ioutil.WriteFile(mkdocsFilePath, out, os.ModePerm)
}
