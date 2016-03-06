// Package client is a library for interacting with the Philips Hue bridge. The client provides basic service discovery
// and authentication for interacting with the bridge.
package client

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

// MultipleBridgesError represents an error when there is more than one Philips Hue bridges is detected.
type MultipleBridgesError struct {
	NumberOfBridges int
}

// Error satisfies the error interface.
func (e *MultipleBridgesError) Error() string {
	return fmt.Sprintf("Expected 1 bridge, found %v", e.NumberOfBridges)
}

// ServiceError represents an error returned from the Philips Hue bridge.
type ServiceError struct {
	StatusCode int
	Status     string
	Message    string
}

// Error satisfies the error interface.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("Received %v %v: %v", e.StatusCode, e.Status, e.Message)
}

// Client represents a client to a Philips Hue bridge.
type Client struct {
	username  string
	connected bool

	client   *http.Client
	endpoint *url.URL
}

// ipRegexp is a regular expression for verifying IP addresses.
const ipRegexp = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

// MeetHueResp represents
type MeetHueResp struct {
	ID         string `json:"id"`
	InternalIP string `json:"internalipaddress"`
}

// NewClient returns a client to a Philips Hue bridge given a username.
func NewClient(username string) (client *Client) {
	return &Client{
		username:  username,
		connected: false,
	}
}

// Connect performs one-time work required to connect to the Philips Hue bridge. This call is idempotent.
func (c *Client) Connect() (err error) {
	// Check to see if the client is already connected. If it is, there is nothing to do.
	if c.connected {
		return nil
	}

	httpClient := &http.Client{}
	r, err := httpClient.Get("https://www.meethue.com/api/nupnp")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
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
		return err
	}
	if len(resp) != 1 {
		return &MultipleBridgesError{NumberOfBridges: len(resp)}
	}
	bridge := resp[0]

	re := regexp.MustCompile(ipRegexp)
	if !re.MatchString(bridge.InternalIP) {
		return &InvalidIPError{IP: bridge.InternalIP}
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue",
		"function": "NewClient",
	}).Debugf("Found bridge %v with IP %v", bridge.ID, bridge.InternalIP)

	endpoint, err := url.Parse("http://" + bridge.InternalIP)
	if err != nil {
		return err
	}

	c.connected = true
	c.client = &http.Client{}
	c.endpoint = endpoint
	return nil
}

// Do sends a command to the to the Philips Hue bridge on behalf of the configured user.
func (c *Client) Do(method string, address string, message []byte, resp interface{}) (err error) {
	// Ensure the client is connected, if not connect it.
	if !c.connected {
		if err = c.Connect(); err != nil {
			return err
		}
	}

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
