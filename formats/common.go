package formats

import (
	_ "github.com/ligurio/testres-db/formats/gotest"
	_ "github.com/ligurio/testres-db/formats/gtester"
	"github.com/ligurio/testres-db/formats/junit"
	_ "github.com/ligurio/testres-db/formats/subunit"
	_ "github.com/ligurio/testres-db/formats/tap13"
	_ "github.com/ligurio/testres-db/formats/trx"
	"log"
	"errors"
	"time"
	"os"
)

type fnParseFile func(filename string) (*TestResult, error)

var Parser = map[string]fnParseFile{
	"xml":     ParseJUnit,
	"tap":     ParseTAP13,
	"subunit": ParseSubUnit,
}

type TestResult struct {
	Name      string
	CreatedAt time.Time
	TestCases []TestCase
}

type TestCase struct {
	Name     string
	Status   TestStatus
	Duration time.Duration
}

type TestStatus int

const (
	StatusFail TestStatus = iota
	StatusPass
	StatusSkip
)

func (s TestStatus) String() string {
	return [...]string{"PASS", "FAIL", "SKIP"}[s]
}

var (
	errUnknownStatus = errors.New("Unknown test status")
)

func ParseJUnit(filename string) (*TestResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	report, err := junit.NewParser(file)
	if err != nil {
		return nil, err
	}

	res := TestResult{}
	for _, suite := range report.Suites {
		for _, testcase := range suite.TestCases {
			test := TestCase{Name: testcase.Name, Duration: 0}
			switch testcase.Status {
			case "Passed":
				test.Status = StatusPass
			case "Failed":
				test.Status = StatusFail
			case "Skipped":
				test.Status = StatusSkip
			}
			if test.Status == 0 {
				log.Println(testcase.Name, errUnknownStatus)
				return nil, errUnknownStatus
			}
			res.TestCases = append(res.TestCases, test)
		}
	}

	return nil, nil
}

func ParseTAP13(filename string) (*TestResult, error) {

	return nil, nil
}

func ParseSubUnit(filename string) (*TestResult, error) {

	return nil, nil
}

func ParseTRX(filename string) (*TestResult, error) {

	return nil, nil
}

func ParseGTester(filename string) (*TestResult, error) {

	return nil, nil
}
