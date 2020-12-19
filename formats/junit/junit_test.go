package junit

import (
	"encoding/json"
	"flag"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testcase struct {
	Filename   string
	XFail      bool
	GoldenFile string
}

/* https://github.com/jstemmer/go-junit-report/blob/master/go-junit-report_test.go */

var testset = []testcase{
	{
		Filename:   "junit-sample-1.xml",
		XFail:      false,
		GoldenFile: "junit-sample-1.json",
	},
	{
		Filename:   "junit-sample-2.xml",
		XFail:      false,
		GoldenFile: "junit-sample-2.json",
	},
	{
		Filename:   "junit-sample-3.xml",
		XFail:      false,
		GoldenFile: "junit-sample-3.json",
	},
	{
		Filename:   "junit-sample-4.xml",
		XFail:      false,
		GoldenFile: "junit-sample-4.json",
	},
	{
		Filename:   "junit-sample-5.xml",
		XFail:      false,
		GoldenFile: "junit-sample-5.json",
	},
	{
		Filename:   "junit-sample-6.xml",
		XFail:      false,
		GoldenFile: "junit-sample-6.json",
	},
	{
		/* JUnit report, that is successfully parsed by Atlassian Bamboo */
		Filename:   "junit-sample-7.xml",
		XFail:      true,
		GoldenFile: "junit-sample-7.json",
	},
	{
		/* JUnit report, that is successfully parsed by Atlassian Bamboo */
		Filename:   "junit-sample-8.xml",
		XFail:      true,
		GoldenFile: "junit-sample-8.json",
	},
}

var update = flag.Bool("update", false, "update .golden files")

func TestParser(t *testing.T) {
	for _, test := range testset {
		if test.XFail {
			t.Logf("XFail: %s", test.Filename)
			continue
		}
		t.Logf("Running: %s", test.Filename)

		file, err := os.Open("testdata/" + test.Filename)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Match structs", test.Filename)

		report, err := NewParser(file)
		if err != nil {
			t.Fatalf("Error parsing %s: %s", test.Filename, err)
		}

		goldenPath := filepath.Join("testdata", test.GoldenFile)
		if *update {
			t.Log("Update", goldenPath)
			data, _ := json.MarshalIndent(report, "", " ")
			if err := ioutil.WriteFile(goldenPath, data, 0644); err != nil {
				t.Fatalf("Failed to update golden file: %s", err)
			}
		}

		rawJson, err := ioutil.ReadFile(goldenPath)
		if err != nil {
			t.Fatal(err.Error())
		}
		var goldenReport JUnitReport
		json.Unmarshal(rawJson, &goldenReport)
		if diff := cmp.Diff(*report, goldenReport); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
