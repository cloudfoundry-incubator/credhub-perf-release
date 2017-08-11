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
	X509Cert string
	X509Key string
}

func (rr *RampedRequest) FireRequests(url, httpVerb, requestBody string) {
	if rr.MinConcurrent >= rr.MaxConcurrent {
		panic("MinConcurrent must be greater than MaxConcurrent")
	}

	if rr.NumberOfRequests < rr.MaxConcurrent {
		panic("Can't have less requests than number of concurrent threads")
	}

	rr.runBenchmark(url, httpVerb, requestBody, rr.NumberOfRequests, rr.MinConcurrent, rr.MaxConcurrent, rr.Step, 0, rr.X509Cert, rr.X509Key)
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
	threshold int,
	x509Cert,
	x509Key string,
) {

	benchmarkData := new(bytes.Buffer)
	for i := lowerConcurrency; i <= upperConcurrency; i += concurrencyStep {
		heyData, benchmarkErr := run(url, httpVerb, requestBody, numRequests, i, threshold, x509Cert, x509Key)
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

func run(url, httpVerb, requestBody string, numRequests, concurrentRequests, rateLimit int, x509Cert, x509Key string) ([]byte, error) {
	fmt.Fprintf(os.Stdout, "Running benchmark with %d requests, %d concurrency, and %d rate limit\n", numRequests, concurrentRequests, rateLimit)
	args := []string{
		"-n", strconv.Itoa(numRequests),
		"-c", strconv.Itoa(concurrentRequests),
		"-q", strconv.Itoa(rateLimit),
		"-T", `application/json`,
		"-d", requestBody,
		"-m", httpVerb,
		"-disable-compression",
		"-disable-keepalive",
		"-o", "csv",
		url,
	}

	heyCmd := exec.Command("hey", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	heyCmd.Stdout = &out
	heyCmd.Stderr = &stderr
	x509UserCert := fmt.Sprintf("X509_USER_CERT=%s", x509Cert)
	x509UserKey := fmt.Sprintf("X509_USER_KEY=%s", x509Key)
	existing := []string{x509UserCert, x509UserKey}
	heyCmd.Env = existing
	err := heyCmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	return []byte(out.String()), nil
}
