package lights

import (
	"fmt"
	"testing"

	"gopkg.in/check.v1"

	"github.com/drombosky/disco-dance-party/hue/message"
	"github.com/drombosky/disco-dance-party/hue/mockHue"
	"github.com/golang/mock/gomock"
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

type LightsSuite struct{}

var _ = check.Suite(&LightsSuite{})

func (s *LightsSuite) TestNewClient(c *check.C) {
	mockHueClient := &mockHue.MockClient{}
	client := NewClient(mockHueClient)
	c.Assert(client, check.DeepEquals, &Client{client: &mockHue.MockClient{}})
}

func (s *LightsSuite) TestGetAll(c *check.C) {
	var resp map[string]message.Light
	var err error

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("GET", "/api/<username>/lights", nil, nil).Return(&TestError{})
	resp, err = client.GetAll()
	c.Assert(err, check.DeepEquals, &TestError{})
	c.Assert(resp, check.IsNil)

	// Success.
	mockHueClient.EXPECT().Do("GET", "/api/<username>/lights", nil, nil).Return(nil)
	resp, err = client.GetAll()
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.IsNil)
}

func (s *LightsSuite) TestGetNew(c *check.C) {
	var resp *message.GetNewResp
	var err error

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("GET", "/api/<username>/lights/new", nil, nil).Return(&TestError{})
	resp, err = client.GetNew()
	c.Assert(err, check.DeepEquals, &TestError{})
	c.Assert(resp, check.IsNil)

	// Success.
	mockHueClient.EXPECT().Do("GET", "/api/<username>/lights/new", nil, nil).Return(nil)
	resp, err = client.GetNew()
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.IsNil)
}

func (s *LightsSuite) TestGet(c *check.C) {
	var resp *message.Light
	var err error

	id := "id"

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("GET", fmt.Sprintf("/api/<username>/lights/%v", id), nil, nil).Return(&TestError{})
	resp, err = client.Get(id)
	c.Assert(err, check.DeepEquals, &TestError{})
	c.Assert(resp, check.IsNil)

	// Success.
	mockHueClient.EXPECT().Do("GET", fmt.Sprintf("/api/<username>/lights/%v", id), nil, nil).Return(nil)
	resp, err = client.Get(id)
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.IsNil)
}

func (s *LightsSuite) TestRename(c *check.C) {
	var err error

	id := "id"
	name := "name"
	message := []byte(`{"name":"name"}`)

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("PUT", fmt.Sprintf("/api/<username>/lights/%v", id), message, nil).Return(&TestError{})
	err = client.Rename(id, name)
	c.Assert(err, check.DeepEquals, &TestError{})

	// Success.
	mockHueClient.EXPECT().Do("PUT", fmt.Sprintf("/api/<username>/lights/%v", id), message, nil).Return(nil)
	err = client.Rename(id, name)
	c.Assert(err, check.IsNil)
}

func (s *LightsSuite) TestSet(c *check.C) {
	var err error

	id := "id"
	newLightState := message.NewLightState{
		BasicState: message.BasicState{
			On:  true,
			Bri: 254,
			Hue: 65535,
			Sat: 128,
		},
	}
	message := []byte(`{"on":true,"bri":254,"hue":65535,"sat":128}`)

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("PUT", fmt.Sprintf("/api/<username>/lights/%v/state", id), message, nil).Return(&TestError{})
	err = client.Set(id, newLightState)
	c.Assert(err, check.DeepEquals, &TestError{})

	// Success.
	mockHueClient.EXPECT().Do("PUT", fmt.Sprintf("/api/<username>/lights/%v/state", id), message, nil).Return(nil)
	err = client.Set(id, newLightState)
	c.Assert(err, check.IsNil)
}

func (s *LightsSuite) TestDelete(c *check.C) {
	var err error

	id := "id"

	mockControl := gomock.NewController(c)
	defer mockControl.Finish()

	mockHueClient := mockHue.NewMockClient(mockControl)
	client := &Client{client: mockHueClient}

	// Error with Philips Hue bridge client.
	mockHueClient.EXPECT().Do("DELETE", fmt.Sprintf("/api/<username>/lights/%v", id), nil, nil).Return(&TestError{})
	err = client.Delete(id)
	c.Assert(err, check.DeepEquals, &TestError{})

	// Success.
	mockHueClient.EXPECT().Do("DELETE", fmt.Sprintf("/api/<username>/lights/%v", id), nil, nil).Return(nil)
	err = client.Delete(id)
	c.Assert(err, check.IsNil)
}
