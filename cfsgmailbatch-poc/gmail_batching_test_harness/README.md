
# gmail batching - proof of concept

This is a PoC for a client to make batch requests to gmail.  
This allows the caller to retrieve multiple messages (a batch) in a single request.

Initially I used scripts to make the calls - see 'original_script_poc' directory.  
Secondly I moved on to making a Go library which does the same thing.