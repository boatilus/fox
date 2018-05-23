package fox

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const faxSID = "FX46b216dea50b3ee395fd534cc6349f5c"

const errorResponseJSON = `{
	"code": 1228,
	"message": "Twilio error message",
	"more_info": "https://url/to/more/info",
	"status": 404
}`

const successReponseJSON = `{}`

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
