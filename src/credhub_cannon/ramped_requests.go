package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type RampedRequest struct {
	MinConcurrent    int
	MaxConcurrent    int
	Step             int
	NumberOfRequests int
	LocalCSV         string
}

func (rr *RampedRequest) FireRequests(url, httpVerb, requestBody string) {
	if rr.MinConcurrent >= rr.MaxConcurrent {
		panic("MinConcurrent must be greater than MaxConcurrent")
	}

	if rr.NumberOfRequests < rr.MaxConcurrent {
		panic("Can't have less requests than number of concurrent threads")
	}

	rr.runBenchmark(url, httpVerb, requestBody, rr.NumberOfRequests, rr.MinConcurrent, rr.MaxConcurrent, rr.Step, 0)
}

func writeFile(path string, data []byte) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Creating csv file error: %s\n", err)
		os.Exit(1)
	}
	_, err = f.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Writing csv data to a file error: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "csv stored locally in file %s\n", path)
}

func (rr *RampedRequest) runBenchmark(
	url,
	httpVerb,
	requestBody string,
	numRequests,
	lowerConcurrency,
	upperConcurrency,
	concurrencyStep,
	threshold int) {

	benchmarkData := new(bytes.Buffer)
	for i := lowerConcurrency; i <= upperConcurrency; i += concurrencyStep {
		heyData, benchmarkErr := run(url, httpVerb, requestBody, numRequests, i, threshold)
		if benchmarkErr != nil {
			fmt.Fprintf(os.Stderr, "%s\n", benchmarkErr)
			os.Exit(1)
		}

		_, writeErr := benchmarkData.Write(heyData)
		if benchmarkErr != nil {
			fmt.Fprintf(os.Stderr, "Buffer error: %s\n", writeErr)
			os.Exit(1)
		}
	}
	println(benchmarkData.Bytes())
	println(rr.LocalCSV)
	if rr.LocalCSV != "" {
		writeFile(rr.LocalCSV, benchmarkData.Bytes())
	}
}

func run(url, httpVerb, requestBody string, numRequests, concurrentRequests, rateLimit int) ([]byte, error) {
	fmt.Fprintf(os.Stdout, "Running benchmark with %d requests, %d concurrency, and %d rate limit\n", numRequests, concurrentRequests, rateLimit)
	args := []string{
		"-n", strconv.Itoa(numRequests),
		"-c", strconv.Itoa(concurrentRequests),
		"-q", strconv.Itoa(rateLimit),
		"-H", `Authorization: ` + token(),
		"-T", `application/json`,
		"-d", requestBody,
		"-m", httpVerb,
		"-disable-compression",
		"-disable-keepalive",
		"-o", "csv",
		url,
	}

	heyData := exec.Command("./hey", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	heyData.Stdout = &out
	heyData.Stderr = &stderr
	err := heyData.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	return []byte(out.String()), nil
}
func token() (string) {
	response, err := exec.Command("credhub", "--token").Output()
	if err != nil {
		return err.Error()
	}
	return string(response)
}
