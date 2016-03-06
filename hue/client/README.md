# client
--
    import "github.com/drombosky/disco-dance-party/hue/client"

Package client is a library for interacting with the Philips Hue bridge. The
client provides basic service discovery and authentication for interacting with
the bridge.

## Usage

#### type Client

```go
type Client struct {
}
```

Client represents a client to a Philips Hue bridge.

#### func  NewClient

```go
func NewClient(username string) (client *Client, err error)
```
NewClient returns a client to a Philips Hue bridge given a username.

#### func (*Client) Connect

```go
func (c *Client) Connect() (err error)
```
Connect performs one-time work required to connect to the Philips Hue bridge.
This call is idempotent.

#### func (*Client) Do

```go
func (c *Client) Do(method string, address string, message []byte, resp interface{}) (err error)
```
Do sends a command to the to the Philips Hue bridge on behalf of the configured
user.

#### type InvalidIPError

```go
type InvalidIPError struct {
	IP string
}
```

InvalidIPError represents an error that occurs when an IP address is invalid.

#### func (*InvalidIPError) Error

```go
func (e *InvalidIPError) Error() string
```
Error satisfies the error interface.

#### type MeetHueResp

```go
type MeetHueResp struct {
	ID         string `json:"id"`
	InternalIP string `json:"internalipaddress"`
}
```

MeetHueResp represents

#### type MultipleBridgesError

```go
type MultipleBridgesError struct {
	NumberOfBridges int
}
```

MultipleBridgesError represents an error when there is more than one Philips Hue
bridges is detected.

#### func (*MultipleBridgesError) Error

```go
func (e *MultipleBridgesError) Error() string
```
Error satisfies the error interface.

#### type ServiceError

```go
type ServiceError struct {
	StatusCode int
	Status     string
	Message    string
}
```

ServiceError represents an error returned from the Philips Hue bridge.

#### func (*ServiceError) Error

```go
func (e *ServiceError) Error() string
```
Error satisfies the error interface.
