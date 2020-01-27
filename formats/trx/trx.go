package trx

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type TRXReport struct {
	XMLName xml.Name `xml:"TestRun"`
	Text    string   `xml:",chardata"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	RunUser string   `xml:"runUser,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Times   struct {
		Text     string `xml:",chardata"`
		Creation string `xml:"creation,attr"`
		Queuing  string `xml:"queuing,attr"`
		Start    string `xml:"start,attr"`
		Finish   string `xml:"finish,attr"`
	} `xml:"Times"`
	Results struct {
		Text           string `xml:",chardata"`
		UnitTestResult []struct {
			Text                     string `xml:",chardata"`
			ExecutionId              string `xml:"executionId,attr"`
			TestId                   string `xml:"testId,attr"`
			TestName                 string `xml:"testName,attr"`
			ComputerName             string `xml:"computerName,attr"`
			Duration                 string `xml:"duration,attr"`
			StartTime                string `xml:"startTime,attr"`
			EndTime                  string `xml:"endTime,attr"`
			TestType                 string `xml:"testType,attr"`
			Outcome                  string `xml:"outcome,attr"`
			TestListId               string `xml:"testListId,attr"`
			RelativeResultsDirectory string `xml:"relativeResultsDirectory,attr"`
			Output                   struct {
				Text      string `xml:",chardata"`
				ErrorInfo struct {
					Text    string `xml:",chardata"`
					Message string `xml:"Message"`
				} `xml:"ErrorInfo"`
			} `xml:"Output"`
		} `xml:"UnitTestResult"`
	} `xml:"Results"`
	TestDefinitions struct {
		Text     string `xml:",chardata"`
		UnitTest []struct {
			Text      string `xml:",chardata"`
			Name      string `xml:"name,attr"`
			Storage   string `xml:"storage,attr"`
			ID        string `xml:"id,attr"`
			Execution struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"Execution"`
			TestMethod struct {
				Text      string `xml:",chardata"`
				CodeBase  string `xml:"codeBase,attr"`
				ClassName string `xml:"className,attr"`
				Name      string `xml:"name,attr"`
			} `xml:"TestMethod"`
		} `xml:"UnitTest"`
	} `xml:"TestDefinitions"`
	TestEntries struct {
		Text      string `xml:",chardata"`
		TestEntry []struct {
			Text        string `xml:",chardata"`
			TestId      string `xml:"testId,attr"`
			ExecutionId string `xml:"executionId,attr"`
			TestListId  string `xml:"testListId,attr"`
		} `xml:"TestEntry"`
	} `xml:"TestEntries"`
	TestLists struct {
		Text     string `xml:",chardata"`
		TestList struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			ID   string `xml:"id,attr"`
		} `xml:"TestList"`
	} `xml:"TestLists"`
}

func NewParser(r io.Reader) (*TRXReport, error) {

	var report = new(TRXReport)

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
