package main

import (
	"flag"
)

func main() {
	numRequests := flag.Int("numRequests", 10000, "number of requests per concurrency step")
	requestType := flag.String("requestType", "get", "type of request: get, set, or interpolate")
	url := flag.String("url", "https://localhost:8844", "credhub url")
	minConcurrent := flag.Int("minConcurrent", 1, "minimum number of concurrent requests")
	maxConcurrent := flag.Int("maxConcurrent", 50, "maximum number of concurrent requests")
	step := flag.Int("step", 1, "interval for concurrent requests")

	//requestType = set or get or interpolate
	//url = https://34.231.67.18:8844
	flag.Parse()

	rampedRequest := &RampedRequest{*minConcurrent, *maxConcurrent, *step, *numRequests, ""}

	const credentialName = "perf-test-json"
	switch *requestType {
	case "set":
		launchSetRequests(rampedRequest, *url, credentialName)
	case "get":
		launchGetRequests(rampedRequest, *url, credentialName)
	case "interpolate":
		launchInterpolateRequests(rampedRequest, *url, credentialName)
	}
}

func launchSetRequests(rampedRequest *RampedRequest, url string, name string) {
	rampedRequest.LocalCSV = "/Users/Pivotal/workspace/setPerfResults.csv"
	requestBody := `{
  "name": "` + name + `",
  "type": "json",
  "value": {
    "key": "value",
    "fancy": { "num": 10 }
  },
  "overwrite": false
}`
	rampedRequest.FireRequests(url+"/api/v1/data", "PUT", requestBody)
}

func launchGetRequests(rampedRequest *RampedRequest, url string, name string) {
	rampedRequest.LocalCSV = "/Users/Pivotal/workspace/getPerfResults.csv"
	rampedRequest.FireRequests(url+"/api/v1/data?name=/"+name, "GET", "")
}

func launchInterpolateRequests(rampedRequest *RampedRequest, url string, name string) {
	rampedRequest.LocalCSV = "/Users/Pivotal/workspace/interpolatePerfResults.csv"

	requestBody := `{
	"pp-config-server": [
	  {
	    "credentials": {
	      "credhub-ref": "((/` + name + `))"
	    },
	    "label": "pp-config-server"
	  }
	]
	}`
	rampedRequest.FireRequests(url+"/api/v1/interpolate", "POST", requestBody)
}
