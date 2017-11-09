package composeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	BxUsSouthApiBase = "https://composebroker-dashboard-public.mybluemix.net/api/2016-07/"
	BxEuGbApiBase    = "https://composebroker-dashboard-public.eu-gb.mybluemix.net/api/2016-07/"
	BxEuDeApiBase    = "https://composebroker-dashboard-public.eu-de.mybluemix.net/api/2016-07/"
)

// Client is a structure that holds session information for the API
type Client struct {
	// The number of times to retry a failing request if the status code is
	// retryable (e.g. for HTTP 429 or 500)
	Retries int
	// The interval to wait between retries. gorequest does not yet support
	// exponential back-off on retries
	RetryInterval time.Duration
	// RetryStatusCodes is the list of status codes to retry for
	RetryStatusCodes []int

	apiBase       string
	apiToken      string
	logger        *log.Logger
	enableLogging bool
}

// NewClient returns a Client for further interaction with the API
func NewClient(apiToken string, apiBase string) (*Client, error) {
	return &Client{
		apiBase:       apiBase,
		apiToken:      apiToken,
		logger:        log.New(ioutil.Discard, "", 0),
		Retries:       5,
		RetryInterval: 3 * time.Second,
		RetryStatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}, nil
}

// SetLogger can enable or disable http logging to and from the Compose
// API endpoint using the provided io.Writer for the provided client.
func (c *Client) SetLogger(enableLogging bool, logger io.Writer) *Client {
	c.logger = log.New(logger, "[composeapi]", log.LstdFlags)
	c.enableLogging = enableLogging
	return c
}

func (c *Client) newRequest(method, targetURL string) *gorequest.SuperAgent {
	return gorequest.New().
		CustomMethod(method, targetURL).
		Set("Authorization", "Bearer "+c.apiToken).
		Set("Content-type", "application/json").
		SetLogger(c.logger).
		SetDebug(c.enableLogging).
		SetCurlCommand(c.enableLogging).
		Retry(c.Retries, c.RetryInterval, c.RetryStatusCodes...)
}

// Link structure for JSON+HAL links
type Link struct {
	HREF      string `json:"href"`
	Templated bool   `json:"templated"`
}

//Errors struct for parsing error returns
type Errors struct {
	Error map[string][]string `json:"errors,omitempty"`
}

func printJSON(jsontext string) {
	var tempholder map[string]interface{}

	if err := json.Unmarshal([]byte(jsontext), &tempholder); err != nil {
		log.Fatal(err)
	}
	indentedjson, err := json.MarshalIndent(tempholder, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(indentedjson))
}

//SetAPIToken overrides the API token
func (c *Client) SetAPIToken(newtoken string) {
	c.apiToken = newtoken
}

//GetJSON Gets JSON string of content at an endpoint
func (c *Client) getJSON(endpoint string) (string, []error) {
	response, body, errs := c.newRequest("GET", c.apiBase+endpoint).End()

	if response.StatusCode != 200 {
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse error - status code %d - body %s", response.StatusCode, response.Body))
		} else {
			errs = append(errs, fmt.Errorf("%v", myerrors.Error))
		}
	}
	return body, errs
}

func (c *Client) getInstanceEndpointPrefix(deploymentid string) (string, []error) {
	if strings.HasPrefix(deploymentid, "bmix-") {
		deploymentid = strings.SplitN(deploymentid, "-", 4)[3]
	}
	body, errs := c.getJSON("instances/" + deploymentid)
	if errs != nil {
		return "", errs
	}

	instance := Instance{}
	err := json.Unmarshal([]byte(body), &instance)
	if err != nil {
		return "", []error{err}
	}

	return "instances/" + deploymentid + "/deployments/" + instance.ResourceID, nil
}
