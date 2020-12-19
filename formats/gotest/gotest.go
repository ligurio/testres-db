package gotest

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

/* https://golang.org/cmd/test2json/ */
type GoTestReport struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func NewParser(r io.Reader) (*GoTestReport, error) {

	var report = new(GoTestReport)

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(buf), &report)
	if err != nil {
		return nil, err
	}
	return report, nil
}
