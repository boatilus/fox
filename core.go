// Package fox implements a wrapper around the Twilio fax API.
package fox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
const DefaultTimeoutDuration = 10 * time.Second

// AccountSID is the Twilio account SID, and should be set prior to calling any methods.
var AccountSID string

// AuthToken is the Twilio auth token, and should be set prior to calling any methods.
var AuthToken string

// TimeoutDuration is the length of time to wait for an HTTP request to finish before timing out.
var TimeoutDuration = DefaultTimeoutDuration

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

func (qt qualityType) String() string {
	switch qt {
	default:
		return ""
	case QualityStandard:
		return "standard"
	case QualityFine:
		return "fine"
	case QualitySuperfine:
		return "superfine"
	}
}

// SendOpts describes the options to use when sending a faxes.
type SendOpts struct {
	// From is phone number to use as the caller id, E.164-formatted. If using a phone number, it must
	// be a Twilio number or a verified outgoing caller id for your account. If sending to a SIP
	// address, this can be any alphanumeric string (plus the characters +, _, ., and -) to use in the
	// From header of the SIP request.
	// From string
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

func (so *SendOpts) urlEncode(data url.Values) {
	//data.Add("From", so.From)
	data.Add("Quality", so.Quality.String())
	data.Add("SipAuthPassword", so.SipAuthPassword)
	data.Add("SipAuthUsername", so.SipAuthUsername)

	if so.StatusCallback != "" {
		data.Add("StatusCallback", so.StatusCallback)
	}

	data.Add("StoreMedia", strconv.FormatBool(so.StoreMedia))

	if so.TTL.Minutes() > 0.0 {
		minutes := so.TTL.Nanoseconds() * int64(1000000000)
		data.Add("Ttl", strconv.FormatInt(minutes, 10))
	}
}

// SendResponse describes the response returned from sending a fax.
type SendResponse struct {
	// AccountSid	is the account from which the fax was sent.
	AccountSid string `json:"account_sid"`
	// APIVersion is the API version used to send the fax, which is always "v1".
	APIVersion string `json:"api_version"`
	// Status is the current status of the fax transmission (typically "queued").
	Status string `json:"status"`
	// SID is the Twilio SID for the created fax.
	SID string `json:"sid"`
	// URL is the fully-qualified reference URL to the fax resource.
	URL string `json:"url"`
	// Direction is always "outbound".
	Direction string `json:"direction"`
	// To	is the phone number or SIP URI of the destination.
	To string `json:"to"`
	// From is the caller ID or SIP.
	From string `json:"from"`
	// Quality is one of "standard", "fine" or "superfine".
	Quality string `json:"quality"`
	// DateCreated is the timestamp at which the fax resource was created.
	DateCreated time.Time `json:"date_created"`
	// DateUpdated is the timestamp at which the fax resource was updated.
	DateUpdated time.Time `json:"date_updated"`
	Links       struct {
		// Media is a fully-qualified reference URL to the fax media resource.
		Media string `json:"media"`
	} `json:"links"`
	// The following fields are present in the response from Twilio, but always null:
	// MediaSid string `json:"media_sid"`
	// PriceUnit string `json:"price_unit"``
	// Price string `json:"price"`
	// Duration int `json:"duration"`
	// NumPages int `json:"num_pages"`
	// MediaURL string `json:"media_url"`
}

// StatusCallback describes the data received from calling a status callback.
type StatusCallback struct {
	// FaxSid is the 34-character unique identifier for the fax.
	FaxSid string
	// AccountSid	is the account from which the fax was sent.
	AccountSid string
	// From is the caller ID or SIP.
	From string
	// To	is the phone number or SIP URI of the destination.
	To string
	// RemoteStationID is the called subscriber identification (CSID) reported by the receiving fax
	// machine.
	RemoteStationID string `json:"RemoteStationId"`
	// FaxStatus is the current status of the fax transmission.
	FaxStatus string
	// APIVersion is the API version used to send the fax, which for this API will be "v1".
	APIVersion string `json:"ApiVersion"`
	// OriginalMediaURL is the original URL passed when sending the fax.
	OriginalMediaURL string `json:"OriginalMediaUrl"`
	// NumPages	is the number of pages sent (only if successful).
	NumPages int
	// MediaURL is the media URL on Twilio's servers that can be used to fetch the original media
	// sent. Note that this URL will expire after 2 hours, but a new URL can be fetched from the
	// instance resource.
	MediaURL string `json:"MediaUrl"`
	// ErrorCode is a Twilio error code that gives more information about a failure, if any.
	ErrorCode int
	// ErrorMessage is a detailed message describing a failure, if any.
	ErrorMessage string
}

// DefaultSendOpts is the default set of options to use in Send.
var DefaultSendOpts = &SendOpts{
	Quality:    QualityFine,
	StoreMedia: true,
}

func init() {
	client = http.Client{
		Timeout: TimeoutDuration,
	}
}

// Send initiates a fax to the specified number.
func Send(to, from, mediaURL string, opts *SendOpts) (*SendResponse, error) {
	if AccountSID == "" || AuthToken == "" {
		return nil, ErrNotAuthenticated
	}

	if opts == nil {
		opts = DefaultSendOpts
	}

	u := url.URL{}
	u.Scheme = scheme
	u.Host = host
	u.Path = fmt.Sprintf("%s/%s", version, endpoint)

	data := url.Values{}
	data.Add("To", to)
	data.Add("From", from)
	data.Add("MediaUrl", mediaURL)
	opts.urlEncode(data)

	r, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	r.SetBasicAuth(AccountSID, AuthToken)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		// do the error thing.
	}

	var sr SendResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}

	return &sr, nil
}
