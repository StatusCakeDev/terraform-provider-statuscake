package statuscake

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTest_Validate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	test := &Test{
		Timeout:      200,
		Confirmation: 100,
		Public:       200,
		Virus:        200,
		TestType:     "FTP",
		RealBrowser:  100,
		TriggerRate:  100,
		CheckRate:    100000,
		CustomHeader: "here be dragons",
		WebsiteName:  "",
		WebsiteURL:   "",
	}

	err := test.Validate()
	require.NotNil(err)

	message := err.Error()
	assert.Contains(message, "WebsiteName is required")
	assert.Contains(message, "WebsiteURL is required")
	assert.Contains(message, "Timeout must be 0 or between 6 and 99")
	assert.Contains(message, "Confirmation must be between 0 and 9")
	assert.Contains(message, "CheckRate must be between 0 and 23999")
	assert.Contains(message, "Public must be 0 or 1")
	assert.Contains(message, "Virus must be 0 or 1")
	assert.Contains(message, "TestType must be HTTP, TCP, or PING")
	assert.Contains(message, "RealBrowser must be 0 or 1")
	assert.Contains(message, "TriggerRate must be between 0 and 59")
	assert.Contains(message, "CustomHeader must be provided as json string")

	test.Timeout = 10
	test.Confirmation = 2
	test.Public = 1
	test.Virus = 1
	test.TestType = "HTTP"
	test.RealBrowser = 1
	test.TriggerRate = 50
	test.CheckRate = 10
	test.WebsiteName = "Foo"
	test.WebsiteURL = "http://example.com"
	test.CustomHeader = `{"test": 15}`
	test.NodeLocations = []string{"foo", "bar"}

	err = test.Validate()
	assert.Nil(err)
}

func TestTest_ToURLValues(t *testing.T) {
	assert := assert.New(t)

	test := &Test{
		TestID:         123,
		Paused:         true,
		WebsiteName:    "Foo Bar",
		CustomHeader:   `{"some":{"json": ["value"]}}`,
		WebsiteURL:     "http://example.com",
		Port:           3000,
		NodeLocations:  []string{"foo", "bar"},
		Timeout:        11,
		PingURL:        "http://example.com/ping",
		Confirmation:   1,
		CheckRate:      500,
		BasicUser:      "myuser",
		BasicPass:      "mypass",
		Public:         1,
		LogoImage:      "http://example.com/logo.jpg",
		Branding:       1,
		WebsiteHost:    "hoster",
		Virus:          1,
		FindString:     "hello",
		DoNotFind:      true,
		TestType:       "HTTP",
		RealBrowser:    1,
		TriggerRate:    50,
		TestTags:       []string{"tag1", "tag2"},
		StatusCodes:    "500",
		EnableSSLAlert: false,
		FollowRedirect: false,
	}

	expected := url.Values{
		"TestID":         {"123"},
		"Paused":         {"1"},
		"WebsiteName":    {"Foo Bar"},
		"WebsiteURL":     {"http://example.com"},
		"CustomHeader":   {`{"some":{"json": ["value"]}}`},
		"Port":           {"3000"},
		"NodeLocations":  {"foo,bar"},
		"Timeout":        {"11"},
		"PingURL":        {"http://example.com/ping"},
		"ContactGroup":   {""},
		"Confirmation":   {"1"},
		"CheckRate":      {"500"},
		"BasicUser":      {"myuser"},
		"BasicPass":      {"mypass"},
		"Public":         {"1"},
		"LogoImage":      {"http://example.com/logo.jpg"},
		"Branding":       {"1"},
		"WebsiteHost":    {"hoster"},
		"Virus":          {"1"},
		"FindString":     {"hello"},
		"DoNotFind":      {"1"},
		"TestType":       {"HTTP"},
		"RealBrowser":    {"1"},
		"TriggerRate":    {"50"},
		"TestTags":       {"tag1,tag2"},
		"StatusCodes":    {"500"},
		"UseJar":         {"0"},
		"PostRaw":        {""},
		"FinalEndpoint":  {""},
		"EnableSSLAlert": {"0"},
		"FollowRedirect": {"0"},
	}

	assert.Equal(expected, test.ToURLValues())

	test.TestID = 0
	delete(expected, "TestID")

	assert.Equal(expected.Encode(), test.ToURLValues().Encode())
}

func TestTests_All(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_all_ok.json",
	}
	tt := newTests(c)
	tests, err := tt.All()
	require.Nil(err)

	assert.Equal("/Tests", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	assert.Len(tests, 2)

	expectedTest := &Test{
		TestID:        100,
		Paused:        false,
		TestType:      "HTTP",
		WebsiteName:   "www 1",
		ContactGroup:  []string{"1"},
		Status:        "Up",
		Uptime:        100,
		NodeLocations: []string{"foo", "bar"},
	}
	assert.Equal(expectedTest, tests[0])

	expectedTest = &Test{
		TestID:        101,
		Paused:        true,
		TestType:      "HTTP",
		WebsiteName:   "www 2",
		ContactGroup:  []string{"2"},
		Status:        "Down",
		Uptime:        0,
		NodeLocations: []string{"foo"},
		TestTags:  	   []string{"test1", "test2"},
	}
	assert.Equal(expectedTest, tests[1])
}

func TestTests_AllWithFilter(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_all_ok.json",
	}

	v := url.Values{}
	v.Set("tags", "test1,test2")
	tt := newTests(c)
	tests, err := tt.AllWithFilter(v)
	require.Nil(err)

	assert.Equal("/Tests", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.NotNil(c.sentRequestValues)
	assert.Len(tests, 1)

	expectedTest := &Test{
		TestID:        101,
		Paused:        true,
		TestType:      "HTTP",
		WebsiteName:   "www 2",
		ContactGroup:  []string{"2"},
		Status:        "Down",
		Uptime:        0,
		NodeLocations: []string{"foo"},
		TestTags:  	   []string{"test1", "test2"},
	}
	assert.Equal(expectedTest, tests[0])
}

func TestTests_Update_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_update_ok.json",
	}
	tt := newTests(c)

	test1 := &Test{
		WebsiteName: "foo",
	}

	test2, err := tt.Update(test1)
	require.Nil(err)

	assert.Equal("/Tests/Update", c.sentRequestPath)
	assert.Equal("PUT", c.sentRequestMethod)
	assert.NotNil(c.sentRequestValues)
	assert.NotNil(test2)

	assert.Equal(1234, test2.TestID)
}

func TestTests_Update_Error(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_update_error.json",
	}
	tt := newTests(c)

	test1 := &Test{
		WebsiteName: "foo",
	}

	test2, err := tt.Update(test1)
	assert.Nil(test2)

	require.NotNil(err)
	assert.IsType(&updateError{}, err)
	assert.Contains(err.Error(), "issue a")
}

func TestTests_Update_ErrorWithSliceOfIssues(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_update_error_slice_of_issues.json",
	}
	tt := newTests(c)

	test1 := &Test{
		WebsiteName: "foo",
	}

	test2, err := tt.Update(test1)
	assert.Nil(test2)

	require.NotNil(err)
	assert.IsType(&updateError{}, err)
	assert.Equal("Required Data is Missing., hello, world", err.Error())
}

func TestTests_Delete_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_delete_ok.json",
	}
	tt := newTests(c)

	err := tt.Delete(1234)
	require.Nil(err)

	assert.Equal("/Tests/Details", c.sentRequestPath)
	assert.Equal("DELETE", c.sentRequestMethod)
	assert.Equal(url.Values{"TestID": {"1234"}}, c.sentRequestValues)
}

func TestTests_Delete_Error(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_delete_error.json",
	}
	tt := newTests(c)

	err := tt.Delete(1234)
	require.NotNil(err)
	assert.Equal("this is an error", err.Error())
}

func TestTests_Detail_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_detail_ok.json",
	}
	tt := newTests(c)

	test, err := tt.Detail(1234)
	require.Nil(err)

	assert.Equal("/Tests/Details", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Equal(url.Values{"TestID": {"1234"}}, c.sentRequestValues)

	assert.Equal(test.TestID, 6735)
	assert.Equal(test.TestType, "HTTP")
	assert.Equal(test.Paused, false)
	assert.Equal(test.WebsiteName, "NL")
	assert.Equal(test.CustomHeader, `{"some":{"json": ["value"]}}`)
	assert.Equal(test.UserAgent, "product/version (comment)")
	assert.Equal(test.ContactGroup, []string{"536"})
	assert.Equal(test.Status, "Up")
	assert.Equal(test.Uptime, 0.0)
	assert.Equal(test.CheckRate, 60)
	assert.Equal(test.Timeout, 40)
	assert.Equal(test.LogoImage, "")
	assert.Equal(test.WebsiteHost, "Various")
	assert.Equal(test.FindString, "")
	assert.Equal(test.DoNotFind, false)
	assert.Equal(test.NodeLocations, []string{"foo", "bar"})
}

type fakeAPIClient struct {
	sentRequestPath   string
	sentRequestMethod string
	sentRequestValues url.Values
	fixture           string
}

func (c *fakeAPIClient) put(path string, v url.Values) (*http.Response, error) {
	return c.all("PUT", path, v)
}

func (c *fakeAPIClient) delete(path string, v url.Values) (*http.Response, error) {
	return c.all("DELETE", path, v)
}

func (c *fakeAPIClient) get(path string, v url.Values) (*http.Response, error) {
	return c.all("GET", path, v)
}

func (c *fakeAPIClient) all(method string, path string, v url.Values) (*http.Response, error) {
	c.sentRequestMethod = method
	c.sentRequestPath = path
	c.sentRequestValues = v

	p := filepath.Join("fixtures", c.fixture)
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	var resp *http.Response
	if len(c.sentRequestValues.Get("tags")) > 0 {
		var storedResponses []Test
		var returnResponses []Test
		byteValue, _ := ioutil.ReadAll(f)
		json.Unmarshal(byteValue, &storedResponses)
		requestedTags := strings.Split(c.sentRequestValues.Get("tags"), ",")

		for _, storedResponse := range storedResponses {
			if len(requestedTags) > len(storedResponse.TestTags) { // if we are requesting more tags than whats stored then there are no matches
				continue
			}

			match := true
			for i, tag := range requestedTags {
				if tag != storedResponse.TestTags[i] {
					match = false
				}
			}

			if match { // we can add it to the response now
				returnResponses = append(returnResponses, storedResponse)
			}
		}

		if len(returnResponses) == 0 {
			return nil, nil
		}

		newByteValue, _ := json.Marshal(&returnResponses)
		resp = &http.Response{
			Body: ioutil.NopCloser(bytes.NewBuffer(newByteValue)),
		}

	} else {
		resp = &http.Response{
			Body: f,
		}
	}

	return resp, nil
}
