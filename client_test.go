package fox

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	accountSID  = "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	authToken   = "ATXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	faxSID      = "FXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	faxMediaURL = "https://www.twilio.com/docs/documents/25/justthefaxmaam.pdf"
)

const errorResponseJSON = `{
	"code": 1228,
	"message": "Twilio error message",
	"more_info": "https://url/to/more/info",
	"status": 404
}`

var to, from string

var c *Client

func init() {
	to = os.Getenv("TO")
	from = os.Getenv("FROM")

	envSID := os.Getenv("ACCOUNT_SID")
	envToken := os.Getenv("AUTH_TOKEN")

	// If the ACCOUNT_SID and AUTH_TOKEN environment variables are set, use them to construct
	// a client and use real endpoints for testing.
	if envSID != "" && envToken != "" {
		accountSID = envSID
		authToken = envToken
	}

	c = NewClient(accountSID, authToken)
}

func makeServer(h http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(h)

	u, _ := url.Parse(server.URL)
	scheme = u.Scheme
	host = u.Host

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	c.HTTPClient = &http.Client{Transport: transport}
	return server
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	sid := "SID"
	token := "TOKEN"

	t.Run("WithOpts", func(t *testing.T) {
		got := NewClient(sid, token, &SendOpts{
			Quality: QualitySuperfine,
		})
		assert.Equal(sid, got.accountSID)
		assert.Equal(token, got.authToken)
		assert.Equal(QualitySuperfine, got.SendOpts.Quality)
	})

	t.Run("NoOpts", func(t *testing.T) {
		got := NewClient(sid, token)
		assert.Equal(sid, got.accountSID)
		assert.Equal(token, got.authToken)
		assert.Equal(DefaultSendOpts, got.SendOpts)
	})
}

func TestClient_buildURL(t *testing.T) {
	assert := assert.New(t)

	t.Run("NoParam", func(t *testing.T) {
		want := fmt.Sprintf("%s://%s/%s/%s", scheme, host, version, endpoint)
		got := c.buildURL("").String()
		assert.Equal(want, got)
	})

	t.Run("WithParam", func(t *testing.T) {
		want := fmt.Sprintf("%s://%s/%s/%s/%s", scheme, host, version, endpoint, "PARAM")
		got := c.buildURL("PARAM").String()
		assert.Equal(want, got)
	})
}

func TestClient_do(t *testing.T) {
	assert := assert.New(t)

	t.Run("OK", func(t *testing.T) {
		server := makeServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("OK"))
		}))
		defer server.Close()

		r, err := http.NewRequest(http.MethodGet, server.URL+"/get-success", nil)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		got, err := c.do(r)
		assert.NoError(err)
		assert.True(bytes.Equal([]byte("OK"), got))
	})

	t.Run("Error", func(t *testing.T) {
		server := makeServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ERROR"))
		}))
		defer server.Close()

		r, err := http.NewRequest(http.MethodGet, server.URL+"/get-error", nil)

		_, err = c.do(r)
		assert.Error(err)
	})
}

func TestClient_Get(t *testing.T) {
	assert := assert.New(t)

	t.Run("OK", func(t *testing.T) {
		server := makeServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(`{
				"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"api_version": "v1",
				"date_created": "2015-07-30T20:00:00Z",
				"date_updated": "2015-07-30T20:00:00Z",
				"direction": "outbound",
				"from": "+15017122661",
				"media_url": "https://www.twilio.com/docs/documents/25/justthefaxmaam.pdf",
				"media_sid": "MEXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"num_pages": null,
				"price": null,
				"price_unit": null,
				"quality": "fine",
				"sid": "FXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"status": "queued",
				"to": "+15558675310",
				"duration": null,
				"links": {
					"media": "https://fax.twilio.com/v1/Faxes/FXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Media"
				},
				"url": "https://fax.twilio.com/v1/Faxes/FXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
			}`))
		}))
		defer server.Close()

		got, err := c.Get(faxSID)
		assert.NoError(err)

		if got == nil {
			t.Error("got is nil")
			t.FailNow()
		}

		assert.Equal("queued", got.Status)
		assert.Equal("fine", got.Quality)
	})

	t.Run("Error", func(t *testing.T) {
		server := makeServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorResponseJSON))
		}))
		defer server.Close()

		_, err := c.Get(faxSID)
		assert.Error(err)
	})
}

func TestClient_Send(t *testing.T) {
	got, err := c.Send(
		to,
		from,
		"http://unec.edu.az/application/uploads/2014/12/pdf-sample.pdf",
	)

	assert.NoError(t, err)
	assert.Equal(t, got.Status, "queued")
}
