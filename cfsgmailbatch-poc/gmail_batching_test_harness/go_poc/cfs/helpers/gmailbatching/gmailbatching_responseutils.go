package gmailbatching

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// parse something like: multipart/mixed; boundary=batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof
func parseResponseBoundary(responseContentType string) (string, error) {
	marker := "boundary="
	idx := strings.Index(responseContentType, marker)
	if idx == -1 {
		return "", fmt.Errorf("could not parse response boundary, header: %v", responseContentType)
	}

	return responseContentType[idx+len(marker):], nil
}

type responseItemBuilder struct {
	httpStatusCode            int
	blankLineCount            int // How many blank lines we have encountered so far in the current batch item block
	contentId                 string
	capturingNextLinesForJson bool // Whether we encountered the first json line so we need to capture all subsequent lines of the batch block to add to the json
	jsonLines                 []string
	errors                    []error
}

func parseResponseBody(responseBoundary string, respBody string) (BatchResult, error) {

	boundarySeparatorPrefix := "--"
	boundarySeparatorLine := boundarySeparatorPrefix + responseBoundary

	batchResult := &batchResult{}
	var currentResponseItemBuilder *responseItemBuilder = nil

	scanner := bufio.NewScanner(strings.NewReader(respBody))
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text() // Do we need strings.TrimSpace here?
		log.Printf("line: %v  (len: %v)", line, len(line))

		if strings.HasPrefix(line, boundarySeparatorPrefix) {

			if !strings.HasPrefix(line, boundarySeparatorLine) {
				log.Fatalf("Error parsing line as it has the wrong boundary string, line: %v", line)
			}

			if strings.HasSuffix(line, "--") {
				// The end of the response
				// Ensure we process any pending item
				if currentResponseItemBuilder != nil {
					batchResult.addToBatchResult(currentResponseItemBuilder)
				}
				break
			} else {
				// Start of a new item
				// Ensure we process any pending item
				if currentResponseItemBuilder != nil {
					batchResult.addToBatchResult(currentResponseItemBuilder)
				}
				currentResponseItemBuilder = &responseItemBuilder{capturingNextLinesForJson: false}

			}

		} else {
			// First line is allowed to be blank
			if firstLine && len(line) == 0 {
				firstLine = false
				continue
			}

			// A 'normal line' inside a batch result
			if currentResponseItemBuilder == nil {
				log.Fatalf("Error parsing line, as there is no current item to build, line: %v", line)
			}
			currentResponseItemBuilder.registerLine(line)
		}

	}

	// Header: Content-Type = [multipart/mixed; boundary=batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof]
	// TODO: parse body of this format:
	//     --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof
	//     Content-Type: application/http
	//     Content-ID: <response-cfs_gmail_item_1>
	//
	//     HTTP/1.1 200 OK
	//     Content-Type: application/json; charset=UTF-8
	//     Vary: Origin
	//     Vary: X-Origin
	//     Vary: Referer
	//
	//     {"id":"15a4a1151a83d7ab","threadId":"15a4a1151a83d7ab","labelIds":["UNREAD","CATEGORY_SOCIAL","INBOX"],"snippet":"A lot has happened on Facebook since you last logged in. Here are some notifications you&#39;ve missed from your friends. David Davids 2 new notifications You have new notifications. A lot has happened","payload":{"mimeType":"multipart/alternative","headers":[{"name":"Date","value":"Thu, 16 Feb 2017 19:15:33 -0800"},{"name":"Subject","value":"David, you have 2 new notifications"},{"name":"From","value":"Facebook \u003cnotification+zj4ttzsayz=y@facebookmail.com\u003e"}]},"sizeEstimate":17176,"historyId":"28357","internalDate":"1487301333000"}
	//     --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof
	//     Content-Type: application/http
	//     Content-ID: <response-cfs_gmail_item_2>
	//
	//     HTTP/1.1 200 OK
	//     Content-Type: application/json; charset=UTF-8
	//     Vary: Origin
	//     Vary: X-Origin
	//     Vary: Referer
	//
	//     {"id":"19281d0d510ff568","threadId":"19281d0d510ff568","labelIds":["UNREAD","INBOX"],"snippet":"CFS Gmail Cleanup was granted access to your Google account connor.sadler.androidemul@gmail.com If you did not grant access, you should check this activity and secure your account. Check activity You","payload":{"mimeType":"multipart/alternative","headers":[{"name":"Date","value":"Sat, 12 Oct 2024 17:39:58 GMT"},{"name":"Subject","value":"Security alert"},{"name":"From","value":"Google \u003cno-reply@accounts.google.com\u003e"}]},"sizeEstimate":12140,"historyId":"185955","internalDate":"1728754798000"}
	//     --batch_2UgdgWCUtu1zQ8O4D_pFmLoGrQG1OQof--

	return batchResult, nil
}

func (rib *responseItemBuilder) registerLine(line string) {
	if len(line) == 0 {
		log.Printf("Parse blank line")
		rib.blankLineCount++
	} else if strings.HasPrefix(line, "HTTP/") {
		log.Printf("Parse HTTP line")
		httpStatusCode, err := parseHttpStatusCode(line)
		if err != nil {
			rib.errors = append(rib.errors, err)
		} else {
			rib.httpStatusCode = httpStatusCode
		}

	} else if strings.HasPrefix(line, "Content-ID") {
		log.Printf("Parse Content-ID line")
		contentId, err := parseContentId(line)
		if err != nil {
			rib.errors = append(rib.errors, err)
		} else {
			rib.contentId = contentId
		}

	} else if strings.HasPrefix(line, "{") {
		log.Printf("Parse json line")
		rib.capturingNextLinesForJson = true
		rib.appendJsonLine(line)
	} else if rib.capturingNextLinesForJson {
		log.Printf("Parse json next line")
		rib.appendJsonLine(line)
	}
}

func (rib *responseItemBuilder) appendJsonLine(line string) {
	rib.jsonLines = append(rib.jsonLines, line)
}

func parseHttpStatusCode(httpLine string) (int, error) {
	// e.g. HTTP/1.1 200 OK
	re := regexp.MustCompile(`^HTTP/1\.1 ([0-9]+) .*$`)
	match := re.FindStringSubmatch(httpLine)
	if len(match) == 0 {
		return 0, fmt.Errorf("could not parse http status code, line: '%v'", httpLine)
	}
	result, err := strconv.Atoi(match[1])
	if err != nil {
		// This should never happen now due to the regex requiring digits in this position
		return 0, fmt.Errorf("could not parse http status code, could not convert '%v' to int, line: '%v'", match[1], httpLine)
	}
	return result, nil
}

func parseContentId(contentIdLine string) (string, error) {
	// e.g. Content-ID: <response-cfs_gmail_item_2>
	re := regexp.MustCompile(`^Content-ID: <(.*)>$`)
	match := re.FindStringSubmatch(contentIdLine)
	if len(match) == 0 {
		return "", fmt.Errorf("could not parse content id, line: '%v'", contentIdLine)
	}
	return match[1], nil
}

type batchResult struct {
	items []batchResultItem
}

type batchResultItem struct {
	httpStatusCode int
	contentId      string
	json           string // may be multiline, usually for errors
}

func (bri *batchResultItem) GetHttpStatusCode() int {
	return bri.httpStatusCode
}

func (bri *batchResultItem) GetContentId() string {
	return bri.contentId
}

// may be multiline, usually for errors
func (bri *batchResultItem) GetJson() string {
	return bri.json
}

func (br *batchResult) Size() int {
	return len(br.items)
}

func (br *batchResult) GetItems() []BatchResultItem {
	result := make([]BatchResultItem, len(br.items))
	for i := 0; i < len(br.items); i++ {
		result[i] = &br.items[i]
	}
	return result
}

func (br *batchResult) addToBatchResult(rib *responseItemBuilder) {
	log.Printf("addToBatchResult, rib has properties: contentId: %v, blankLineCount: %v, json line count: %v", rib.contentId, rib.blankLineCount, len(rib.jsonLines))

	fullJson := strings.Join(rib.jsonLines, "\n")
	br.items = append(br.items, batchResultItem{httpStatusCode: rib.httpStatusCode, contentId: rib.contentId, json: fullJson})
}
