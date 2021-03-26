package backends

import (
	"os"
	"testing"
)

func TestSyncCircleCI(t *testing.T) {
	t.Log("Basic test with spotify/helios project")
	username := os.Getenv("CIRCLECI_USERNAME")
	token := os.Getenv("CIRCLECI_TOKEN")
	if username == "" || token == "" {
		t.Skip("No CIRCLECI_USERNAME and CIRCLECI_TOKEN.")
	}
	backend := Backend{Type: "circleci", Base: "https://circleci.com/", Name: "spotify/helios",
		Project: "spotify/helios", Branch: "master", Username: username, Secret: token}

	buildsNumber := 5
	httpClient := NewAPIClient()
	builds, err := SyncCircleCI(httpClient, &backend, buildsNumber)
	if builds == nil || err != nil {
		t.Failed()
	}
}
