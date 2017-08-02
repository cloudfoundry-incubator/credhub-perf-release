package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry-incubator/credhub_cannon"

	"github.com/onsi/gomega/gexec"
	"testing"
)

func NewRampedRequest(args Args) {
	rampedRequest := &RampedRequest{
		args.MinConcurrent,
		args.MaxConcurrent,
		args.Step,
		args.NumberOfRequests,
		args.LocalCSV}
	rampedRequest.FireRequests(args.URL, "POST", args.RequestBody)
}

type Args struct {
	MinConcurrent    int
	MaxConcurrent    int
	Step             int
	NumberOfRequests int
	LocalCSV         string
	URL              string
	RequestBody      string
}

func TestRampedRequest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RampedRequest Suite")
}

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
