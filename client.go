package fox

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// DefaultTimeoutDuration is the default length of time for a Client to wait for an HTTP request to
// complete before timing out.
const DefaultTimeoutDuration = 10 * time.Second

// Client describes an encapsulation of an HTTP client, send options and Twilio credentials.
type Client struct {
	HTTPClient      *http.Client
	TimeoutDuration time.Duration
	SendOpts        *SendOpts
	accountSID      string
	authToken       string
}

// NewClient constructs a new Client given a Twilio account SID, auth token and an optional
// pointer to a SendOpts object. If no argument is supplied for sendOpts, the default send options
// are used.
//
// By default, the HTTP client sets its request timeout duration to DefaultTimeDuration. To
// override, assign a new time.Duration value to HTTPClient.Timeout.
func NewClient(accountSID, authToken string, sendOpts ...*SendOpts) *Client {
	c := Client{
		HTTPClient: &http.Client{
			Timeout: DefaultTimeoutDuration,
		},
		accountSID: accountSID,
		authToken:  authToken,
	}

	if len(sendOpts) > 0 {
		c.SendOpts = sendOpts[0]
	} else {
		c.SendOpts = DefaultSendOpts
	}

	return &c
}

// Get retrieves the data for a single fax instance by its SID, or an error of the type
// ErrorResponse.
func (c *Client) Get(sid string) (*SendResponse, error) {
	if c.accountSID == "" || c.authToken == "" {
		return nil, ErrNotAuthenticated
	}
	if sid == "" {
		return nil, ErrMissingSID
	}

	u := c.buildURL(sid)

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.do(r)
	if err != nil {
		return nil, err
	}

	var sr SendResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}

	return &sr, nil
}

// Send initiates a fax to the specified number. The arguments for the to and from numbers are
// expected to be in the E.164 format, and the media URL argument is expected to be a
// fully-qualified, publicly-accessible URL. It returns the response received from Twilio, or
// an error of the type ErrorResponse.
func (c *Client) Send(to, from, mediaURL string, sendOpts ...*SendOpts) (*SendResponse, error) {
	if c.accountSID == "" || c.authToken == "" {
		return nil, ErrNotAuthenticated
	}
	if to == "" {
		return nil, ErrMissingToNumber
	}
	if from == "" {
		return nil, ErrMissingFromNumber
	}
	if mediaURL == "" {
		return nil, ErrMissingMediaURL
	}

	var opts *SendOpts
	if len(sendOpts) > 0 {
		opts = sendOpts[0]
	} else {
		opts = c.SendOpts
	}

	u := c.buildURL("")

	data := url.Values{}
	data.Add("To", to)
	data.Add("From", from)
	data.Add("MediaUrl", mediaURL)
	opts.urlEncode(data)

	r, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	body, err := c.do(r)
	if err != nil {
		return nil, err
	}

	var sr SendResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}

	return &sr, nil
}

func (c *Client) buildURL(param string) *url.URL {
	u := url.URL{}
	u.Scheme = scheme
	u.Host = host
	u.Path = path.Join(version, endpoint, param)
	return &u
}

// do performs the actual request, setting authentication credentials and returning either a success
// response body as a byte slice or an error of type ErrorResponse.
func (c *Client) do(r *http.Request) ([]byte, error) {
	r.SetBasicAuth(c.accountSID, c.authToken)

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Twilio returns 201 CREATED for fax resources created succesfully via a POST request and 200 OK
	// when retrieving resources via a GET request. All other status codes indicate an error, in which
	// the response body is described by ErrorResponse.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, &errRes
	}

	return body, nil
}
