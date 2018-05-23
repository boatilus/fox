package fox

import (
	"fmt"
	"net/url"
	"testing"

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

func TestSendOpts_URLEncode(t *testing.T) {
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
