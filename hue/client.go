package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type InvalidIpError struct {
	Ip string
}

func (e *InvalidIpError) Error() string {
	return fmt.Sprintf("%v is not a valid IP address", e.Ip)
}

type MultipleHubsError struct {
	NumberOfHubs int
}

func (e *MultipleHubsError) Error() string {
	return fmt.Sprintf("Expected 1 hub, found %v", e.NumberOfHubs)
}

type ServiceError struct {
	StatusCode int
	Status     string
	Message    string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("Received %v %v: %v", e.StatusCode, e.Status, e.Message)
}

type Client struct {
	client   *http.Client
	endpoint *url.URL
	username string
}

const ipRegexp = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

type MeetHueResp struct {
	Id         string `json:"id"`
	InternalIp string `json:"internalipaddress"`
}

func NewClient(username string) (client *Client, err error) {
	httpClient := &http.Client{}
	r, err := httpClient.Get("https://www.meethue.com/api/nupnp")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
		"host":     "https://www.meethue.com/api/nupnp",
		"method":   "GET",
		"request":  string(body),
	}).Debugf("Recieved %v from %v", string(body), "https://www.meethue.com/api/nupnp")

	resp := []MeetHueResp{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if len(resp) != 1 {
		return nil, &MultipleHubsError{NumberOfHubs: len(resp)}
	}
	hub := resp[0]

	re := regexp.MustCompile(ipRegexp)
	if !re.MatchString(hub.InternalIp) {
		return nil, &InvalidIpError{Ip: hub.InternalIp}
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
	}).Debugf("Found hub %v with IP %v", hub.Id, hub.InternalIp)

	endpoint, err := url.Parse("http://" + hub.InternalIp)
	if err != nil {
		return nil, err
	}

	return &Client{client: &http.Client{}, endpoint: endpoint, username: username}, nil
}

func (c *Client) Do(method string, address string, message []byte, resp interface{}) (err error) {
	// Get the URL for the resource.
	url := c.endpoint
	url.Path = strings.Replace(address, "<username>", c.username, -1)
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
		"host":     url.String(),
		"method":   method,
		"request":  string(message),
	}).Debugf("%v %v %v", method, url.String(), string(message))

	// Create the http request.
	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the http request.
	r, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Get the body content.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
		"host":     url.String(),
		"method":   method,
		"response": string(body),
	}).Debugf("%v %v %v", method, url.String(), string(body))

	// Check for errors.
	if r.StatusCode != 200 {
		return &ServiceError{StatusCode: r.StatusCode, Status: r.Status, Message: string(body)}
	}

	// Return the result.
	return json.Unmarshal(body, resp)
}
