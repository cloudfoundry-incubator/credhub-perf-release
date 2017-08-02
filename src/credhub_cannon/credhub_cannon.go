package main

import (
	"os"
	"strconv"
)

func main() {
	numRequests, _ := strconv.Atoi(os.Args[1])
	requestType := os.Args[2]
	url := os.Args[3]

	//requestType = set or get or interpolate
	//url = https://34.231.67.18:8844

	rampedRequest := &RampedRequest{1, 15, 1, numRequests, ""}

	const credentialName = "perf-test-json"
	switch requestType {
	case "set":
		launchSetRequests(rampedRequest, url, credentialName)
	case "get":
		launchGetRequests(rampedRequest, url, credentialName)
	case "interpolate":
		launchInterpolateRequests(rampedRequest, url, credentialName)
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
