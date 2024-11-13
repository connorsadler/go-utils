package gmailbatching

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type GmailBatchingProcessor interface {
	AddToBatch(msgId string)
	GetBatch() (BatchResult, error)
}

type gmailBatchingProcessor struct {
	client *http.Client
	msgIds []string
}

type BatchResult interface {
	Size() int
	GetItems() []BatchResultItem
}

type BatchResultItem interface {
	GetHttpStatusCode() int
	GetContentId() string
	GetJson() string // may be multiline, usually for errors
}

func NewGmailBatchingProcessor(client *http.Client) GmailBatchingProcessor {
	return &gmailBatchingProcessor{client, make([]string, 0)}
}

func (gbp *gmailBatchingProcessor) AddToBatch(msgId string) {
	gbp.msgIds = append(gbp.msgIds, msgId)
}

func (gbp *gmailBatchingProcessor) GetBatch() (BatchResult, error) {
	log.Printf(">>> GetBatch")

	log.Printf("gbp.msgIds: %v", gbp.msgIds)

	url := "https://www.googleapis.com/batch/gmail/v1"
	// Use this to test with a proxy, to see the outgoing https request body and headers
	//url := "http://localhost:9090"

	boundaryString := "cfsgbh_boundary"
	contentType := "multipart/mixed; boundary=\"" + boundaryString + "\""
	body := createBatchRequestBody(gbp.msgIds, boundaryString)
	reader := strings.NewReader(body)
	resp, err := gbp.client.Post(url, contentType, reader)
	if err != nil {
		return nil, fmt.Errorf("error in GetBatch: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error in GetBatch: %v", err)
	}

	responseBoundaryString, err := parseResponseBoundary(resp.Header["Content-Type"][0])
	if err != nil {
		return nil, fmt.Errorf("error in GetBatch: %v", err)
	}
	log.Printf("Response body: %v", string(respBody))
	batchResult, err := parseResponseBody(responseBoundaryString, string(respBody))
	if err != nil {
		return nil, fmt.Errorf("error in GetBatch: %v", err)
	}

	log.Printf("<<< GetBatch")

	return batchResult, nil
}
