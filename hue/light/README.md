# light
--
    import "github.com/drombosky/disco-dance-party/hue/light"

Package light is a library for interacting with the lights connected to a
Philips Hue bridge. Commands include getting and setting light attributes as
well as discovering and deleting lights.

## Usage

#### type Client

```go
type Client struct {
}
```

Client represents a client to control lights via the Hue hub.

#### func  NewClient

```go
func NewClient(hueClient *hue.Client) (client *Client, err error)
```
NewClient takes a *hue.Client and returns a client for interacting with lights.

#### func (*Client) Delete

```go
func (c *Client) Delete(id string) (err error)
```
Delete deletes a light from the bridge.

#### func (*Client) Get

```go
func (c *Client) Get(id string) (resp *Light, err error)
```
Get gets the attributes and state of a given light.

#### func (*Client) GetAll

```go
func (c *Client) GetAll() (resp map[string]Light, err error)
```
GetAll gets a list of all lights that have been discovered by the bridge.

#### func (*Client) GetNew

```go
func (c *Client) GetNew() (resp *GetNewResp, err error)
```
GetNew gets a list of lights that were discovered the last time a search for new
lights was performed. The list of new lights is always deleted when a new search
is started.

#### func (*Client) Rename

```go
func (c *Client) Rename(id string, name string) (err error)
```
Rename is used to rename lights. A light can have its name changed when in any
state, including when it is unreachable or off.

#### func (*Client) SetState

```go
func (c *Client) SetState(id string, state NewState) (err error)
```
SetState allows the user to turn the light on and off, modify the hue and
effects.

#### type GetNewResp

```go
type GetNewResp struct {
	// Returns “active” if a scan is currently on-going, “none” if a scan has not been performed since the bridge was
	// powered on, or else the date and time that the last scan was completed in ISO 8601:2004 format
	// (YYYY-MM-DDThh:mm:ss).
	LastScan string `json:"lastscan"`
}
```

GetNewResp represents ...

#### type Light

```go
type Light struct {
	// The state of the light as reported by the Hue hub.
	State State `json:"state,omitempty"`
	// A fixed name describing the type of light e.g. “Extended color light”.
	Type string `json:"type,omitempty"`
	// A unique, editable name given to the light.
	Name string `json:"name,omitempty"`
	// The hardware model of the light.
	ModelID string `json:"modelid,omitempty"`
	// Unique id of the device. The MAC address of the device with a unique endpoint id in the form:
	// AA:BB:CC:DD:EE:FF:00:11-XX
	UniqueID string `json:"uniqueid,omitempty"`
	// The manufacturer name.
	ManifacturerName string
	// Unique ID of the luminaire the light is a part of in the format: AA:BB:CC:DD-XX-YY. AA:BB:, ... represents the hex
	// of the luminaireid, XX the lightsource position (incremental but may contain gaps) and YY the lightpoint position
	// (index of light in luminaire group). A gap in the lightpoint position indicates an incomplete luminaire (light
	// search required to discover missing light points in this case).
	LuminaireUniqueID string
	// An identifier for the software version running on the light.
	SwVersion string `json:"swversion,omitempty"`
	// This parameter is reserved for future functionality. As from 1.11 point symbols are no longer returned.
	PointSymbol map[string]string `json:"pointsymbol,omitempty"`
}
```

Light represents the complete state of a light including the light's state type,
name, model ID, and software version.

#### type NewState

```go
type NewState struct {
	// The basic state provided during sets and returned during gets.
	UserState
	// The duration of the transition from the light’s current state to the new state. This is given as a multiple of
	// 100ms and defaults to 4 (400ms). For example, setting transitiontime:10 will make the transition last 1 second.
	TransitionTime int `json:"transitiontime,omitempty"`
	// Increments or decrements the value of the brightness.bri_inc is ignored if the bri attribute is provided. Any
	// ongoing bri transition is stopped. Setting a value of 0 also stops any ongoing transition. The bridge will return
	// the bri value after the increment is performed.
	BriInc int `json:"bri_inc,omitempty"`
	// Increments or decrements the value of the sat.sat_inc is ignored if the sat attribute is provided. Any ongoing
	// sat transition is stopped. Setting a value of 0 also stops any ongoing transition. The bridge will return the sat
	// value after the increment is performed.
	SatInc int `json:"sat_inc,omitempty"`
	// Increments or decrements the value of the hue. hue_inc is ignored if the hue attribute is provided. Any ongoing
	// color transition is stopped. Setting a value of 0 also stops any ongoing transition. The bridge will return the hue
	// value after the increment is performed.
	// Note if the resulting values are < 0 or > 65535 the result is wrapped. For example {"hue_inc": 1} on a hue value
	// of 65535 results in a hue of 0. {"hue_inc": -2} on a hue value of 0 results in a hue of 65534.
	HueInc int `json:"hue_inc,omitempty"`
	// Increments or decrements the value of the ct. ct_inc is ignored if the ct attribute is provided. Any ongoing color
	// transition is stopped. Setting a value of 0 also stops any ongoing transition. The bridge will return the ct value
	// after the increment is performed.
	CtInc int `json:"ct_inc,omitempty"`
	// Increments or decrements the value of the xy. xy_inc is ignored if the xy attribute is provided. Any ongoing color
	// transition is stopped. Setting a value of 0 also stops any ongoing transition. Will stop at it's gamut boundaries.
	// The bridge will return the xy value after the increment is performed. Max value [0.5, 0.5].
	XyInc [2]float64 `json:"xy_inc,omitempty"`
}
```

NewState represents the new state of the light to be provided to the Hue hub.

#### type State

```go
type State struct {
	// The basic light state provided during sets and returned during gets.
	UserState
	// Indicates the color mode in which the light is working, this is the last command type it received. Values are “hs”
	// for Hue and Saturation, “xy” for XY and “ct” for Color Temperature. This parameter is only present when the light
	// supports at least one of the values.
	Colormode string `json:"colormode,omitempty"`
	// Indicates if a light can be reached by the bridge.
	Reachable bool `json:"reachable,omitempty"`
}
```

State represents the state of the light as reported by the Hue hub.

#### type UserState

```go
type UserState struct {
	// On/Off state of the light. On=true, Off=false
	On bool `json:"on,ompitempty"`
	// Brightness of the light. This is a scale from the minimum brightness the light is capable of, 1, to the maximum
	// capable brightness, 254.
	Bri int `json:"bri,omitempty"`
	// Hue of the light. This is a wrapping value between 0 and 65535. Both 0 and 65535 are red, 25500 is green and 46920
	// is blue.
	Hue int `json:"hue,omitempty"`
	// Saturation of the light. 254 is the most saturated (colored) and 0 is the least saturated (white).
	Sat int `json:"sat,omitempty"`
	// The x and y coordinates of a color in CIE color space. The first entry is the x coordinate and the second entry is
	// the y coordinate. Both x and y are between 0 and 1.
	Xy [2]float64 `json:"xy,omitempty"`
	// The Mired Color temperature of the light. 2012 connected lights are capable of 153 (6500K) to 500 (2000K).
	Ct int `json:"ct,omitempty"`
	// The alert effect, which is a temporary change to the bulb’s state. This can take one of the following values:
	//   “none” – The light is not performing an alert effect.
	//   “select” – The light is performing one breathe cycle.
	//   “lselect” – The light is performing breathe cycles for 15 seconds or until an "alert": "none" command is
	//   received.
	// Note that this contains the last alert sent to the light and not its current state. i.e. After the breathe cycle
	// has finished the bridge does not reset the alert to "none".
	Alert string `json:"alert,omitempty"`
	// The dynamic effect of the light, can either be “none” or “colorloop”. If set to colorloop, the light will cycle
	// through all hues using the current brightness and saturation settings.
	Effect string `json:"effect,omitempty"`
}
```

UserState represents the basic light state provided during sets and returned
during gets.
