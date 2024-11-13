package gmailbatching

import (
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

func TestCreateBatchRequestBody_HappyPath(t *testing.T) {

	boundaryString := "myboundarystring"
	msgIds := []string{"msg1", "msg2"}
	body := createBatchRequestBody(msgIds, boundaryString)

	expectedBody := strings.TrimPrefix(dedent.Dedent(`
        --myboundarystring
        Content-Type: application/http
        Content-ID: <cfs_gmail_item_1>

        GET /gmail/v1/users/me/messages/msg1?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false

        --myboundarystring
        Content-Type: application/http
        Content-ID: <cfs_gmail_item_2>

        GET /gmail/v1/users/me/messages/msg2?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false

        --myboundarystring--
        `), "\n") // TrimPrefix with "\n" to remove the leading newline

	assert.Equal(t, expectedBody, body)
}
