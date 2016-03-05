// Package hue contains the interfaces for interacting with the Philip Hue bridge.
package hue

import (
	"github.com/drombosky/disco-dance-party/hue/message"
)

// Client represents an interface for interacting with the Philips Hue bridge.
type Client interface {
	Do(method string, address string, message []byte, resp interface{}) (err error)
}

// Lights represents an interface for a client to control lights via the Hue bridge.
type Lights interface {
	// GetAll gets a list of all lights that have been discovered by the Philips Hue bridge.
	GetAll() (resp map[string]message.Light, err error)
	// GetNew gets a list of lights that were discovered the last time a search for new lights was performed. The list of
	// new lights is always deleted when a new search is started.
	GetNew() (resp *message.GetNewResp, err error)
	// Get gets the attributes and state of a given light.
	Get(id string) (resp *message.Light, err error)
	// Rename is used to rename lights. A light can have its name changed when in any state, including when it is
	// unreachable or off.
	Rename(id, name string) (err error)
	// Set allows the user to turn the light on and off, modify the hue and effects.
	Set(id string, state message.NewLightState) (err error)
	// Delete deletes a light from the Philips Hue bridge.
	Delete(id string) (err error)
}
