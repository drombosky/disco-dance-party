// Package hue is a library for interacting with the Philips Hue bridge. The client provides basic service discovery and
// authentication for interacting with the bridge.
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

// InvalidIPError represents an error that occurs when an IP address is invalid.
type InvalidIPError struct {
	IP string
}

// Error satisfies the error interface.
func (e *InvalidIPError) Error() string {
	return fmt.Sprintf("%v is not a valid IP address", e.IP)
}

// MultipleHubsError represents an error when there is more than one hub detected.
type MultipleHubsError struct {
	NumberOfHubs int
}

// Error satisfies the error interface.
func (e *MultipleHubsError) Error() string {
	return fmt.Sprintf("Expected 1 hub, found %v", e.NumberOfHubs)
}

// ServiceError represents an error returned from the hub.
type ServiceError struct {
	StatusCode int
	Status     string
	Message    string
}

// Error satisfies the error interface.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("Received %v %v: %v", e.StatusCode, e.Status, e.Message)
}

// Client represents a client to a Hue hub.
type Client struct {
	client   *http.Client
	endpoint *url.URL
	username string
}

// ipRegexp is a regular expression for verifying IP addresses.
const ipRegexp = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

// MeetHueResp represents
type MeetHueResp struct {
	ID         string `json:"id"`
	InternalIP string `json:"internalipaddress"`
}

// NewClient returns a client to a Hue hub given a username.
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
	if !re.MatchString(hub.InternalIP) {
		return nil, &InvalidIPError{IP: hub.InternalIP}
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
	}).Debugf("Found hub %v with IP %v", hub.ID, hub.InternalIP)

	endpoint, err := url.Parse("http://" + hub.InternalIP)
	if err != nil {
		return nil, err
	}

	return &Client{client: &http.Client{}, endpoint: endpoint, username: username}, nil
}

// Do sends a command to the to the Hue hub on behalf of the configured user.
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
