package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	numRequests := flag.Int("numRequests", 10000, "number of requests per concurrency step")
	requestType := flag.String("requestType", "get", "type of request: get, set, or interpolate")
	url := flag.String("url", "https://localhost:8844", "credhub url")
	minConcurrent := flag.Int("minConcurrent", 1, "minimum number of concurrent requests")
	maxConcurrent := flag.Int("maxConcurrent", 50, "maximum number of concurrent requests")
	step := flag.Int("step", 1, "interval for concurrent requests")
	x509UserCert := flag.String("x509Cert", "", "path to mtls x509 certificate")
	x509UserKey := flag.String("x509Key", "", "path to mtls x509 key")

	//requestType = set or get or interpolate
	//url = https://34.231.67.18:8844
	flag.Parse()

	x509Cert := *x509UserCert
	x509Key := *x509UserKey

	if x509Cert == "" {
		os.Stderr.WriteString("Please enter x509Cert")
		os.Exit(1)
	} else if x509Key == "" {
		os.Stderr.WriteString("Please enter x509Key")
		os.Exit(1)
	}

	fmt.Println("beginning [" + *requestType + "] tests on instance: " + *url)
	fmt.Printf("concurrency from %v to %v, step by %v", *minConcurrent, *maxConcurrent, *step)
	fmt.Printf("%v requests per step", *numRequests)

	rampedRequest := &RampedRequest{*minConcurrent, *maxConcurrent, *step, *numRequests, "", x509Cert, x509Key}

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
	rampedRequest.LocalCSV = "/var/vcap/sys/log/credhub_cannon/setPerfResults-" + time.Now().UTC().Format("20060102150405") + ".csv"
	requestBody := `{
  "name": "/c/p-spring-cloud-services/circuit-breaker/5c9073f9-677b-4eb7-8c95-4b89d66d2890/credential-json",
  "type": "json",
  "value": {
  "dashboard": "https://hystrix-cd5e1b33-2989-414c-8d96-ebe456a0905b.diego.example.com.com",
  "amqp": {
    "uris": [
      "amqp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44/53d2200b-10b3-4e8e-97f6-f3b5c6318e86"
    ],
    "uri": "amqp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44/53d2200b-10b3-4e8e-97f6-f3b5c6318e86",
    "http_api_uri": "https://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@rabbitmq-management.diego.example.com.com/api/",
    "vhost": "53d2200b-10b3-4e8e-97f6-f3b5c6318e86",
    "http_api_uris": [
      "https://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@rabbitmq-management.diego.example.com.com/api/"
    ],
    "ssl": false,
    "dashboard_url": "https://rabbitmq-management.diego.example.com.com/#/login/f962c55f-bc96-4765-b642-0683ab67e9a2/p9fivu9u09v58r8db6r2087pt",
    "password": "p9fivu9u09v58r8db6r2087pt",
    "protocols": {
      "stomp": {
        "uris": [
          "stomp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:61613"
        ],
        "vhost": "53d2200b-10b3-4e8e-97f6-f3b5c6318e86",
        "username": "f962c55f-bc96-4765-b642-0683ab67e9a2",
        "password": "p9fivu9u09v58r8db6r2087pt",
        "port": 61613,
        "host": "10.0.3.44",
        "hosts": [
          "10.0.3.44"
        ],
        "ssl": false,
        "uri": "stomp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:61613"
      },
      "mqtt": {
        "uris": [
          "mqtt://53d2200b-10b3-4e8e-97f6-f3b5c6318e86%3Af962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:1883"
        ],
        "uri": "mqtt://53d2200b-10b3-4e8e-97f6-f3b5c6318e86%3Af962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:1883",
        "ssl": false,
        "hosts": [
          "10.0.3.44"
        ],
        "host": "10.0.3.44",
        "port": 1883,
        "password": "p9fivu9u09v58r8db6r2087pt",
        "username": "53d2200b-10b3-4e8e-97f6-f3b5c6318e86:f962c55f-bc96-4765-b642-0683ab67e9a2"
      },
      "management": {
        "uris": [
          "http://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:15672/api/"
        ],
        "path": "/api/",
        "ssl": false,
        "hosts": [
          "10.0.3.44"
        ],
        "password": "p9fivu9u09v58r8db6r2087pt",
        "username": "f962c55f-bc96-4765-b642-0683ab67e9a2",
        "port": 15672,
        "host": "10.0.3.44",
        "uri": "http://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:15672/api/"
      },
      "amqp": {
        "uris": [
          "amqp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:5672/53d2200b-10b3-4e8e-97f6-f3b5c6318e86"
        ],
        "vhost": "53d2200b-10b3-4e8e-97f6-f3b5c6318e86",
        "username": "f962c55f-bc96-4765-b642-0683ab67e9a2",
        "password": "p9fivu9u09v58r8db6r2087pt",
        "port": 5672,
        "host": "10.0.3.44",
        "hosts": [
          "10.0.3.44"
        ],
        "ssl": false,
        "uri": "amqp://f962c55f-bc96-4765-b642-0683ab67e9a2:p9fivu9u09v58r8db6r2087pt@10.0.3.44:5672/53d2200b-10b3-4e8e-97f6-f3b5c6318e86"
      }
    },
    "username": "f962c55f-bc96-4765-b642-0683ab67e9a2",
    "hostname": "10.0.3.44",
    "hostnames": [
      "10.0.3.44"
    ]
  },
  "stream": "https://turbine-cd5e1b33-2989-414c-8d96-ebe456a0905b.diego.example.com.com"
},
  "overwrite": true
}`
	rampedRequest.FireRequests(url+"/api/v1/data", "PUT", requestBody)
}

func launchGetRequests(rampedRequest *RampedRequest, url string, name string) {
	rampedRequest.LocalCSV = "/var/vcap/sys/log/credhub_cannon/getPerfResults-" + time.Now().UTC().Format("20060102150405") + ".csv"
	rampedRequest.FireRequests(url+"/api/v1/data?name=/c/p-spring-cloud-services/circuit-breaker/5c9073f9-677b-4eb7-8c95-4b89d66d2890/credential-json", "GET", "")
}

func launchInterpolateRequests(rampedRequest *RampedRequest, url string, name string) {
	rampedRequest.LocalCSV = "/var/vcap/sys/log/credhub_cannon/interpolatePerfResults-" + time.Now().UTC().Format("20060102150405") + ".csv"

	requestBody := `{
  "p-circuit-breaker-dashboard": [
    {
      "tags": [
        "circuit-breaker",
        "hystrix-amqp",
        "spring-cloud"
      ],
      "name": "circuit-breaker",
      "plan": "standard",
      "provider": null,
      "label": "p-circuit-breaker-dashboard",
      "volume_mounts": [],
      "syslog_drain_url": null,
      "credentials": {
        "credhub-ref": "((/c/p-spring-cloud-services/circuit-breaker/5c9073f9-677b-4eb7-8c95-4b89d66d2890/credential-json))"
      }
    }
  ],
  "p-service-registry": [
    {
      "tags": [
        "eureka",
        "discovery",
        "registry",
        "spring-cloud"
      ],
      "name": "service-registry",
      "plan": "standard",
      "provider": null,
      "label": "p-service-registry",
      "volume_mounts": [],
      "syslog_drain_url": null,
      "credentials": {
        "credhub-ref": "((/c/p-spring-cloud-services/service-registry/5c9073f9-677b-4eb7-8c95-4b89d66d2891/credential-json))"
      }
    }
  ],
  "p-config-server": [
    {
      "tags": [
        "configuration",
        "spring-cloud"
      ],
      "name": "config-server",
      "plan": "standard",
      "provider": null,
      "label": "p-config-server",
      "volume_mounts": [],
      "syslog_drain_url": null,
      "credentials": {
        "credhub-ref": "((/c/p-spring-cloud-services/config-server/5c9073f9-677b-4eb7-8c95-4b89d66d2892/credential-json))"
      }
    }
  ]
}`
	rampedRequest.FireRequests(url+"/api/v1/interpolate", "POST", requestBody)
}
