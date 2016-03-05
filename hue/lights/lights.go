// Package lights is a library for interacting with the lights connected to a Philips Hue bridge. Commands include
// getting and setting light attributes as well as discovering and deleting lights.
package lights

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/drombosky/disco-dance-party/hue"
	"github.com/drombosky/disco-dance-party/hue/message"
)

// Client represents a client to control lights via the Philips Hue bridge.
type Client struct {
	client hue.Client
}

// NewClient takes a *hue.Client and returns a client for interacting with lights.
func NewClient(hueClient hue.Client) (client *Client, err error) {
	return &Client{client: hueClient}, nil
}

// GetAll gets a list of all lights that have been discovered by the Philips Hue bridge.
func (c *Client) GetAll() (resp map[string]message.Light, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) GetNew",
	}).Debugf("Get all")
	if err = c.client.Do("GET", "/api/<username>/lights", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetNew gets a list of lights that were discovered the last time a search for new lights was performed. The list of
// new lights is always deleted when a new search is started.
func (c *Client) GetNew() (resp *message.GetNewResp, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) GetNew",
	}).Debugf("Get new")
	if err = c.client.Do("GET", "/api/<username>/lights/new", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get gets the attributes and state of a given light.
func (c *Client) Get(id string) (resp *message.Light, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) Get",
	}).Debugf("Get %v", id)
	if err = c.client.Do("GET", fmt.Sprintf("/api/<username>/lights/%v", id), nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Rename is used to rename lights. A light can have its name changed when in any state, including when it is
// unreachable or off.
func (c *Client) Rename(id, name string) (err error) {
	type Body struct {
		Name string `json:"name"`
	}
	message, err := json.Marshal(Body{Name: name})
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue/light",
		"function": "(c *Client) Rename",
		"request":  string(message),
	}).Debugf("Rename %v to %v", id, name)

	if err = c.client.Do("PUT", fmt.Sprintf("/api/<username>/lights/%v", id), message, nil); err != nil {
		return err
	}
	return nil
}

// Set allows the user to turn the light on and off, modify the hue and effects.
func (c *Client) Set(id string, state message.NewLightState) (err error) {
	message, err := json.Marshal(state)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue/light",
		"function": "(c *Client) SetState",
		"request":  string(message),
	}).Debugf("Set state for %v to %v", id, string(message))

	if err = c.client.Do("PUT", fmt.Sprintf("/api/<username>/lights/%v/state", id), message, nil); err != nil {
		return err
	}
	return nil
}

// Delete deletes a light from the Philips Hue bridge.
func (c *Client) Delete(id string) (err error) {
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue/light",
		"function": "(c *Client) Delete",
	}).Debugf("Delete %v", id)

	if err = c.client.Do("PUT", fmt.Sprintf("/api/<username>/lights/%v/state", id), nil, nil); err != nil {
		return err
	}
	return nil
}
