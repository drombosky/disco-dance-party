package light

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/drombosky/disco-dance-party/hue"
)

type Client struct {
	client *hue.Client
}

func NewClient(hueClient *hue.Client) (client *Client, err error) {
	return &Client{client: hueClient}, nil
}

type UserState struct {
	On             bool      `json:"on,ompitempty"`
	Bri            int       `json:"bri,omitempty"`
	Hue            int       `json:"hue,omitempty"`
	Sat            int       `json:"sat,omitempty"`
	Xy             []float64 `json:"xy,omitempty"`
	Ct             int       `json:"ct,omitempty"`
	Alert          string    `json:"alert,omitempty"`
	Effect         string    `json:"effect,omitempty"`
	TransitionTime int       `json:"transitiontime,omitempty"`
}

type State struct {
	UserState
	Colormode string `json:"colormode,omitempty"`
	Reachable bool   `json:"reachable,omitempty"`
}

type Light struct {
	State       State             `json:"state,omitempty"`
	Type        string            `json:"type,omitempty"`
	Name        string            `json:"name,omitempty"`
	ModelId     string            `json:"modelid,omitempty"`
	SwVersion   string            `json:"swversion,omitempty"`
	PointSymbol map[string]string `json:"pointsymbol,omitempty"`
}

func (c *Client) GetAll() (resp map[string]Light, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) GetNew",
	}).Debugf("Get all")
	if err = c.client.Do("GET", "/api/<username>/lights", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type GetNewResp struct {
	Lights   map[string]Light
	LastScan string `json:"lastscan"`
}

func (c *Client) GetNew() (resp *GetNewResp, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) GetNew",
	}).Debugf("Get new")
	if err = c.client.Do("GET", "/api/<username>/lights/new", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Get(id int) (resp *Light, err error) {
	log.WithFields(log.Fields{
		"package": "github.com/drombosky/disco-dance-party/hue/light",
		"method":  "(c *Client) Get",
	}).Debugf("Get %v", id)
	if err = c.client.Do("GET", fmt.Sprintf("/api/<username>/lights/%v", id), nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type SetAttributesResp struct {
}

func (c *Client) Rename(id int, name string) (resp *SetAttributesResp, err error) {
	type Body struct {
		Name string `json:"name"`
	}
	message, err := json.Marshal(Body{Name: name})
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue/light",
		"function": "(c *Client) Rename",
		"request":  string(message),
	}).Debugf("Rename %v to %v", id, name)

	if err = c.client.Do("PUT", fmt.Sprintf("/api/<username>/lights/%v", id), message, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type SetStateResp struct {
}

func (c *Client) SetState(id int, state UserState) (resp *SetStateResp, err error) {
	message, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"package":  "github.com/drombosky/disco-dance-party/hue/light",
		"function": "(c *Client) SetState",
		"request":  string(message),
	}).Debugf("Set state for %v to %v", id, string(message))

	if err = c.client.Do("PUT", fmt.Sprintf("/api/<username>/lights/%v/state", id), message, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
