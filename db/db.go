package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ligurio/testres-db/formats"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DB struct {
	db   *sql.DB
	path string
}

const projectsTbl = `
CREATE TABLE IF NOT EXISTS "project" (
	project_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT UNIQUE NOT NULL DEFAULT '',
	UNIQUE (name)
);

CREATE UNIQUE INDEX name_index ON project(name);
`

const buildsTbl = `
CREATE TABLE IF NOT EXISTS "build" (
	build_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT NOT NULL DEFAULT '',
	project_id INTEGER,
	created TEXT,
	UNIQUE (name)
);
`
const testsTbl = `
CREATE TABLE IF NOT EXISTS "test" (
	test_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT NOT NULL DEFAULT '',
	build_id TEXT NOT NULL DEFAULT '',
	FOREIGN KEY(build_id) REFERENCES build(build_id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE (name)
);
`

const statusesTbl = `
CREATE TABLE IF NOT EXISTS "status" (
	status_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT,
	UNIQUE (name)
);

INSERT INTO status(name) VALUES ('PASSED'), ('FAILED'), ('SKIPPED');
`

const testrunsTbl = `
CREATE TABLE IF NOT EXISTS "testrun" (
	testrun_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	test_id INTEGER,
	build_id INTEGER,
	status_id INTEGER,
	duration TEXT,
	FOREIGN KEY(test_id) REFERENCES test(test_id) ON DELETE CASCADE ON UPDATE CASCADE
	FOREIGN KEY(status_id) REFERENCES status(status_id) ON DELETE CASCADE ON UPDATE CASCADE
	FOREIGN KEY(build_id) REFERENCES build(build_id) ON DELETE CASCADE ON UPDATE CASCADE
);
`

func (db *DB) Init() error {
	for _, sql := range []string{projectsTbl, buildsTbl, testsTbl, statusesTbl, testrunsTbl} {
		if _, err := db.db.Exec(sql); err != nil {
			return errors.New(fmt.Sprintf("SQL error with %s\n%s", sql, err.Error()))
		}
	}
	return nil
}

func (db *DB) SetPath(p string) {
	db.path = p
}

func (db *DB) GetPath() string {
	return db.path
}

func (db *DB) Open() error {
	var (
		err    error
		dbPath string
	)
	if db.GetPath() != "" {
		dbPath = db.GetPath()
	} else {
		dbPath = "testres.sqlite"
	}
	db.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) GetProjectIdByName(project string) (int64, error) {
	var projectId int64
	row := db.db.QueryRow("SELECT project_id FROM project WHERE name=$1", project)
	err := row.Scan(&projectId)
	if err != nil {
		return 0, err
	}

	return projectId, nil
}

func (db *DB) AddProject(project string) (int64, error) {
	var projectId int64
	stmt, err := db.db.Prepare("INSERT INTO project (name) values(?)")
	defer stmt.Close()
	res, err := stmt.Exec(project)
	if err != nil {
		return 0, err
	}
	projectId, err = res.LastInsertId()

	return projectId, err
}

func (db *DB) AddBuild(project_id int64, build string, created time.Time) (int64, error) {
	var buildId int64
	stmt, err := db.db.Prepare("INSERT INTO build (name, project_id, created) values(?, ?, ?)")
	defer stmt.Close()
	res, err := stmt.Exec(build, project_id, created)
	if err != nil {
		return 0, err
	}
	buildId, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return buildId, err
}

func (db *DB) GetBuildIdByName(build string) (int64, error) {
	var buildId int64
	row := db.db.QueryRow("SELECT build_id FROM build WHERE name=$1", build)
	err := row.Scan(&buildId)
	if err != nil {
		return 0, err
	}

	return buildId, err
}

func (db *DB) AddTest(testname string) (int64, error) {
	var testId int64
	stmt, _ := db.db.Prepare("INSERT INTO test (name) values(?)")
	defer stmt.Close()
	res, err := stmt.Exec(testname)
	if err != nil {
		return 0, err
	}
	testId, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return testId, nil
}

func (db *DB) GetTestIdByName(testname string) (int64, error) {
	var testId int64
	row := db.db.QueryRow("SELECT test_id FROM test WHERE name=$1", testname)
	err := row.Scan(&testId)
	if err != nil {
		return 0, err
	}

	return testId, nil
}

func (db *DB) GetTimeOfLatestBuild(project string) (time.Time, error) {
	project_id, err := db.GetProjectIdByName(project)
	if err != nil {
		return time.Time{}, err
	}
	row := db.db.QueryRow("SELECT created FROM build WHERE project_id=$1", project_id)
	var created time.Time
	err = row.Scan(&created)
	if err != nil {
		return time.Time{}, err
	}

	return created, nil
}

func (db *DB) AddResults(project string, results *[]formats.TestResult) error {
	stmt, err := db.db.Prepare("INSERT INTO testrun(test_id, build_id, status_id, duration) VALUES(?, ?, ?, ?)")
	defer stmt.Close()

	for _, build := range *results {
		var buildId int64
		buildId, err = db.GetBuildIdByName(build.Name)
		if err != nil {
			var projectId int64
			projectId, err = db.GetProjectIdByName(project)
			if err != nil {
				projectId, err = db.AddProject(project)
				if err != nil {
					return err
				}
			}
			buildId, err = db.AddBuild(projectId, build.Name, build.CreatedAt)
			if err != nil {
				return err
			}
		}
		for _, testcase := range build.TestCases {
			var testId int64
			testId, err = db.GetTestIdByName(testcase.Name)
			if err != nil {
				testId, err = db.AddTest(testcase.Name)
				if err != nil {
					return err
				}
			}
			_, err = stmt.Exec(testId, buildId, testcase.Status, testcase.Duration)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
