package file

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	serverURL, teardown := serveFixturesContent()
	defer teardown()

	uri := serverURL + "/data.txt"

	bytes, err := Download(uri)
	require.NoError(t, err)

	assert.Equal(t, mustReadFile(filepath.Join(".", "fixtures", "data.txt")), bytes)
}

func TestDownload_Fail(t *testing.T) {
	serverURL, teardown := serveFixturesContent()
	defer teardown()

	uri := serverURL + "/missing.txt"

	_, err := Download(uri)

	require.Error(t, err)
	assert.Regexp(t, `failed to download "http://127.0.0.1:\d+/missing.txt": 404 Not Found`, err.Error())
}

func mustReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func serveFixturesContent() (string, func()) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./fixtures")))

	server := httptest.NewServer(mux)

	return server.URL, server.Close
}
