package backends

// Documentation: https://circleci.com/docs/2.0/artifacts/

import (
	"github.com/jszwedko/go-circleci"
	"github.com/ligurio/testres-db/formats"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// limit a number of builds that should be processed
// -1 means all available builds
const buildsNumber = 10
const isDebug = false

func SyncCircleCI(client *http.Client, b *Backend) (*[]formats.TestResult, error) {
	project_path := strings.Split(b.Project, "/")
	if len(project_path) != 2 {
		log.Println("Perhaps wrong project name specified")
		return nil, nil
	}

	account := project_path[0]
	repo := project_path[1]

	connection := &circleci.Client{Token: b.Secret, HTTPClient: client, Debug: isDebug}
	builds, err := connection.ListRecentBuildsForProject(account, repo, b.Branch, "", buildsNumber, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Found %d builds\n", len(builds))
	if len(builds) == 0 {
		log.Println("no builds found")
		return nil, nil
	}

	status := map[string]formats.TestStatus{
		"success": formats.StatusPass,
		"skipped": formats.StatusSkip,
		"failed":  formats.StatusFail,
	}

	results := make([]formats.TestResult, len(builds))
	for _, build := range builds {
		log.Printf("Found build: %d, status %s\n", build.BuildNum, build.Status)
		metadata, err := connection.ListTestMetadata(account, repo, build.BuildNum)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(metadata) == 0 {
			log.Printf("\tno tests\n")
			continue
		} else {
			log.Printf("\ttests %d\n", len(metadata))
		}

		artifacts, err := connection.ListBuildArtifacts(account, repo, build.BuildNum)
		for _, artifact := range artifacts {
			log.Println("Found artifact:", artifact.URL)
		}

		var testcases []formats.TestCase
		for _, test := range metadata {
			var testcase formats.TestCase
			log.Println("Found test:", test.Result, test.Name, test.RunTime)
			testcase = formats.TestCase{Name: test.Name, Status: status[test.Result]}
			testcases = append(testcases, testcase)
		}
		results = append(results, formats.TestResult{Name: strconv.Itoa(build.BuildNum), TestCases: testcases})
	}

	return &results, nil
}
