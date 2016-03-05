# hue
--
    import "github.com/drombosky/disco-dance-party/hue"

Package hue is a library for interacting with the Philips Hue bridge. The client
provides basic service discovery and authentication for interacting with the
bridge.

## Usage

#### type Client

```go
type Client struct {
}
```

Client represents a client to a Hue hub.

#### func  NewClient

```go
func NewClient(username string) (client *Client, err error)
```
NewClient returns a client to a Hue hub given a username.

#### func (*Client) Do

```go
func (c *Client) Do(method string, address string, message []byte, resp interface{}) (err error)
```
Do sends a command to the to the Hue hub on behalf of the configured user.

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

#### type MultipleHubsError

```go
type MultipleHubsError struct {
	NumberOfHubs int
}
```

MultipleHubsError represents an error when there is more than one hub detected.

#### func (*MultipleHubsError) Error

```go
func (e *MultipleHubsError) Error() string
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

ServiceError represents an error returned from the hub.

#### func (*ServiceError) Error

```go
func (e *ServiceError) Error() string
```
Error satisfies the error interface.
