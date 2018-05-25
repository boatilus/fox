// Package fox implements a simple wrapper for the Twilio programmatic fax API. It implements all
// the functions associated with the "Faxes" endpoint, but to keep the library tight, does not
// faciliate, for example, E.164 phone number parsing and validation, or handling Twilio status
// callbacks.
//
// To get started, construct a new Client with your Twilio account SID and auth token:
//
// 		c := fox.NewClient("YOUR_TWILIO_ACCOUNT_SID", "YOUR_TWILIO_AUTH_TOKEN")
//
// Optionally, you can also pass a pointer to a SendOptions object to NewClient to specify custom
// send options (to, for example, tell Twilio *not* to store fax media):
//
//   opts := SendOpts{StoreMedia: false}
//   c := fox.NewClient("YOUR_TWILIO_ACCOUNT_SID", "YOUR_TWILIO_AUTH_TOKEN", &opts)
//
// The Get, List and Send methods on the returned Client are used to make the API calls as described
// by Twilio's API reference.
package fox

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var (
	scheme = "https" // the API is always accessed over HTTPS
	host   = "fax.twilio.com"
)

const (
	version  = "v1" // pins this package to API v1
	endpoint = "Faxes"
)

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

type statusType int

const (
	// StatusQueued indicates that the fax is queued, waiting for processing.
	StatusQueued statusType = iota
	// StatusProcessing indicates that the fax is being downloaded, uploaded, or transcoded into a
	// different format.
	StatusProcessing
	// StatusSending indicates that the fax is in the process of being sent.
	StatusSending
	// StatusDelivered indicates that he fax has been successfuly delivered.
	StatusDelivered
	// StatusReceiving indicates that the fax is in the process of being received.
	StatusReceiving
	// StatusReceived indicates that the fax has been successfully received.
	StatusReceived
	// StatusNoAnswer indicates that the outbound fax failed because the other end did not pick up.
	StatusNoAnswer
	// StatusBusy indicates that the outbound fax failed because the other side sent back a busy
	// signal.
	StatusBusy
	// StatusFailed indicates that the fax failed to send or receive.
	StatusFailed
	// StatusCanceled indicates that the fax was canceled, either by using the REST API, or rejected
	// by TwiML.
	StatusCanceled
)

func (st statusType) String() string {
	switch st {
	default:
		return ""
	case StatusQueued:
		return "queued"
	case StatusProcessing:
		return "processing"
	case StatusSending:
		return "sending"
	case StatusDelivered:
		return "delivered"
	case StatusReceiving:
		return "receiving"
	case StatusReceived:
		return "received"
	case StatusNoAnswer:
		return "no-answer"
	case StatusBusy:
		return "busy"
	case StatusFailed:
		return "failed"
	case StatusCanceled:
		return "canceled"
	}
}

// ListOpts describes the options to use when listing faxes.
type ListOpts struct {
	// DateCreatedAfter filters the returned list to only include faxes created after the supplied
	// date.
	DateCreatedAfter time.Time
	// DateCreatedOnOrBefore filters the returned list to only include faxes created on or before the
	// supplied date.
	DateCreatedOnOrBefore time.Time
	// From filters the returned list to only include faxes sent from the supplied number, given in
	// E.164 format.
	From string
	// To filters the returned list to only include faxes sent to the supplied number, given in E.164
	// format.
	To string
}

// urlEncode adds ListOpts fields to a url.Values map using standard param=value URL encoding.
func (lo *ListOpts) urlEncode(data url.Values) {
	if !lo.DateCreatedAfter.IsZero() {
		data.Add("DateCreatedAfter", lo.DateCreatedAfter.Format(time.RFC3339))
	}
	if !lo.DateCreatedOnOrBefore.IsZero() {
		data.Add("DateCreatedOnOrBefore", lo.DateCreatedOnOrBefore.Format(time.RFC3339))
	}
	if lo.From != "" {
		data.Add("From", lo.From)
	}
	if lo.To != "" {
		data.Add("To", lo.To)
	}
}

// SendOpts describes the options to use when sending a fax.
type SendOpts struct {
	// Quality is a quality value, one of QualityStandard, QualityFine or QualitySuperfine.
	Quality qualityType
	// SIPAuthPassword is the password to use for authentication when sending to a SIP address.
	SIPAuthPassword string
	// SIPAuthUsername is the username to use for authentication when sending to a SIP address.
	SIPAuthUsername string
	// StatusCallback is a status callback URL that will receive a GET or POST request when the status
	// of the fax changes.
	StatusCallback string
	// StoreMedia specifies whether or not to store a copy of the sent media on Twilio's servers for
	// later retrieval.
	StoreMedia bool
	// TTLMinutes is the duration, in minutes, from when a fax was initiated should Twilio attempt to
	// send the fax.
	TTLMinutes int
}

// urlEncode adds SendOpts fields to a url.Values map using standard param=value URL encoding.
func (so *SendOpts) urlEncode(data url.Values) {
	data.Add("Quality", so.Quality.String())

	if so.SIPAuthPassword != "" {
		data.Add("SipAuthPassword", so.SIPAuthPassword)
	}
	if so.SIPAuthUsername != "" {
		data.Add("SipAuthUsername", so.SIPAuthUsername)
	}
	if so.StatusCallback != "" {
		data.Add("StatusCallback", so.StatusCallback)
	}

	data.Add("StoreMedia", strconv.FormatBool(so.StoreMedia))

	if so.TTLMinutes > 0 {
		data.Add("Ttl", strconv.FormatInt(int64(so.TTLMinutes), 10))
	}
}

// DefaultSendOpts is the default set of options to use for Client.Send. It mirrors the defaults
// specified by Twilio.
var DefaultSendOpts = &SendOpts{
	Quality:    QualityFine,
	StoreMedia: true,
}

// ErrorResponse describes Twilio's error response.
type ErrorResponse struct {
	// Code is the unique Twilio error code.
	Code int `json:"code"`
	// Message is a descriptive error message.
	Message string `json:"message"`
	// MoreInfo is a link to the Twilio documentation for the error code.
	MoreInfo string `json:"more_info"`
	// Status is the HTTP status code for this error.
	Status int `json:"status"`
}

// Error satisfies the error interface.
func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("fox: error %v (Twilio error %v): %s", err.Status, err.Code, err.Message)
}

// Meta describes the metadata object component of a ListResponse
type Meta struct {
	FirstPageURL    string `json:"first_page_url"`
	Key             string `json:"key"`
	NextPageURL     string `json:"next_page_url,omitempty"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	PreviousPageURL string `json:"previous_page_url,omitempty"`
	URL             string `json:"url"`
}

// ListResponse describes the success response returned from listing faxes.
type ListResponse struct {
	Faxes []SendResponse `json:"faxes"`
	Meta  Meta           `json:"meta"`
}

// SendResponse describes the success response returned from sending a fax.
type SendResponse struct {
	// AccountSid	is the unique SID identifier of the account from which the fax was sent.
	AccountSid string `json:"account_sid"`
	// APIVersion is the API version used to send the fax, which is always "v1".
	APIVersion string `json:"api_version"`
	// Status is the current status of the fax transmission (typically "queued").
	Status string `json:"status"`
	// SID is the 34-character string that uniquely identifies this fax.
	SID string `json:"sid"`
	// URL is the fully-qualified reference URL to the fax resource.
	URL string `json:"url"`
	// Direction is the transmission direction of this fax. Always "outbound".
	Direction string `json:"direction"`
	// To	is the phone number or SIP URI of the destination.
	To string `json:"to"`
	// From is the number the fax was sent from, in E.164 format, or the SIP From display name.
	From string `json:"from"`
	// Quality is one of "standard", "fine" or "superfine".
	Quality string `json:"quality"`
	// DateCreated is the timestamp at which the fax resource was created.
	DateCreated time.Time `json:"date_created"`
	// DateUpdated is the timestamp at which the fax resource was updated.
	DateUpdated time.Time `json:"date_updated"`
	// Links is a dictionary of URL links to nested resources of this fax.
	Links struct {
		// Media is a fully-qualified reference URL to the fax media resource.
		Media string `json:"media"`
	} `json:"links"`
	// MediaSid string `json:"media_sid"`
	// PriceUnit is the currency unit of the Price. E.g., "USD".
	PriceUnit string `json:"price_unit"`
	Price     string `json:"price"`
	// Duration is the time taken to transmit the fax, in seconds.
	Duration int    `json:"duration"`
	NumPages int    `json:"num_pages"`
	MediaURL string `json:"media_url"`
}

// StatusCallbackResponse describes the response received from calling a status callback.
type StatusCallbackResponse struct {
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
