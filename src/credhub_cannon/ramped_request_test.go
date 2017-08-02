package main_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("RampedRequest", func() {
	var (
		testServer *ghttp.Server
		bodyChan   chan []byte
		runnerArgs Args
	)

	Context("when correct arguments are used", func() {
		BeforeEach(func() {
			testServer = ghttp.NewUnstartedServer()
			handler := ghttp.CombineHandlers(
				func(rw http.ResponseWriter, req *http.Request) {
					Expect(req.Host).To(Equal(strings.TrimPrefix(testServer.URL(), "http://")))
				},
				ghttp.RespondWith(http.StatusOK, nil),
			)
			testServer.AppendHandlers(handler)
			testServer.AllowUnhandledRequests = true
			testServer.Start()

			bodyChan = make(chan []byte, 3)

			runnerArgs = Args{
				MinConcurrent:    5,
				MaxConcurrent:    10,
				Step:             5,
				NumberOfRequests: 10,
				LocalCSV:         "/Users/pivotal/workspace/testResults.csv",
				URL:              testServer.URL(),
				RequestBody: `"{
  \"name\": \"test_credential\" ,
  \"type\": \"json\",
  \"value\": {
    \"key\": \"value\",
    \"fancy\": { \"num\": 10 }
  },
  \"overwrite\": true,
}"`,
			}
		})

		JustBeforeEach(func() {
			NewRampedRequest(runnerArgs)

		})

		AfterEach(func() {
			testServer.Close()
			close(bodyChan)
		})

		It("ramps up throughput over multiple tests", func() {
			Expect(testServer.ReceivedRequests()).To(HaveLen(20))
		})

		Context("when local-csv is specified", func() {
			var dir string
			BeforeEach(func() {
				var err error
				dir, err = ioutil.TempDir("", "test")
				Expect(err).NotTo(HaveOccurred())
				runnerArgs.LocalCSV = dir

				header := make(http.Header)
				header.Add("Content-Type", "application/json")

			})
			It("stores the csv locally", func() {

				checkFiles := func() int {
					files, err := ioutil.ReadDir(dir)
					Expect(err).ToNot(HaveOccurred())
					fileCount := 0
					for _, file := range files {
						if strings.Contains(file.Name(), "csv") {
							Expect(file.Size()).ToNot(BeZero())
							fileCount++
						}
					}
					return fileCount
				}
				Eventually(checkFiles).Should(Equal(1))
				Expect(os.RemoveAll(dir)).To(Succeed())
			})
		})

	})

	Context("when incorrect arguments are used", func() {
		Context("when minConccurent is greater than maxConcurent", func() {
			It("panics", func() {
				runnerArgs = Args{
					MinConcurrent:    10,
					MaxConcurrent:    5,
					Step:             5,
					NumberOfRequests: 10,
					LocalCSV:         "/Users/pivotal/workspace",
					URL:              "example.com",
					RequestBody: "",
				}

				panics := func() {
					NewRampedRequest(runnerArgs)
				}

				Expect(panics).To(Panic())
			})
		})

		Context("when concurrency is higher than number of requests", func() {
			It("panics", func() {
				runnerArgs = Args{
					MinConcurrent:    5,
					MaxConcurrent:    10,
					Step:             5,
					NumberOfRequests: 1,
					LocalCSV:         "/Users/pivotal/workspace",
					URL:              "example.com",
					RequestBody: "",
				}

				panics := func() {
					NewRampedRequest(runnerArgs)
				}

				Expect(panics).To(Panic())
			})
		})
	})
})
