package menu

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const menuCSSFileName = "structor-menu.css"

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
