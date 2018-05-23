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

func TestSend(t *testing.T) {
	res, err := Send(
		to,
		from,
		"http://unec.edu.az/application/uploads/2014/12/pdf-sample.pdf",
		DefaultSendOpts,
	)

	assert.NoError(t, err)
	t.Logf("%#v", res)
}
