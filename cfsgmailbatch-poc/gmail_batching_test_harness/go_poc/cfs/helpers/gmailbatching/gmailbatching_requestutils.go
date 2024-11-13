package gmailbatching

import (
	"fmt"
	"strings"
)

func createBatchRequestBody(msgIds []string, boundaryString string) string {
	// Produces a body like this:
	//     --cfsgbh_boundary
	//     Content-Type: application/http
	//     Content-ID: <cfs_gmail_item_1>
	//
	//     GET /gmail/v1/users/me/messages/15a4a1151a83d7ab?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false
	//
	//     --cfsgbh_boundary
	//     Content-Type: application/http
	//     Content-ID: <cfs_gmail_item_2>
	//
	//     GET /gmail/v1/users/me/messages/19281d0d510ff568?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false
	//
	//     --cfsgbh_boundary--

	var sb strings.Builder

	for i, msgId := range msgIds {
		sb.WriteString("--")
		sb.WriteString(boundaryString)
		sb.WriteString("\n")
		sb.WriteString("Content-Type: application/http\n")
		sb.WriteString(fmt.Sprintf("Content-ID: <cfs_gmail_item_%v>\n", i+1))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("GET /gmail/v1/users/me/messages/%v?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false\n", msgId))
		sb.WriteString("\n")
	}
	sb.WriteString("--")
	sb.WriteString(boundaryString)
	sb.WriteString("--")
	sb.WriteString("\n")

	return sb.String()
}
