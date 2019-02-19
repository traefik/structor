package menu

import (
	"github.com/containous/structor/manifest"
)

func editManifest(manif map[string]interface{}, versionJsFile string, versionCSSFile string) {
	// Append menu JS file
	manifest.AppendExtraJs(versionJsFile, manif)

	// Append menu CSS file
	manifest.AppendExtraCSS(versionCSSFile, manif)

	// reset site URL
	manif["site_url"] = ""
}
