package gmailbatching

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResponseBoundary_HappyPath(t *testing.T) {

	val, err := parseResponseBoundary("multipart/mixed; boundary=batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof")

	assert.Equal(t, nil, err)
	assert.Equal(t, "batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof", val)
}

func TestParseResponseBoundary_Error(t *testing.T) {

	val, err := parseResponseBoundary("XXX")

	assert.ErrorContains(t, err, "could not parse response boundary, header: XXX")
	assert.Equal(t, "", val)
}

func TestParseResponseBody(t *testing.T) {
	responseBody := `
        --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof
        Content-Type: application/http
        Content-ID: <response-cfs_gmail_item_1>

        HTTP/1.1 200 OK
        Content-Type: application/json; charset=UTF-8
        Vary: Origin
        Vary: X-Origin
        Vary: Referer

        {"id":"15a4a1151a83d7ab","threadId":"15a4a1151a83d7ab","labelIds":["UNREAD","CATEGORY_SOCIAL","INBOX"],"snippet":"A lot has happened on Facebook since you last logged in. Here are some notifications you&#39;ve missed from your friends. David Davids 2 new notifications You have new notifications. A lot has happened","payload":{"mimeType":"multipart/alternative","headers":[{"name":"Date","value":"Thu, 16 Feb 2017 19:15:33 -0800"},{"name":"Subject","value":"David, you have 2 new notifications"},{"name":"From","value":"Facebook \u003cnotification+zj4ttzsayz=y@facebookmail.com\u003e"}]},"sizeEstimate":17176,"historyId":"28357","internalDate":"1487301333000"}
        --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof
        Content-Type: application/http
        Content-ID: <response-cfs_gmail_item_2>

        HTTP/1.1 200 OK
        Content-Type: application/json; charset=UTF-8
        Vary: Origin
        Vary: X-Origin
        Vary: Referer

        {"id":"19281d0d510ff568","threadId":"19281d0d510ff568","labelIds":["UNREAD","INBOX"],"snippet":"CFS Gmail Cleanup was granted access to your Google account connor.sadler.androidemul@gmail.com If you did not grant access, you should check this activity and secure your account. Check activity You","payload":{"mimeType":"multipart/alternative","headers":[{"name":"Date","value":"Sat, 12 Oct 2024 17:39:58 GMT"},{"name":"Subject","value":"Security alert"},{"name":"From","value":"Google \u003cno-reply@accounts.google.com\u003e"}]},"sizeEstimate":12140,"historyId":"185955","internalDate":"1728754798000"}
        --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof--
    `
	responseBody = fixupMultilineResponseBody(t, responseBody, "batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof")

	result, err := parseResponseBody("batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof", responseBody)

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, result.Size())
	resultItems := result.GetItems()
	assert.Equal(t, 200, resultItems[0].GetHttpStatusCode())
	assert.Equal(t, "response-cfs_gmail_item_1", resultItems[0].GetContentId())
	assert.True(t, strings.HasPrefix(resultItems[0].GetJson(), `{"id":"15a4a1151a83d7ab","threadId":"15a4a1151a83d7ab"`))
	assert.True(t, strings.HasSuffix(resultItems[0].GetJson(), `"internalDate":"1487301333000"}`))
	assert.Equal(t, 200, resultItems[1].GetHttpStatusCode())
	assert.Equal(t, "response-cfs_gmail_item_2", resultItems[1].GetContentId())
	assert.True(t, strings.HasPrefix(resultItems[1].GetJson(), `{"id":"19281d0d510ff568","threadId":"19281d0d510ff568"`))
	assert.True(t, strings.HasSuffix(resultItems[1].GetJson(), `"internalDate":"1728754798000"}`))

	// TODO: More checks!!!
}

// TODO: Parse body where there is an error - for example, a content id line is badly formatted
// TODO: or the boundary separator being wrong

func TestParseResponseBody_ErrorItemWithMultilineJson(t *testing.T) {
	responseBody := `
        --batch_0-znMOV-dw53PPMWPtc3FX3mc7Ym7JHZ
        Content-Type: application/http
        Content-ID: <response-cfs_gmail_item_3>

        HTTP/1.1 404 Not Found
        Vary: Origin
        Vary: X-Origin
        Vary: Referer
        Content-Type: application/json; charset=UTF-8

        {
        "error": {
            "code": 404,
            "message": "Requested entity was not found.",
            "errors": [
            {
                "message": "Requested entity was not found.",
                "domain": "global",
                "reason": "notFound"
            }
            ],
            "status": "NOT_FOUND"
        }
        }

        --batch_0-znMOV-dw53PPMWPtc3FX3mc7Ym7JHZ--
    `
	responseBody = fixupMultilineResponseBody(t, responseBody, "batch_0-znMOV-dw53PPMWPtc3FX3mc7Ym7JHZ")

	result, err := parseResponseBody("batch_0-znMOV-dw53PPMWPtc3FX3mc7Ym7JHZ", responseBody)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, result.Size())
	resultItems := result.GetItems()
	assert.Equal(t, 404, resultItems[0].GetHttpStatusCode())
	assert.Equal(t, "response-cfs_gmail_item_3", resultItems[0].GetContentId())
	assert.True(t, strings.HasPrefix(resultItems[0].GetJson(), "{\n\"error\": {"))
	// Note: the json does not include a trailing newline, as the blank line after it is not included in the json lines, it's a blank line
	assert.True(t, strings.HasSuffix(resultItems[0].GetJson(), "\"status\": \"NOT_FOUND\"\n}\n}"))

}

func TestParseHttpStatusCode_HappyPath(t *testing.T) {
	result, err := parseHttpStatusCode("HTTP/1.1 200 OK")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, result)
}

func TestParseHttpStatusCode_Error_BadInteger(t *testing.T) {
	result, err := parseHttpStatusCode("HTTP/1.1 xxx OK")
	assert.ErrorContains(t, err, "could not parse http status code, line: 'HTTP/1.1 xxx OK'")
	assert.Equal(t, 0, result)
}

func TestParseHttpStatusCode_Error(t *testing.T) {
	result, err := parseHttpStatusCode("YYY")
	assert.ErrorContains(t, err, "could not parse http status code, line: 'YYY'")
	assert.Equal(t, 0, result)
}

func TestParseContentId_HappyPath(t *testing.T) {
	result, err := parseContentId("Content-ID: <response-cfs_gmail_item_2>")
	assert.Equal(t, nil, err)
	assert.Equal(t, "response-cfs_gmail_item_2", result)
}

func TestParseContentId_Error(t *testing.T) {
	result, err := parseContentId("Content-ID: >")
	assert.ErrorContains(t, err, "could not parse content id, line: 'Content-ID: >'")
	assert.Equal(t, "", result)
}
