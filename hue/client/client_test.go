package client

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

// TestError is a made up error that can be used in mocks to ensure proper logic flow.
type TestError struct {
	Message string
}

// Error implements the error interface.
func (e *TestError) Error() string {
	if e.Message == "" {
		return "Test Error"
	}
	return "Test Error " + e.Message
}

type ClientSuite struct{}

var _ = check.Suite(&ClientSuite{})

func (s *ClientSuite) TestInvalidIPError(c *check.C) {
	e := &InvalidIPError{IP: "badIP"}
	c.Assert(e.Error(), check.Equals, "badIP is not a valid IP address")
}

func (s *ClientSuite) TestMultipleBridgesError(c *check.C) {
	e := &MultipleBridgesError{NumberOfBridges: 2}
	c.Assert(e.Error(), check.Equals, "Expected 1 bridge, found 2")
}

func (s *ClientSuite) TestServiceError(c *check.C) {
	e := &ServiceError{
		StatusCode: 500,
		Status:     "status",
		Message:    "message",
	}
	c.Assert(e.Error(), check.Equals, "Received 500 status: message")
}
