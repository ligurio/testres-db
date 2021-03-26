package backends

import (
	"os"
	"testing"
)

func TestSyncGitLab(t *testing.T) {
	t.Log("Basic test with cmocka/cmocka project")
	username := os.Getenv("GITLAB_USERNAME")
	token := os.Getenv("GITLAB_TOKEN")
	if username == "" || token == "" {
		t.Skip("No GITLAB_USERNAME and GITLAB_TOKEN.")
	}
	backend := Backend{Type: "gitlab", Base: "https://gitlab.com/", Name: "cmocka/cmocka",
		Project: "cmocka/cmocka", Branch: "master", Username: username, Secret: token}

	buildsNumber := 5
	httpClient := NewAPIClient()
	builds, err := SyncGitLab(httpClient, &backend, buildsNumber)
	if builds == nil || err != nil {
		t.Failed()
	}
}
