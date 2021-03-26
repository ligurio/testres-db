package backends

import (
	"context"
	"github.com/bndr/gojenkins"
	"github.com/ligurio/testres-db/formats"
	"log"
	"net/http"
)

func SyncJenkins(client *http.Client, b *Backend, buildsNumber int) (*[]formats.TestResult, error) {
	var jenkins *gojenkins.Jenkins
	jenkins = gojenkins.CreateJenkins(client, b.Base, b.Username, b.Secret)
	ctx := context.Background()
	_, err := jenkins.Init(ctx)
	if err != nil {
		return nil, err
	}

	jobBuilds, err := jenkins.GetAllBuildIds(ctx, b.Pipeline)
	if err != nil {
		return nil, err
	}

	if buildsNumber != -1 && len(jobBuilds) > buildsNumber {
		jobBuilds = jobBuilds[:buildsNumber]
	}
	results := make([]formats.TestResult, len(jobBuilds))
	for _, jobBuild := range jobBuilds {
		buildNum, err := jenkins.GetBuild(ctx, b.Pipeline, jobBuild.Number)
		if err != nil {
			return &results, err
		}
		log.Println(jobBuild.URL, buildNum.GetResult())
		TestResult, err := buildNum.GetResultSet(ctx)
		if err != nil {
			return &results, err
		}
		var testcases []formats.TestCase
		for _, suite := range TestResult.Suites {
			var testcase formats.TestCase
			for _, test := range suite.Cases {
				testcase = formats.TestCase{Name: test.Name}
			}
			testcases = append(testcases, testcase)
		}
		buildInfo := buildNum.Info()
		var result = formats.TestResult{Name: buildInfo.ID, TestCases: testcases}
		results = append(results, result)
	}

	return &results, nil
}
