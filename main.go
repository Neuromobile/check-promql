package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type promResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	Result     []map[string]interface{} `json:"result"`
	ResultType string                   `json:"result_type"`
}

var (
	host              = flag.String("host", "", "Prometheus server host")
	port              = flag.String("port", "9090", "Prometheus server port")
	authBasicUser     = flag.String("auth-basic-user", "", "Prometheus server basic user")
	authBasicPassword = flag.String("auth-basic-password", "", "Prometheus server basic password")
	ssl               = flag.Bool("ssl", false, "Prometheus server SSL?")
	query             = flag.String("query", "", "Query that be executed in prometheus")
	critical          = flag.Float64("critical", 0.0, "Critical if value is greater than")
	warning           = flag.Float64("warning", 0.0, "Warning if value is greater than")
	lessThan          = flag.Bool("lt", false, "Change whether value is less than check")
)

const (
	criticalStatus = 2
	warningStatus  = 1
	okStatus       = 0
)

func init() {
	flag.Parse()

	if *query == "" {
		fmt.Println("You need to pass the prometheus query that be executed")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *host == "" {
		fmt.Println("You need to pass the host")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	promURL := "http://"
	if *ssl {
		promURL = "https://"
	}
	if *authBasicUser != "" && *authBasicPassword != "" {
		promURL += *authBasicUser + ":" + *authBasicPassword + "@"
	}
	promURL += *host + ":" + *port + "/api/v1/query?query=" + *query

	URL, err := url.Parse(promURL)
	if err != nil {
		log.Fatal(err)
	}

	queryString, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}

	URL.RawQuery = queryString.Encode()

	resp, err := http.Get(URL.String())
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to query prometheus. Got %d status code\n", resp.StatusCode)
		os.Exit(-1)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	var response promResponse

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Fatal(err)
	}

	if len(response.Data.Result) == 0 {
		fmt.Printf("No data received in the response\n")
		os.Exit(-1)
	}

	status := okStatus
	for _, r := range response.Data.Result {
		if *lessThan {
			switch value, _ := strconv.ParseFloat(r["value"].([]interface{})[1].(string), 64); {
			case value < *critical:
				fmt.Printf("CRITICAL: '%s' value is %.4f and got %.4f\n", r["metric"], *critical, value)
				status = criticalStatus
			case value < *warning:
				fmt.Printf("WARNING: '%s' value is %.4f and got %.4f\n", r["metric"], *warning, value)
				if status < warningStatus {
					status = warningStatus
				}
			}
		} else {
			switch value, _ := strconv.ParseFloat(r["value"].([]interface{})[1].(string), 64); {
			case value >= *critical:
				fmt.Printf("CRITICAL: '%s' value is %.4f and got %.4f\n", r["metric"], *critical, value)
				status = criticalStatus
			case value >= *warning:
				fmt.Printf("WARNING: '%s' value is %.4f and got %.4f\n", r["metric"], *warning, value)
				if status < warningStatus {
					status = warningStatus
				}
			}
		}
	}

	if status == okStatus {
		fmt.Printf("OK: All values OK!")
	}
	os.Exit(status)
}
