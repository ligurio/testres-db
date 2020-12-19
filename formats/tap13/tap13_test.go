package tap13

import (
	"testing"
)

type tcase struct {
	name string
	skip bool
}

var testset = []tcase{

	{name: "tap-sample-01.tap", skip: false},
	{name: "tap-sample-02.tap", skip: false},
	{name: "tap-sample-03.tap", skip: false},
	{name: "tap-sample-04.tap", skip: false},
	{name: "tap-sample-05.tap", skip: false},
	{name: "tap-sample-06.tap", skip: false},
	{name: "tap-sample-07.tap", skip: false},
	{name: "tap-sample-08.tap", skip: false},
	// CRIU tests
	{name: "tap-sample-09.tap", skip: false},
	// GIT tests
	{name: "tap-sample-10.tap", skip: false},
	// LibVirt tests (libvirt-tck-f10-broken.txt)
	{name: "tap-sample-11.tap", skip: false},
	// LibVirt tests (libvirt-tck-f10-fixed.txt)
	{name: "tap-sample-12.tap", skip: false},
	/*
		TODO:
		- pytest-tap https://github.com/python-tap/pytest-tap
		- postgresql project
	*/
}

func TestNewParser(t *testing.T) {

}
