package gmailbatchingtestharness

import (
	"fmt"
	"go_poc/cfs/helpers"
	"go_poc/cfs/helpers/gmailbatching"
	"go_poc/cfs/helpers/gmailhelper"
	"log"
	"slices"

	"golang.org/x/exp/maps"
	"google.golang.org/api/gmail/v1"
)

//
// WARNING: Public code
//          Don't include any message ids in this code, just in case
//
// Run with:
//    Not currently supported - code supplied as-is
//

func MainCode() {

	fmt.Println("gmailbatchingtestharness - mainCode")

	gmh := gmailhelper.NewGmailHelper()
	gmh.InitFlags()
	gmh.InitCredentialsAndClient()

	// v0.3 - get message with batching
	client := gmh.GetClient()
	//client := http.DefaultClient
	//cfsutils.InstallLoggingRoundTripper(http.DefaultClient)
	gbp := gmailbatching.NewGmailBatchingProcessor(client)
	gbp.AddToBatch("xxxxxxxxxxxxd7ab")
	gbp.AddToBatch("xxxxxxxxxxxxf568")
	gbp.AddToBatch("xxxxxxxxxxxx1111") // Bad msgId

	gbp.AddToBatch("xxxxxxxxxxxxca02") // extra 1
	gbp.AddToBatch("xxxxxxxxxxxxd7ec")
	gbp.AddToBatch("xxxxxxxxxxxx3599")
	gbp.AddToBatch("xxxxxxxxxxxxc13f")
	gbp.AddToBatch("xxxxxxxxxxxx257f")
	gbp.AddToBatch("xxxxxxxxxxxxd55e")
	gbp.AddToBatch("xxxxxxxxxxxxf3d7")
	gbp.AddToBatch("xxxxxxxxxxxx1a01")
	gbp.AddToBatch("xxxxxxxxxxxx6818")
	gbp.AddToBatch("xxxxxxxxxxxxd0e1") // 10
	gbp.AddToBatch("xxxxxxxxxxxx95c0")
	gbp.AddToBatch("xxxxxxxxxxxx8489")
	gbp.AddToBatch("xxxxxxxxxxxxc482")

	batchResult, err := gbp.GetBatch()
	if err != nil {
		log.Fatal(err)
	}

	log.Println()
	log.Printf("batchResult size: %v", batchResult.Size())
	statusCodeToBri := make(map[int][]gmailbatching.BatchResultItem)
	for _, bri := range batchResult.GetItems() {
		statusCodeToBri[bri.GetHttpStatusCode()] = append(statusCodeToBri[bri.GetHttpStatusCode()], bri)
	}

	httpStatusCodes := make([]int, 0, len(statusCodeToBri))
	httpStatusCodes = append(httpStatusCodes, maps.Keys(statusCodeToBri)...)
	slices.Sort(httpStatusCodes)

	log.Println()
	for _, k := range httpStatusCodes {
		v := statusCodeToBri[k]
		log.Printf("http code: %v - results len: %v", k, len(v))
	}
}

func getSampleMessage(srv *gmail.Service) {
	//msgId := "xxxxxxxxxxxx6ec0" // BT thing
	msgId := "xxxxxxxxxxxxd7ab" // Facebook thing - David, you have 2 new notifications

	msg := gmailhelper.GetMessageFromGmailServiceEx(srv, msgId, true)
	log.Printf("msg: %v", msg)
	log.Printf("  payload body: %v", msg.Payload.Body)

	msgAsJson := helpers.ConvToJsonStringIndent(msg)
	log.Printf("msgAsJson: %v", msgAsJson)
	outputFilename := fmt.Sprintf("go_poc_gmail_batching_test_harness__msgId_%v_output.json", msgId)
	helpers.WriteToFile(outputFilename, msgAsJson)

	log.Printf("message written to output file: %v", outputFilename)
}
