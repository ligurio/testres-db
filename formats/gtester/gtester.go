/*
GLib provides a testing framework.
There are two parts: a set of functions in GLib that you can use when writing
unit tests, and a test runner tool named gtester.

GTester can generate a machine-readable summary of the test results. The report
uses a GTester-specific XML format.

https://wiki.gnome.org/Projects/GLib/GTester

*/

package gtester

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type GTesterReport struct {
	XMLName    xml.Name `xml:"gtester"`
	Text       string   `xml:",chardata"`
	Testbinary []struct {
		Text   string `xml:",chardata"`
		Path   string `xml:"path,attr"`
		Binary struct {
			Text string `xml:",chardata"`
			File string `xml:"file,attr"`
		} `xml:"binary"`
		RandomSeed string `xml:"random-seed"`
		Testcase   []struct {
			Text     string `xml:",chardata"`
			Path     string `xml:"path,attr"`
			Skipped  string `xml:"skipped,attr"`
			Duration string `xml:"duration"`
			Status   struct {
				Text       string `xml:",chardata"`
				ExitStatus string `xml:"exit-status,attr"`
				NForks     string `xml:"n-forks,attr"`
				Result     string `xml:"result,attr"`
			} `xml:"status"`
			Message string `xml:"message"`
			Error   string `xml:"error"`
		} `xml:"testcase"`
		Duration string `xml:"duration"`
	} `xml:"testbinary"`
}

func NewParser(r io.Reader) (*GTesterReport, error) {

	var report = new(GTesterReport)

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
