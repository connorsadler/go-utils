
## cfshttplogging - Test harness

This is a test harness for cfshttplogging, showing how to use that module.  

___   
>TL;DR - copy the code from "cfshttplogging_testharness.go", and remember to "go get" the module.  
___


Brief instructions are simply:

1. go get the module as a dependency for you module.  
Drop to a command line inside the module that wishes to use cfshttplogging, and run:
```
go get github.com/connorsadler/go-utils/cfshttplogging
```

2. Add an import like this:
```
import "github.com/connorsadler/go-utils/cfshttplogging"
```

3. In the setup part of your code, call something like this:
```
cfshttplogging.InstallLoggingRoundTripper(http.DefaultClient)
```

The argument must be a *http.Client, and can be the default one as above.

4. Use that http client to make outgoing http calls, and you'll see logging in stdout, like this:
```
[v0.1 LoggingRoundTripper] >>> Sending request to https://api.sampleapis.com/beers/ale/1
[v0.1 LoggingRoundTripper] >>> Headers: TODO
[v0.1 LoggingRoundTripper] >>> Body: [no body]
[v0.1 LoggingRoundTripper] <<< Received 200 OK response
[v0.1 LoggingRoundTripper] <<< Body: {"price":"$16.99","name":"Founders All Day IPA","rating":{"average":4.411243509154233,"reviews":453},"image":"https://www.totalwine.com/media/sys_master/twmmedia/h00/h94/11891416367134.png","id":1}
Done
2024/10/11 13:05:25 Response body: {"price":"$16.99","name":"Founders All Day IPA","rating":{"average":4.411243509154233,"reviews":453},"image":"https://www.totalwine.com/media/sys_master/twmmedia/h00/h94/11891416367134.png","id":1}
```


### Local Development

If you change the cfshttplogging code locally, use this line in the test harness go.mod:
```
replace github.com/connorsadler/go-utils/cfshttplogging => /XXX/go-utils/cfshttplogging
```
... where XXX is the full path to where you have cloned 'go-utils'

Then your local test harness will use your local (new/changed) code.