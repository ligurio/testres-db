package xunit

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type XUnitReport struct {
	XMLName  xml.Name `xml:"assemblies"`
	Text     string   `xml:",chardata"`
	Assembly struct {
		Text          string `xml:",chardata"`
		TestFramework string `xml:"test-framework,attr"`
		Collection    struct {
			Text string `xml:",chardata"`
			Test struct {
				Text    string `xml:",chardata"`
				Name    string `xml:"name,attr"`
				Type    string `xml:"type,attr"`
				Method  string `xml:"method,attr"`
				Time    string `xml:"time,attr"`
				Result  string `xml:"result,attr"`
				Output  string `xml:"output"`
				Failure struct {
					Text          string `xml:",chardata"`
					ExceptionType string `xml:"exception-type,attr"`
					Message       string `xml:"message"`
					StackTrace    string `xml:"stack-trace"`
				} `xml:"failure"`
			} `xml:"test"`
		} `xml:"collection"`
	} `xml:"assembly"`
} 


func NewParser(r io.Reader) (*XUnitReport, error) {

	var report = new(XUnitReport)

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
