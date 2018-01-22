package menu

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func editManifest(mkdocsFilePath string, versionJsFile string) error {
	bytes, err := ioutil.ReadFile(mkdocsFilePath)
	if err != nil {
		return errors.Wrap(err, "error when reading MkDocs Manifest.")
	}

	manif := make(map[string]interface{})

	err = yaml.Unmarshal(bytes, manif)
	if err != nil {
		return errors.Wrap(err, "error when unmarshal MkDocs Manifest.")
	}

	var extraJs []interface{}
	if val, ok := manif["extra_javascript"]; ok {
		extraJs = val.([]interface{})
	}
	extraJs = append(extraJs, versionJsFile)
	manif["extra_javascript"] = extraJs
	manif["site_url"] = ""

	out, err := yaml.Marshal(manif)
	if err != nil {
		return errors.Wrap(err, "error when marshal MkDocs Manifest.")
	}

	return ioutil.WriteFile(mkdocsFilePath, out, os.ModePerm)
}
