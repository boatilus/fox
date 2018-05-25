package fox

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse_Error(t *testing.T) {
	in := ErrorResponse{
		Code:    12228,
		Message: "Twilio error message",
		Status:  404,
	}

	want := "fox: error 404 (Twilio error 12228): Twilio error message"
	got := in.Error()
	assert.Equal(t, want, got)
}

func TestListOpts_urlEncode(t *testing.T) {
	in := ListOpts{
		DateCreatedAfter:      time.Now().Add(time.Hour * 4),
		DateCreatedOnOrBefore: time.Now(),
		From: from,
		To:   to,
	}

	data := url.Values{}
	in.urlEncode(data)

	got := data.Encode()
	want := fmt.Sprintf(
		"DateCreatedAfter=%s&DateCreatedOnOrBefore=%s&From=%s&To=%s",
		url.QueryEscape(in.DateCreatedAfter.Format(time.RFC3339)),
		url.QueryEscape(in.DateCreatedOnOrBefore.Format(time.RFC3339)),
		url.QueryEscape(from),
		url.QueryEscape(to),
	)

	assert.Equal(t, want, got)
}

func TestSendOpts_urlEncode(t *testing.T) {
	in := SendOpts{
		Quality:         QualitySuperfine,
		SIPAuthPassword: "password",
		SIPAuthUsername: "username",
		StatusCallback:  "callback",
		StoreMedia:      true,
		TTLMinutes:      10,
	}

	data := url.Values{}
	in.urlEncode(data)

	got := data.Encode()
	want := fmt.Sprintf(
		"Quality=%s&SipAuthPassword=%s&SipAuthUsername=%s&StatusCallback=%s&StoreMedia=%v&Ttl=%v",
		in.Quality.String(),
		in.SIPAuthPassword,
		in.SIPAuthUsername,
		in.StatusCallback,
		in.StoreMedia,
		in.TTLMinutes,
	)

	assert.Equal(t, want, got)
}
