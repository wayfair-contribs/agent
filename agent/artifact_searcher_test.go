package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/buildkite/agent/v3/api"
	"github.com/buildkite/agent/v3/logger"
	"github.com/stretchr/testify/assert"
)

func TestArtifactSearcherConnectsToEndpoint(t *testing.T) {
	defer os.Remove("llamas.txt")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.RequestURI() {
		case `/builds/my-build/artifacts/search?query=llamas.txt&scope=my-build&state=finished`:
			fmt.Fprint(rw, `[{
				"id": "4600ac5c-5a13-4e92-bb83-f86f218f7b32",
				"file_size": 3,
				"absolute_path": "llamas.txt",
				"path": "llamas.txt",
				"url": "http://example.com/download"
			}]`)
		default:
			fmt.Println(req.URL.String())
			http.Error(rw, "Not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	ac := api.NewClient(logger.Discard, api.Config{
		Endpoint: server.URL,
		Token:    `llamasforever`,
	})

	s := NewArtifactSearcher(logger.Discard, ac, "my-build")

	artifacts, err := s.Search("llamas.txt", "my-build", false, false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*api.Artifact{{
		ID:           "4600ac5c-5a13-4e92-bb83-f86f218f7b32",
		Path:         "llamas.txt",
		AbsolutePath: "llamas.txt",
		FileSize:     3,
		URL:          "http://example.com/download",
	}}, artifacts)
}
