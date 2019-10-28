package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Download Downloads a file.
func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("failed to download %q: %s", url, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
