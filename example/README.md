# ocpp1.6-go

Open Charge Point Protocol (OCPP) version 1.6 in Go.

## Example

Example shows implementation ocppj version 1.6 for central system.

### Folder structure
- callbacks.go - Includes handlers for each OCPP request (Implementation DB logic)
- configs.json - File to specify list of chargers for the demo in JSON format
- simplequeue.go - Simple messages queue and charger objects for the demo only
- README.md - this file

#### Docker
To spin contaienr on docker, please use commands below:
```bash
docker build -t ocpp16:latest -f Dockerfile .
docker run --rm --name ocpp16-example -p "9033:8080" ocpp16:latest
```


## Central System Example

To use library in your project, you must implement the callbacks with your business logic, as shown below:

```go
import (
	"time"
	"net/http"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/ocpp16-go/messages"
)

type OCPPHandlers struct {
    // ... Add required variables for your implementation
}

func (cs *OCPPHandlers) BootNotificationRequestHandler (callMessage messages.CallMessage) (string, error, bool) {

    // ... Implement your business logic

	// Create CallResult message
	callMessageResponse := messages.CallResultMessageWithParam (
		callMessage.UniqueID,
		bootNotificationResp.GetPayload(),
	)

	return callMessageResponse.ToString(), nil, WEBSOCKET_KEEP_OPEN
}

// further callbacks... 
```
### Requirements for the design
Name of the methods for the call requests has to in in the format: action + "RequestHandler", as example "HeartbeatRequestHandler".
For the responses handlers have to be in the formar: action + "ResponseHandler", as example "AuthorizeResponseHandler".
The error handler named "OCPPErrorHandler".