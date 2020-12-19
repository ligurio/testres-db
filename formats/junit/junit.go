/*
<testsuites>        => the aggregated result of all junit testfiles
  <testsuite>       => the output from a single TestSuite
    <properties>    => the defined properties at test execution
      <property>    => name/value pair for a single property
      ...
    </properties>
    <error></error> => optional information, in place of a test case - normally if the tests in the suite could not be found etc.
    <testcase>      => the results from executing a test method
      <system-out>  => data written to System.out during the test run
      <system-err>  => data written to System.err during the test run
      <skipped/>    => test was skipped
      <failure>     => test failed
      <error>       => test encountered an error
    </testcase>
    ...
  </testsuite>
  ...
</testsuites>
*/

package junit

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitReport struct {
	XMLName  xml.Name `xml:"testsuites"`
	Disabled int      `xml:"disabled,attr,omitempty"`
	Errors   int      `xml:"errors,attr,omitempty"`
	Failures int      `xml:"failures,attr,omitempty"`
	Name     string   `xml:"name,attr,omitempty"`
	Time     float64  `xml:"time,attr,omitempty"`
	Tests    int      `xml:"tests,attr,omitempty"`

	Suites []JUnitTestsuite `xml:"testsuite"`
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestsuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Disabled  int      `xml:"disabled,attr,omitempty"`
	Errors    int      `xml:"errors,attr,omitempty"`
	Failures  int      `xml:"failures,attr,omitempty"`
	Hostname  string   `xml:"hostname,attr,omitempty"`
	Id        string   `xml:"id,attr,omitempty"`
	Name      string   `xml:"name,attr,omitempty"`
	Package   string   `xml:"package,attr,omitempty"`
	Skipped   int      `xml:"skipped,attr,omitempty"`
	Tests     int      `xml:"tests,attr,omitempty"`
	Time      float64  `xml:"time,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`

	Properties []JUnitProperty `xml:"property->property"`
	TestCases  []JUnitTestcase `xml:"testcase"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	XMLName xml.Name `xml:"property"`
	name    string   `xml:"name,attr"`
	value   string   `xml:"value,attr"`
}

// JUnitTestCase is a single test case with its result.
type JUnitTestcase struct {
	XMLName    xml.Name      `xml:"testcase"`
	Assertions string        `xml:"assertions,attr,omitempty"`
	Classname  string        `xml:"classname,attr,omitempty"`
	Error      *JUnitFailure `xml:"error,omitempty"`
	Failure    *JUnitFailure `xml:"failure,omitempty"`
	Name       string        `xml:"name,attr,omitempty"`
	Skipped    int           `xml:"skipped,attr,omitempty"`
	Status     string        `xml:"status,attr,omitempty"`
	Systemout  InnerResult   `xml:"system-out,omitempty"`
	Systemerr  InnerResult   `xml:"system-err,omitempty"`
	Time       float64       `xml:"time,attr,omitempty"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Value   string `xml:",innerxml"`
	Type    string `xml:"type,attr,omitempty"`
	Message string `xml:"message,attr,omitempty"`
}

type InnerResult struct {
	Value string `xml:",innerxml"`
}

func NewParser(r io.Reader) (*JUnitReport, error) {

	var report = new(JUnitReport)

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal([]byte(buf), &report)
	if err != nil {
		return nil, err
	}
	return report, nil
}
