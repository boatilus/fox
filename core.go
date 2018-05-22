// Package fox implements a wrapper around the Twilio fax API.
package fox

import (
	"net/http"
	"net/url"
	"time"
)

const (
	scheme   = "https" // the API is always accessed over HTTPS
	host     = "fax.twilio.com"
	version  = "v1" // pins this package to API v1
	endpoint = "Faxes"
)

// DefaultTimeoutDuration is the default length of time to wait for an HTTP request to finish before
// timing out.
const DefaultTimeoutDuration time.Time = 10 * time.Second

// AccountSID is the Twilio account SID, and should be set prior to calling any methods.
var AccountSID string

// AuthToken is the Twilio auth token, and should be set prior to calling any methods.
var AuthToken string

// TimeoutDuration is the length of time to wait for an HTTP request to finish before timing out.
var TimeoutDuration time.Time = DefaultTimeoutDuration

var client http.Client

type qualityType int

const (
	// QualityStandard is a low quality (204x98) fax resolution that should be supported by all
	// devices.
	QualityStandard qualityType = iota
	// QualityFine is a medium quality (204x196) fax resolution; this quality boasts wide device
	// support.
	QualityFine
	// QualitySuperfine is a high quality (204x392) fax resolution; this quality my not be supported
	// by many devices.
	QualitySuperfine
)

// SendOptions describes the options to use when sending a faxes.
type SendOptions struct {
	// From is phone number to use as the caller id, E.164-formatted. If using a phone number, it must
	// be a Twilio number or a verified outgoing caller id for your account. If sending to a SIP
	// address, this can be any alphanumeric string (plus the characters +, _, ., and -) to use in the
	// From header of the SIP request.
	From string
	// Quality is a quality value, one of QualityStandard, QualityFine or QualitySuperfine.
	Quality qualityType
	// SipAuthPassword is the password to use for authentication when sending to a SIP address.
	SipAuthPassword string
	// SipAuthUsername is the username to use for authentication when sending to a SIP address.
	SipAuthUsername string
	// StatusCallback is a status callback URL that will receive a POST when the status of the fax
	// changes.
	StatusCallback string
	// StoreMedia specifies whether or not to store a copy of the sent media on Twilio's servers for
	// later retrieval.
	StoreMedia bool
	// TTL is the duration from when a fax was initiated should Twilio attempt to send the fax.
	// Twilio observes only the minutes length component of the duration.
	TTL time.Duration
}

type SendResponse struct {
	// fields
}

// DefaultSendOptions is the default set of options to use in Send.
var DefaultSendOptions = &SendOptions{
	Quality:    QualityFine,
	StoreMedia: true,
}

func init() {
	client = http.Client{
		Timeout: TimeoutDuration,
	}
}

// Send initiates a fax to the specified number.
func Send(to, mediaURL string, opts *SendOptions) (*SendResponse, error) {
	if AccountSID == "" || AuthToken == "" {
		return ErrNotAuthenticated
	}

	if opts == nil {
		opts = DefaultSendOptions
	}

	url := url.URL{}
	url.Scheme = scheme
	url.Host = host
	url.Path = fmt.Sprintf("%s/%s", version, endpoint)

	r, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return err
	}

	res, err := client.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		// do the thing
	}

}
