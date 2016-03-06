# lights
--
    import "github.com/drombosky/disco-dance-party/hue/lights"

Package lights is a library for interacting with the lights connected to a
Philips Hue bridge. Commands include getting and setting light attributes as
well as discovering and deleting lights.

## Usage

#### type Client

```go
type Client struct {
}
```

Client represents a client to control lights via the Philips Hue bridge.

#### func  NewClient

```go
func NewClient(hueClient hue.Client) (client *Client)
```
NewClient takes a *hue.Client and returns a client for interacting with lights.

#### func (*Client) Delete

```go
func (c *Client) Delete(id string) (err error)
```
Delete deletes a light from the Philips Hue bridge.

#### func (*Client) Get

```go
func (c *Client) Get(id string) (resp *message.Light, err error)
```
Get gets the attributes and state of a given light.

#### func (*Client) GetAll

```go
func (c *Client) GetAll() (resp map[string]message.Light, err error)
```
GetAll gets a list of all lights that have been discovered by the Philips Hue
bridge.

#### func (*Client) GetNew

```go
func (c *Client) GetNew() (resp *message.GetNewResp, err error)
```
GetNew gets a list of lights that were discovered the last time a search for new
lights was performed. The list of new lights is always deleted when a new search
is started.

#### func (*Client) Rename

```go
func (c *Client) Rename(id, name string) (err error)
```
Rename is used to rename lights. A light can have its name changed when in any
state, including when it is unreachable or off.

#### func (*Client) Set

```go
func (c *Client) Set(id string, state message.NewLightState) (err error)
```
Set allows the user to turn the light on and off, modify the hue and effects.
