
## GMail batching - Go Proof of Concept

This is my Go POC for gmail batching calls.  

### Overview
Rather than making multiple 'GET message' calls, GMail allows us to make multiple GET calls in a batch.  
Unfortunately the Go Gmail SDK does not support batching, so we need to add the batching functionality ourselves.  

Main code requirements are:
1. Construct a batch request packet to send to Google
2. Make the call to Google
3. Parse the response

### Google Gmail Documentation

Gmail batching docs:  
https://developers.google.com/gmail/api/guides/batch#:~:text=The%20Gmail%20API%20supports%20batching,lot%20of%20data%20to%20upload


### Quickstart Steps

See this file for a details/code:
gmail_batching_test_harness/quickstart_gmail_batching_test_harness.go

Basic steps are as follows:

1. Create a gmail helper with:
```
gmh := gmailhelper.NewGmailHelper()
gmh.InitFlags()
gmh.InitCredentialsAndClient()
```

2. Create a client with:
```
client := gmh.GetClient()
```

3. Setup a batching call like this:
```
gbp := gmailbatching.NewGmailBatchingProcessor(client)
gbp.AddToBatch("xxxxxxxxxxxxd7ab")
```

4. Retrieve the batch with:
```
gbp.GetBatch()
```

5. Check the results - see quickstart_gmail_batching_test_harness.go for details

