package db

import (
	"fmt"
	"github.com/ligurio/testres-db/formats"
	"os"
	"testing"
	"time"
)

const (
	dbPath      = "testres-test.db"
	projectName = "project #1"
	buildName   = "#1"
	testName    = "test #1"
)

func setUpDB(t *testing.T) *DB {
	db := &DB{}
	db.SetPath(dbPath)
	err := db.Open()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = db.Init()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return db
}

func tearDownDB(t *testing.T, db *DB) {
	db.Close()
	err := os.Remove(db.GetPath())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestAddProject(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	_, err := db.AddProject(projectName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetProjectIdByName_NotFound(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	_, err := db.GetProjectIdByName(projectName)
	if err == nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetProjectIdByName(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	projectId, err := db.AddProject(projectName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	id, err := db.GetProjectIdByName(projectName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if id != projectId {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetBuildIdByName_NotFound(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	_, err := db.GetBuildIdByName(buildName)
	if err == nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetBuildIdByName(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	projectId, err := db.AddProject(projectName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	buildId, err := db.AddBuild(projectId, buildName, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	id, err := db.GetBuildIdByName(buildName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if id != buildId {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetTestIdByName_NotFound(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	_, err := db.GetTestIdByName(testName)
	if err == nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetTestIdByName(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	testId, err := db.AddTest(testName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	id, err := db.GetTestIdByName(testName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if id != testId {
		t.Log(err)
		t.FailNow()
	}
}

func TestAddTest(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)

	_, err := db.AddTest(testName)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func setUpResults(t *testing.T) []formats.TestResult {
	testcases := make([]formats.TestCase, 20)
	t1 := time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2017, time.February, 16, 0, 0, 0, 0, time.UTC)
	duration := t2.Sub(t1)
	for i := range testcases {
		testcases[i].Name = fmt.Sprintf("test %d", i)
		testcases[i].Status = formats.StatusFail
		testcases[i].Duration = duration
	}
	testcases[3].Status = formats.StatusPass
	testcases[6].Status = formats.StatusSkip
	testcases[8].Status = formats.StatusPass
	testcases[12].Status = formats.StatusPass
	testcases[17].Status = formats.StatusSkip
	results := make([]formats.TestResult, 2)
	results[0] = formats.TestResult{Name: buildName, TestCases: testcases, CreatedAt: time.Now().Format(time.RFC3339)}
	results[1] = formats.TestResult{Name: "build #2", TestCases: testcases, CreatedAt: time.Now().Format(time.RFC3339)}

	return results
}

func TestAddResults(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	results := setUpResults(t)
	err := db.AddResults(projectName, &results)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
