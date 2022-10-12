package file

import (
	"fmt"
	"io"
	"net/http"
)

// Download Downloads a file.
func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("failed to download %q: %s", url, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
