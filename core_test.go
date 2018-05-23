package fox

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var to, from string

func init() {
	AccountSID = os.Getenv("ACCOUNT_SID")
	AuthToken = os.Getenv("AUTH_TOKEN")
	to = os.Getenv("TO")
	from = os.Getenv("FROM")
}

func TestGet(t *testing.T) {
	got, err := Get("FX46b216dea50b3ee395fd534cc6349f5c")
	assert.NoError(t, err)
	t.Logf("%#v", got)
}

func TestSend(t *testing.T) {
	got, err := Send(
		to,
		from,
		"http://unec.edu.az/application/uploads/2014/12/pdf-sample.pdf",
		DefaultSendOpts,
	)

	assert.NoError(t, err)
	assert.Equal(t, got.Status, "queued")
}
