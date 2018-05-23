package fox

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c *Client
var to, from string

func init() {
	c = NewClient(os.Getenv("ACCOUNT_SID"), os.Getenv("AUTH_TOKEN"))
	to = os.Getenv("TO")
	from = os.Getenv("FROM")
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

func TestClient_BuildURL(t *testing.T) {
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

func TestClient_Get(t *testing.T) {
	got, err := c.Get(faxSID)
	assert.NoError(t, err)
	t.Logf("%#v", got)
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
